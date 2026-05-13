# PowerShell 脚本：本地交叉编译 RSK API (linux/amd64)
# 用法: .\scripts\build-local.ps1 [-Version "v1.2.3"]
# 默认版本: dev
param(
    [string]$Version = "dev"
)

Write-Host "=== 本地交叉编译 dujiao-api (linux/amd64) ===" -ForegroundColor Cyan

# 确保 bin 目录存在
New-Item -ItemType Directory -Force -Path "bin" | Out-Null

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -trimpath -tags release `
  -ldflags "-s -w -X github.com/dujiao-next/internal/version.Version=$Version" `
  -o bin\dujiao-api .\cmd\server

if ($LASTEXITCODE -eq 0) {
    Write-Host "=== 编译完成 ===" -ForegroundColor Green
    Get-ChildItem bin\dujiao-api | Select-Object Length, LastWriteTime
} else {
    Write-Host "=== 编译失败 ===" -ForegroundColor Red
    exit $LASTEXITCODE
}
