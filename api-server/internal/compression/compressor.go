package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
)

// Compressor defines the interface for different compression algorithms
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(compressed []byte) ([]byte, error)
	Algorithm() string
}

// ZstdCompressor implements Zstandard compression
type ZstdCompressor struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

// NewZstdCompressor creates a new Zstd compressor
func NewZstdCompressor() (*ZstdCompressor, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %w", err)
	}

	return &ZstdCompressor{
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (z *ZstdCompressor) Compress(data []byte) ([]byte, error) {
	return z.encoder.EncodeAll(data, nil), nil
}

func (z *ZstdCompressor) Decompress(compressed []byte) ([]byte, error) {
	return z.decoder.DecodeAll(compressed, nil)
}

func (z *ZstdCompressor) Algorithm() string {
	return "zstd"
}

// GzipCompressor implements Gzip compression
type GzipCompressor struct {
	level int // compression level
}

// NewGzipCompressor creates a new Gzip compressor with specified level
func NewGzipCompressor(level int) *GzipCompressor {
	return &GzipCompressor{level: level}
}

func (g *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, g.level)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *GzipCompressor) Decompress(compressed []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from gzip reader: %w", err)
	}

	return result, nil
}

func (g *GzipCompressor) Algorithm() string {
	return "gzip"
}

// LZ4Compressor implements LZ4 compression
type LZ4Compressor struct{}

// NewLZ4Compressor creates a new LZ4 compressor
func NewLZ4Compressor() *LZ4Compressor {
	return &LZ4Compressor{}
}

func (l *LZ4Compressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := lz4.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to lz4 writer: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close lz4 writer: %w", err)
	}

	return buf.Bytes(), nil
}

func (l *LZ4Compressor) Decompress(compressed []byte) ([]byte, error) {
	reader := lz4.NewReader(bytes.NewReader(compressed))
	
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from lz4 reader: %w", err)
	}

	return result, nil
}

func (l *LZ4Compressor) Algorithm() string {
	return "lz4"
}

// OpenZLCompressor implements OpenZL-based compression (interface only, actual implementation in Rust)
type OpenZLCompressor struct {
	modelID string
}

// NewOpenZLCompressor creates a new OpenZL compressor
func NewOpenZLCompressor(modelID string) *OpenZLCompressor {
	return &OpenZLCompressor{modelID: modelID}
}

func (o *OpenZLCompressor) Compress(data []byte) ([]byte, error) {
	// This would call into the Rust OpenZL implementation
	// For now, we'll use Zstd as a fallback
	// The actual OpenZL FFI call would be here
	compressor, err := NewZstdCompressor()
	if err != nil {
		return nil, err
	}
	return compressor.Compress(data)
}

func (o *OpenZLCompressor) Decompress(compressed []byte) ([]byte, error) {
	// This would call into the Rust OpenZL implementation
	// For now, we'll use Zstd as a fallback
	// The actual OpenZL FFI call would be here
	compressor, err := NewZstdCompressor()
	if err != nil {
		return nil, err
	}
	return compressor.Decompress(compressed)
}

func (o *OpenZLCompressor) Algorithm() string {
	return "openzl:" + o.modelID
}

// Pool-based Compressor that reuses resources
type PooledCompressor struct {
	compressor Compressor
	pool       sync.Pool
}

// NewPooledCompressor creates a compressor with internal resource pooling
func NewPooledCompressor(compressor Compressor) *PooledCompressor {
	return &PooledCompressor{
		compressor: compressor,
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (p *PooledCompressor) Compress(data []byte) ([]byte, error) {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer p.pool.Put(buf)

	writer, err := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close: %w", err)
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

func (p *PooledCompressor) Decompress(compressed []byte) ([]byte, error) {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer p.pool.Put(buf)

	buf.Write(compressed)

	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	return result, nil
}

func (p *PooledCompressor) Algorithm() string {
	return fmt.Sprintf("pooled_%s", p.compressor.Algorithm())
}

// CompressionStrategy determines the best compression algorithm for a given data type
type CompressionStrategy interface {
	SelectCompressor(dataType string) Compressor
	GetCompressor(algorithm string) Compressor
	RegisterCompressor(algorithm string, compressor Compressor)
}

// DefaultCompressionStrategy provides default algorithm selection
type DefaultCompressionStrategy struct {
	compressors map[string]Compressor
}

// NewDefaultCompressionStrategy creates a new strategy with default compressors
func NewDefaultCompressionStrategy() *DefaultCompressionStrategy {
	zstdComp, _ := NewZstdCompressor()
	
	return &DefaultCompressionStrategy{
		compressors: map[string]Compressor{
			"zstd":     zstdComp,
			"gzip":     NewGzipCompressor(gzip.BestSpeed),
			"lz4":      NewLZ4Compressor(),
			"openzl":   NewOpenZLCompressor("default"),
			"pooled":   NewPooledCompressor(zstdComp),
		},
	}
}

func (d *DefaultCompressionStrategy) SelectCompressor(dataType string) Compressor {
	switch dataType {
	case "logs":
		return d.compressors["lz4"]  // Fast compression for logs
	case "metrics":
		return d.compressors["zstd"] // High ratio for metrics
	case "traces":
		return d.compressors["lz4"]  // Fast for traces
	case "anomaly_data":
		return d.compressors["zstd"] // High ratio for anomaly data
	default:
		return d.compressors["gzip"] // Default
	}
}

func (d *DefaultCompressionStrategy) GetCompressor(algorithm string) Compressor {
	return d.compressors[algorithm]
}

func (d *DefaultCompressionStrategy) RegisterCompressor(algorithm string, compressor Compressor) {
	d.compressors[algorithm] = compressor
}

// CompressionManager orchestrates compression across the application
type CompressionManager struct {
	strategy CompressionStrategy
}

// NewCompressionManager creates a new compression manager
func NewCompressionManager(strategy CompressionStrategy) *CompressionManager {
	if strategy == nil {
		strategy = NewDefaultCompressionStrategy()
	}
	
	return &CompressionManager{
		strategy: strategy,
	}
}

// CompressData compresses data using the appropriate algorithm
func (c *CompressionManager) CompressData(data []byte, dataType string) ([]byte, string, error) {
	compressor := c.strategy.SelectCompressor(dataType)
	if compressor == nil {
		return nil, "", fmt.Errorf("no compressor found for data type: %s", dataType)
	}

	compressedData, err := compressor.Compress(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compress %s data: %w", dataType, err)
	}

	return compressedData, compressor.Algorithm(), nil
}

// DecompressData decompresses data using the specified algorithm
func (c *CompressionManager) DecompressData(compressedData []byte, algorithm string) ([]byte, error) {
	compressor := c.strategy.GetCompressor(algorithm)
	if compressor == nil {
		return nil, fmt.Errorf("no compressor found for algorithm: %s", algorithm)
	}

	decompressedData, err := compressor.Decompress(compressedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data with %s: %w", algorithm, err)
	}

	return decompressedData, nil
}

// GetCompressionRatio returns the compression ratio for the given data
func (c *CompressionManager) GetCompressionRatio(original []byte, compressed []byte) float64 {
	if len(original) == 0 {
		return 0
	}
	return float64(len(compressed)) / float64(len(original))
}