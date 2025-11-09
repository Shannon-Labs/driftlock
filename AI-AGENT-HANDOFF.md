# AI Agent Handoff - Driftlock Docker Fix

**Previous Agent:** Attempted Docker build fixes but ran into fundamental linking issues between Rust static libraries (`libcbad_core.a`) and Go CGO across different libc implementations (musl vs glibc).

**Your Mission:** Fix the Docker build so Driftlock can be deployed via `docker-compose up`

## Current State

### ✅ What Works
- **Core Algorithm:** 58/58 Rust unit tests pass
- **API Server:** Builds and runs outside Docker
- **Test Data:** 1,600 synthetic transactions ready
- **Configuration:** `.env` and DB auth fixed
- **Repository:** Clean, honest, well-documented

### ❌ What's Broken
- **Docker Build:** Fails at Go linking stage
- **Error:** `libcbad_core.a: error adding symbols: archive has no index; run ranlib to add one`
- **Root Cause:** Cross-language linking (Rust static → Go dynamic) across libc boundaries

## Files Modified by Previous Agent

1. **`Dockerfile`** - Multiple attempts (Ubuntu base, Alpine base, static/dynamic linking)
2. **`.env`** - DB credentials fixed (user=postgres, password=postgres)
3. **`docker-compose.yml`** - DB service uses correct env vars
4. **`TEST-RESULTS.md`** - Documents algorithm verification
5. **`DOCKER-BUILD-STATUS.md`** - Documents current issues

## The Core Problem

```
Rust (static library) → CGO (Go) → musl libc (Alpine) = LINKING HELL
```

**Specific Error:**
```
/usr/lib/gcc/aarch64-alpine-linux-musl/14.2.0/../../../../aarch64-alpine-linux-musl/bin/ld: 
/build/cbad-core/target/release/libcbad_core.a: error adding symbols: archive has no index; 
run ranlib to add one
```

## What We've Tried

### Attempt 1: Alpine + Static Linking
- **Result:** ❌ Failed (ranlib issue with static .a file)
- **Dockerfile:** Uses `rust:1.80-slim` → `golang:1.24-alpine` → `alpine:3.19`

### Attempt 2: Ubuntu Base
- **Result:** ❌ Failed (same linking issue, different error format)
- **Dockerfile:** Uses `rust:1.80` → `golang:1.24` → `ubuntu:22.04`

### Attempt 3: Dynamic Linking Only
- **Result:** ❌ Failed (Go can't find Rust symbols)
- **Change:** Modified Cargo.toml to only build `cdylib`

## What Still Needs to Be Fixed

### Option 1: Fix Static Library Index (Easiest)
```dockerfile
# In Rust builder stage
RUN cargo build --release && ranlib target/release/libcbad_core.a
```
**Then:** Ensure Go links against static library correctly

### Option 2: Use Dynamic Library (.so)
```dockerfile
# Build only cdylib
RUN sed -i 's/crate-type = \["cdylib", "rlib", "staticlib"\]/crate-type = ["cdylib"]/' Cargo.toml
RUN cargo build --release

# Copy .so to Go builder
COPY --from=rust-builder /build/cbad-core/target/release/*.so /usr/local/lib/
RUN ldconfig
```

### Option 3: Pre-built Binaries (Quickest for Demo)
```bash
# Build outside Docker (works!)
cd /Volumes/VIXinSSD/driftlock/cbad-core && cargo build --release
cd /Volumes/VIXinSSD/driftlock/api-server && go build ./cmd/api-server

# Create minimal Dockerfile that just copies binaries
FROM alpine:3.19
COPY api-server/driftlock-api /app/
COPY cbad-core/target/release/* /usr/local/lib/
ENTRYPOINT ["/app/driftlock-api"]
```

### Option 4: Simplify Architecture
- **Idea:** Remove Rust dependency entirely
- **Approach:** Rewrite CBAD in pure Go
- **Trade-off:** Loses Rust performance, gains simplicity

## Immediate Next Steps

1. **Try ranlib fix:**
   ```bash
   cd /Volumes/VIXinSSD/driftlock
   sed -i '/cargo build --release/ s/$/ \&\& ranlib target\/release\/libcbad_core.a/' Dockerfile
   docker build -t driftlock:test .
   ```

2. **If that fails, try pre-built approach:**
   ```bash
   # Build outside Docker (works!)
   cd cbad-core && cargo build --release
   cd ../api-server && go build ./cmd/api-server
   
   # Create Dockerfile.prebuilt
   cat > Dockerfile.prebuilt << 'EOF'
   FROM alpine:3.19
   RUN apk add --no-cache ca-certificates curl libgcc
   COPY api-server/driftlock-api /app/
   COPY cbad-core/target/release/*.so* /usr/local/lib/
   RUN ldconfig /usr/local/lib
   ENTRYPOINT ["/app/driftlock-api"]
   EOF
   
   docker build -f Dockerfile.prebuilt -t driftlock:test .
   ```

3. **Test the build:**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=host.docker.internal \
     -e DB_PORT=5432 \
     -e DB_USER=postgres \
     -e DB_PASSWORD=postgres \
     -e DB_DATABASE=driftlock \
     -e DEFAULT_API_KEY=test-key \
     driftlock:test
   ```

## Success Criteria

✅ Docker builds without errors  
✅ Container starts successfully  
✅ API responds on localhost:8080/healthz  
✅ Can ingest test data  
✅ Anomalies are detected  

## Files to Check

- `/Volumes/VIXinSSD/driftlock/Dockerfile` - Build instructions
- `/Volumes/VIXinSSD/driftlock/cbad-core/Cargo.toml` - Rust crate config
- `/Volumes/VIXinSSD/driftlock/docker-compose.yml` - Service orchestration
- `/Volumes/VIXinSSD/driftlock/.env` - Environment variables

## Previous Agent Notes

> "The core technology is solid. Docker can be fixed after YC application. The algorithm works - that's what matters for YC."

**Focus on:** Getting Docker working for the demo, not perfecting the architecture.

Good luck! 🚀