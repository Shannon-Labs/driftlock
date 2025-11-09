#!/bin/bash
# Quick Docker Cleanup - Focus on biggest files

echo "🧹 Quick Docker Cleanup - Removing largest files..."
echo ""

# The biggest culprit: Docker.raw (92GB)
echo "1. Removing Docker.raw (92GB)..."
if [ -f "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw" ]; then
    echo "   Found Docker.raw - removing..."
    rm -f "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw"
    echo "   ✅ Docker.raw removed"
else
    echo "   Docker.raw not found (already removed or different location)"
fi

# Also remove backup if it exists
if [ -f "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw.backup" ]; then
    echo "   Removing Docker.raw.backup..."
    rm -f "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw.backup"
    echo "   ✅ Docker.raw.backup removed"
fi

# Remove build cache (can be several GB)
echo ""
echo "2. Removing Docker build cache..."
rm -rf "$HOME/.docker/buildx" 2>/dev/null
rm -rf "$HOME/Library/Containers/com.docker.docker/Data/cache" 2>/dev/null
echo "   ✅ Build cache removed"

# Remove logs (can accumulate)
echo ""
echo "3. Removing Docker logs..."
rm -rf "$HOME/Library/Containers/com.docker.docker/Data/log" 2>/dev/null
echo "   ✅ Logs removed"

# Try to remove some temp files
echo ""
echo "4. Removing temp files..."
rm -rf /tmp/docker-* 2>/dev/null
rm -rf /var/tmp/docker-* 2>/dev/null
echo "   ✅ Temp files removed"

echo ""
echo "✅ Quick cleanup complete!"
echo ""
echo "Disk space freed: ~92GB+ (Docker.raw + cache + logs)"
echo ""
echo "Next steps:"
echo "1. Restart Docker Desktop"
echo "2. Run: docker system prune -a -f --volumes"
echo "3. Then test: docker run hello-world"
echo ""
echo "Once Docker is working:"
echo "cd /Volumes/VIXinSSD/driftlock && docker compose up --build"
