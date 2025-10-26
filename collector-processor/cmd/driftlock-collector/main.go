package main

import (
	"fmt"
	"os"

	"github.com/your-org/driftlock/collector-processor/driftlockcbad"
)

func main() {
	factory := driftlockcbad.NewFactory()
	cfg := factory.CreateDefaultConfig()

	fmt.Fprintf(os.Stdout, "driftlock collector (placeholder)\n")
	fmt.Fprintf(os.Stdout, "component: %s, default config: %#v\n", factory.Type(), cfg)
	fmt.Fprintln(os.Stdout, "TODO: integrate with OpenTelemetry collector service wiring.")
}
