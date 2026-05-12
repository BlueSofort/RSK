#!/usr/bin/env bash
# VPS 端部署脚本
# 用法: ssh user@vps 'cd /opt/rsk && ./deploy.sh'
# 或:   ssh user@vps 'cd /opt/rsk && ./deploy.sh api'

set -euo pipefail

cd /opt/rsk
TARGET="${1:-all}"

pull() {
  local dir="$1"
  echo "=== Pulling $dir ==="
  (cd "/opt/rsk/$dir" && git pull)
}

build_api() {
  echo "=== Building api ==="
  cd /opt/rsk
  docker build --pull -t rsk-api:latest -f RSK-main/Dockerfile RSK-main/
}

build_user() {
  echo "=== Building user ==="
  cd /opt/rsk
  docker build -t rsk-user:latest -f RSK-user/Dockerfile.prod RSK-user/
}

build_admin() {
  echo "=== Building admin ==="
  cd /opt/rsk
  docker build -t rsk-admin:latest -f RSK-admin/Dockerfile.prod RSK-admin/
}

case "$TARGET" in
  all)
    pull RSK-main
    pull RSK-user
    pull RSK-admin
    build_api
    build_user
    build_admin
    ;;
  api)
    pull RSK-main
    build_api
    ;;
  user)
    pull RSK-user
    build_user
    ;;
  admin)
    pull RSK-admin
    build_admin
    ;;
  *)
    echo "用法: ./deploy.sh [all|api|user|admin]"
    exit 1
    ;;
esac

echo "=== Restarting ==="
docker-compose up -d --remove-orphans

echo "=== Cleaning old images ==="
docker image prune -f

echo "Done."
