#!/bin/bash
set -e

# Sync docs/ to landing-page/public/docs/
# This ensures the frontend always serves the latest documentation.

echo "Syncing documentation to frontend..."

# Ensure destination exists
mkdir -p landing-page/public/docs

# Rsync with delete to mirror structure (excluding hidden files/git)
# We use cp -R for simplicity if rsync isn't desired, but rsync is better.
# Since this is dev env, cp is safer/easier to read.

# Clear old docs in public to remove stale files
rm -rf landing-page/public/docs/*

# Copy new structure
cp -R docs/* landing-page/public/docs/

echo "Documentation synced to landing-page/public/docs/"
