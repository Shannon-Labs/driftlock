#!/bin/bash
set -e

REPO_ROOT="/Volumes/VIXinSSD/driftlock"
ACTIVE_LOG=$(ls -t "$REPO_ROOT/logs/crypto-api-test-*.log" 2>/dev/null | head -1)

if [ -z "$ACTIVE_LOG" ]; then
    echo "No active log file found"
    exit 1
fi

echo "🚀 Anomaly alert monitor started - watching: $ACTIVE_LOG"
echo "Will alert immediately when anomaly detected"
echo ""

# Track last reported anomalies
LAST_ANOMALY_COUNT=0
LAST_BATCH_COUNT=0
START_TIME=$(date +%s)

while true; do
    # Get current counts
    CURRENT_BATCHES=$(grep -c "Batch sent" "$ACTIVE_LOG" 2>/dev/null || echo "0")
    CURRENT_ANOMALIES=$(grep -c "ANOMALY DETECTED\|anomalies detected" "$ACTIVE_LOG" 2>/dev/null || echo "0")
    
    # Check for new anomalies
    if [ "$CURRENT_ANOMALIES" -gt "$LAST_ANOMALY_COUNT" ]; then
        echo ""
        echo "═══════════════════════════════════════════════════════════"
        echo "🚨 START LOOM NOW! ANOMALY DETECTED! 🚨"
        echo "═══════════════════════════════════════════════════════════"
        echo "Time: $(date '+%Y-%m-%d %H:%M:%S')"
        echo ""
        echo "📊 Anomaly Summary:"
        tail -n 20 "$ACTIVE_LOG" | grep -E "(ANOMALY DETECTED|anomalies detected|Anomaly:|🎯 Anomaly)" | tail -5
        echo ""
        echo "📈 Recent Context:"
        tail -n 40 "$ACTIVE_LOG" | grep -A 10 -B 10 -E "(ANOMALY DETECTED|anomalies detected)" | tail -20 || true
        echo ""
        echo "📹 START YOUR LOOM RECORDING NOW!"
        echo "═══════════════════════════════════════════════════════════"
        echo ""
        say "Anomaly detected! Start Loom now!" 2>/dev/null || true
        afplay /System/Library/Sounds/Glass.aiff 2>/dev/null || true
        
        LAST_ANOMALY_COUNT=$CURRENT_ANOMALIES
    fi
    
    # Check for 401 errors
    AUTH_ERRORS=$(grep -c "401\|unauthenticated" "$ACTIVE_LOG" 2>/dev/null || echo "0")
    if [ "$AUTH_ERRORS" -gt 0 ]; then
        echo "⚠️  AUTH ERROR DETECTED - Please update DRIFTLOCK_API_KEY in .env"
        echo "401 errors found in log, test may be failing"
        echo ""
    fi
    
    # Check if test is still running
    TEST_PID=$(pgrep -f "api_crypto_test" 2>/dev/null || echo "")
    MONITOR_PID=$(pgrep -f "monitor_anomalies" 2>/dev/null || echo "")
    
    if [ -z "$TEST_PID" ]; then
        echo "⚠️  Test process not running - restarting..."
        sleep 5
        ACTIVE_LOG=$(ls -t "$REPO_ROOT/logs/crypto-api-test-*.log" 2>/dev/null | head -1)
        if [ -n "$ACTIVE_LOG" ]; then
            export CRYPTO_SOURCE=binance
            set -a && source "$REPO_ROOT/.env" 2>/dev/null && set +a
            export DRIFTLOCK_API_URL=${DRIFTLOCK_API_URL:-https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1}
            cd "$REPO_ROOT"
            nohup python3 scripts/api_crypto_test_sensitive.py --api-key "$DRIFTLOCK_API_KEY" --api-url "$DRIFTLOCK_API_URL" >> "$ACTIVE_LOG" 2>&1 &
            echo "✅ Test restarted with PID: $!"
        fi
    fi
    
    if [ -z "$MONITOR_PID" ]; then
        echo "⚠️  Monitor process not running - restarting..."
        sleep 5
        cd "$REPO_ROOT"
        nohup bash scripts/monitor_anomalies.sh >> /dev/null 2>&1 &
        echo "✅ Monitor restarted with PID: $!"
    fi
    
    # Periodic status every 15 minutes (900 seconds)
    ELAPSED=$(( $(date +%s) - START_TIME ))
    if [ $((ELAPSED % 900)) -lt 30 ] && [ "$ELAPSED" -gt 60 ]; then
        echo ""
        echo "⏰ STATUS UPDATE - $(date '+%H:%M:%S')"
        echo "   Runtime: $((ELAPSED / 60)) minutes"
        echo "   Batches: $CURRENT_BATCHES | Anomalies: $CURRENT_ANOMALIES"
        echo "   Log: $ACTIVE_LOG"
        echo "   Test PID: $(pgrep -f api_crypto_test 2>/dev/null | head -1 || echo 'N/A')"
        echo ""
    fi
    
    LAST_BATCH_COUNT=$CURRENT_BATCHES
    sleep 5
done