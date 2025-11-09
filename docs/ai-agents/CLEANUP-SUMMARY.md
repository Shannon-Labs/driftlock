# Docker Cleanup Summary

## Large Files to Remove

### 1. Docker.raw (92GB) ⭐ BIGGEST
**Location:** `~/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw`
**What it is:** Main Docker disk image
**Safe to delete:** Yes, Docker will recreate it on next start

### 2. Docker.raw.backup (92GB) ⭐ BIGGEST
**Location:** `~/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw.backup`
**What it is:** Backup of Docker disk image
**Safe to delete:** Yes, this is just a backup

### 3. Build Cache (5-10GB)
**Location:** `~/.docker/buildx/`
**What it is:** Docker build cache
**Safe to delete:** Yes, will be rebuilt as needed

### 4. Docker Logs (1-5GB)
**Location:** `~/Library/Containers/com.docker.docker/Data/log/`
**What it is:** Docker daemon and container logs
**Safe to delete:** Yes, just logs

### 5. Application Support (500MB-2GB)
**Location:** `~/Library/Application Support/Docker Desktop/`
**What it is:** Docker Desktop application data
**Safe to delete:** Yes, will be recreated

## Total Space: ~200GB+

## Quick Cleanup Commands

```bash
# Run the quick cleanup script
cd /Volumes/VIXinSSD/driftlock
./QUICK-CLEANUP.sh

# Or run manually:
# Remove Docker.raw (92GB)
rm ~/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw

# Remove build cache
rm -rf ~/.docker/buildx
rm -rf ~/Library/Containers/com.docker.docker/Data/cache

# Remove logs
rm -rf ~/Library/Containers/com.docker.docker/Data/log
```

## After Cleanup

1. **Restart Docker Desktop**
2. **Prune everything:** `docker system prune -a -f --volumes`
3. **Test:** `docker run hello-world`
4. **Run Driftlock:** `docker compose up --build`

## Expected Results

- **Disk space freed:** ~200GB
- **Docker performance:** Much faster
- **Build times:** Reduced (no cache)
- **Driftlock:** Should run without disk space errors
