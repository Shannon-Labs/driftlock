//! Crypto Streaming Anomaly Detection Demo
//!
//! Demonstrates how to use CBAD for real-time streaming anomaly detection
//! on cryptocurrency price data.
//!
//! Run with: cargo run --example crypto_stream --release
//!
//! For live data, you'd connect to a WebSocket like:
//! - Binance: wss://stream.binance.com:9443/ws/btcusdt@trade
//! - Coinbase: wss://ws-feed.exchange.coinbase.com

use cbad_core::anomaly::{AnomalyConfig, AnomalyDetector};
use cbad_core::window::WindowConfig;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::path::Path;
use std::time::Instant;

/// Simulated crypto trade event (matches Binance trade format)
#[derive(Debug, Clone, Serialize, Deserialize)]
struct CryptoTrade {
    #[serde(rename = "s")]
    symbol: String,
    #[serde(rename = "p")]
    price: String,
    #[serde(rename = "q")]
    quantity: String,
    #[serde(rename = "T")]
    trade_time: u64,
    #[serde(rename = "m")]
    is_buyer_maker: bool,
}

/// Configuration for streaming detector
#[derive(Debug, Clone)]
pub struct StreamConfig {
    /// Number of events to establish baseline
    pub baseline_events: usize,
    /// Events per detection window
    pub window_events: usize,
    /// NCD threshold for anomaly (lower = more sensitive)
    pub ncd_threshold: f64,
    /// Whether to require statistical significance
    pub require_significance: bool,
}

impl Default for StreamConfig {
    fn default() -> Self {
        Self {
            baseline_events: 100,
            window_events: 30,
            ncd_threshold: 0.25,
            require_significance: true,
        }
    }
}

/// Streaming anomaly detector for crypto data
pub struct CryptoStreamDetector {
    detector: AnomalyDetector,
    event_count: usize,
    anomaly_count: usize,
    last_anomaly_ncd: f64,
}

impl CryptoStreamDetector {
    pub fn new(config: StreamConfig) -> Result<Self, String> {
        let anomaly_config = AnomalyConfig {
            window_config: WindowConfig {
                baseline_size: config.baseline_events,
                window_size: config.window_events,
                hop_size: config.window_events / 2,
                max_capacity: config.baseline_events + config.window_events * 3,
                ..Default::default()
            },
            permutation_count: 100,
            ncd_threshold: config.ncd_threshold,
            require_statistical_significance: config.require_significance,
            ..Default::default()
        };

        let detector = AnomalyDetector::new(anomaly_config)
            .map_err(|e| format!("Failed to create detector: {:?}", e))?;

        Ok(Self {
            detector,
            event_count: 0,
            anomaly_count: 0,
            last_anomaly_ncd: 0.0,
        })
    }

    /// Process a trade event and return anomaly status
    pub fn process_trade(&mut self, trade: &CryptoTrade) -> Option<AnomalyAlert> {
        // Serialize trade to bytes
        let data = serde_json::to_vec(trade).unwrap_or_default();

        // Add to detector
        if self.detector.add_data(data).is_err() {
            return None;
        }

        self.event_count += 1;

        // Check for anomalies
        if !self.detector.is_ready().unwrap_or(false) {
            return None;
        }

        if let Ok(Some(result)) = self.detector.detect_anomaly() {
            if result.is_anomaly {
                self.anomaly_count += 1;
                self.last_anomaly_ncd = result.metrics.ncd;

                return Some(AnomalyAlert {
                    event_index: self.event_count,
                    symbol: trade.symbol.clone(),
                    price: trade.price.clone(),
                    ncd: result.metrics.ncd,
                    p_value: result.metrics.p_value,
                    confidence: result.confidence_level,
                    severity: if result.metrics.ncd > 0.8 {
                        "HIGH"
                    } else if result.metrics.ncd > 0.5 {
                        "MEDIUM"
                    } else {
                        "LOW"
                    },
                });
            }
        }

        None
    }

    /// Get current stats
    pub fn stats(&self) -> StreamStats {
        StreamStats {
            events_processed: self.event_count,
            anomalies_detected: self.anomaly_count,
            anomaly_rate: if self.event_count > 0 {
                self.anomaly_count as f64 / self.event_count as f64
            } else {
                0.0
            },
            last_ncd: self.last_anomaly_ncd,
        }
    }
}

#[derive(Debug)]
pub struct AnomalyAlert {
    pub event_index: usize,
    pub symbol: String,
    pub price: String,
    pub ncd: f64,
    pub p_value: f64,
    pub confidence: f64,
    pub severity: &'static str,
}

#[derive(Debug)]
pub struct StreamStats {
    pub events_processed: usize,
    pub anomalies_detected: usize,
    pub anomaly_rate: f64,
    pub last_ncd: f64,
}

fn main() {
    println!("\n=== Crypto Streaming Anomaly Detection Demo ===\n");

    // Demo 1: Simulated live stream with Terra Luna data
    println!("📊 Demo 1: Terra Luna Price Collapse Detection\n");
    demo_terra_luna();

    // Demo 2: Show how you'd integrate with a real WebSocket
    println!("\n📡 Demo 2: WebSocket Integration Pattern\n");
    show_websocket_pattern();

    // Demo 3: Multi-stream detection (multiple trading pairs)
    println!("\n🔀 Demo 3: Multi-Stream Detection\n");
    demo_multi_stream();

    println!("\n=== Demo Complete ===\n");
}

fn demo_terra_luna() {
    // Load Terra Luna historical data
    let path = Path::new(env!("CARGO_MANIFEST_DIR")).join("../test-data/terra_luna/terra-luna.csv");

    if !path.exists() {
        println!("   ⚠️  Terra Luna data not found, using synthetic data");
        demo_synthetic_crash();
        return;
    }

    let file = File::open(&path).expect("Failed to open file");
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().filter_map(Result::ok).collect();

    println!(
        "   Loading {} price points from Terra Luna crash...",
        lines.len()
    );

    // Create detector with crypto-optimized settings
    let config = StreamConfig {
        baseline_events: 100, // ~4 days of hourly data
        window_events: 24,    // ~1 day window
        ncd_threshold: 0.20,  // Sensitive to changes
        require_significance: true,
    };

    let mut detector = CryptoStreamDetector::new(config).expect("Failed to create detector");

    let start = Instant::now();
    let mut alerts: Vec<(usize, f64)> = Vec::new();

    // Process each line as a simulated trade
    for (i, line) in lines.iter().enumerate() {
        // Parse CSV: timestamp,date,price
        let parts: Vec<&str> = line.split(',').collect();
        if parts.len() < 3 || parts[0] == "timestamp" {
            continue; // Skip header
        }

        let trade = CryptoTrade {
            symbol: "LUNA/USD".to_string(),
            price: parts[2].to_string(), // price is third column
            quantity: "1.0".to_string(),
            trade_time: parts[0].parse().unwrap_or(i as u64),
            is_buyer_maker: false,
        };

        if let Some(alert) = detector.process_trade(&trade) {
            alerts.push((i, alert.ncd));
            if alerts.len() <= 5 {
                println!(
                    "   🚨 ALERT at index {}: {} @ ${} | NCD={:.3} | Severity: {}",
                    alert.event_index, alert.symbol, alert.price, alert.ncd, alert.severity
                );
            }
        }
    }

    let elapsed = start.elapsed();
    let stats = detector.stats();

    println!("\n   📈 Results:");
    println!("      - Events processed: {}", stats.events_processed);
    println!("      - Anomalies detected: {}", stats.anomalies_detected);
    println!("      - Anomaly rate: {:.2}%", stats.anomaly_rate * 100.0);
    if !alerts.is_empty() {
        println!("      - First alert at index: {}", alerts[0].0);
        println!(
            "      - Peak NCD: {:.3}",
            alerts.iter().map(|(_, ncd)| *ncd).fold(0.0_f64, f64::max)
        );
    }
    println!("      - Processing time: {:?}", elapsed);
    println!(
        "      - Throughput: {:.0} events/sec",
        stats.events_processed as f64 / elapsed.as_secs_f64()
    );
}

fn demo_synthetic_crash() {
    use cbad_core::anomaly::synthetic::{generate_anomaly, generate_baseline, AnomalyType};

    let config = StreamConfig {
        baseline_events: 100,
        window_events: 30,
        ncd_threshold: 0.20,
        require_significance: false,
    };

    let mut detector = CryptoStreamDetector::new(config).expect("Failed to create detector");

    // Generate stable baseline (normal market)
    let baseline = generate_baseline(150, 42);
    for (i, event) in baseline.iter().enumerate() {
        let trade = CryptoTrade {
            symbol: "SYNTH/USD".to_string(),
            price: format!("{:.2}", 100.0 + (i as f64 * 0.01)),
            quantity: "1.0".to_string(),
            trade_time: i as u64,
            is_buyer_maker: i % 2 == 0,
        };

        // Also feed raw bytes
        let _ = detector.detector.add_data(event.clone());
    }

    println!("   Baseline established with 150 normal events");

    // Inject crash pattern
    let crash_events = generate_anomaly(50, AnomalyType::VolumeSpike, 43);
    let mut crash_detected = false;

    for (i, event) in crash_events.iter().enumerate() {
        let trade = CryptoTrade {
            symbol: "SYNTH/USD".to_string(),
            price: format!("{:.2}", 100.0 - (i as f64 * 2.0)), // Price dropping fast
            quantity: format!("{:.2}", 100.0 + (i as f64 * 10.0)), // Volume spike
            trade_time: (150 + i) as u64,
            is_buyer_maker: true, // All sells
        };

        if let Some(alert) = detector.process_trade(&trade) {
            if !crash_detected {
                println!(
                    "   🚨 CRASH DETECTED at index {}: {} dropped to ${} | NCD={:.3}",
                    alert.event_index, alert.symbol, alert.price, alert.ncd
                );
                crash_detected = true;
            }
        }
    }

    if crash_detected {
        println!("   ✅ Successfully detected synthetic crash pattern!");
    } else {
        println!("   ⚠️  Crash not detected - try adjusting thresholds");
    }
}

fn show_websocket_pattern() {
    println!("   Example WebSocket integration for Binance:");
    println!();
    println!("   ```rust");
    println!("   use tokio_tungstenite::connect_async;");
    println!("   use futures_util::StreamExt;");
    println!("");
    println!("   async fn stream_btc() -> Result<(), Box<dyn std::error::Error>> {{");
    println!("       let url = \"wss://stream.binance.com:9443/ws/btcusdt@trade\";");
    println!("       let (ws_stream, _) = connect_async(url).await?;");
    println!("       let (_, mut read) = ws_stream.split();");
    println!("");
    println!("       let mut detector = CryptoStreamDetector::new(StreamConfig::default())?;");
    println!("");
    println!("       while let Some(msg) = read.next().await {{");
    println!("           if let Ok(text) = msg?.into_text() {{");
    println!("               let trade: CryptoTrade = serde_json::from_str(&text)?;");
    println!("");
    println!("               if let Some(alert) = detector.process_trade(&trade) {{");
    println!("                   println!(\"ANOMALY: {{}} @ {{}} - NCD={{}}\",");
    println!("                       alert.symbol, alert.price, alert.ncd);");
    println!("               }}");
    println!("           }}");
    println!("       }}");
    println!("       Ok(())");
    println!("   }}");
    println!("   ```");
    println!();
    println!("   The streaming detector maintains state across calls,");
    println!("   so you get continuous anomaly detection on the live feed.");
}

fn demo_multi_stream() {
    // Demonstrate monitoring multiple trading pairs
    let pairs = ["BTC/USD", "ETH/USD", "SOL/USD"];

    println!(
        "   Monitoring {} trading pairs simultaneously...",
        pairs.len()
    );
    println!();

    for pair in pairs {
        // Each pair gets its own detector
        let config = StreamConfig {
            baseline_events: 50,
            window_events: 20,
            ncd_threshold: 0.25,
            require_significance: true,
        };

        let detector = CryptoStreamDetector::new(config);

        match detector {
            Ok(_) => println!("   ✅ {} detector initialized", pair),
            Err(e) => println!("   ❌ {} failed: {}", pair, e),
        }
    }

    println!();
    println!("   In production, each detector runs in its own async task:");
    println!();
    println!("   ```rust");
    println!("   let handles: Vec<_> = pairs.iter().map(|pair| {{");
    println!("       let pair = pair.to_string();");
    println!("       tokio::spawn(async move {{");
    println!("           let mut detector = CryptoStreamDetector::new(config)?;");
    println!("           let mut stream = connect_to_exchange(&pair).await?;");
    println!("");
    println!("           while let Some(trade) = stream.next().await {{");
    println!("               if let Some(alert) = detector.process_trade(&trade) {{");
    println!("                   alert_channel.send(alert).await?;");
    println!("               }}");
    println!("           }}");
    println!("           Ok::<_, Error>(())");
    println!("       }})");
    println!("   }}).collect();");
    println!("");
    println!("   futures::future::join_all(handles).await;");
    println!("   ```");
}
