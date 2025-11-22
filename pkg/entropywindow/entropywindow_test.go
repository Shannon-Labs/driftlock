package entropywindow

import (
	"strings"
	"testing"
)

func TestAnalyzerDetectsHighEntropyLines(t *testing.T) {
	analyzer, err := NewAnalyzer(Config{BaselineLines: 10, Threshold: 0.2, CompressionAlgorithm: "zstd"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { analyzer.Close() })

	baseline := []string{
		"INFO service=auth status=ok",
		"INFO service=auth latency=21",
		"INFO service=payments latency=10",
		"INFO service=api latency=15",
		"INFO service=api latency=13",
		"INFO service=api latency=14",
		"INFO service=api latency=16",
		"INFO service=api latency=11",
		"INFO service=api latency=15",
		"INFO service=api latency=12",
	}

	for _, line := range baseline {
		res := analyzer.Process(line)
		if res.Ready && res.IsAnomaly {
			t.Fatalf("baseline produced anomaly: %+v", res)
		}
	}

	spike := "ERROR !!! random payload " + randomBlob(128)
	res := analyzer.Process(spike)
	if !res.Ready {
		t.Fatalf("analyzer not ready after baseline")
	}
	if !res.IsAnomaly {
		t.Fatalf("expected anomaly, got %+v", res)
	}
}

func TestAnalyzerWarmsBeforeScoring(t *testing.T) {
	analyzer, err := NewAnalyzer(Config{BaselineLines: 3, Threshold: 0.2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer analyzer.Close()

	inputs := []string{"a", "b", "c", "d"}
	for idx, line := range inputs {
		res := analyzer.Process(line)
		if idx < 2 && res.Ready {
			t.Fatalf("result ready too early: %+v", res)
		}
	}
}

func randomBlob(length int) string {
	var out strings.Builder
	for i := 0; i < length; i++ {
		out.WriteByte(byte(32 + (i*73)%90))
	}
	return out.String()
}
