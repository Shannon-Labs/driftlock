#!/bin/bash
# Fix build issues for local development

set -e

echo "🔧 Fixing build issues..."
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CBAD_LIB_DIR="$ROOT_DIR/cbad-core/target/release"

echo "Root directory: $ROOT_DIR"
echo "Library directory: $CBAD_LIB_DIR"
echo ""

# Check if Rust library exists
if [ ! -f "$CBAD_LIB_DIR/libcbad_core.dylib" ] && [ ! -f "$CBAD_LIB_DIR/libcbad_core.so" ]; then
    echo -e "${YELLOW}⚠ Rust library not found. Building...${NC}"
    cd "$ROOT_DIR/cbad-core"
    cargo build --release
    cd "$ROOT_DIR"
fi

# For macOS: Copy to standard location (requires user to run with appropriate permissions)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "macOS detected. Setting up library paths..."
    
    # Create local lib directory
    mkdir -p "$ROOT_DIR/lib"
    
    # Copy libraries to local lib directory
    if [ -f "$CBAD_LIB_DIR/libcbad_core.dylib" ]; then
        cp "$CBAD_LIB_DIR/libcbad_core.dylib" "$ROOT_DIR/lib/"
        echo -e "${GREEN}✓ Copied libcbad_core.dylib to lib/${NC}"
    fi
    
    if [ -f "$CBAD_LIB_DIR/libcbad_core.a" ]; then
        cp "$CBAD_LIB_DIR/libcbad_core.a" "$ROOT_DIR/lib/"
        echo -e "${GREEN}✓ Copied libcbad_core.a to lib/${NC}"
    fi
    
    # Export environment variables
    export CGO_ENABLED=1
    export CGO_LDFLAGS="-L$ROOT_DIR/lib -lcbad_core"
    export DYLD_LIBRARY_PATH="$ROOT_DIR/lib:$DYLD_LIBRARY_PATH"
    export LD_LIBRARY_PATH="$ROOT_DIR/lib:$LD_LIBRARY_PATH"
    
    echo ""
    echo "Environment variables set:"
    echo "  CGO_ENABLED=$CGO_ENABLED"
    echo "  CGO_LDFLAGS=$CGO_LDFLAGS"
    echo "  DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH"
    echo ""
    echo "To build, run:"
    echo "  cd collector-processor"
    echo "  go build -o ../bin/driftlock-http-test ./cmd/driftlock-http"
    
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Linux detected. Setting up library paths..."
    
    # Create local lib directory
    mkdir -p "$ROOT_DIR/lib"
    
    # Copy libraries to local lib directory
    if [ -f "$CBAD_LIB_DIR/libcbad_core.so" ]; then
        cp "$CBAD_LIB_DIR/libcbad_core.so" "$ROOT_DIR/lib/"
        echo -e "${GREEN}✓ Copied libcbad_core.so to lib/${NC}"
    fi
    
    if [ -f "$CBAD_LIB_DIR/libcbad_core.a" ]; then
        cp "$CBAD_LIB_DIR/libcbad_core.a" "$ROOT_DIR/lib/"
        echo -e "${GREEN}✓ Copied libcbad_core.a to lib/${NC}"
    fi
    
    # Export environment variables
    export CGO_ENABLED=1
    export CGO_LDFLAGS="-L$ROOT_DIR/lib -lcbad_core"
    export LD_LIBRARY_PATH="$ROOT_DIR/lib:$LD_LIBRARY_PATH"
    
    echo ""
    echo "Environment variables set:"
    echo "  CGO_ENABLED=$CGO_ENABLED"
    echo "  CGO_LDFLAGS=$CGO_LDFLAGS"
    echo "  LD_LIBRARY_PATH=$LD_LIBRARY_PATH"
fi

echo ""
echo -e "${GREEN}✓ Build environment configured${NC}"
echo ""
echo "Note: For Docker builds, you may need to restart Docker Desktop if experiencing I/O errors."

