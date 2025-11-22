module github.com/driftlock/driftlock/cmd/driftlock-mcp

go 1.24.1

replace github.com/driftlock/driftlock/pkg/entropywindow => ../../pkg/entropywindow

require github.com/driftlock/driftlock/pkg/entropywindow v0.0.0-00010101000000-000000000000

require github.com/klauspost/compress v1.17.9 // indirect
