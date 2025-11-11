#!/bin/bash

# 编译脚本 - 编译所有服务的二进制文件

set -e

# 检测并设置 Go 路径
if ! command -v go &> /dev/null; then
    # 尝试常见的 Go 安装路径
    if [ -f "/usr/local/go/bin/go" ]; then
        export PATH=$PATH:/usr/local/go/bin
    elif [ -f "/opt/homebrew/bin/go" ]; then
        export PATH=$PATH:/opt/homebrew/bin
    elif [ -f "$HOME/go/bin/go" ]; then
        export PATH=$PATH:$HOME/go/bin
    else
        echo "Error: Go not found. Please install Go or add it to your PATH."
        exit 1
    fi
fi

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 创建 bin 目录
BIN_DIR="$PROJECT_ROOT/bin"
mkdir -p "$BIN_DIR"

echo "Building binaries..."
echo "Output directory: $BIN_DIR"
echo ""

# 编译 STT 服务
echo "Building STT server..."
cd "$PROJECT_ROOT"
go build -o "$BIN_DIR/stt-server" ./cmd/stt-server
if [ $? -eq 0 ]; then
    echo "✓ STT server built successfully"
    ls -lh "$BIN_DIR/stt-server"
else
    echo "✗ Failed to build STT server"
    exit 1
fi
echo ""

# 编译 TTS 服务
echo "Building TTS server..."
go build -o "$BIN_DIR/tts-server" ./cmd/tts-server
if [ $? -eq 0 ]; then
    echo "✓ TTS server built successfully"
    ls -lh "$BIN_DIR/tts-server"
else
    echo "✗ Failed to build TTS server"
    exit 1
fi
echo ""

# 编译统一服务
echo "Building unified speech server..."
go build -o "$BIN_DIR/speech-server" ./cmd/speech-server
if [ $? -eq 0 ]; then
    echo "✓ Speech server built successfully"
    ls -lh "$BIN_DIR/speech-server"
else
    echo "✗ Failed to build speech server"
    exit 1
fi
echo ""

echo "=========================================="
echo "All binaries built successfully!"
echo "=========================================="
echo ""
echo "Binaries location: $BIN_DIR"
echo ""
echo "Usage:"
echo "  STT server:     $BIN_DIR/stt-server"
echo "  TTS server:     $BIN_DIR/tts-server"
echo "  Unified server: $BIN_DIR/speech-server"
echo ""

