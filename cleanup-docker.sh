#!/bin/bash
# Docker Cleanup Script for Driftlock
# Removes large Docker files to free up disk space

echo "🧹 Cleaning up Docker files..."
echo ""

# Stop Docker first (if it's running)
echo "1. Stopping Docker..."
osascript -e 'quit app "Docker"' 2>/dev/null
sleep 5

# Function to remove with sudo if needed
remove_with_sudo() {
    local path="$1"
    if [ -e "$path" ]; then
        echo "   Removing: $path"
        sudo rm -rf "$path" 2>/dev/null || rm -rf "$path" 2>/dev/null || echo "   ⚠️  Could not remove: $path"
    fi
}

# 1. Remove Docker.raw (the biggest file - 92GB)
echo "2. Removing Docker.raw (92GB)..."
remove_with_sudo "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw"
remove_with_sudo "$HOME/Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw.backup"

# 2. Remove Docker build cache
echo "3. Removing Docker build cache..."
remove_with_sudo "$HOME/.docker/buildx"
remove_with_sudo "$HOME/Library/Containers/com.docker.docker/Data/cache"

# 3. Remove Docker log files
echo "4. Removing Docker logs..."
remove_with_sudo "$HOME/Library/Containers/com.docker.docker/Data/log"

# 4. Remove Docker application support
echo "5. Removing Docker application support..."
remove_with_sudo "$HOME/Library/Application Support/Docker Desktop"
remove_with_sudo "$HOME/Library/Preferences/com.docker.docker.plist"

# 5. Remove system Docker files (if they exist)
echo "6. Removing system Docker files..."
remove_with_sudo "/var/lib/docker/images"
remove_with_sudo "/var/lib/docker/volumes"
remove_with_sudo "/var/lib/docker/overlay2"

# 6. Remove Docker temp files
echo "7. Removing Docker temp files..."
remove_with_sudo "/tmp/docker-*"
remove_with_sudo "/var/tmp/docker-*"

echo ""
echo "✅ Docker cleanup complete!"
echo ""
echo "Next steps:"
echo "1. Restart Docker Desktop"
echo "2. Run: docker system prune -a -f --volumes"
echo "3. Then: cd /Volumes/VIXinSSD/driftlock && docker compose up --build"
echo ""
echo "Disk space freed: Check with 'df -h /Users'"
