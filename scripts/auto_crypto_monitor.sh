#!/bin/bash
# Master orchestration script for automated crypto anomaly detection
# Manages 4-hour live session with Kraken WebSocket + Driftlock entropy detector

set -euo pipefail

# Configuration
DURATION="${DURATION:-14400}"  # 4 hours in seconds
TEST_MODE="${TEST_MODE:-false}"
KRAKEN_PAIR="${KRAKEN_PAIR:-BTC/USD}"
BASELINE_LINES="${BASELINE_LINES:-120}"
COMPRESSION_ALGO="${COMPRESSION_ALGO:-zstd}"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --duration)
            DURATION="$2"
            shift 2
            ;;
        --test-mode)
            TEST_MODE=true
            shift
            ;;
        --pair)
            KRAKEN_PAIR="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BLUE='\033[0;34m'
NC='\033[0m'

# Session setup
SESSION_TIMESTAMP=$(date +%s)
SESSION_DIR="logs/session_${SESSION_TIMESTAMP}"
mkdir -p "$SESSION_DIR"

# Log files
RAW_LOG="$SESSION_DIR/kraken_raw.ndjson"
ANOMALY_LOG="$SESSION_DIR/kraken_anomalies.ndjson"
ALERT_LOG="logs/anomaly_alerts.log"
STATUS_LOG="$SESSION_DIR/status.log"
SYNTHETIC_LOG="$SESSION_DIR/synthetic.log"
LOOM_LOG="$SESSION_DIR/loom.log"
MASTER_LOG="$SESSION_DIR/master.log"

# PID files
STREAMER_PID_FILE="$SESSION_DIR/streamer.pid"
DETECTOR_PID_FILE="$SESSION_DIR/detector.pid"
ALERTER_PID_FILE="$SESSION_DIR/alerter.pid"
STATUS_PID_FILE="$SESSION_DIR/status.pid"
SYNTHETIC_PID_FILE="$SESSION_DIR/synthetic.pid"
LOOM_PID_FILE="$SESSION_DIR/loom.pid"
RECORDER_PID_FILE="$SESSION_DIR/recorder.pid"

# Named pipes
KRAKEN_PIPE="/tmp/kraken_${SESSION_TIMESTAMP}.pipe"

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}🛑 Shutting down session...${NC}" | tee -a "$MASTER_LOG"
    
    # Kill all child processes
    for pid_file in "$STREAMER_PID_FILE" "$DETECTOR_PID_FILE" "$ALERTER_PID_FILE" "$STATUS_PID_FILE" "$SYNTHETIC_PID_FILE" "$LOOM_PID_FILE" "$RECORDER_PID_FILE"; do
        if [ -f "$pid_file" ]; then
            pid=$(cat "$pid_file")
            if kill -0 "$pid" 2>/dev/null; then
                echo "   Stopping process $pid..." | tee -a "$MASTER_LOG"
                kill "$pid" 2>/dev/null || true
            fi
            rm -f "$pid_file"
        fi
    done
    
    # Clean up named pipes
    rm -f "$KRAKEN_PIPE"
    
    # Final statistics
    echo -e "\n${CYAN}📊 SESSION SUMMARY${NC}" | tee -a "$MASTER_LOG"
    echo "   Duration: $(($(date +%s) - SESSION_TIMESTAMP))s" | tee -a "$MASTER_LOG"
    echo "   Session Dir: $SESSION_DIR" | tee -a "$MASTER_LOG"
    
    if [ -f "$RAW_LOG" ]; then
        echo "   Total Trades: $(wc -l < "$RAW_LOG" | tr -d ' ')" | tee -a "$MASTER_LOG"
    fi
    
    if [ -f "$ANOMALY_LOG" ]; then
        echo "   Total Anomalies: $(grep -c '"anomaly".*true\|"is_anomaly".*true\|"detected".*true' "$ANOMALY_LOG" 2>/dev/null || echo "0")" | tee -a "$MASTER_LOG"
    fi
    
    echo -e "${GREEN}✅ Session complete${NC}" | tee -a "$MASTER_LOG"
    exit 0
}

trap cleanup SIGINT SIGTERM EXIT

# Banner
echo -e "${BLUE}"
echo "═══════════════════════════════════════════════════════════"
echo "  🚀 AUTOMATED CRYPTO ANOMALY DETECTION SESSION"
echo "═══════════════════════════════════════════════════════════"
echo -e "${NC}"
echo "   Pair: $KRAKEN_PAIR"
echo "   Duration: $((DURATION / 3600))h $((DURATION % 3600 / 60))m"
echo "   Baseline: $BASELINE_LINES lines"
echo "   Compression: $COMPRESSION_ALGO"
echo "   Session Dir: $SESSION_DIR"
echo "   Started: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
if command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${GREEN}   🎥 Automatic screen recording enabled${NC}"
else
    echo -e "${YELLOW}   📹 Loom will be prompted on anomaly detection${NC}"
fi
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "" | tee "$MASTER_LOG"

# Export environment for child processes
export SESSION_DIR
export RAW_LOG
export ANOMALY_LOG
export ALERT_LOG
export STATUS_LOG
export SYNTHETIC_LOG
export LOOM_LOG
export STREAMER_PID_FILE
export KRAKEN_PAIR

# Create named pipe
rm -f "$KRAKEN_PIPE"
mkfifo "$KRAKEN_PIPE"

# Step 1: Start Kraken WebSocket streamer
echo -e "${CYAN}[1/6] Starting Kraken WebSocket streamer...${NC}" | tee -a "$MASTER_LOG"
python3 -u scripts/stream_kraken_ws.py 2>"$SESSION_DIR/streamer.stderr.log" | \
    tee -a "$RAW_LOG" > "$KRAKEN_PIPE" &

STREAMER_PID=$!
echo "$STREAMER_PID" > "$STREAMER_PID_FILE"
echo "   PID: $STREAMER_PID" | tee -a "$MASTER_LOG"
sleep 2

# Verify streamer is running
if ! kill -0 "$STREAMER_PID" 2>/dev/null; then
    echo -e "${RED}   ❌ Streamer failed to start${NC}" | tee -a "$MASTER_LOG"
    cat "$SESSION_DIR/streamer.stderr.log"
    exit 1
fi
echo -e "   ${GREEN}✅ Streamer running${NC}" | tee -a "$MASTER_LOG"
echo ""

# Step 2: Start Driftlock detector
echo -e "${CYAN}[2/6] Starting Driftlock entropy detector...${NC}" | tee -a "$MASTER_LOG"
cat "$KRAKEN_PIPE" | \
    ./bin/driftlock scan --stdin --follow --format ndjson \
        --baseline-lines "$BASELINE_LINES" \
        --algo "$COMPRESSION_ALGO" \
        --show-all 2>"$SESSION_DIR/detector.stderr.log" | \
    tee -a "$ANOMALY_LOG" | \
    tee >(./scripts/alert_on_anomaly.sh 2>"$SESSION_DIR/alerter.stderr.log" &
         echo $! > "$ALERTER_PID_FILE") > /dev/null &

DETECTOR_PID=$!
echo "$DETECTOR_PID" > "$DETECTOR_PID_FILE"
echo "   PID: $DETECTOR_PID" | tee -a "$MASTER_LOG"
sleep 2

if ! kill -0 "$DETECTOR_PID" 2>/dev/null; then
    echo -e "${RED}   ❌ Detector failed to start${NC}" | tee -a "$MASTER_LOG"
    cat "$SESSION_DIR/detector.stderr.log"
    exit 1
fi
echo -e "   ${GREEN}✅ Detector running${NC}" | tee -a "$MASTER_LOG"

# Wait for alerter to start
sleep 2
if [ -f "$ALERTER_PID_FILE" ]; then
    ALERTER_PID=$(cat "$ALERTER_PID_FILE")
    if kill -0 "$ALERTER_PID" 2>/dev/null; then
        echo "   Alerter PID: $ALERTER_PID" | tee -a "$MASTER_LOG"
        echo -e "   ${GREEN}✅ Alerter running${NC}" | tee -a "$MASTER_LOG"
    fi
fi
echo ""

# Step 3: Start status reporter
echo -e "${CYAN}[3/6] Starting status reporter (15-minute intervals)...${NC}" | tee -a "$MASTER_LOG"
./scripts/status_reporter.sh 2>"$SESSION_DIR/status.stderr.log" &
STATUS_PID=$!
echo "$STATUS_PID" > "$STATUS_PID_FILE"
echo "   PID: $STATUS_PID" | tee -a "$MASTER_LOG"
echo -e "   ${GREEN}✅ Status reporter running${NC}" | tee -a "$MASTER_LOG"
echo ""

# Step 4: Start synthetic anomaly injector
if [ "$TEST_MODE" = false ]; then
    echo -e "${CYAN}[4/6] Starting synthetic anomaly injector (30-min threshold)...${NC}" | tee -a "$MASTER_LOG"
    ./scripts/synthetic_anomaly.sh 2>"$SESSION_DIR/synthetic.stderr.log" &
    SYNTHETIC_PID=$!
    echo "$SYNTHETIC_PID" > "$SYNTHETIC_PID_FILE"
    echo "   PID: $SYNTHETIC_PID" | tee -a "$MASTER_LOG"
    echo -e "   ${GREEN}✅ Synthetic injector running${NC}" | tee -a "$MASTER_LOG"
else
    echo -e "${YELLOW}[4/6] Skipping synthetic injector (test mode)${NC}" | tee -a "$MASTER_LOG"
fi
echo ""

# Step 5: Start auto screen recorder (if ffmpeg available)
if command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${CYAN}[5/7] Starting auto screen recorder...${NC}" | tee -a "$MASTER_LOG"
    ./scripts/auto_screen_recorder.sh 2>"$SESSION_DIR/recorder.stderr.log" &
    RECORDER_PID=$!
    echo "$RECORDER_PID" > "$RECORDER_PID_FILE"
    echo "   PID: $RECORDER_PID" | tee -a "$MASTER_LOG"
    echo -e "   ${GREEN}✅ Auto recorder running (60s per anomaly)${NC}" | tee -a "$MASTER_LOG"
else
    echo -e "${YELLOW}[5/7] Skipping auto recorder (ffmpeg not found)${NC}" | tee -a "$MASTER_LOG"
    echo "   Install with: brew install ffmpeg" | tee -a "$MASTER_LOG"
fi
echo ""

# Step 6: Start Loom controller
echo -e "${CYAN}[6/7] Starting Loom controller...${NC}" | tee -a "$MASTER_LOG"
./scripts/loom_controller.sh 2>"$SESSION_DIR/loom.stderr.log" &
LOOM_PID=$!
echo "$LOOM_PID" > "$LOOM_PID_FILE"
echo "   PID: $LOOM_PID" | tee -a "$MASTER_LOG"
echo -e "   ${GREEN}✅ Loom controller running${NC}" | tee -a "$MASTER_LOG"
echo ""

# Step 7: Monitor session
echo -e "${CYAN}[7/7] Session active - monitoring for ${DURATION}s...${NC}" | tee -a "$MASTER_LOG"
echo -e "${YELLOW}   Press Ctrl+C to stop early${NC}"
echo ""

# Wait for duration or until interrupted
END_TIME=$(($(date +%s) + DURATION))

while [ "$(date +%s)" -lt "$END_TIME" ]; do
    # Check all processes are still running
    if ! kill -0 "$STREAMER_PID" 2>/dev/null; then
        echo -e "${RED}⚠️  Streamer died - restarting...${NC}" | tee -a "$MASTER_LOG"
        
        python3 -u scripts/stream_kraken_ws.py 2>>"$SESSION_DIR/streamer.stderr.log" | \
            tee -a "$RAW_LOG" > "$KRAKEN_PIPE" &
        
        STREAMER_PID=$!
        echo "$STREAMER_PID" > "$STREAMER_PID_FILE"
        echo "   New PID: $STREAMER_PID" | tee -a "$MASTER_LOG"
    fi
    
    sleep 30
done

echo -e "\n${GREEN}⏱️  Session duration reached - shutting down...${NC}" | tee -a "$MASTER_LOG"
