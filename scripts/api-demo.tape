Output screenshots/api-demo-demo.gif
Set Shell zsh
Set FontSize 28
Set Width 1200
Set Height 900
Set TypingSpeed 120ms

Type 'cd /Volumes/VIXinSSD/driftlock'
Enter
Type 'export API_KEY="dlk_24f0460f-bbff-4c18-8709-b7fcc8fcdf4c.gd5u763q5whkng7fuafhosgkdby43ws3"'
Enter
Type 'export BASE_URL="http://localhost:8080"'
Enter
Type 'export DATABASE_URL="postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable"'
Enter
Type 'curl -s $BASE_URL/healthz | jq'
Enter
Sleep 4s
Type 'jq ".[0:600]" test-data/financial-demo.json | curl -sS -X POST -H "X-Api-Key: $API_KEY" -H "Content-Type: application/json" --data @- $BASE_URL/v1/detect | jq ".anomalies[0]"'
Enter
Sleep 8s
Type 'psql "$DATABASE_URL" -c "SELECT id, stream_id, ncd, p_value::numeric, confidence::numeric FROM anomalies ORDER BY detected_at DESC LIMIT 5;"'
Enter
Sleep 6s
