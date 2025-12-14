package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// convertTransactionsToNDJSON converts CSV or JSONL transaction datasets into
// per-event NDJSON ready to wrap into { "events": [...] } requests for Driftlock.
// It aims to stay dependency-free (stdlib only) and supports light coercion
// (timestamps, numbers, booleans).
func main() {
	input := flag.String("input", "", "Path to CSV or JSONL input file")
	output := flag.String("output", "", "Path to NDJSON output file")
	limit := flag.Int("limit", 1000, "Maximum number of records to emit")
	tsField := flag.String("timestamp", "timestamp", "Timestamp column/field name (CSV)")
	timeLayout := flag.String("time-layout", time.RFC3339, "Go time layout for parsing the timestamp field")
	dataset := flag.String("dataset", "", "Optional preset for known datasets (e.g. fraud)")
	flag.Parse()

	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}

	inExt := strings.ToLower(filepath.Ext(*input))
	switch inExt {
	case ".csv":
		if err := convertCSV(*input, *output, *limit, *tsField, *timeLayout, *dataset); err != nil {
			fail(err)
		}
	case ".jsonl", ".ndjson":
		if err := copyNDJSON(*input, *output, *limit); err != nil {
			fail(err)
		}
	default:
		fail(fmt.Errorf("unsupported input extension %q (expected .csv or .jsonl/.ndjson)", inExt))
	}

	fmt.Printf("Wrote NDJSON to %s\n", *output)
}

func convertCSV(input, output string, limit int, tsField, timeLayout, dataset string) error {
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	row := 0
	for {
		if limit > 0 && row >= limit {
			break
		}

		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read row %d: %w", row+1, err)
		}
		if len(record) != len(headers) {
			return fmt.Errorf("row %d: header/data length mismatch (%d vs %d)", row+1, len(headers), len(record))
		}

		event := map[string]interface{}{}
		for i, h := range headers {
			event[h] = coerce(record[i])
		}

		if preset := strings.ToLower(dataset); preset == "fraud" {
			if v, ok := event["is_fraud"].(float64); ok {
				event["label"] = v == 1
			}
		}

		// Normalize timestamp
		if raw, ok := event[tsField]; ok {
			switch val := raw.(type) {
			case string:
				if ts, err := time.Parse(timeLayout, val); err == nil {
					event["timestamp"] = ts.UTC().Format(time.RFC3339Nano)
				}
			case float64:
				// Treat numeric timestamps as seconds offset from now.
				base := time.Now().Add(-time.Duration(val) * time.Second)
				event["timestamp"] = base.UTC().Format(time.RFC3339Nano)
			}
		}

		b, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal row %d: %w", row+1, err)
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
		row++
	}

	return nil
}

func copyNDJSON(input, output string, limit int) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	sc := bufio.NewScanner(in)
	w := bufio.NewWriter(out)
	defer w.Flush()

	line := 0
	for sc.Scan() {
		if limit > 0 && line >= limit {
			break
		}
		if _, err := w.Write(sc.Bytes()); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
		line++
	}
	return sc.Err()
}

func coerce(val string) interface{} {
	if val == "" {
		return val
	}
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	lower := strings.ToLower(val)
	if lower == "true" || lower == "false" {
		return lower == "true"
	}
	// Unquote if it looks like a quoted string with doubled quotes (CSV)
	if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		return strings.Trim(val, "\"")
	}
	return val
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
