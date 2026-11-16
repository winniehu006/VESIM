# Builds a self-contained VESIM executable and portable zip for Windows.
# Usage: .\scripts\build-windows.ps1 [-Version v2.1.35] [-SkipFrontend]
param(
  [string]$Version = "v2.1.35",
  [string]$GoBin = "",
  [switch]$SkipFrontend
)

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot
Set-Location $root

if ($GoBin) { $env:PATH = "$GoBin;$env:PATH" }
$go = Get-Command go -ErrorAction SilentlyContinue
if (-not $go) {
  throw "go was not found in PATH. Install Go 1.26+ or add its bin directory to PATH."
}
$npm = Get-Command npm -ErrorAction SilentlyContinue
if (-not $npm) {
  throw "npm was not found in PATH. Install Node.js 20+ first."
}

if (-not $SkipFrontend) {
  npm ci --prefix web
  if ($LASTEXITCODE -ne 0) { throw "npm ci failed" }
  npm run build --prefix web
  if ($LASTEXITCODE -ne 0) { throw "frontend build failed" }

  if (Test-Path internal/web/dist) { Remove-Item -Recurse -Force internal/web/dist }
  Copy-Item -Recurse web/dist internal/web/dist
}

$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

$buildTime = (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
$ldflags = "-s -w -X 'github.com/winniehu006/vesim/internal/global.Version=$Version' -X 'github.com/winniehu006/vesim/internal/global.BuildTime=$buildTime'"

New-Item -ItemType Directory -Force -Path dist | Out-Null
go build -trimpath -buildvcs=false -tags "with_utls nomsgpack" -ldflags $ldflags -o dist/vesim.exe ./cmd/vesim
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

$stage = "dist/vesim_windows_amd64"
if (Test-Path $stage) { Remove-Item -Recurse -Force $stage }
New-Item -ItemType Directory -Force -Path "$stage/config" | Out-Null

Copy-Item dist/vesim.exe $stage
Copy-Item config/config.example.yaml "$stage/config/config.example.yaml"
Copy-Item packaging/windows/start-vesim.bat $stage
Copy-Item LICENSE, THIRD_PARTY_NOTICES.md $stage
New-Item -ItemType Directory -Force -Path "$stage/data" | Out-Null

$zip = "dist/vesim_${Version}_windows_amd64.zip"
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path "$stage/*" -DestinationPath $zip

Write-Host "Built $zip"
