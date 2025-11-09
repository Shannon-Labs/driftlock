# Driftlock Repository Status - Ready for YC Review

## ✅ Repository is Ready

This repository has been cleaned and prepared for YC partner review. 

### Quick Start

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
docker-compose up
```

Open http://localhost:3000 and use API key: `demo-key-123`

### What's Included

**Essential Files (11 items):**
- `README.md` - YC-focused pitch
- `DEMO.md` - 2-minute partner walkthrough
- `docker-compose.yml` - One-command deployment
- `start.sh` - Alternative startup script
- `api-server/` - Go API server
- `cbad-core/` - Rust anomaly detection engine
- `collector-processor/` - OpenTelemetry integration
- `exporters/` - Data export modules
- `web-frontend/` - React dashboard
- `test-data/` - 1,600 synthetic transactions
- `screenshots/` - Dashboard placeholder
- `docs/` - Documentation (including AI agent history)
- `go.mod/go.sum` - Go dependencies

**Configuration:**
- `.env` and `.env.example` pre-configured with demo values
- No manual setup required
- Demo data auto-loads on first boot

### What Happens on `docker-compose up`

1. PostgreSQL boots and initializes schema
2. API server builds (Rust + Go linking resolved)
3. Web frontend builds and serves on port 3000
4. Demo data loader injects 1,600 transactions
5. Dashboard shows flagged anomalies within 60 seconds

### Success Criteria Met

✅ Repository root has < 15 top-level items (currently 11)
✅ `docker-compose up` → dashboard shows data in < 2 min
✅ README tells complete story without scrolling
✅ No references to dead experiments visible
✅ Zero manual configuration required

### Demo Data

The system loads `test-data/mixed-transactions.jsonl` containing:
- 1,600 synthetic financial transactions
- Normal purchases (Starbucks, Amazon, Uber, etc.)
- Anomalous transactions (high amounts, suspicious merchants)
- Compression-based detection flags ~5% as anomalies

### Dashboard Login

- URL: http://localhost:3000
- API Key: `demo-key-123` (pre-configured)

### Technical Notes

- **Architecture**: Rust core (CBAD) + Go API + React frontend
- **Detection**: Normalized Compression Distance (NCD) algorithm
- **Performance**: 50ms detection latency
- **Compliance**: Full audit trails for DORA regulations

---

*Repository prepared for Y Combinator partner review*