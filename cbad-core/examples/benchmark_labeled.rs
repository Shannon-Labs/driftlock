//! Benchmark CBAD against labeled datasets with precision/recall metrics
//!
//! Run with: cargo run --example benchmark_labeled --release
//!
//! This tests whether CBAD actually detects the labeled anomalies in:
//! - Jailbreak prompts (labeled as jailbreak vs regular)
//! - Hallucination detection (labeled as correct vs hallucinated)
//! - Fraud detection (labeled as fraud vs normal)

use cbad_core::anomaly::{AnomalyConfig, AnomalyDetector};
use cbad_core::window::WindowConfig;
use serde::Deserialize;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::path::Path;

#[derive(Debug, Default)]
struct BenchmarkResult {
    true_positives: usize,
    false_positives: usize,
    true_negatives: usize,
    false_negatives: usize,
    total_events: usize,
}

impl BenchmarkResult {
    fn precision(&self) -> f64 {
        let denom = self.true_positives + self.false_positives;
        if denom == 0 {
            0.0
        } else {
            self.true_positives as f64 / denom as f64
        }
    }

    fn recall(&self) -> f64 {
        let denom = self.true_positives + self.false_negatives;
        if denom == 0 {
            0.0
        } else {
            self.true_positives as f64 / denom as f64
        }
    }

    fn f1(&self) -> f64 {
        let p = self.precision();
        let r = self.recall();
        if p + r == 0.0 {
            0.0
        } else {
            2.0 * p * r / (p + r)
        }
    }

    fn accuracy(&self) -> f64 {
        let correct = self.true_positives + self.true_negatives;
        if self.total_events == 0 {
            0.0
        } else {
            correct as f64 / self.total_events as f64
        }
    }

    fn print(&self, name: &str) {
        println!("\n📊 {} Results:", name);
        println!("   Total events: {}", self.total_events);
        println!(
            "   True Positives:  {} (anomalies correctly detected)",
            self.true_positives
        );
        println!(
            "   False Positives: {} (normal flagged as anomaly)",
            self.false_positives
        );
        println!(
            "   True Negatives:  {} (normal correctly passed)",
            self.true_negatives
        );
        println!(
            "   False Negatives: {} (anomalies missed)",
            self.false_negatives
        );
        println!("   ---");
        println!("   Precision: {:.1}%", self.precision() * 100.0);
        println!("   Recall:    {:.1}%", self.recall() * 100.0);
        println!("   F1 Score:  {:.1}%", self.f1() * 100.0);
        println!("   Accuracy:  {:.1}%", self.accuracy() * 100.0);
    }
}

fn main() {
    println!("\n=== CBAD Labeled Dataset Benchmark ===\n");
    println!("Testing whether CBAD detects known anomalies in labeled datasets.\n");

    // Test 1: Jailbreak prompt detection
    benchmark_jailbreak();

    // Test 2: Hallucination detection
    benchmark_hallucination();

    // Test 3: Fraud detection
    benchmark_fraud();

    println!("\n=== Benchmark Complete ===\n");
}

/// Benchmark jailbreak prompt detection
///
/// Strategy: Train on regular prompts, test on jailbreak prompts
/// Expected: Jailbreak prompts should be detected as anomalies
fn benchmark_jailbreak() {
    println!("🔒 JAILBREAK PROMPT DETECTION");
    println!("   Training: Regular user prompts");
    println!("   Testing: Mix of regular + jailbreak prompts");

    let base_path = Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets/jailbreak/data/prompts");

    let regular_path = base_path.join("regular_prompts_2023_12_25.csv");
    let jailbreak_path = base_path.join("jailbreak_prompts_2023_12_25.csv");

    if !regular_path.exists() || !jailbreak_path.exists() {
        println!("   ⚠️  Dataset not found at {:?}", base_path);
        return;
    }

    // Load prompts (need more data for baseline)
    let regular_prompts = load_csv_column(&regular_path, "prompt", 2000);
    let jailbreak_prompts = load_csv_column(&jailbreak_path, "prompt", 500);

    println!("   Loaded {} regular prompts", regular_prompts.len());
    println!("   Loaded {} jailbreak prompts", jailbreak_prompts.len());

    if regular_prompts.len() < 100 || jailbreak_prompts.len() < 50 {
        println!("   ⚠️  Not enough data");
        return;
    }

    // Create detector trained on regular prompts
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

    let detector = AnomalyDetector::new(config).expect("create detector");

    // Train on regular prompts (first 300)
    for prompt in regular_prompts.iter().take(300) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
    }

    // Now test on remaining regular (should NOT be anomalies) and jailbreak (SHOULD be anomalies)
    let mut result = BenchmarkResult::default();

    // Test regular prompts (label = 0 = normal)
    for prompt in regular_prompts.iter().skip(300).take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.false_positives += 1; // Normal flagged as anomaly
                } else {
                    result.true_negatives += 1; // Normal correctly passed
                }
            }
        }
    }

    // Test jailbreak prompts (label = 1 = anomaly)
    for prompt in jailbreak_prompts.iter().take(100) {
        let _ = detector.add_data(prompt.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.true_positives += 1; // Jailbreak detected
                } else {
                    result.false_negatives += 1; // Jailbreak missed
                }
            }
        }
    }

    result.print("Jailbreak Detection");
}

/// Benchmark hallucination detection
///
/// Strategy: Train on correct answers, test on hallucinated answers
fn benchmark_hallucination() {
    println!("\n🤥 HALLUCINATION DETECTION");
    println!("   Training: Correct QA answers");
    println!("   Testing: Mix of correct + hallucinated answers");

    let qa_path = Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets/halueval/data/qa_data.json");

    if !qa_path.exists() {
        println!("   ⚠️  Dataset not found at {:?}", qa_path);
        return;
    }

    // Load QA pairs
    #[derive(Deserialize)]
    struct QaPair {
        question: String,
        right_answer: String,
        hallucinated_answer: String,
    }

    let file = File::open(&qa_path).expect("open file");
    let reader = BufReader::new(file);
    let mut correct_answers = Vec::new();
    let mut hallucinated_answers = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i >= 500 {
            break;
        }
        if let Ok(line) = line {
            if let Ok(qa) = serde_json::from_str::<QaPair>(&line) {
                let correct = format!("Q: {} A: {}", qa.question, qa.right_answer);
                let hallu = format!("Q: {} A: {}", qa.question, qa.hallucinated_answer);
                correct_answers.push(correct);
                hallucinated_answers.push(hallu);
            }
        }
    }

    println!("   Loaded {} correct answers", correct_answers.len());
    println!(
        "   Loaded {} hallucinated answers",
        hallucinated_answers.len()
    );

    if correct_answers.len() < 100 {
        println!("   ⚠️  Not enough data");
        return;
    }

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

    let detector = AnomalyDetector::new(config).expect("create detector");

    // Train on correct answers
    for answer in correct_answers.iter().take(200) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
    }

    let mut result = BenchmarkResult::default();

    // Test remaining correct answers (should NOT be anomalies)
    for answer in correct_answers.iter().skip(200).take(100) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.false_positives += 1;
                } else {
                    result.true_negatives += 1;
                }
            }
        }
    }

    // Test hallucinated answers (SHOULD be anomalies)
    for answer in hallucinated_answers.iter().skip(200).take(100) {
        let _ = detector.add_data(answer.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.true_positives += 1;
                } else {
                    result.false_negatives += 1;
                }
            }
        }
    }

    result.print("Hallucination Detection");
}

/// Benchmark fraud detection
///
/// Strategy: Train on normal transactions, test on fraud
fn benchmark_fraud() {
    println!("\n💳 FRAUD DETECTION");
    println!("   Training: Normal transactions");
    println!("   Testing: Mix of normal + fraud transactions");

    let fraud_path = Path::new(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmark-datasets/../driftlock-archives/test-data/fraud/fraud_data.csv");

    // Try alternate path
    let fraud_path = if !fraud_path.exists() {
        Path::new("/Volumes/VIXinSSD/driftlock-archives/test-data/fraud/fraud_data.csv")
            .to_path_buf()
    } else {
        fraud_path
    };

    if !fraud_path.exists() {
        println!("   ⚠️  Dataset not found at {:?}", fraud_path);
        return;
    }

    // Load transactions
    let file = File::open(&fraud_path).expect("open file");
    let reader = BufReader::new(file);
    let mut normal_txns = Vec::new();
    let mut fraud_txns = Vec::new();

    for (i, line) in reader.lines().enumerate() {
        if i == 0 {
            continue;
        } // Skip header
        if normal_txns.len() >= 1000 && fraud_txns.len() >= 300 {
            break;
        }

        if let Ok(line) = line {
            let parts: Vec<&str> = line.split(',').collect();
            if parts.len() < 15 {
                continue;
            }

            // is_fraud is last column
            let is_fraud = parts.last().map(|s| s.trim()).unwrap_or("0");

            // Create text representation
            let text = format!(
                "merchant={} category={} amount={} city={} state={}",
                parts.get(1).unwrap_or(&""),
                parts.get(2).unwrap_or(&""),
                parts.get(3).unwrap_or(&""),
                parts.get(4).unwrap_or(&""),
                parts.get(5).unwrap_or(&""),
            );

            if is_fraud == "1" && fraud_txns.len() < 300 {
                fraud_txns.push(text);
            } else if is_fraud == "0" && normal_txns.len() < 1000 {
                normal_txns.push(text);
            }
        }
    }

    println!("   Loaded {} normal transactions", normal_txns.len());
    println!("   Loaded {} fraud transactions", fraud_txns.len());

    if normal_txns.len() < 200 || fraud_txns.len() < 50 {
        println!("   ⚠️  Not enough data");
        return;
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

    let detector = AnomalyDetector::new(config).expect("create detector");

    // Train on normal transactions
    for txn in normal_txns.iter().take(500) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
    }

    let mut result = BenchmarkResult::default();

    // Test remaining normal transactions (should NOT be anomalies)
    for txn in normal_txns.iter().skip(500).take(200) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.false_positives += 1;
                } else {
                    result.true_negatives += 1;
                }
            }
        }
    }

    // Test fraud transactions (SHOULD be anomalies)
    for txn in fraud_txns.iter().take(200) {
        let _ = detector.add_data(txn.as_bytes().to_vec());
        result.total_events += 1;

        if detector.is_ready().unwrap_or(false) {
            if let Ok(Some(detection)) = detector.detect_anomaly() {
                if detection.is_anomaly {
                    result.true_positives += 1;
                } else {
                    result.false_negatives += 1;
                }
            }
        }
    }

    result.print("Fraud Detection");
}

/// Load a specific column from a CSV file using proper CSV parsing
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

    // Find column index
    let headers = match reader.headers() {
        Ok(h) => h.clone(),
        Err(_) => return Vec::new(),
    };

    let col_idx = headers.iter().position(|h| h.trim() == column);
    let col_idx = match col_idx {
        Some(i) => i,
        None => return Vec::new(),
    };

    // Read data
    let mut result = Vec::new();
    for record in reader.records().take(max_rows) {
        if let Ok(record) = record {
            if let Some(val) = record.get(col_idx) {
                let val = val.trim();
                if !val.is_empty() && val.len() > 10 {
                    // Skip very short prompts
                    result.push(val.to_string());
                }
            }
        }
    }

    result
}
