# Benchmark Harness

The benchmark tools evaluate throughput, latency, CPU, and memory utilization for the processor and CBAD core library.

## Planned Metrics

### Performance Benchmarks
- Events per second at varying window sizes
- p50 / p95 propagation latency through the OTel pipeline
- CPU and memory budgets for structured vs unstructured data paths
- Compression algorithm performance comparison
- FFI boundary overhead measurement

### Accuracy Benchmarks
- False positive/negative rates across different data types
- Statistical significance test validation
- Anomaly detection sensitivity analysis
- Cross-platform result consistency

### Scalability Benchmarks
- Multi-threaded processing performance
- Memory usage scaling with window size
- Distributed deployment performance
- Network overhead in clustered setups

## Benchmark Suites

### CBAD Core Benchmarks (Rust)
```bash
cd cbad-core
cargo bench
```

Measures:
- Compression adapter performance
- Metrics calculation throughput
- Permutation test efficiency
- Memory allocation patterns

### API Server Benchmarks (Go)
```bash
go test -bench=. ./api-server/...
```

Measures:
- HTTP endpoint latency
- Event ingestion throughput
- Database query performance
- WebSocket streaming capacity

### End-to-End Benchmarks
```bash
make benchmark
```

Measures:
- Full pipeline latency
- System resource utilization
- Integration overhead
- Real-world workload simulation

## Performance Targets

### Throughput Targets
- CBAD Core: >10k events/second on 4-core hardware
- API Server: >1k requests/second with <100ms p95 latency
- Collector Processor: >5k events/second with <400ms p95 latency

### Resource Targets
- Memory usage: <2GB for 1M event windows
- CPU utilization: <80% under normal load
- Disk I/O: Minimal during normal operation
- Network overhead: <10% of payload size

### Accuracy Targets
- False positive rate: <1% on synthetic benchmarks
- Reproducibility: 100% deterministic across runs
- Statistical significance: Proper p-value calibration
- Cross-platform consistency: Identical results across architectures

## Usage

### Running Benchmarks

```bash
# Full benchmark suite
make benchmark

# Component-specific benchmarks
make benchmark-cbad
make benchmark-api
make benchmark-e2e

# Performance regression testing
make benchmark-compare
```

### Analyzing Results

```bash
# Generate performance report
make benchmark-report

# Compare with baseline
make benchmark-baseline

# Export metrics for CI
make benchmark-export
```

### Continuous Benchmarking

CI pipeline includes:
- Performance regression detection
- Resource usage monitoring
- Accuracy validation
- Cross-platform consistency checks

## Contributing

When adding new benchmarks:
1. Follow the naming convention: `BenchmarkComponentFeature`
2. Use deterministic data sets with documented seeds
3. Include baseline performance expectations
4. Document resource requirements and constraints
5. Update this README with new benchmark descriptions