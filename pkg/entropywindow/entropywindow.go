package entropywindow

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
)

const (
	defaultBaselineLines = 400
	defaultThreshold     = 0.35
	defaultAlgo          = "zstd"
	maxEntropyBits       = 8.0 // byte alphabet
)

// Config controls the sliding-window entropy analysis.
type Config struct {
	BaselineLines        int
	Threshold            float64
	CompressionAlgorithm string
	MinLineLength        int
	IgnoreEmptyLines     bool
}

// Result represents the analysis output for a single record.
type Result struct {
	Sequence                 int     `json:"sequence"`
	Line                     string  `json:"line"`
	Entropy                  float64 `json:"entropy"`
	BaselineEntropy          float64 `json:"baseline_entropy"`
	CompressionRatio         float64 `json:"compression_ratio"`
	BaselineCompressionRatio float64 `json:"baseline_compression_ratio"`
	EntropyDelta             float64 `json:"entropy_delta"`
	CompressionDelta         float64 `json:"compression_delta"`
	Score                    float64 `json:"score"`
	IsAnomaly                bool    `json:"is_anomaly"`
	Ready                    bool    `json:"ready"`
	Reason                   string  `json:"reason"`
}

// Analyzer maintains sliding baseline statistics and scores new lines.
type Analyzer struct {
	cfg            Config
	samples        []sample
	idx            int
	count          int
	entropySum     float64
	compressionSum float64
	sequence       int
	compressor     compressor
	mu             sync.Mutex
}

type sample struct {
	entropy  float64
	compress float64
}

// NewAnalyzer returns an Analyzer configured with sane defaults.
func NewAnalyzer(cfg Config) (*Analyzer, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	comp, err := newCompressor(cfg.CompressionAlgorithm)
	if err != nil {
		return nil, err
	}
	return &Analyzer{
		cfg:        cfg,
		samples:    make([]sample, cfg.BaselineLines),
		compressor: comp,
	}, nil
}

// Close releases compressor resources.
func (a *Analyzer) Close() error {
	if closer, ok := a.compressor.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// Process analyzes a single line and updates the sliding baseline.
func (a *Analyzer) Process(line string) Result {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.sequence++
	clean := sanitizeLine(line, a.cfg.IgnoreEmptyLines)
	res := Result{
		Sequence: a.sequence,
		Line:     clean,
	}

	metrics := a.measure(clean)
	res.Entropy = metrics.entropy
	res.CompressionRatio = metrics.compress

	if a.count > 0 {
		res.BaselineEntropy = a.entropySum / float64(a.count)
		res.BaselineCompressionRatio = a.compressionSum / float64(a.count)
	}

	ready := a.count >= a.cfg.BaselineLines
	res.Ready = ready

	if ready {
		res.EntropyDelta = normalizeEntropyDelta(res.Entropy, res.BaselineEntropy)
		res.CompressionDelta = normalizeCompressionDelta(res.CompressionRatio, res.BaselineCompressionRatio)
		res.Score = clamp((res.EntropyDelta + res.CompressionDelta) / 2)
		res.IsAnomaly = res.Score >= a.cfg.Threshold && len(clean) >= a.cfg.MinLineLength
		if res.IsAnomaly {
			res.Reason = fmt.Sprintf("ΔH=%.2f ΔCR=%.2f (score=%.2f)", res.EntropyDelta, res.CompressionDelta, res.Score)
		} else {
			res.Reason = "Within baseline variance"
		}
	} else {
		res.Reason = fmt.Sprintf("Warming baseline (%d/%d)", a.count, a.cfg.BaselineLines)
	}

	a.push(metrics)
	return res
}

// ProcessJSON canonicalizes JSON input before analysis.
func (a *Analyzer) ProcessJSON(raw string) (Result, error) {
	var payload interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return Result{}, err
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return Result{}, err
	}
	return a.Process(string(normalized)), nil
}

func (a *Analyzer) measure(line string) sample {
	entropy := shannonEntropy(line)
	compress, _ := a.compressor.ratio([]byte(line))
	return sample{entropy: entropy, compress: compress}
}

func (a *Analyzer) push(s sample) {
	if len(a.samples) == 0 {
		return
	}
	if a.count < len(a.samples) {
		a.samples[a.count] = s
		a.count++
		a.entropySum += s.entropy
		a.compressionSum += s.compress
		return
	}
	old := a.samples[a.idx]
	a.entropySum += s.entropy - old.entropy
	a.compressionSum += s.compress - old.compress
	a.samples[a.idx] = s
	a.idx = (a.idx + 1) % len(a.samples)
}

func sanitizeLine(line string, skipEmpty bool) string {
	clean := strings.TrimRight(line, "\r\n")
	if skipEmpty && strings.TrimSpace(clean) == "" {
		return ""
	}
	return clean
}

func normalizeEntropyDelta(value, baseline float64) float64 {
	if baseline <= 0 {
		return clamp(value / maxEntropyBits)
	}
	return clamp(math.Abs(value-baseline) / maxEntropyBits)
}

func normalizeCompressionDelta(value, baseline float64) float64 {
	if baseline <= 0 {
		return clamp(value)
	}
	return clamp(math.Abs(value-baseline) / math.Max(baseline, 1e-9))
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func validateConfig(cfg *Config) error {
	if cfg.BaselineLines <= 0 {
		cfg.BaselineLines = defaultBaselineLines
	}
	if cfg.Threshold <= 0 {
		cfg.Threshold = defaultThreshold
	}
	if cfg.CompressionAlgorithm == "" {
		cfg.CompressionAlgorithm = defaultAlgo
	}
	if cfg.MinLineLength <= 0 {
		cfg.MinLineLength = 12
	}
	switch cfg.CompressionAlgorithm {
	case "zstd", "gzip":
	default:
		return fmt.Errorf("unsupported compression algorithm %q", cfg.CompressionAlgorithm)
	}
	return nil
}

type compressor interface {
	ratio([]byte) (float64, error)
}

type zstdCompressor struct {
	enc *zstd.Encoder
}

func newCompressor(algo string) (compressor, error) {
	switch algo {
	case "zstd":
		enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
		if err != nil {
			return nil, err
		}
		return &zstdCompressor{enc: enc}, nil
	case "gzip":
		return gzipCompressor{}, nil
	default:
		return nil, errors.New("unsupported compression algorithm")
	}
}

func (z *zstdCompressor) ratio(data []byte) (float64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	compressed := z.enc.EncodeAll(data, make([]byte, 0, len(data)))
	return float64(len(compressed)) / float64(len(data)), nil
}

func (z *zstdCompressor) Close() error {
	return z.enc.Close()
}

type gzipCompressor struct{}

func (gzipCompressor) ratio(data []byte) (float64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return 0, err
	}
	if _, err := w.Write(data); err != nil {
		return 0, err
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	return float64(buf.Len()) / float64(len(data)), nil
}

func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	var counts [256]int
	for i := 0; i < len(s); i++ {
		counts[s[i]]++
	}
	length := float64(len(s))
	entropy := 0.0
	for _, count := range counts {
		if count == 0 {
			continue
		}
		p := float64(count) / length
		entropy -= p * math.Log2(p)
	}
	return entropy
}
