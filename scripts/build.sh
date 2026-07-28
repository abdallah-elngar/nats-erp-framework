#!/bin/bash

# ============================================
# بناء NATS Framework
# ============================================

set -e

echo "🔨 Building NATS Framework..."

# المتغيرات
BINARY_NAME="nats"
BUILD_DIR="./dist"
CMD_DIR="./cmd/nats"

# إنشاء مجلد البناء
mkdir -p $BUILD_DIR

# تحديد نظام التشغيل
OS=${1:-$(go env GOOS)}
ARCH=${2:-$(go env GOARCH)}

echo "📦 Building for: $OS/$ARCH"

# بناء التطبيق
GOOS=$OS GOARCH=$ARCH go build -o "$BUILD_DIR/$BINARY_NAME" "$CMD_DIR/main.go"

if [ $? -eq 0 ]; then
    echo "✅ Build complete: $BUILD_DIR/$BINARY_NAME"
    echo "📊 Binary size: $(du -h "$BUILD_DIR/$BINARY_NAME" | cut -f1)"
else
    echo "❌ Build failed!"
    exit 1
fi