package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/driftlock/driftlock/pkg/entropywindow"
	"github.com/spf13/cobra"
)

var (
	scanFormat     string
	scanFollow     bool
	scanThreshold  float64
	scanBaseline   int
	scanAlgo       string
	scanOutput     string
	scanStdin      bool
	scanShowAll    bool
	scanMinLineLen int
	scanAudio      bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [file]",
	Short: "Analyze raw logs or text streams for entropy spikes",
	Long: `Stream raw log lines through Driftlock's entropy window engine.

Examples:
  # Pipe journalctl output into the analyzer
  journalctl -u api.service | driftlock scan --format raw --threshold 0.4

  # Follow a file and emit NDJSON anomalies with audio alerts
  driftlock scan /var/log/app.log --follow --output ndjson --audio
`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringVar(&scanFormat, "format", "raw", "Input format: raw|ndjson|json")
	scanCmd.Flags().BoolVar(&scanFollow, "follow", false, "Keep the file handle open and stream appended lines")
	scanCmd.Flags().Float64Var(&scanThreshold, "threshold", 0.35, "Anomaly score threshold (0-1)")
	scanCmd.Flags().IntVar(&scanBaseline, "baseline-lines", 400, "Number of historical lines kept in the baseline window")
	scanCmd.Flags().StringVar(&scanAlgo, "algo", "zstd", "Compression algorithm: zstd|gzip")
	scanCmd.Flags().StringVar(&scanOutput, "output", "ndjson", "Output format: ndjson|pretty")
	scanCmd.Flags().BoolVar(&scanStdin, "stdin", false, "Force reading from STDIN even if a file path is provided")
	scanCmd.Flags().BoolVar(&scanShowAll, "show-all", false, "Emit every line instead of anomalies only")
	scanCmd.Flags().IntVar(&scanMinLineLen, "min-line-length", 12, "Minimum line length required before evaluating anomalies")
	scanCmd.Flags().BoolVar(&scanAudio, "audio", false, "Play a system sound when an anomaly is detected (macOS only)")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	format := strings.ToLower(scanFormat)
	if format != "raw" && format != "ndjson" && format != "json" {
		return fmt.Errorf("unsupported format %q", scanFormat)
	}
	if scanFollow && format == "json" {
		return errors.New("--follow is not supported for JSON payloads")
	}

	useStdin := scanStdin || len(args) == 0
	if !useStdin && len(args) == 0 {
		return errors.New("provide a file path or pipe data via --stdin")
	}

	var reader io.Reader
	var closer io.Closer
	if useStdin {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode()&os.ModeCharDevice) != 0 && len(args) == 0 {
			return errors.New("no STDIN data detected; pass a file or pipe data")
		}
		reader = os.Stdin
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		reader = file
		closer = file
	}
	if closer != nil {
		defer closer.Close()
	}

	analyzer, err := entropywindow.NewAnalyzer(entropywindow.Config{
		BaselineLines:        scanBaseline,
		Threshold:            scanThreshold,
		CompressionAlgorithm: scanAlgo,
		MinLineLength:        scanMinLineLen,
	})
	if err != nil {
		return err
	}
	defer analyzer.Close()

	emitter := newResultEmitter(scanOutput, scanShowAll, scanAudio)
	var handler func(string) error
	switch format {
	case "json":
		return processJSON(reader, analyzer, emitter)
	case "ndjson":
		handler = func(raw string) error {
			if strings.TrimSpace(raw) == "" {
				return nil
			}
			res, err := analyzer.ProcessJSON(raw)
			if err != nil {
				return err
			}
			return emitter.Emit(res)
		}
	default:
		handler = func(raw string) error {
			res := analyzer.Process(raw)
			return emitter.Emit(res)
		}
	}

	return streamLines(reader, handler, scanFollow)
}

func processJSON(r io.Reader, analyzer *entropywindow.Analyzer, emitter resultEmitter) error {
	decoder := json.NewDecoder(r)
	for {
		var payload interface{}
		if err := decoder.Decode(&payload); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		res := analyzer.Process(string(encoded))
		if err := emitter.Emit(res); err != nil {
			return err
		}
	}
}

func streamLines(r io.Reader, handler func(string) error, follow bool) error {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if err := handler(line); err != nil {
				return err
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) && follow {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

type resultEmitter struct {
	format    string
	showAll   bool
	playAudio bool
}

func newResultEmitter(format string, showAll bool, playAudio bool) resultEmitter {
	format = strings.ToLower(format)
	if format != "pretty" {
		format = "ndjson"
	}
	return resultEmitter{format: format, showAll: showAll, playAudio: playAudio}
}

func (re resultEmitter) Emit(res entropywindow.Result) error {
	if !res.IsAnomaly && !re.showAll {
		return nil
	}

	if res.IsAnomaly && re.playAudio {
		go playAlert()
	}

	switch re.format {
	case "pretty":
		line := res.Line
		if len(line) > 120 {
			line = line[:120] + "…"
		}
		status := "OK"
		if res.IsAnomaly {
			status = "ALERT"
		}
		fmt.Printf("[%s] #%d score=%.2f entropy=%.2f comp=%.2f :: %s\n",
			status, res.Sequence, res.Score, res.Entropy, res.CompressionRatio, line)
	default:
		b, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func playAlert() {
	// Simple fire-and-forget audio alert for macOS
	// Uses the system "Ping" sound or similar
	cmd := exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
	_ = cmd.Run()
}
