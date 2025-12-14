//! Statistical utilities for CBAD
//!
//! This module provides online statistics computation using Welford's algorithm
//! and approximate quantile estimation using reservoir sampling.

use serde::{Deserialize, Serialize};

/// Welford's online algorithm for computing running mean and variance.
///
/// This allows incremental updates without storing all historical values,
/// providing O(1) memory complexity regardless of the number of samples.
#[derive(Clone, Copy, Debug, Default, Serialize, Deserialize)]
pub struct Welford {
    n: u64,
    mean: f64,
    m2: f64,
}

impl Welford {
    /// Create a new Welford statistics tracker
    pub fn new() -> Self {
        Self {
            n: 0,
            mean: 0.0,
            m2: 0.0,
        }
    }

    /// Update the statistics with a new value
    pub fn update(&mut self, x: f64) {
        self.n += 1;
        let delta = x - self.mean;
        self.mean += delta / (self.n as f64);
        let delta2 = x - self.mean;
        self.m2 += delta * delta2;
    }

    /// Get the running mean
    pub fn mean(&self) -> f64 {
        self.mean
    }

    /// Get the sample variance (using n-1 for unbiased estimation)
    pub fn variance(&self) -> f64 {
        if self.n < 2 {
            0.0
        } else {
            self.m2 / ((self.n - 1) as f64)
        }
    }

    /// Get the sample standard deviation
    pub fn std(&self) -> f64 {
        self.variance().sqrt()
    }

    /// Get the number of samples seen
    pub fn count(&self) -> u64 {
        self.n
    }

    /// Reset the statistics
    pub fn reset(&mut self) {
        self.n = 0;
        self.mean = 0.0;
        self.m2 = 0.0;
    }
}

/// Calculate Shannon entropy in bits per byte.
///
/// Returns a value between 0.0 (completely predictable) and 8.0 (maximum entropy).
pub fn entropy_bits_per_byte(data: &[u8]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }

    let mut counts = [0usize; 256];
    for &b in data {
        counts[b as usize] += 1;
    }

    let n = data.len() as f64;
    let mut h = 0.0f64;
    for c in counts.iter().copied() {
        if c == 0 {
            continue;
        }
        let p = (c as f64) / n;
        h -= p * p.log2();
    }
    h // max 8.0 for uniform distribution
}

/// Fixed-capacity reservoir for approximate quantile estimation.
///
/// Uses a circular buffer to maintain the most recent samples,
/// providing approximate quantile computation without storing all history.
#[derive(Clone, Debug)]
pub struct SimpleQuantile {
    buf: Vec<f64>,
    cap: usize,
    next_idx: usize,
}

impl SimpleQuantile {
    /// Create a new quantile estimator with the specified capacity
    pub fn new_with_cap(cap: usize) -> Self {
        Self {
            buf: Vec::with_capacity(cap),
            cap,
            next_idx: 0,
        }
    }

    /// Add a value to the reservoir
    pub fn add(&mut self, x: f64) {
        if self.buf.len() < self.cap {
            self.buf.push(x);
        } else {
            self.buf[self.next_idx] = x;
            self.next_idx = (self.next_idx + 1) % self.cap;
        }
    }

    /// Get the approximate quantile (0.0 to 1.0)
    pub fn quantile(&self, q: f64) -> Option<f64> {
        if self.buf.is_empty() {
            return None;
        }
        let mut v = self.buf.clone();
        v.sort_by(|a, b| a.partial_cmp(b).unwrap());
        let n = v.len();
        let idx = ((q.clamp(0.0, 1.0)) * ((n - 1) as f64)).round() as usize;
        Some(v[idx])
    }

    /// Get the number of samples in the reservoir
    pub fn len(&self) -> usize {
        self.buf.len()
    }

    /// Check if the reservoir is empty
    pub fn is_empty(&self) -> bool {
        self.buf.is_empty()
    }

    /// Clear the reservoir
    pub fn clear(&mut self) {
        self.buf.clear();
        self.next_idx = 0;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_welford_basic() {
        let mut w = Welford::new();
        let data = vec![2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0];
        for x in &data {
            w.update(*x);
        }

        assert!((w.mean() - 5.0).abs() < 0.001);
        assert!((w.variance() - 4.571).abs() < 0.01);
        assert_eq!(w.count(), 8);
    }

    #[test]
    fn test_entropy() {
        // Uniform distribution should have high entropy
        let uniform: Vec<u8> = (0..=255).collect();
        let entropy = entropy_bits_per_byte(&uniform);
        assert!(entropy > 7.9); // Should be close to 8.0

        // Single repeated byte should have 0 entropy
        let constant = vec![0u8; 100];
        let entropy = entropy_bits_per_byte(&constant);
        assert!(entropy < 0.001);
    }

    #[test]
    fn test_simple_quantile() {
        let mut sq = SimpleQuantile::new_with_cap(100);
        for i in 0..100 {
            sq.add(i as f64);
        }

        assert!((sq.quantile(0.5).unwrap() - 50.0).abs() < 1.0);
        assert!((sq.quantile(0.0).unwrap() - 0.0).abs() < 0.001);
        assert!((sq.quantile(1.0).unwrap() - 99.0).abs() < 0.001);
    }
}
