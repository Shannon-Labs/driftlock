module github.com/Shannon-Labs/driftlock

go 1.24.1

replace github.com/Shannon-Labs/driftlock => .

replace github.com/Shannon-Labs/driftlock/cbad-core => ./cbad-core

replace github.com/Shannon-Labs/driftlock/collector-processor => ./collector-processor

replace github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad => ./collector-processor/driftlockcbad

replace github.com/Shannon-Labs/driftlock/pkg/version => ./pkg/version

replace github.com/Shannon-Labs/driftlock/api-server => ./api-server

require github.com/lib/pq v1.10.9
