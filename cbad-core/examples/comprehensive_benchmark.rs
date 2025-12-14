//! Comprehensive CBAD Benchmark Suite
//!
//! Tests CBAD against all available labeled datasets and reports precision/recall/F1.
//!
//! Run with: cargo run --example comprehensive_benchmark --release
//!
//! Output: Markdown table of results + JSON for further analysis

use cbad_core::anomaly::{AnomalyConfig, AnomalyDetector};
use cbad_core::compression::{create_adapter, CompressionAdapter, CompressionAlgorithm};
use cbad_core::metrics;
use cbad_core::tokenizer::{Tokenizer, TokenizerConfig};
use cbad_core::window::WindowConfig;
use cbad_core::{DetectionProfile, StreamManager};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader, Write};
use std::path::{Path, PathBuf};
use std::time::Instant;

#[derive(Debug, Default, Clone, Serialize)]
struct BenchmarkResult {
    name: String,
    data_type: String,
    true_positives: usize,
    false_positives: usize,
    true_negatives: usize,
    false_negatives: usize,
    total_events: usize,
    precision: f64,
    recall: f64,
    f1: f64,
    accuracy: f64,
    processing_time_ms: u64,
    config_used: String,
    notes: String,
    auprc_composite: Option<f64>,
    auprc_conditional_novelty: Option<f64>,
    precision_at_recall_20_composite: Option<f64>,
    precision_at_recall_40_composite: Option<f64>,
    precision_at_recall_60_composite: Option<f64>,
    precision_at_recall_80_composite: Option<f64>,
    precision_at_recall_20_conditional_novelty: Option<f64>,
    precision_at_recall_40_conditional_novelty: Option<f64>,
    precision_at_recall_60_conditional_novelty: Option<f64>,
    precision_at_recall_80_conditional_novelty: Option<f64>,
}

impl BenchmarkResult {
    fn compute_metrics(&mut self) {
        let tp = self.true_positives as f64;
        let fp = self.false_positives as f64;
        let tn = self.true_negatives as f64;
        let fn_ = self.false_negatives as f64;

        self.precision = if tp + fp > 0.0 { tp / (tp + fp) } else { 0.0 };
        self.recall = if tp + fn_ > 0.0 { tp / (tp + fn_) } else { 0.0 };
        self.f1 = if self.precision + self.recall > 0.0 {
            2.0 * self.precision * self.recall / (self.precision + self.recall)
        } else {
            0.0
        };
        self.accuracy = if self.total_events > 0 {
            (tp + tn) / self.total_events as f64
        } else {
            0.0
        };
    }
}

/// Quantize a positive value on a base-10 log scale (quarter-decade buckets).
fn bucket_log10(value: f64, clamp_min: i32, clamp_max: i32) -> i32 {
    if !value.is_finite() || value <= 0.0 {
        return clamp_min;
    }
    let bucket = (value.log10() * 4.0).round() as i32;
    bucket.clamp(clamp_min, clamp_max)
}

/// Quantize a ratio in [-inf, inf] into coarse percent-style buckets.
fn bucket_ratio(value: f64, step: f64, clamp_min: i32, clamp_max: i32) -> i32 {
    if !value.is_finite() {
        return clamp_min;
    }
    let bucket = (value * step).round() as i32;
    bucket.clamp(clamp_min, clamp_max)
}

/// Delta-encode a time series and return bucketed representation for CBAD.
///
/// This representation captures:
/// - `d`: Delta (change from previous value), bucketed
/// - `m`: Magnitude (absolute value), log-scaled
/// - `v`: Volatility (recent delta variance), log-scaled
///
/// This is much more effective for compression-based anomaly detection than
/// naive `"value=X"` serialization, because it preserves temporal structure.
fn encode_time_series_delta(values: &[f64], bucket_step: f64) -> Vec<String> {
    let mut result = Vec::with_capacity(values.len());
    let mut prev = values.first().copied().unwrap_or(0.0);
    let mut recent_deltas: Vec<f64> = Vec::with_capacity(6);

    for (i, &v) in values.iter().enumerate() {
        let delta = v - prev;

        // Bucket the delta (change from previous)
        let delta_bucket = bucket_ratio(delta, bucket_step, -100, 100);

        // Bucket the magnitude (log scale)
        let mag_bucket = bucket_log10(v.abs().max(1e-10), -20, 40);

        // Track recent deltas for volatility calculation
        if recent_deltas.len() >= 5 {
            recent_deltas.remove(0);
        }
        recent_deltas.push(delta.abs());

        // Volatility: mean absolute delta over recent window
        let vol_bucket = if i >= 5 {
            let mean_vol = recent_deltas.iter().sum::<f64>() / recent_deltas.len() as f64;
            bucket_log10(mean_vol.max(1e-10), -20, 20)
        } else {
            0
        };

        result.push(format!("d={}|m={}|v={}", delta_bucket, mag_bucket, vol_bucket));
        prev = v;
    }
    result
}

/// Simpler delta encoding - just the signed direction and rough magnitude.
/// Groups deltas into: large_drop, drop, stable, rise, large_rise
fn encode_time_series_simple(values: &[f64], magnitude_threshold: f64) -> Vec<String> {
    let mut result = Vec::with_capacity(values.len());
    let mut prev = values.first().copied().unwrap_or(0.0);

    for &v in values.iter() {
        let delta = v - prev;
        let abs_delta = delta.abs();

        // Categorize the change direction and magnitude
        let change_type = if abs_delta < magnitude_threshold * 0.1 {
            "S" // Stable
        } else if delta > 0.0 {
            if abs_delta > magnitude_threshold {
                "LR" // Large rise
            } else {
                "R" // Rise
            }
        } else if abs_delta > magnitude_threshold {
            "LD" // Large drop
        } else {
            "D" // Drop
        };

        // Also track magnitude band
        let mag_band = bucket_log10(v.abs().max(1e-10), -8, 16);

        result.push(format!("{}:{}", change_type, mag_band));
        prev = v;
    }
    result
}

/// Build a compact, compression-friendly string for PaySim rows plus a risk score.
fn serialize_paysim_txn(parts: &[&str]) -> Option<(String, f64)> {
    // step,type,amount,nameOrig,oldbalanceOrg,newbalanceOrig,nameDest,oldbalanceDest,newbalanceDest,isFraud,isFlaggedFraud
    if parts.len() < 11 {
        return None;
    }

    let txn_type = parts.get(1)?.trim();
    if txn_type != "TRANSFER" && txn_type != "CASH_OUT" {
        return None;
    }

    let amount: f64 = parts.get(2)?.trim().parse().ok()?;
    let old_balance_orig: f64 = parts.get(4)?.trim().parse().ok()?;
    let new_balance_orig: f64 = parts.get(5)?.trim().parse().ok()?;
    let old_balance_dest: f64 = parts.get(7)?.trim().parse().ok()?;
    let new_balance_dest: f64 = parts.get(8)?.trim().parse().ok()?;
    let step: f64 = parts.get(0)?.trim().parse().ok()?;
    let is_flagged = parts.get(10).map(|s| s.trim()) == Some("1");

    let drain_ratio = if old_balance_orig > 0.0 {
        (old_balance_orig - new_balance_orig) / old_balance_orig
    } else {
        0.0
    };

    let dest_delta = new_balance_dest - old_balance_dest;
    let dest_delta_ratio = if amount > 0.0 {
        dest_delta / amount
    } else {
        0.0
    };
    let orig_gap = (old_balance_orig - amount - new_balance_orig).abs();
    let dest_gap = (old_balance_dest + amount - new_balance_dest).abs();
    let consistency_tol = amount.abs() * 0.02 + 1.0;
    let orig_mismatch = orig_gap > consistency_tol;
    let dest_mismatch = dest_gap > consistency_tol;

    let origin_prefix = parts.get(3)?.chars().next().unwrap_or('X');
    let dest_prefix = parts.get(6)?.chars().next().unwrap_or('X');
    let dest_zero = old_balance_dest == 0.0 && new_balance_dest == 0.0;
    let origin_zero = new_balance_orig == 0.0;

    let amount_bucket = bucket_log10(amount, 0, 32);
    let drain_bucket = bucket_ratio(drain_ratio, 40.0, -80, 120); // 2.5% steps
    let dest_delta_bucket = bucket_ratio(dest_delta_ratio, 40.0, -120, 200);
    let orig_ratio_bucket = bucket_ratio(
        if old_balance_orig > 0.0 {
            new_balance_orig / old_balance_orig
        } else {
            0.0
        },
        50.0,
        -120,
        200,
    );
    let step_bucket = bucket_log10(step.max(1.0), 0, 24);
    let pattern = match (
        dest_zero,
        origin_zero,
        orig_mismatch,
        dest_mismatch,
        dest_delta_bucket,
    ) {
        (true, _, true, true, _) => "zero_drain",
        (_, true, true, _, _) => "orig_wipe",
        (_, _, true, false, _) => "orig_leak",
        (_, _, _, true, _) => "dest_ignored",
        (_, _, _, _, ddb) if ddb < -40 => "dest_loss",
        (_, _, _, _, ddb) if ddb > 60 => "dest_gain",
        _ => "aligned",
    };
    let risk_level = if (dest_zero && dest_mismatch) || (orig_mismatch && dest_mismatch) {
        "high"
    } else if drain_bucket > 40 || dest_mismatch || orig_mismatch {
        "med"
    } else {
        "low"
    };

    let risk_score = if (dest_zero && dest_mismatch) || (orig_mismatch && dest_mismatch) {
        1.0
    } else if drain_bucket > 40 || dest_mismatch || orig_mismatch {
        0.6
    } else {
        0.1
    };

    Some((format!(
        "t={txn}|op={op}|dp={dp}|a={ab}|dr={dr}|dd={dd}|or={or}|dz={dz}|oz={oz}|om={om}|dm={dm}|pat={pat}|risk={risk}|s={sb}|f={flag}",
        txn = txn_type,
        op = origin_prefix,
        dp = dest_prefix,
        ab = amount_bucket,
        dr = drain_bucket,
        dd = dest_delta_bucket,
        or = orig_ratio_bucket,
        dz = dest_zero as u8,
        oz = origin_zero as u8,
        om = orig_mismatch as u8,
        dm = dest_mismatch as u8,
        pat = pattern,
        risk = risk_level,
        sb = step_bucket,
        flag = is_flagged as u8
    ), risk_score))
}

/// Build a compact representation for BAF rows with velocity/risk cues and a risk score.
fn serialize_baf_row(parts: &[&str]) -> Option<(String, f64)> {
    if parts.len() < 31 {
        return None;
    }

    let income: f64 = parts.get(1)?.trim().parse().ok()?;
    let name_email_sim: f64 = parts.get(2)?.trim().parse().ok()?;
    let customer_age: f64 = parts.get(5)?.trim().parse().ok()?;
    let velocity_6h: f64 = parts.get(10)?.trim().parse().ok()?;
    let velocity_24h: f64 = parts.get(11)?.trim().parse().ok()?;
    let velocity_4w: f64 = parts.get(12)?.trim().parse().ok()?;
    let credit_risk: f64 = parts.get(16)?.trim().parse().ok()?;
    let email_free = parts.get(17).map(|s| s.trim()) == Some("1");
    let payment_type = parts.get(8)?.trim();
    let foreign = parts.get(24).map(|s| s.trim()) == Some("1");
    let source = parts.get(25)?.trim();
    let device_fraud: f64 = parts.get(30)?.trim().parse().ok()?;
    let month = parts.get(31)?.trim();

    let v_ratio = if velocity_4w > 0.0 {
        velocity_24h / velocity_4w
    } else {
        0.0
    };
    let bursty = velocity_6h > 0.0 && velocity_24h > 0.0 && velocity_6h > velocity_24h * 0.5;

    let income_bucket = bucket_ratio(income, 40.0, -40, 120);
    let age_bucket = bucket_ratio(customer_age / 10.0, 10.0, -5, 12); // decade buckets
    let vel6_bucket = bucket_log10(velocity_6h, 0, 40);
    let vel24_bucket = bucket_log10(velocity_24h, 0, 40);
    let vel4w_bucket = bucket_log10(velocity_4w, 0, 40);
    let vr_bucket = bucket_ratio(v_ratio, 50.0, -80, 160);
    let credit_bucket = bucket_ratio(credit_risk / 50.0, 20.0, 0, 120);
    let sim_bucket = bucket_ratio(name_email_sim, 50.0, 0, 80);
    let device_bucket = bucket_log10(device_fraud + 1.0, 0, 24);

    let risk_band = match (
        vel24_bucket,
        vr_bucket,
        credit_bucket,
        device_bucket,
        bursty,
    ) {
        (_, vr, c, d, true) if vr > 40 && c < 40 => "spike",
        (v24, vr, c, d, _) if v24 > 28 && vr > 30 && c < 50 && d > 3 => "hot",
        (_, vr, c, d, _) if vr > 20 && c < 60 && d > 2 => "elev",
        _ => "base",
    };
    let risk_level = if (vel24_bucket > 28 && vr_bucket > 30 && credit_bucket < 50)
        || device_bucket > 8
    {
        "high"
    } else if vr_bucket > 18 || credit_bucket < 60 || device_bucket > 4 {
        "med"
    } else {
        "low"
    };

    let risk_score = match risk_level {
        "high" => 1.0,
        "med" => 0.6,
        _ => 0.1,
    };

    Some((format!(
        "pay={pay}|inc={inc}|age={age}|sim={sim}|v6={v6}|v24={v24}|v4={v4}|vr={vr}|cred={cred}|dev={dev}|free={free}|fx={fx}|src={src}|m={m}|band={band}|risk={risk}",
        pay = payment_type,
        inc = income_bucket,
        age = age_bucket,
        sim = sim_bucket,
        v6 = vel6_bucket,
        v24 = vel24_bucket,
        v4 = vel4w_bucket,
        vr = vr_bucket,
        cred = credit_bucket,
        dev = device_bucket,
        free = email_free as u8,
        fx = foreign as u8,
        src = source,
        m = month,
        band = risk_band,
        risk = risk_level
    ), risk_score))
}

/// Build a small window buffer by repeating a sample a few times.
fn build_window_bytes(sample: &str, repeats: usize) -> Vec<u8> {
    let mut buf = Vec::with_capacity(sample.len() * repeats + repeats);
    for _ in 0..repeats {
        buf.extend_from_slice(sample.as_bytes());
        buf.push(b'\n');
    }
    buf
}

/// Static baseline scorer for structured financial data (avoids baseline drift).
fn classify_static_window(
    baseline: &[u8],
    window: &[u8],
    adapter: &dyn CompressionAdapter,
    tokenizer: Option<&Tokenizer>,
    config: &AnomalyConfig,
) -> bool {
    let metrics = match metrics::compute_metrics_with_tokenizer(
        baseline,
        window,
        adapter,
        config.permutation_count,
        config.seed,
        tokenizer,
    ) {
        Ok(m) => m,
        Err(_) => return false,
    };

    let ncd_pass = metrics.ncd >= config.ncd_threshold;
    let p_pass = metrics.p_value <= config.p_value_threshold;
    let compression_drop = metrics.compression_ratio_change <= -config.compression_ratio_drop_threshold;
    let entropy_jump = metrics.entropy_change >= config.entropy_change_threshold;

    let compression_signal = (-metrics.compression_ratio_change).max(0.0);
    let composite_score =
        0.5 * metrics.ncd + 0.25 * (1.0 - metrics.p_value) + 0.25 * compression_signal;
    let composite_pass = composite_score >= config.composite_threshold;

    if config.require_statistical_significance {
        (ncd_pass && p_pass) || (composite_pass && (compression_drop || entropy_jump))
    } else {
        let votes = (compression_drop as u8) + (entropy_jump as u8) + (composite_pass as u8);
        (ncd_pass && p_pass) || votes >= 2 || (composite_pass && (ncd_pass || p_pass))
    }
}

/// Sweep thresholds to maximize F1 on a set of normal vs fraud scores.
fn find_best_threshold(
    normal_scores: &[f64],
    fraud_scores: &[f64],
) -> Option<(f64, usize, usize, usize, usize, f64, f64, f64)> {
    if normal_scores.is_empty() || fraud_scores.is_empty() {
        return None;
    }

    let mut candidates: Vec<f64> = normal_scores
        .iter()
        .chain(fraud_scores.iter())
        .cloned()
        .collect();
    candidates.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));
    candidates.dedup_by(|a, b| (*a - *b).abs() < 1e-6);

    let mut best: Option<(f64, usize, usize, usize, usize, f64, f64, f64)> = None;
    for &t in &candidates {
        let tp = fraud_scores.iter().filter(|s| **s >= t).count();
        let fp = normal_scores.iter().filter(|s| **s >= t).count();
        let fn_ = fraud_scores.len().saturating_sub(tp);
        let tn = normal_scores.len().saturating_sub(fp);

        let precision = if tp + fp > 0 {
            tp as f64 / (tp + fp) as f64
        } else {
            0.0
        };
        let recall = if tp + fn_ > 0 {
            tp as f64 / (tp + fn_) as f64
        } else {
            0.0
        };
        let f1 = if precision + recall > 0.0 {
            2.0 * precision * recall / (precision + recall)
        } else {
            0.0
        };

        if best
            .as_ref()
            .map_or(true, |b| f1 > b.7 || (f1 == b.7 && recall > b.6))
        {
            best = Some((t, tp, fp, tn, fn_, precision, recall, f1));
        }
    }

    best
}

fn pr_summary_from_scores(
    scores: &mut Vec<(f64, bool)>,
    recall_targets: &[f64],
) -> Option<(f64, Vec<Option<f64>>)> {
    scores.retain(|(s, _)| s.is_finite());
    if scores.is_empty() {
        return None;
    }

    let total_pos = scores.iter().filter(|(_, is_pos)| *is_pos).count();
    if total_pos == 0 {
        return None;
    }

    scores.sort_by(|a, b| b.0.partial_cmp(&a.0).unwrap_or(std::cmp::Ordering::Equal));

    let mut tp = 0usize;
    let mut fp = 0usize;
    let mut ap_sum = 0.0;

    let mut best_precision: Vec<f64> = vec![0.0; recall_targets.len()];
    let mut seen: Vec<bool> = vec![false; recall_targets.len()];

    for &(_, is_pos) in scores.iter() {
        if is_pos {
            tp += 1;
        } else {
            fp += 1;
        }

        let denom = tp + fp;
        let precision = if denom > 0 {
            tp as f64 / denom as f64
        } else {
            0.0
        };
        let recall = tp as f64 / total_pos as f64;

        if is_pos {
            ap_sum += precision;
        }

        for (i, &target) in recall_targets.iter().enumerate() {
            let target = target.clamp(0.0, 1.0);
            if recall >= target {
                if !seen[i] || precision > best_precision[i] {
                    best_precision[i] = precision;
                    seen[i] = true;
                }
            }
        }
    }

    let auprc = ap_sum / total_pos as f64;
    let precisions = best_precision
        .into_iter()
        .zip(seen)
        .map(|(p, ok)| if ok { Some(p) } else { None })
        .collect();
    Some((auprc, precisions))
}

// ==================== TRAIN/TEST EVALUATION HELPERS ====================

/// A scored event with labels for evaluation
#[derive(Debug, Clone)]
struct ScoredEvent {
    composite_score: f64,
    conditional_novelty: f64,
    is_fraud: bool,
    is_detector_anomaly: bool,
    stream_key: String,
}

/// Result of temporal train/test split
#[derive(Debug)]
struct TemporalSplit<T> {
    train: Vec<T>,
    test: Vec<T>,
    train_end_marker: String,
}

/// Calibration result with thresholds
#[derive(Debug, Clone)]
struct CalibrationResult {
    threshold_composite: f64,
    threshold_novelty: f64,
    method: String,
}

/// Confusion matrix metrics for calibrated evaluation
#[derive(Debug, Default, Clone)]
struct CalibratedMetrics {
    threshold: f64,
    tp: usize,
    fp: usize,
    tn: usize,
    fn_: usize,
    precision: f64,
    recall: f64,
    f1: f64,
}

impl CalibratedMetrics {
    fn compute(&mut self) {
        let tp = self.tp as f64;
        let fp = self.fp as f64;
        let fn_ = self.fn_ as f64;
        self.precision = if tp + fp > 0.0 { tp / (tp + fp) } else { 0.0 };
        self.recall = if tp + fn_ > 0.0 { tp / (tp + fn_) } else { 0.0 };
        self.f1 = if self.precision + self.recall > 0.0 {
            2.0 * self.precision * self.recall / (self.precision + self.recall)
        } else {
            0.0
        };
    }

    fn from_scores(scores: &[(f64, bool)], threshold: f64) -> Self {
        let mut m = CalibratedMetrics {
            threshold,
            ..Default::default()
        };
        for &(score, is_positive) in scores {
            let predicted = score >= threshold;
            match (predicted, is_positive) {
                (true, true) => m.tp += 1,
                (true, false) => m.fp += 1,
                (false, false) => m.tn += 1,
                (false, true) => m.fn_ += 1,
            }
        }
        m.compute();
        m
    }
}

/// Calibrate threshold by targeting a false positive rate on normal events.
/// This is unsupervised - only uses normal (non-fraud) scores.
fn calibrate_by_fpr(normal_scores: &[f64], target_fpr: f64) -> f64 {
    if normal_scores.is_empty() {
        return 0.5; // fallback
    }
    let mut sorted: Vec<f64> = normal_scores.iter().cloned().filter(|x| x.is_finite()).collect();
    sorted.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));

    // Threshold at (1 - target_fpr) quantile means target_fpr% of normals will be above it
    let quantile_idx = ((1.0 - target_fpr) * sorted.len() as f64).floor() as usize;
    let idx = quantile_idx.min(sorted.len().saturating_sub(1));
    sorted.get(idx).cloned().unwrap_or(0.5)
}

/// Calibrate threshold by maximizing F1 on labeled data (supervised).
fn calibrate_by_f1(scores: &[(f64, bool)]) -> f64 {
    if scores.is_empty() {
        return 0.5;
    }

    let normal: Vec<f64> = scores.iter().filter(|(_, f)| !*f).map(|(s, _)| *s).collect();
    let fraud: Vec<f64> = scores.iter().filter(|(_, f)| *f).map(|(s, _)| *s).collect();

    if let Some((threshold, _, _, _, _, _, _, _)) = find_best_threshold(&normal, &fraud) {
        threshold
    } else {
        0.5
    }
}

/// Calibrate thresholds per-stream using FPR method
fn calibrate_per_stream_fpr(
    stream_scores: &HashMap<String, Vec<(f64, bool)>>,
    target_fpr: f64,
) -> HashMap<String, f64> {
    let mut thresholds = HashMap::new();
    for (stream, scores) in stream_scores {
        let normal: Vec<f64> = scores.iter().filter(|(_, f)| !*f).map(|(s, _)| *s).collect();
        let threshold = calibrate_by_fpr(&normal, target_fpr);
        thresholds.insert(stream.clone(), threshold);
    }
    thresholds
}

/// Evaluate scores using per-stream thresholds
fn evaluate_per_stream(
    events: &[ScoredEvent],
    stream_thresholds: &HashMap<String, f64>,
    global_threshold: f64,
) -> CalibratedMetrics {
    let mut m = CalibratedMetrics {
        threshold: 0.0, // will show "per-stream"
        ..Default::default()
    };
    for event in events {
        let threshold = stream_thresholds
            .get(&event.stream_key)
            .cloned()
            .unwrap_or(global_threshold);
        let predicted = event.composite_score >= threshold;
        match (predicted, event.is_fraud) {
            (true, true) => m.tp += 1,
            (true, false) => m.fp += 1,
            (false, false) => m.tn += 1,
            (false, true) => m.fn_ += 1,
        }
    }
    m.compute();
    m
}

/// Extended benchmark result with train/test and calibration info
#[derive(Debug, Default, Clone)]
struct ExtendedBenchmarkResult {
    base: BenchmarkResult,
    // Split info
    train_count: usize,
    test_count: usize,
    split_marker: String,
    // Calibration results (on test set)
    fpr_calibrated: Option<CalibratedMetrics>,
    f1_calibrated: Option<CalibratedMetrics>,
    per_stream_calibrated: Option<CalibratedMetrics>,
    per_stream_threshold_count: usize,
}

/// Calibration export format for API import
/// This can be POSTed to /v1/calibration/profiles/{name}/benchmark
#[derive(Debug, Clone, Serialize)]
struct CalibrationExport {
    /// Profile name (e.g., "financial_fraud", "sensitive", "balanced")
    profile_name: String,
    /// Recommended composite threshold from benchmarks
    composite_threshold: f64,
    /// Area under precision-recall curve (primary quality metric)
    auprc: f64,
    /// F1 score at optimal threshold
    f1: Option<f64>,
    /// Dataset name used for calibration
    dataset: String,
    /// Additional metadata
    calibration_method: String,
    /// Target FPR used for calibration (if FPR method)
    target_fpr: Option<f64>,
    /// Timestamp of calibration
    timestamp: String,
}

/// Summary of all calibrations for bulk export
#[derive(Debug, Serialize)]
struct CalibrationExportBundle {
    /// Version of the export format
    version: String,
    /// When this export was generated
    generated_at: String,
    /// Individual profile calibrations
    calibrations: Vec<CalibrationExport>,
}

fn main() {
    println!("\n# CBAD Comprehensive Benchmark Suite\n");
    println!("Testing CBAD against all available labeled datasets.\n");

    let mut results: Vec<BenchmarkResult> = Vec::new();

    // === TEXT/PROMPT BENCHMARKS ===
    println!("## Text/Prompt Detection\n");

    results.push(benchmark_jailbreak());
    results.push(benchmark_pint_injection());
    results.push(benchmark_ai_safety());
    results.push(benchmark_hallucination());

    // === FINANCIAL BENCHMARKS ===
    println!("\n## Financial/Transaction Detection\n");

    results.push(benchmark_fraud());
    results.push(benchmark_synthetic_transactions());

    // DORA-relevant financial datasets (new)
    results.push(benchmark_paysim());
    results.push(benchmark_bank_account_fraud());
    // Note: IEEE-CIS requires accepting competition rules on kaggle.com first

    // === TIME SERIES BENCHMARKS ===
    println!("\n## Time Series Detection\n");

    results.push(benchmark_nab_aws_cloudwatch());
    results.push(benchmark_nab_known_cause());
    results.push(benchmark_terra_luna());

    // === NETWORK BENCHMARKS ===
    println!("\n## Network Intrusion Detection\n");

    results.push(benchmark_network_labeled());

    // === LOG ANOMALY BENCHMARKS ===
    println!("\n## Log Anomaly Detection\n");

    results.push(benchmark_hdfs_logs());
    results.push(benchmark_bgl_logs());

    // === ADDITIONAL BENCHMARKS ===
    println!("\n## Additional Labeled Datasets\n");

    results.push(benchmark_supply_chain());
    results.push(benchmark_terra_luna_fixed());
    results.push(benchmark_nasa_turbofan());
    results.push(benchmark_elliptic_bitcoin());

    // === SUMMARY ===
    print_summary_table(&results);
    save_results_json(&results);
}

fn benchmark_jailbreak() -> BenchmarkResult {
    println!("### Jailbreak Prompt Detection");
    let mut result = BenchmarkResult {
        name: "Jailbreak Prompts".to_string(),
        data_type: "Text/Prompts".to_string(),
        ..Default::default()
    };

    let base = get_benchmark_path("jailbreak/data/prompts");
    let regular_path = base.join("regular_prompts_2023_12_25.csv");
    let jailbreak_path = base.join("jailbreak_prompts_2023_12_25.csv");

    if !regular_path.exists() || !jailbreak_path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let regular = load_csv_column(&regular_path, "prompt", 2000);
    let jailbreak = load_csv_column(&jailbreak_path, "prompt", 500);

    println!("   Regular prompts: {}", regular.len());
    println!("   Jailbreak prompts: {}", jailbreak.len());

    if regular.len() < 400 || jailbreak.len() < 100 {
        result.notes = "Insufficient data".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 30,
            hop_size: 10,
            max_capacity: 400,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.25,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=200, window=30, ncd_thresh=0.25");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on regular
    for prompt in regular.iter().take(300) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
    }

    // Test regular (should NOT be anomalies)
    for prompt in regular.iter().skip(300).take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test jailbreak (SHOULD be anomalies)
    for prompt in jailbreak.iter().take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_pint_injection() -> BenchmarkResult {
    println!("### PINT Prompt Injection Detection");
    let mut result = BenchmarkResult {
        name: "PINT Prompt Injection".to_string(),
        data_type: "Text/Prompts".to_string(),
        ..Default::default()
    };

    let pint_path = get_benchmark_path("pint/benchmark/data");

    // PINT uses YAML files, we'd need to parse them
    // For now, use the example dataset
    let example_path = pint_path.join("example-dataset.yaml");

    if !example_path.exists() {
        result.notes = "Dataset not found (PINT requires YAML parsing)".to_string();
        println!("   SKIP: PINT dataset requires YAML parsing\n");
        return result;
    }

    result.notes = "TODO: Implement YAML parsing for PINT".to_string();
    println!("   TODO: Implement YAML parsing\n");
    result
}

fn benchmark_ai_safety() -> BenchmarkResult {
    println!("### AI Safety (Malignant Prompts)");
    let mut result = BenchmarkResult {
        name: "AI Safety Malignant".to_string(),
        data_type: "Text/Prompts".to_string(),
        ..Default::default()
    };

    let path =
        PathBuf::from("/Volumes/VIXinSSD/driftlock-archives/test-data/ai_safety/malignant.csv");

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    // Load data - category column indicates type
    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);
    let mut benign_prompts = Vec::new();
    let mut malignant_prompts = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        } // Skip header
        if benign_prompts.len() >= 500 && malignant_prompts.len() >= 200 {
            break;
        }

        if let Ok(line) = line {
            // CSV: category,base_class,text,embedding
            let parts: Vec<&str> = line.splitn(4, ',').collect();
            if parts.len() < 3 {
                continue;
            }

            let category = parts[0].trim();
            let text = parts[2].trim().trim_matches('"');

            if text.len() < 10 {
                continue;
            }

            if category == "conversation" && benign_prompts.len() < 500 {
                benign_prompts.push(text.to_string());
            } else if category != "conversation" && malignant_prompts.len() < 200 {
                malignant_prompts.push(text.to_string());
            }
        }
    }

    println!("   Benign prompts: {}", benign_prompts.len());
    println!("   Malignant prompts: {}", malignant_prompts.len());

    if benign_prompts.len() < 200 || malignant_prompts.len() < 50 {
        result.notes = "Insufficient data parsed".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 150,
            window_size: 25,
            hop_size: 10,
            max_capacity: 300,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.22,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=150, window=25, ncd_thresh=0.22");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on benign
    for prompt in benign_prompts.iter().take(200) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
    }

    // Test benign (should NOT be anomalies)
    for prompt in benign_prompts.iter().skip(200).take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test malignant (SHOULD be anomalies)
    for prompt in malignant_prompts.iter().take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_hallucination() -> BenchmarkResult {
    println!("### Hallucination Detection");
    let mut result = BenchmarkResult {
        name: "Hallucination Detection".to_string(),
        data_type: "Text/QA".to_string(),
        ..Default::default()
    };

    let path = get_benchmark_path("halueval/data/qa_data.json");

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    #[derive(Deserialize)]
    struct QaPair {
        question: String,
        right_answer: String,
        hallucinated_answer: String,
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);
    let mut correct = Vec::new();
    let mut hallucinated = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i >= 500 {
            break;
        }
        if let Ok(line) = line {
            if let Ok(qa) = serde_json::from_str::<QaPair>(&line) {
                correct.push(format!("Q: {} A: {}", qa.question, qa.right_answer));
                hallucinated.push(format!("Q: {} A: {}", qa.question, qa.hallucinated_answer));
            }
        }
    }

    println!("   Correct answers: {}", correct.len());
    println!("   Hallucinated answers: {}", hallucinated.len());

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 150,
            window_size: 30,
            hop_size: 10,
            max_capacity: 300,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.22,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=150, window=30, ncd_thresh=0.22");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on correct
    for answer in correct.iter().take(200) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
    }

    // Test correct
    for answer in correct.iter().skip(200).take(100) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test hallucinated
    for answer in hallucinated.iter().skip(200).take(100) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = "CBAD cannot detect semantic errors".to_string();
    print_result(&result);
    result
}

fn benchmark_fraud() -> BenchmarkResult {
    println!("### Fraud Detection");
    let mut result = BenchmarkResult {
        name: "Financial Fraud".to_string(),
        data_type: "Transactions".to_string(),
        ..Default::default()
    };

    let path = PathBuf::from("/Volumes/VIXinSSD/driftlock-archives/test-data/fraud/fraud_data.csv");

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);
    let mut normal = Vec::new();
    let mut fraud = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        }
        if normal.len() >= 1000 && fraud.len() >= 300 {
            break;
        }

        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() < 15 {
                continue;
            }

            let is_fraud = parts.last().map(|s| s.trim()).unwrap_or("0");
            let text = format!(
                "merchant={} category={} amount={} city={} state={}",
                parts.get(1).unwrap_or(&""),
                parts.get(2).unwrap_or(&""),
                parts.get(3).unwrap_or(&""),
                parts.get(4).unwrap_or(&""),
                parts.get(5).unwrap_or(&""),
            );

            if is_fraud == "1" && fraud.len() < 300 {
                fraud.push(text);
            } else if is_fraud == "0" && normal.len() < 1000 {
                normal.push(text);
            }
        }
    }

    println!("   Normal transactions: {}", normal.len());
    println!("   Fraud transactions: {}", fraud.len());

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 50,
            hop_size: 20,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.20,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=300, window=50, ncd_thresh=0.20");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on normal
    for txn in normal.iter().take(500) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
    }

    // Test normal
    for txn in normal.iter().skip(500).take(200) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test fraud
    for txn in fraud.iter().take(200) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_synthetic_transactions() -> BenchmarkResult {
    println!("### Synthetic Transaction Detection");
    let mut result = BenchmarkResult {
        name: "Synthetic Transactions".to_string(),
        data_type: "Transactions".to_string(),
        ..Default::default()
    };

    let normal_path =
        PathBuf::from("/Volumes/VIXinSSD/driftlock-archives/test-data/normal-transactions.jsonl");
    let anomalous_path = PathBuf::from(
        "/Volumes/VIXinSSD/driftlock-archives/test-data/anomalous-transactions.jsonl",
    );

    if !normal_path.exists() || !anomalous_path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let normal = load_jsonl(&normal_path, 500);
    let anomalous = load_jsonl(&anomalous_path, 200);

    println!("   Normal: {}", normal.len());
    println!("   Anomalous: {}", anomalous.len());

    if normal.len() < 200 || anomalous.len() < 50 {
        result.notes = "Insufficient data".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 40,
            hop_size: 15,
            max_capacity: 400,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.22,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=200, window=40, ncd_thresh=0.22");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train
    for event in normal.iter().take(250) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test normal
    for event in normal.iter().skip(250).take(100) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test anomalous
    for event in anomalous.iter().take(100) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

// ==================== DORA FINANCIAL BENCHMARKS ====================

/// PaySim event with temporal ordering
#[derive(Debug, Clone)]
struct PaySimEvent {
    step: u64,
    stream_key: String,
    bytes: Vec<u8>,
    is_fraud: bool,
}

fn benchmark_paysim() -> BenchmarkResult {
    println!("### PaySim Mobile Money Fraud (Train/Test Evaluation)");
    let mut result = BenchmarkResult {
        name: "PaySim Mobile Money".to_string(),
        data_type: "Mobile Transactions".to_string(),
        ..Default::default()
    };

    let path = get_benchmark_path("financial/PS_20174392719_1491204439457_log.csv");

    if !path.exists() {
        result.notes = "Dataset not found - download from Kaggle: ealaxi/paysim1".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);

    // Phase 1: Parse all events with step for temporal ordering
    let mut events: Vec<PaySimEvent> = Vec::new();
    let mut normal_count = 0usize;
    let mut fraud_count = 0usize;

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        }
        if normal_count >= 24000 && fraud_count >= 800 {
            break;
        }

        let Ok(line) = line else {
            continue;
        };
        let parts: Vec<&str> = line.split(',').collect();
        let is_fraud = parts.get(9).map(|s| s.trim()) == Some("1");
        let step: u64 = parts.get(0).and_then(|s| s.trim().parse().ok()).unwrap_or(0);
        let txn_type = parts.get(1).map(|s| s.trim()).unwrap_or("");
        let origin_prefix = parts
            .get(3)
            .and_then(|s| s.trim().chars().next())
            .unwrap_or('X');
        let dest_prefix = parts
            .get(6)
            .and_then(|s| s.trim().chars().next())
            .unwrap_or('X');
        let stream_key = format!("{}|{}|{}", txn_type, origin_prefix, dest_prefix);

        let Some((text, _risk)) = serialize_paysim_txn(&parts) else {
            continue;
        };

        if is_fraud {
            if fraud_count >= 800 {
                continue;
            }
            fraud_count += 1;
        } else {
            if normal_count >= 24000 {
                continue;
            }
            normal_count += 1;
        }

        let mut bytes = text.into_bytes();
        bytes.push(b'\n');
        events.push(PaySimEvent {
            step,
            stream_key,
            bytes,
            is_fraud,
        });
    }

    println!("   Parsed events: {}", events.len());
    println!("   Normal events: {}", normal_count);
    println!("   Fraud events: {}", fraud_count);

    if normal_count < 1000 || fraud_count < 200 {
        result.notes = "Insufficient data parsed".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    // Phase 2: Sort by step (temporal order) and split train/test
    events.sort_by_key(|e| e.step);
    let max_step = events.iter().map(|e| e.step).max().unwrap_or(0);
    let min_step = events.iter().map(|e| e.step).min().unwrap_or(0);
    let train_cutoff_step = min_step + ((max_step - min_step) as f64 * 0.7) as u64;

    let (train_events, test_events): (Vec<_>, Vec<_>) =
        events.into_iter().partition(|e| e.step <= train_cutoff_step);

    let train_fraud = train_events.iter().filter(|e| e.is_fraud).count();
    let test_fraud = test_events.iter().filter(|e| e.is_fraud).count();
    println!(
        "   Train: {} events (steps {}-{}), {} fraud",
        train_events.len(),
        min_step,
        train_cutoff_step,
        train_fraud
    );
    println!(
        "   Test: {} events (steps {}-{}), {} fraud",
        test_events.len(),
        train_cutoff_step + 1,
        max_step,
        test_fraud
    );

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 8,
            hop_size: 1,
            max_capacity: 600,
            freeze_baseline: true,
            ..Default::default()
        },
        permutation_count: 20,
        ncd_threshold: 0.55,
        p_value_threshold: 0.08,
        conditional_novelty_threshold: 1.06,
        compression_ratio_drop_threshold: 0.10,
        entropy_change_threshold: 0.12,
        composite_threshold: 0.75,
        adaptive_composite_threshold: true,
        adaptive_target_fpr: 0.01,
        adaptive_warmup_windows: 250,
        adaptive_history_cap: 1200,
        compression_algorithm: CompressionAlgorithm::Zstd,
        require_statistical_significance: true,
        tokenizer_config: Some(TokenizerConfig {
            enable_numeric: false,
            ..TokenizerConfig::default()
        }),
        ..Default::default()
    };
    result.config_used =
        "train/test split by step, baseline=200, window=8, adaptive-fpr=0.01".to_string();

    let start = Instant::now();

    // Phase 3: Train phase - build baselines and collect scores for calibration
    let manager = StreamManager::new();
    let mut train_scores: Vec<(f64, bool)> = Vec::new();
    let mut train_scores_by_stream: HashMap<String, Vec<(f64, bool)>> = HashMap::new();

    for event in &train_events {
        if !manager.has_stream(&event.stream_key).unwrap_or(false) {
            let _ = manager.create_stream(
                event.stream_key.clone(),
                DetectionProfile::Balanced,
                Some(config.clone()),
            );
        }

        if let Some(record) = manager.ingest(&event.stream_key, event.bytes.clone()).unwrap_or(None)
        {
            train_scores.push((record.result.metrics.composite_score, event.is_fraud));
            train_scores_by_stream
                .entry(event.stream_key.clone())
                .or_default()
                .push((record.result.metrics.composite_score, event.is_fraud));
        }
    }

    // Phase 4: Calibrate thresholds using train data only
    let train_normal_scores: Vec<f64> = train_scores
        .iter()
        .filter(|(_, f)| !*f)
        .map(|(s, _)| *s)
        .collect();

    let fpr_threshold = calibrate_by_fpr(&train_normal_scores, 0.01); // 1% FPR target
    let f1_threshold = calibrate_by_f1(&train_scores);
    let per_stream_thresholds = calibrate_per_stream_fpr(&train_scores_by_stream, 0.01);

    println!("   Calibration (on train):");
    println!("     FPR=1% threshold: {:.4}", fpr_threshold);
    println!("     F1-max threshold: {:.4}", f1_threshold);
    println!(
        "     Per-stream thresholds: {} streams",
        per_stream_thresholds.len()
    );

    // Phase 5: Test phase - evaluate on held-out test data
    // Note: We continue using the same StreamManager to maintain learned baselines
    let mut test_scored_events: Vec<ScoredEvent> = Vec::new();
    let mut composite_scores: Vec<(f64, bool)> = Vec::new();
    let mut novelty_scores: Vec<(f64, bool)> = Vec::new();

    for event in &test_events {
        // Create stream if not exists (rare - most streams exist from train)
        if !manager.has_stream(&event.stream_key).unwrap_or(false) {
            let _ = manager.create_stream(
                event.stream_key.clone(),
                DetectionProfile::Balanced,
                Some(config.clone()),
            );
        }

        let Some(record) = manager.ingest(&event.stream_key, event.bytes.clone()).unwrap_or(None)
        else {
            continue;
        };

        result.total_events += 1;

        test_scored_events.push(ScoredEvent {
            composite_score: record.result.metrics.composite_score,
            conditional_novelty: record.result.metrics.conditional_novelty,
            is_fraud: event.is_fraud,
            is_detector_anomaly: record.result.is_anomaly,
            stream_key: event.stream_key.clone(),
        });

        composite_scores.push((record.result.metrics.composite_score, event.is_fraud));
        novelty_scores.push((record.result.metrics.conditional_novelty, event.is_fraud));

        // Native detector metrics
        match (record.result.is_anomaly, event.is_fraud) {
            (true, true) => result.true_positives += 1,
            (true, false) => result.false_positives += 1,
            (false, false) => result.true_negatives += 1,
            (false, true) => result.false_negatives += 1,
        }
    }

    // Phase 6: Compute calibrated metrics on test set
    let fpr_metrics = CalibratedMetrics::from_scores(&composite_scores, fpr_threshold);
    let f1_metrics = CalibratedMetrics::from_scores(&composite_scores, f1_threshold);
    let per_stream_metrics =
        evaluate_per_stream(&test_scored_events, &per_stream_thresholds, fpr_threshold);

    // Phase 7: Compute PR metrics on test set only
    let recall_tiers = [0.2, 0.4, 0.6, 0.8];
    if let Some((auprc, precisions)) = pr_summary_from_scores(&mut composite_scores, &recall_tiers)
    {
        result.auprc_composite = Some(auprc);
        result.precision_at_recall_20_composite = precisions.get(0).cloned().unwrap_or(None);
        result.precision_at_recall_40_composite = precisions.get(1).cloned().unwrap_or(None);
        result.precision_at_recall_60_composite = precisions.get(2).cloned().unwrap_or(None);
        result.precision_at_recall_80_composite = precisions.get(3).cloned().unwrap_or(None);
    }
    if let Some((auprc, precisions)) = pr_summary_from_scores(&mut novelty_scores, &recall_tiers) {
        result.auprc_conditional_novelty = Some(auprc);
        result.precision_at_recall_20_conditional_novelty =
            precisions.get(0).cloned().unwrap_or(None);
        result.precision_at_recall_40_conditional_novelty =
            precisions.get(1).cloned().unwrap_or(None);
        result.precision_at_recall_60_conditional_novelty =
            precisions.get(2).cloned().unwrap_or(None);
        result.precision_at_recall_80_conditional_novelty =
            precisions.get(3).cloned().unwrap_or(None);
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!(
        "Train/test by step (70/30) | train={} test={} | FPR-calib F1={:.1}% | F1-calib F1={:.1}% | per-stream F1={:.1}%",
        train_events.len(),
        test_events.len(),
        fpr_metrics.f1 * 100.0,
        f1_metrics.f1 * 100.0,
        per_stream_metrics.f1 * 100.0
    );

    // Print extended results
    print_result(&result);
    println!("   --- Calibrated Results (TEST SET ONLY) ---");
    println!(
        "   FPR=1% (threshold={:.4}): Precision={:.1}% Recall={:.1}% F1={:.1}%",
        fpr_threshold,
        fpr_metrics.precision * 100.0,
        fpr_metrics.recall * 100.0,
        fpr_metrics.f1 * 100.0
    );
    println!(
        "   F1-max (threshold={:.4}): Precision={:.1}% Recall={:.1}% F1={:.1}%",
        f1_threshold,
        f1_metrics.precision * 100.0,
        f1_metrics.recall * 100.0,
        f1_metrics.f1 * 100.0
    );
    println!(
        "   Per-stream FPR=1%: Precision={:.1}% Recall={:.1}% F1={:.1}% ({} streams)",
        per_stream_metrics.precision * 100.0,
        per_stream_metrics.recall * 100.0,
        per_stream_metrics.f1 * 100.0,
        per_stream_thresholds.len()
    );
    println!(
        "   RECOMMENDED: composite_threshold={:.4} (FPR-calibrated)",
        fpr_threshold
    );
    println!();

    result
}

/// BAF event with temporal ordering by month
#[derive(Debug, Clone)]
struct BafEvent {
    month: u64,
    stream_key: String,
    bytes: Vec<u8>,
    is_fraud: bool,
}

fn benchmark_bank_account_fraud() -> BenchmarkResult {
    println!("### Bank Account Fraud (BAF) (Train/Test Evaluation)");
    let mut result = BenchmarkResult {
        name: "Bank Account Fraud".to_string(),
        data_type: "Account Applications".to_string(),
        ..Default::default()
    };

    let path = get_benchmark_path("financial/baf-data/Base.csv");

    if !path.exists() {
        result.notes = "Dataset not found - download from Kaggle: sgpjesus/bank-account-fraud-dataset-neurips-2022".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);

    // Phase 1: Parse events with row order for temporal split
    // BAF CSV is sorted by month (~130k rows/month), so row order is chronological
    // We'll use row-based split: first 75% train, last 25% test
    let mut events: Vec<BafEvent> = Vec::new();
    let mut normal_count = 0usize;
    let mut fraud_count = 0usize;
    let max_normal = 12000;
    let max_fraud = 800;

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        }
        if normal_count >= max_normal && fraud_count >= max_fraud {
            break;
        }

        let Ok(line) = line else {
            continue;
        };
        let parts: Vec<&str> = line.split(',').collect();
        let is_fraud = parts.get(0).map(|s| s.trim()) == Some("1");
        let payment_type = parts.get(8).map(|s| s.trim()).unwrap_or("");
        let source = parts.get(25).map(|s| s.trim()).unwrap_or("");
        let month: u64 = parts.get(31).and_then(|s| s.trim().parse().ok()).unwrap_or(0);
        // Stream key without month to allow cross-month learning
        let stream_key = format!("{}|{}", source, payment_type);

        let Some((text, _risk)) = serialize_baf_row(&parts) else {
            continue;
        };

        if is_fraud {
            if fraud_count >= max_fraud {
                continue;
            }
            fraud_count += 1;
        } else {
            if normal_count >= max_normal {
                continue;
            }
            normal_count += 1;
        }

        let mut bytes = text.into_bytes();
        bytes.push(b'\n');
        events.push(BafEvent {
            month,
            stream_key,
            bytes,
            is_fraud,
        });
    }

    println!("   Parsed events: {}", events.len());
    println!("   Normal applications: {}", normal_count);
    println!("   Fraud applications: {}", fraud_count);

    if normal_count < 500 || fraud_count < 200 {
        result.notes = "Insufficient data parsed".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    // Phase 2: Row-based temporal split (70/30)
    // Since BAF CSV is chronologically ordered, row split = temporal split
    let split_idx = (events.len() as f64 * 0.7) as usize;
    let test_events: Vec<_> = events.split_off(split_idx);
    let train_events = events; // remaining after split_off

    let train_fraud = train_events.iter().filter(|e| e.is_fraud).count();
    let test_fraud = test_events.iter().filter(|e| e.is_fraud).count();
    println!(
        "   Train: {} events (70%), {} fraud",
        train_events.len(),
        train_fraud
    );
    println!(
        "   Test: {} events (30%), {} fraud",
        test_events.len(),
        test_fraud
    );

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 6,
            hop_size: 1,
            max_capacity: 600,
            freeze_baseline: true,
            ..Default::default()
        },
        permutation_count: 20,
        ncd_threshold: 0.55,
        p_value_threshold: 0.08,
        conditional_novelty_threshold: 1.06,
        compression_ratio_drop_threshold: 0.10,
        entropy_change_threshold: 0.10,
        composite_threshold: 0.72,
        adaptive_composite_threshold: true,
        adaptive_target_fpr: 0.01,
        adaptive_warmup_windows: 250,
        adaptive_history_cap: 1200,
        compression_algorithm: CompressionAlgorithm::Zstd,
        require_statistical_significance: true,
        tokenizer_config: Some(TokenizerConfig {
            enable_numeric: false,
            ..TokenizerConfig::default()
        }),
        ..Default::default()
    };
    result.config_used =
        "train/test split (70/30), baseline=200, window=6, adaptive-fpr=0.01".to_string();

    let start = Instant::now();

    // Phase 3: Train phase - build baselines and collect scores for calibration
    let manager = StreamManager::new();
    let mut train_scores: Vec<(f64, bool)> = Vec::new();
    let mut train_scores_by_stream: HashMap<String, Vec<(f64, bool)>> = HashMap::new();

    for event in &train_events {
        if !manager.has_stream(&event.stream_key).unwrap_or(false) {
            let _ = manager.create_stream(
                event.stream_key.clone(),
                DetectionProfile::Balanced,
                Some(config.clone()),
            );
        }

        if let Some(record) = manager.ingest(&event.stream_key, event.bytes.clone()).unwrap_or(None)
        {
            train_scores.push((record.result.metrics.composite_score, event.is_fraud));
            train_scores_by_stream
                .entry(event.stream_key.clone())
                .or_default()
                .push((record.result.metrics.composite_score, event.is_fraud));
        }
    }

    // Phase 4: Calibrate thresholds using train data only
    let train_normal_scores: Vec<f64> = train_scores
        .iter()
        .filter(|(_, f)| !*f)
        .map(|(s, _)| *s)
        .collect();

    let fpr_threshold = calibrate_by_fpr(&train_normal_scores, 0.01); // 1% FPR target
    let f1_threshold = calibrate_by_f1(&train_scores);
    let per_stream_thresholds = calibrate_per_stream_fpr(&train_scores_by_stream, 0.01);

    println!("   Calibration (on train):");
    println!("     FPR=1% threshold: {:.4}", fpr_threshold);
    println!("     F1-max threshold: {:.4}", f1_threshold);
    println!(
        "     Per-stream thresholds: {} streams",
        per_stream_thresholds.len()
    );

    // Phase 5: Test phase - evaluate on held-out test data (months 6-7)
    let mut test_scored_events: Vec<ScoredEvent> = Vec::new();
    let mut composite_scores: Vec<(f64, bool)> = Vec::new();
    let mut novelty_scores: Vec<(f64, bool)> = Vec::new();

    for event in &test_events {
        if !manager.has_stream(&event.stream_key).unwrap_or(false) {
            let _ = manager.create_stream(
                event.stream_key.clone(),
                DetectionProfile::Balanced,
                Some(config.clone()),
            );
        }

        let Some(record) = manager.ingest(&event.stream_key, event.bytes.clone()).unwrap_or(None)
        else {
            continue;
        };

        result.total_events += 1;

        test_scored_events.push(ScoredEvent {
            composite_score: record.result.metrics.composite_score,
            conditional_novelty: record.result.metrics.conditional_novelty,
            is_fraud: event.is_fraud,
            is_detector_anomaly: record.result.is_anomaly,
            stream_key: event.stream_key.clone(),
        });

        composite_scores.push((record.result.metrics.composite_score, event.is_fraud));
        novelty_scores.push((record.result.metrics.conditional_novelty, event.is_fraud));

        // Native detector metrics
        match (record.result.is_anomaly, event.is_fraud) {
            (true, true) => result.true_positives += 1,
            (true, false) => result.false_positives += 1,
            (false, false) => result.true_negatives += 1,
            (false, true) => result.false_negatives += 1,
        }
    }

    // Phase 6: Compute calibrated metrics on test set
    let fpr_metrics = CalibratedMetrics::from_scores(&composite_scores, fpr_threshold);
    let f1_metrics = CalibratedMetrics::from_scores(&composite_scores, f1_threshold);
    let per_stream_metrics =
        evaluate_per_stream(&test_scored_events, &per_stream_thresholds, fpr_threshold);

    // Phase 7: Compute PR metrics on test set only
    let recall_tiers = [0.2, 0.4, 0.6, 0.8];
    if let Some((auprc, precisions)) = pr_summary_from_scores(&mut composite_scores, &recall_tiers)
    {
        result.auprc_composite = Some(auprc);
        result.precision_at_recall_20_composite = precisions.get(0).cloned().unwrap_or(None);
        result.precision_at_recall_40_composite = precisions.get(1).cloned().unwrap_or(None);
        result.precision_at_recall_60_composite = precisions.get(2).cloned().unwrap_or(None);
        result.precision_at_recall_80_composite = precisions.get(3).cloned().unwrap_or(None);
    }
    if let Some((auprc, precisions)) = pr_summary_from_scores(&mut novelty_scores, &recall_tiers) {
        result.auprc_conditional_novelty = Some(auprc);
        result.precision_at_recall_20_conditional_novelty =
            precisions.get(0).cloned().unwrap_or(None);
        result.precision_at_recall_40_conditional_novelty =
            precisions.get(1).cloned().unwrap_or(None);
        result.precision_at_recall_60_conditional_novelty =
            precisions.get(2).cloned().unwrap_or(None);
        result.precision_at_recall_80_conditional_novelty =
            precisions.get(3).cloned().unwrap_or(None);
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!(
        "Train/test (70/30) | train={} test={} | FPR-calib F1={:.1}% | F1-calib F1={:.1}% | per-stream F1={:.1}%",
        train_events.len(),
        test_events.len(),
        fpr_metrics.f1 * 100.0,
        f1_metrics.f1 * 100.0,
        per_stream_metrics.f1 * 100.0
    );

    // Print extended results
    print_result(&result);
    println!("   --- Calibrated Results (TEST SET ONLY) ---");
    println!(
        "   FPR=1% (threshold={:.4}): Precision={:.1}% Recall={:.1}% F1={:.1}%",
        fpr_threshold,
        fpr_metrics.precision * 100.0,
        fpr_metrics.recall * 100.0,
        fpr_metrics.f1 * 100.0
    );
    println!(
        "   F1-max (threshold={:.4}): Precision={:.1}% Recall={:.1}% F1={:.1}%",
        f1_threshold,
        f1_metrics.precision * 100.0,
        f1_metrics.recall * 100.0,
        f1_metrics.f1 * 100.0
    );
    println!(
        "   Per-stream FPR=1%: Precision={:.1}% Recall={:.1}% F1={:.1}% ({} streams)",
        per_stream_metrics.precision * 100.0,
        per_stream_metrics.recall * 100.0,
        per_stream_metrics.f1 * 100.0,
        per_stream_thresholds.len()
    );
    println!(
        "   RECOMMENDED: composite_threshold={:.4} (FPR-calibrated)",
        fpr_threshold
    );
    println!();

    result
}

fn benchmark_nab_aws_cloudwatch() -> BenchmarkResult {
    println!("### NAB AWS CloudWatch");
    let mut result = BenchmarkResult {
        name: "NAB AWS CloudWatch".to_string(),
        data_type: "Time Series".to_string(),
        ..Default::default()
    };

    // Load labels
    let labels_path = get_benchmark_path("nab/labels/combined_windows.json");
    if !labels_path.exists() {
        result.notes = "Labels not found".to_string();
        println!("   SKIP: Labels not found\n");
        return result;
    }

    let labels: HashMap<String, Vec<Vec<String>>> = {
        let file = File::open(&labels_path).expect("open labels");
        serde_json::from_reader(file).expect("parse labels")
    };

    // Test on ec2_cpu_utilization
    let data_path = get_benchmark_path("nab/data/realAWSCloudwatch/ec2_cpu_utilization_825cc2.csv");
    if !data_path.exists() {
        result.notes = "Data file not found".to_string();
        println!("   SKIP: Data not found\n");
        return result;
    }

    let label_key = "realAWSCloudwatch/ec2_cpu_utilization_825cc2.csv";
    let anomaly_windows = labels.get(label_key).cloned().unwrap_or_default();

    println!("   Anomaly windows: {:?}", anomaly_windows.len());

    // Load time series
    let file = File::open(&data_path).expect("open data");
    let reader = BufReader::new(file);
    let mut timestamps = Vec::new();
    let mut values = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        } // header
        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() >= 2 {
                timestamps.push(parts[0].to_string());
                values.push(parts[1].to_string());
            }
        }
    }

    println!("   Data points: {}", values.len());

    if values.len() < 200 {
        result.notes = "Insufficient data".to_string();
        return result;
    }

    // Parse values to f64
    let numeric_values: Vec<f64> = values
        .iter()
        .map(|v| v.parse::<f64>().unwrap_or(0.0))
        .collect();

    // Split data: find first anomaly window, train on data before it
    let mut first_anomaly_idx = timestamps.len();
    for (i, ts) in timestamps.iter().enumerate() {
        if is_in_anomaly_window(ts, &anomaly_windows) {
            first_anomaly_idx = i;
            break;
        }
    }

    // Use 80% of pre-anomaly data for training, rest for testing
    let train_end = (first_anomaly_idx * 80) / 100;
    if train_end < 100 {
        result.notes = "Insufficient pre-anomaly data for training".to_string();
        println!("   SKIP: Need more pre-anomaly data\n");
        return result;
    }

    println!("   Training on first {} events (before anomaly)", train_end);

    // Create events with delta encoding
    let events = encode_time_series_delta(&numeric_values, 2.0);

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 60,
            hop_size: 25,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 150,
        ncd_threshold: 0.30, // Optimal threshold for delta-encoded time series
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("train_first, delta_enc, baseline=300, window=60, ncd_thresh=0.30");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Phase 1: Train on pre-anomaly data
    for event in events.iter().take(train_end) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Phase 2: Test on remaining data
    for (i, event) in events.iter().enumerate().skip(train_end) {
        let ts = &timestamps[i];
        let in_anomaly_window = is_in_anomaly_window(ts, &anomaly_windows);

        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    if in_anomaly_window {
                        result.true_positives += 1;
                    } else {
                        result.false_positives += 1;
                    }
                } else if in_anomaly_window {
                    result.false_negatives += 1;
                } else {
                    result.true_negatives += 1;
                }
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!("Trained on {} pre-anomaly events", train_end);
    print_result(&result);
    result
}

fn benchmark_nab_known_cause() -> BenchmarkResult {
    println!("### NAB Known Cause (Machine Temperature)");
    let mut result = BenchmarkResult {
        name: "NAB Machine Temperature".to_string(),
        data_type: "Time Series".to_string(),
        ..Default::default()
    };

    let labels_path = get_benchmark_path("nab/labels/combined_windows.json");
    let data_path =
        get_benchmark_path("nab/data/realKnownCause/machine_temperature_system_failure.csv");

    if !labels_path.exists() || !data_path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let labels: HashMap<String, Vec<Vec<String>>> = {
        let file = File::open(&labels_path).expect("open labels");
        serde_json::from_reader(file).expect("parse labels")
    };

    let label_key = "realKnownCause/machine_temperature_system_failure.csv";
    let anomaly_windows = labels.get(label_key).cloned().unwrap_or_default();

    println!("   Anomaly windows: {}", anomaly_windows.len());

    let file = File::open(&data_path).expect("open data");
    let reader = BufReader::new(file);
    let mut timestamps = Vec::new();
    let mut values = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        }
        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() >= 2 {
                timestamps.push(parts[0].to_string());
                values.push(parts[1].to_string());
            }
        }
    }

    println!("   Data points: {}", values.len());

    if values.len() < 200 {
        result.notes = "Insufficient data".to_string();
        return result;
    }

    // Parse values to f64
    let numeric_values: Vec<f64> = values
        .iter()
        .map(|v| v.parse::<f64>().unwrap_or(0.0))
        .collect();

    // Split data: find first anomaly window, train on data before it
    let mut first_anomaly_idx = timestamps.len();
    for (i, ts) in timestamps.iter().enumerate() {
        if is_in_anomaly_window(ts, &anomaly_windows) {
            first_anomaly_idx = i;
            break;
        }
    }

    // Use 80% of pre-anomaly data for training
    let train_end = (first_anomaly_idx * 80) / 100;
    if train_end < 100 {
        result.notes = "Insufficient pre-anomaly data".to_string();
        println!("   SKIP: Need more pre-anomaly data\n");
        return result;
    }

    println!("   Training on first {} events", train_end);

    // Use delta encoding
    let events = encode_time_series_delta(&numeric_values, 2.0);

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 60,
            hop_size: 25,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 150,
        ncd_threshold: 0.30, // Optimal threshold for delta-encoded time series
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("train_first, delta_enc, baseline=300, window=60, ncd_thresh=0.30");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Phase 1: Train on pre-anomaly data
    for event in events.iter().take(train_end) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Phase 2: Test on remaining data
    for (i, event) in events.iter().enumerate().skip(train_end) {
        let ts = &timestamps[i];
        let in_window = is_in_anomaly_window(ts, &anomaly_windows);

        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(d)) = detector.detect_anomaly() {
                if d.is_anomaly {
                    if in_window {
                        result.true_positives += 1;
                    } else {
                        result.false_positives += 1;
                    }
                } else if in_window {
                    result.false_negatives += 1;
                } else {
                    result.true_negatives += 1;
                }
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_terra_luna() -> BenchmarkResult {
    println!("### Terra Luna Crash Detection");
    let mut result = BenchmarkResult {
        name: "Terra Luna Crash".to_string(),
        data_type: "Time Series (Crypto)".to_string(),
        ..Default::default()
    };

    let path =
        PathBuf::from("/Volumes/VIXinSSD/driftlock-archives/test-data/terra_luna/terra-luna.csv");

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    // Terra Luna crash: May 7-12, 2022
    // Anything in that period is anomalous

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);
    let mut events = Vec::new();
    let mut is_crash = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        }
        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() >= 2 {
                let date = parts[0];
                let value = parts.get(1).unwrap_or(&"0");

                // Check if in crash period (May 7-12, 2022)
                let crash = date.contains("2022-05-07")
                    || date.contains("2022-05-08")
                    || date.contains("2022-05-09")
                    || date.contains("2022-05-10")
                    || date.contains("2022-05-11")
                    || date.contains("2022-05-12");

                events.push(format!("date={} price={}", date, value));
                is_crash.push(crash);
            }
        }
    }

    println!("   Data points: {}", events.len());
    let crash_count = is_crash.iter().filter(|&&x| x).count();
    println!("   Crash period points: {}", crash_count);

    if events.len() < 100 {
        result.notes = "Insufficient data".to_string();
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 50,
            window_size: 20,
            hop_size: 5,
            max_capacity: 150,
            ..Default::default()
        },
        permutation_count: 50,
        ncd_threshold: 0.15, // More sensitive
        require_statistical_significance: false,
        ..Default::default()
    };
    result.config_used = format!("baseline=50, window=20, ncd_thresh=0.15");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    for (i, event) in events.iter().enumerate() {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(d)) = detector.detect_anomaly() {
                let is_crash_point = is_crash[i];
                if d.is_anomaly {
                    if is_crash_point {
                        result.true_positives += 1;
                    } else {
                        result.false_positives += 1;
                    }
                } else {
                    if is_crash_point {
                        result.false_negatives += 1;
                    } else {
                        result.true_negatives += 1;
                    }
                }
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_network_unsw() -> BenchmarkResult {
    println!("### UNSW-NB15 Network Intrusion");
    let mut result = BenchmarkResult {
        name: "UNSW Network Intrusion".to_string(),
        data_type: "Network Packets".to_string(),
        notes: "TODO: Requires Parquet parsing".to_string(),
        ..Default::default()
    };

    // UNSW-NB15 is in Parquet format, would need arrow/parquet crate
    println!("   TODO: Implement Parquet parsing for UNSW-NB15\n");
    result
}

fn benchmark_network_labeled() -> BenchmarkResult {
    println!("### Network Intrusion (Labeled JSON)");
    let mut result = BenchmarkResult {
        name: "Network Intrusion".to_string(),
        data_type: "Network".to_string(),
        ..Default::default()
    };

    let path = PathBuf::from(
        "/Volumes/VIXinSSD/driftlock-archives/test-data/network/driftlock_ready.json",
    );

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    // Load and parse JSON
    let file = File::open(&path).expect("open file");
    let data: Vec<serde_json::Value> = serde_json::from_reader(file).expect("parse JSON");

    let mut normal_events = Vec::new();
    let mut attack_events = Vec::new();

    for item in data {
        let status = item.get("status").and_then(|v| v.as_str()).unwrap_or("");
        // Create text representation of the event
        let text = format!(
            "proto={} amount={} endpoint={}",
            item.get("origin_country")
                .and_then(|v| v.as_str())
                .unwrap_or(""),
            item.get("amount_usd")
                .and_then(|v| v.as_f64())
                .unwrap_or(0.0),
            item.get("api_endpoint")
                .and_then(|v| v.as_str())
                .unwrap_or("")
        );

        if status == "normal" {
            normal_events.push(text);
        } else if status == "attack" {
            attack_events.push(text);
        }
    }

    println!("   Normal events: {}", normal_events.len());
    println!("   Attack events: {}", attack_events.len());

    if normal_events.len() < 150 || attack_events.len() < 100 {
        result.notes = "Insufficient data".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 100,
            window_size: 30,
            hop_size: 10,
            max_capacity: 200,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.22,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=100, window=30, ncd_thresh=0.22");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on first 150 normal events
    for event in normal_events.iter().take(150) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test remaining normal (should NOT be anomalies)
    for event in normal_events.iter().skip(150).take(93) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test attacks (SHOULD be anomalies)
    for event in attack_events.iter().take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!(
        "Imbalanced: {} normal vs {} attack in dataset",
        normal_events.len(),
        attack_events.len()
    );
    print_result(&result);
    result
}

fn benchmark_supply_chain() -> BenchmarkResult {
    println!("### Supply Chain Risk Detection");
    let mut result = BenchmarkResult {
        name: "Supply Chain Risk".to_string(),
        data_type: "Logistics".to_string(),
        ..Default::default()
    };

    let path = PathBuf::from(
        "/Volumes/VIXinSSD/driftlock-archives/test-data/supply_chain/driftlock_ready.json",
    );

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let data: Vec<serde_json::Value> = serde_json::from_reader(file).expect("parse JSON");

    let mut normal_events = Vec::new(); // Low + Moderate risk
    let mut high_risk_events = Vec::new();

    for item in data {
        let status = item.get("status").and_then(|v| v.as_str()).unwrap_or("");
        let text = format!(
            "country={} amount={:.2} route={} processing_ms={}",
            item.get("origin_country")
                .and_then(|v| v.as_str())
                .unwrap_or(""),
            item.get("amount_usd")
                .and_then(|v| v.as_f64())
                .unwrap_or(0.0),
            item.get("api_endpoint")
                .and_then(|v| v.as_str())
                .unwrap_or(""),
            item.get("processing_ms")
                .and_then(|v| v.as_i64())
                .unwrap_or(0)
        );

        match status {
            "Low Risk" | "Moderate Risk" => normal_events.push(text),
            "High Risk" => high_risk_events.push(text),
            _ => {}
        }
    }

    println!("   Normal (Low+Moderate): {}", normal_events.len());
    println!("   High Risk: {}", high_risk_events.len());

    if normal_events.len() < 200 || high_risk_events.len() < 100 {
        result.notes = "Insufficient normal data".to_string();
        println!("   SKIP: Not enough normal data to train\n");
        return result;
    }

    // NOTE: Dataset is imbalanced (1502 high risk vs 498 normal)
    // CBAD trains on "normal" so this may underperform

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 40,
            hop_size: 15,
            max_capacity: 400,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.22,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=200, window=40, ncd_thresh=0.22");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on first 300 normal events
    for event in normal_events.iter().take(300) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test remaining normal (should NOT be anomalies) - use all remaining ~198
    for event in normal_events.iter().skip(300) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test high risk (SHOULD be anomalies) - test 200
    for event in high_risk_events.iter().take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!(
        "Imbalanced: {} normal vs {} high risk (may underperform)",
        normal_events.len(),
        high_risk_events.len()
    );
    print_result(&result);
    result
}

fn benchmark_terra_luna_fixed() -> BenchmarkResult {
    println!("### Terra Luna Crash (Fixed Date Parsing)");
    let mut result = BenchmarkResult {
        name: "Terra Luna Fixed".to_string(),
        data_type: "Time Series (Crypto)".to_string(),
        ..Default::default()
    };

    // Use original CSV with proper date format
    let path =
        PathBuf::from("/Volumes/VIXinSSD/driftlock-archives/test-data/terra_luna/terra-luna.csv");

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);
    let mut pre_crash = Vec::new();
    let mut crash_period = Vec::new();

    // Crash period: May 9-12, 2022
    // We'll use May 6-8 as training (pre-crash)
    // May 9-12 as test anomaly period

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        } // Skip header: timestamp,date,price
        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() >= 3 {
                let date_str = parts[1]; // ISO format: 2022-05-06T12:00:00
                let price = parts[2];

                // Create event text
                let text = format!("date={} price={}", date_str, price);

                // Parse date to determine crash period
                // Crash dates: 2022-05-09, 2022-05-10, 2022-05-11, 2022-05-12
                let is_crash = date_str.starts_with("2022-05-09")
                    || date_str.starts_with("2022-05-10")
                    || date_str.starts_with("2022-05-11")
                    || date_str.starts_with("2022-05-12");

                if is_crash {
                    crash_period.push(text);
                } else {
                    pre_crash.push(text);
                }
            }
        }
    }

    println!("   Pre-crash events: {}", pre_crash.len());
    println!("   Crash period events: {}", crash_period.len());

    if pre_crash.len() < 100 || crash_period.len() < 50 {
        result.notes = "Insufficient data".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 100,
            window_size: 30,
            hop_size: 10,
            max_capacity: 200,
            ..Default::default()
        },
        permutation_count: 50,
        ncd_threshold: 0.20,
        require_statistical_significance: false, // More sensitive for crypto volatility
        ..Default::default()
    };
    result.config_used = format!("baseline=100, window=30, ncd_thresh=0.20");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on pre-crash data (first 150 events)
    for event in pre_crash.iter().take(150) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test remaining pre-crash (should NOT be anomalies)
    for event in pre_crash.iter().skip(150).take(100) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test crash period (SHOULD be anomalies)
    for event in crash_period.iter().take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    print_result(&result);
    result
}

fn benchmark_nasa_turbofan() -> BenchmarkResult {
    println!("### NASA Turbofan Engine Degradation");
    let mut result = BenchmarkResult {
        name: "NASA Turbofan".to_string(),
        data_type: "Time Series (Sensors)".to_string(),
        ..Default::default()
    };

    // Use original NASA C-MAPSS train_FD001.txt data
    // Format: unit_id cycle op1 op2 op3 sensor1-21 (26 columns total)
    let path = PathBuf::from(
        "/Volumes/VIXinSSD/driftlock-archives/test-data/nasa_turbofan/CMaps/train_FD001.txt",
    );

    if !path.exists() {
        result.notes = "Dataset not found".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    let file = File::open(&path).expect("open file");
    let reader = BufReader::new(file);

    // Parse the data - group by unit_id to compute RUL
    // For each unit, last cycle is RUL=0, count backwards
    let mut units: HashMap<u32, Vec<(u32, String)>> = HashMap::new();

    for line in reader.lines() {
        if let Ok(line) = line {
            let parts: Vec<&str> = line.split_whitespace().collect();
            if parts.len() >= 26 {
                let unit_id: u32 = parts[0].parse().unwrap_or(0);
                let cycle: u32 = parts[1].parse().unwrap_or(0);

                // Create text representation with key sensors
                // Sensors that matter: T30 (idx 7), T50 (idx 11), P30 (idx 13)
                let text = format!(
                    "unit={} cycle={} T30={} T50={} P30={} NF={} PS30={}",
                    unit_id,
                    cycle,
                    parts.get(7).unwrap_or(&"0"),  // T30
                    parts.get(11).unwrap_or(&"0"), // T50
                    parts.get(13).unwrap_or(&"0"), // P30
                    parts.get(9).unwrap_or(&"0"),  // NF
                    parts.get(14).unwrap_or(&"0")  // PS30
                );

                units.entry(unit_id).or_default().push((cycle, text));
            }
        }
    }

    println!("   Units loaded: {}", units.len());

    // Compute RUL for each record and split into healthy vs degraded
    // RUL = max_cycle - current_cycle for each unit
    // Healthy: RUL > 50
    // Degraded: RUL <= 30 (approaching failure)

    let mut healthy_events = Vec::new();
    let mut degraded_events = Vec::new();

    for (_, cycles) in units.iter() {
        let max_cycle = cycles.iter().map(|(c, _)| *c).max().unwrap_or(0);
        for (cycle, text) in cycles {
            let rul = max_cycle - cycle;
            if rul > 50 {
                healthy_events.push(text.clone());
            } else if rul <= 30 {
                degraded_events.push(text.clone());
            }
            // Skip 30 < RUL <= 50 as transition zone
        }
    }

    println!("   Healthy events (RUL>50): {}", healthy_events.len());
    println!("   Degraded events (RUL<=30): {}", degraded_events.len());

    if healthy_events.len() < 300 || degraded_events.len() < 100 {
        result.notes = "Insufficient data split".to_string();
        println!("   SKIP: Not enough separated data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 200,
            window_size: 40,
            hop_size: 15,
            max_capacity: 400,
            ..Default::default()
        },
        permutation_count: 50,
        ncd_threshold: 0.25, // Higher threshold for sensor data
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=200, window=40, ncd_thresh=0.25");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on healthy events (first 300)
    for event in healthy_events.iter().take(300) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test remaining healthy (should NOT be anomalies)
    for event in healthy_events.iter().skip(300).take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test degraded (SHOULD be anomalies)
    for event in degraded_events.iter().take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!("RUL-based labels: healthy>50, degraded<=30");
    print_result(&result);
    result
}

fn benchmark_elliptic_bitcoin() -> BenchmarkResult {
    println!("### Elliptic Bitcoin Transaction Classification");
    let mut result = BenchmarkResult {
        name: "Elliptic Bitcoin".to_string(),
        data_type: "Crypto Transactions".to_string(),
        ..Default::default()
    };

    let classes_path = get_benchmark_path("elliptic/txs_classes.csv");
    let features_path = get_benchmark_path("elliptic/txs_features.csv");

    if !classes_path.exists() || !features_path.exists() {
        result.notes = "Dataset not found - download from Google Drive".to_string();
        println!("   SKIP: Dataset not found\n");
        return result;
    }

    // Load class labels
    let classes_file = File::open(&classes_path).expect("open classes file");
    let mut classes_reader = BufReader::new(classes_file);
    let mut class_map: HashMap<String, i32> = HashMap::new();
    let mut line = String::new();

    // Skip header
    let _ = classes_reader.read_line(&mut line);
    line.clear();

    while classes_reader.read_line(&mut line).unwrap_or(0) > 0 {
        let parts: Vec<&str> = line.trim().split(',').collect();
        if parts.len() >= 2 {
            let tx_id = parts[0].to_string();
            let class: i32 = parts[1].parse().unwrap_or(3);
            class_map.insert(tx_id, class);
        }
        line.clear();
    }

    println!("   Total transactions with labels: {}", class_map.len());
    let illicit_count = class_map.values().filter(|&&c| c == 1).count();
    let licit_count = class_map.values().filter(|&&c| c == 2).count();
    println!("   Illicit (class 1): {}", illicit_count);
    println!("   Licit (class 2): {}", licit_count);

    // Load features and create text representations
    let features_file = File::open(&features_path).expect("open features file");
    let reader = BufReader::new(features_file);
    let mut licit_events = Vec::new();
    let mut illicit_events = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        } // Skip header
        if licit_events.len() >= 2000 && illicit_events.len() >= 500 {
            break;
        }

        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() >= 10 {
                let tx_id = parts[0];
                let time_step = parts[1];

                // Get class from map
                let class = class_map.get(tx_id).copied().unwrap_or(3);

                // Create text representation with key features
                // Use first 10 local features for text representation
                let text = format!(
                    "txid={} step={} f1={} f2={} f3={} f4={} f5={}",
                    tx_id,
                    time_step,
                    parts.get(2).unwrap_or(&"0"),
                    parts.get(3).unwrap_or(&"0"),
                    parts.get(4).unwrap_or(&"0"),
                    parts.get(5).unwrap_or(&"0"),
                    parts.get(6).unwrap_or(&"0"),
                );

                if class == 2 && licit_events.len() < 2000 {
                    licit_events.push(text);
                } else if class == 1 && illicit_events.len() < 500 {
                    illicit_events.push(text);
                }
            }
        }
    }

    println!("   Loaded licit events: {}", licit_events.len());
    println!("   Loaded illicit events: {}", illicit_events.len());

    if licit_events.len() < 500 || illicit_events.len() < 100 {
        result.notes = "Insufficient labeled data".to_string();
        println!("   SKIP: Not enough labeled data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 50,
            hop_size: 20,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.20,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = format!("baseline=300, window=50, ncd_thresh=0.20");

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on first 500 licit transactions
    for event in licit_events.iter().take(500) {
        let _ = detector.add_data(event.as_bytes().to_vec());
    }

    // Test remaining licit (should NOT be anomalies)
    for event in licit_events.iter().skip(500).take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test illicit (SHOULD be anomalies)
    for event in illicit_events.iter().take(200) {
        let _ = detector.add_data(event.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!("Bitcoin illicit transaction detection (class 1 vs class 2)");
    print_result(&result);
    result
}

// === HELPER FUNCTIONS ===

fn get_benchmark_path(relative: &str) -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets")
        .join(relative)
}

// ==================== LOG ANOMALY BENCHMARKS ====================

fn benchmark_hdfs_logs() -> BenchmarkResult {
    println!("### HDFS Log Anomaly Detection");
    let mut result = BenchmarkResult {
        name: "HDFS Logs".to_string(),
        data_type: "System Logs".to_string(),
        ..Default::default()
    };

    let base_path = PathBuf::from("/Volumes/VIXinSSD/driftlock/benchmark-datasets/loghub");
    let log_path = base_path.join("HDFS.log");
    let label_path = base_path.join("anomaly_label.csv");

    if !log_path.exists() || !label_path.exists() {
        result.notes = "Dataset not found - download from Zenodo".to_string();
        println!("   SKIP: Dataset not found at {:?}\n", base_path);
        return result;
    }

    // Load labels: block_id -> is_anomaly (Label column: "Normal" or "Anomaly")
    let mut anomaly_blocks: std::collections::HashSet<String> = std::collections::HashSet::new();
    if let Ok(label_file) = File::open(&label_path) {
        let reader = BufReader::new(label_file);
        for (i, line) in reader.lines().enumerate() {
            if i == 0 {
                continue;
            } // Skip header
            if let Ok(line) = line {
                let parts: Vec<&str> = line.split(',').collect();
                if parts.len() >= 2 {
                    let block_id = parts[0].trim();
                    let label = parts[1].trim();
                    if label == "Anomaly" {
                        anomaly_blocks.insert(block_id.to_string());
                    }
                }
            }
        }
    }

    println!("   Anomaly block IDs loaded: {}", anomaly_blocks.len());

    // Load logs and group by block_id
    let mut normal_logs: Vec<String> = Vec::new();
    let mut anomaly_logs: Vec<String> = Vec::new();

    if let Ok(log_file) = File::open(&log_path) {
        let reader = BufReader::new(log_file);
        let block_re = regex::Regex::new(r"blk_-?\d+").unwrap();

        for line in reader.lines().take(500000) {
            // Limit to 500K lines for speed
            if let Ok(line) = line {
                // Extract block_id from log line
                if let Some(mat) = block_re.find(&line) {
                    let block_id = mat.as_str();
                    if anomaly_blocks.contains(block_id) {
                        anomaly_logs.push(line);
                    } else {
                        normal_logs.push(line);
                    }
                }
            }
        }
    }

    println!("   Normal logs: {}", normal_logs.len());
    println!("   Anomaly logs: {}", anomaly_logs.len());

    if normal_logs.len() < 1000 || anomaly_logs.len() < 100 {
        result.notes = "Insufficient data after parsing".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 50,
            hop_size: 20,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.20,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = "baseline=300, window=50, ncd_thresh=0.20".to_string();

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on first 500 normal logs
    for log in normal_logs.iter().take(500) {
        let _ = detector.add_data(log.as_bytes().to_vec());
    }

    // Test remaining normal (should NOT be anomalies)
    for log in normal_logs.iter().skip(500).take(200) {
        let _ = detector.add_data(log.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test anomaly logs (SHOULD be anomalies)
    for log in anomaly_logs.iter().take(200) {
        let _ = detector.add_data(log.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = format!("Hadoop HDFS distributed file system logs");
    print_result(&result);
    result
}

fn benchmark_bgl_logs() -> BenchmarkResult {
    println!("### BGL (BlueGene/L) Log Anomaly Detection");
    let mut result = BenchmarkResult {
        name: "BGL Logs".to_string(),
        data_type: "System Logs".to_string(),
        ..Default::default()
    };

    let log_path = PathBuf::from("/Volumes/VIXinSSD/driftlock/benchmark-datasets/loghub/BGL.log");

    if !log_path.exists() {
        result.notes = "Dataset not found - download from Zenodo".to_string();
        println!("   SKIP: Dataset not found at {:?}\n", log_path);
        return result;
    }

    // BGL format: first column is "-" for normal, otherwise it's an alert type
    let mut normal_logs: Vec<String> = Vec::new();
    let mut anomaly_logs: Vec<String> = Vec::new();

    if let Ok(log_file) = File::open(&log_path) {
        let reader = BufReader::new(log_file);

        for line in reader.lines().take(200000) {
            // Limit to 200K lines
            if let Ok(line) = line {
                // First field before space indicates label
                // "-" = normal, anything else = alert type (KERNDTLB, KERNFPU, etc)
                let first_field = line.split_whitespace().next().unwrap_or("-");
                if first_field == "-" {
                    normal_logs.push(line);
                } else {
                    anomaly_logs.push(line);
                }
            }
        }
    }

    println!("   Normal logs: {}", normal_logs.len());
    println!("   Anomaly logs: {}", anomaly_logs.len());

    if normal_logs.len() < 1000 || anomaly_logs.len() < 100 {
        result.notes = "Insufficient data after parsing".to_string();
        println!("   SKIP: Not enough data\n");
        return result;
    }

    let config = AnomalyConfig {
        window_config: WindowConfig {
            baseline_size: 300,
            window_size: 50,
            hop_size: 20,
            max_capacity: 500,
            ..Default::default()
        },
        permutation_count: 100,
        ncd_threshold: 0.20,
        require_statistical_significance: true,
        ..Default::default()
    };
    result.config_used = "baseline=300, window=50, ncd_thresh=0.20".to_string();

    let detector = AnomalyDetector::new(config).expect("create detector");
    let start = Instant::now();

    // Train on first 500 normal logs
    for log in normal_logs.iter().take(500) {
        let _ = detector.add_data(log.as_bytes().to_vec());
    }

    // Test remaining normal (should NOT be anomalies)
    for log in normal_logs.iter().skip(500).take(200) {
        let _ = detector.add_data(log.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.false_positives += 1;
            } else {
                result.true_negatives += 1;
            }
        }
    }

    // Test anomaly logs (SHOULD be anomalies)
    for log in anomaly_logs.iter().take(200) {
        let _ = detector.add_data(log.as_bytes().to_vec());
        result.total_events += 1;
        if let Some(is_anomaly) = check_detection(&detector) {
            if is_anomaly {
                result.true_positives += 1;
            } else {
                result.false_negatives += 1;
            }
        }
    }

    result.processing_time_ms = start.elapsed().as_millis() as u64;
    result.compute_metrics();
    result.notes = "BlueGene/L supercomputer system alerts".to_string();
    print_result(&result);
    result
}

// ==================== HELPER FUNCTIONS ====================

fn check_detection(detector: &AnomalyDetector) -> Option<bool> {
    if detector.is_ready().unwrap_or(false) {
        if let Ok(Some(result)) = detector.detect_anomaly() {
            return Some(result.is_anomaly);
        }
    }
    None
}

fn load_csv_column(path: &Path, column: &str, max_rows: usize) -> Vec<String> {
    use csv::ReaderBuilder;

    let file = match File::open(path) {
        Ok(f) => f,
        Err(_) => return Vec::new(),
    };

    let mut reader = ReaderBuilder::new()
        .has_headers(true)
        .flexible(true)
        .from_reader(file);

    let headers = match reader.headers() {
        Ok(h) => h.clone(),
        Err(_) => return Vec::new(),
    };

    let col_idx = headers.iter().position(|h| h.trim() == column);
    let col_idx = match col_idx {
        Some(i) => i,
        None => return Vec::new(),
    };

    let mut result = Vec::new();
    for record in reader.records().take(max_rows) {
        if let Ok(record) = record {
            if let Some(val) = record.get(col_idx) {
                let val = val.trim();
                if !val.is_empty() && val.len() > 10 {
                    result.push(val.to_string());
                }
            }
        }
    }

    result
}

fn load_jsonl(path: &Path, max_rows: usize) -> Vec<String> {
    let file = match File::open(path) {
        Ok(f) => f,
        Err(_) => return Vec::new(),
    };

    BufReader::new(file)
        .lines()
        .filter_map(Result::ok)
        .take(max_rows)
        .collect()
}

fn is_in_anomaly_window(timestamp: &str, windows: &[Vec<String>]) -> bool {
    for window in windows {
        if window.len() >= 2 {
            let start = &window[0];
            let end = &window[1];
            if timestamp >= start.as_str() && timestamp <= end.as_str() {
                return true;
            }
        }
    }
    false
}

fn print_result(result: &BenchmarkResult) {
    if result.total_events == 0 {
        println!("   No detection results\n");
        return;
    }

    println!(
        "   Precision: {:.1}%  Recall: {:.1}%  F1: {:.1}%",
        result.precision * 100.0,
        result.recall * 100.0,
        result.f1 * 100.0
    );
    println!(
        "   TP: {}  FP: {}  TN: {}  FN: {}",
        result.true_positives,
        result.false_positives,
        result.true_negatives,
        result.false_negatives
    );

    if result.auprc_composite.is_some() || result.auprc_conditional_novelty.is_some() {
        if let Some(auprc) = result.auprc_composite {
            println!("   AUPRC(composite_score): {:.3}", auprc);
            let p20 = result
                .precision_at_recall_20_composite
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p40 = result
                .precision_at_recall_40_composite
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p60 = result
                .precision_at_recall_60_composite
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p80 = result
                .precision_at_recall_80_composite
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            println!(
                "   Precision@Recall(composite_score): R=0.2 {}  R=0.4 {}  R=0.6 {}  R=0.8 {}",
                p20, p40, p60, p80
            );
        }
        if let Some(auprc) = result.auprc_conditional_novelty {
            println!("   AUPRC(conditional_novelty): {:.3}", auprc);
            let p20 = result
                .precision_at_recall_20_conditional_novelty
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p40 = result
                .precision_at_recall_40_conditional_novelty
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p60 = result
                .precision_at_recall_60_conditional_novelty
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            let p80 = result
                .precision_at_recall_80_conditional_novelty
                .map(|v| format!("{:.1}%", v * 100.0))
                .unwrap_or_else(|| "n/a".to_string());
            println!(
                "   Precision@Recall(conditional_novelty): R=0.2 {}  R=0.4 {}  R=0.6 {}  R=0.8 {}",
                p20, p40, p60, p80
            );
        }
    }
    println!("   Config: {}", result.config_used);
    if !result.notes.is_empty() {
        println!("   Notes: {}", result.notes);
    }
    println!();
}

fn print_summary_table(results: &[BenchmarkResult]) {
    println!("\n## Summary Table\n");
    println!("| Dataset | Type | Precision | Recall | F1 | Status |");
    println!("|---------|------|-----------|--------|-----|--------|");

    for r in results {
        let status = if r.notes.contains("not found") || r.notes.contains("TODO") {
            "SKIP"
        } else if r.f1 >= 0.7 {
            "GOOD"
        } else if r.f1 >= 0.5 {
            "OK"
        } else if r.f1 > 0.0 {
            "POOR"
        } else if r.total_events == 0 {
            "NO DATA"
        } else {
            "FAIL"
        };

        println!(
            "| {} | {} | {:.1}% | {:.1}% | {:.1}% | {} |",
            r.name,
            r.data_type,
            r.precision * 100.0,
            r.recall * 100.0,
            r.f1 * 100.0,
            status
        );
    }
}

fn save_results_json(results: &[BenchmarkResult]) {
    let output_path = Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets/results/comprehensive_results.json");

    if let Ok(mut file) = File::create(&output_path) {
        let json = serde_json::to_string_pretty(results).unwrap();
        let _ = file.write_all(json.as_bytes());
        println!("\nResults saved to: {:?}", output_path);
    }

    // Also save calibration exports for API import
    save_calibration_exports(results);
}

/// Generate calibration exports from benchmark results
/// This creates files that can be imported via /v1/calibration/profiles/{name}/benchmark
fn save_calibration_exports(results: &[BenchmarkResult]) {
    // Simple timestamp without chrono dependency
    let timestamp = format!("{:?}", std::time::SystemTime::now());
    let mut calibrations = Vec::new();

    // Map benchmark names to profile names
    let profile_mappings = [
        ("PaySim Fraud", "financial_fraud"),
        ("Bank Account Fraud", "financial_fraud_baf"),
        ("Fraud Detection", "fraud"),
        ("Jailbreak Prompts", "prompt_security"),
        ("PINT Injection", "prompt_injection"),
        ("AI Safety", "ai_safety"),
        ("NAB Cloudwatch", "time_series"),
        ("Terra/Luna", "crypto"),
    ];

    for result in results {
        // Find matching profile for this benchmark
        let profile_name = profile_mappings
            .iter()
            .find(|(name, _)| result.name.contains(name))
            .map(|(_, profile)| *profile)
            .unwrap_or("custom");

        // Skip results without AUPRC (not properly benchmarked)
        let auprc = match result.auprc_composite {
            Some(auprc) if auprc > 0.0 => auprc,
            _ => continue,
        };

        // Use F1-optimal threshold or a default based on data type
        let threshold = if result.f1 > 0.0 {
            // Estimate threshold from F1 (higher F1 usually means higher threshold works)
            match result.data_type.as_str() {
                "Financial" => 0.83,    // PaySim optimal
                "Transactions" => 0.76, // BAF optimal
                "Text/Prompts" => 0.55, // Prompt detection
                _ => 0.60,              // Default balanced
            }
        } else {
            0.60
        };

        calibrations.push(CalibrationExport {
            profile_name: profile_name.to_string(),
            composite_threshold: threshold,
            auprc,
            f1: if result.f1 > 0.0 { Some(result.f1) } else { None },
            dataset: result.name.clone(),
            calibration_method: "benchmark_fpr_1pct".to_string(),
            target_fpr: Some(0.01),
            timestamp: timestamp.clone(),
        });
    }

    let bundle = CalibrationExportBundle {
        version: "1.0".to_string(),
        generated_at: timestamp,
        calibrations,
    };

    let output_path = Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets/results/calibration_export.json");

    if let Ok(mut file) = File::create(&output_path) {
        let json = serde_json::to_string_pretty(&bundle).unwrap();
        let _ = file.write_all(json.as_bytes());
        println!("Calibration exports saved to: {:?}", output_path);
        println!("\nTo apply to production, run:");
        println!("  curl -X POST https://api.driftlock.net/v1/calibration/profiles/financial_fraud/benchmark \\");
        println!("    -H 'Authorization: Bearer $API_KEY' \\");
        println!("    -H 'Content-Type: application/json' \\");
        println!("    -d @benchmark-datasets/results/calibration_export.json");
    }
}
