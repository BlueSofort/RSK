#!/usr/bin/env bash
# 本地交叉编译 RSK API (linux/amd64)
# 用法: ./scripts/build-local.sh [版本号]
# 示例: ./scripts/build-local.sh v1.2.3
# 默认: dev
set -euo pipefail

APP_VERSION="${1:-dev}"

cd "$(dirname "$0")/.."

echo "=== 本地交叉编译 dujiao-api (linux/amd64) ==="
mkdir -p bin

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -tags release \
  -ldflags="-s -w -X github.com/dujiao-next/internal/version.Version=${APP_VERSION}" \
  -o bin/dujiao-api ./cmd/server

echo "=== 编译完成 ==="
ls -lh bin/dujiao-api
file bin/dujiao-api
