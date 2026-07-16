# FlowForge CLI — Windows 安装脚本 (PowerShell)
# Usage: irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
# Options: -Version <ver> -Prefix <dir>

param(
    [string]$Version = "latest",
    [string]$Prefix = "$HOME\.flowforge"
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$AppName = "flowforge"
$ReleasesBase = "https://github.com/coderbiq/flowforge/releases"
$BinDir = "$Prefix\bin"

# ── 架构检测 ─────────────────────────────────────────

$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ((Get-CimInstance Win32_Processor).Architecture -eq 9) {
        "windows-arm64"
    } else {
        "windows-amd64"
    }
} else {
    throw "32-bit Windows is not supported"
}

Write-Host "Detected: $Arch"

# ── 版本获取 ─────────────────────────────────────────

if ($Version -eq "latest") {
    try {
        $manifest = Invoke-RestMethod -Uri "$ReleasesBase/latest/download/manifest.json" -ErrorAction Stop
        $Version = $manifest.version
    } catch {
        throw "Failed to fetch latest version: $_"
    }
}

# ── manifest 获取 ────────────────────────────────────

function Get-ArtifactInfo {
    param($ManifestUrl, $Platform)

    try {
        $manifest = Invoke-RestMethod -Uri $ManifestUrl -ErrorAction Stop
    } catch {
        return $null
    }

    $artifact = $manifest.artifacts | Where-Object { $_.platform -eq $Platform } | Select-Object -First 1
    if (-not $artifact) {
        return $null
    }

    return @{ url = $artifact.url; sha256 = $artifact.sha256 }
}

$manifestUrl = "$ReleasesBase/download/$Version/manifest.json"
$artifactInfo = Get-ArtifactInfo -ManifestUrl $manifestUrl -Platform $Arch

if (-not $artifactInfo) {
    throw "Failed to find artifact for $Arch version $Version"
}

$DownloadUrl = $artifactInfo.url
$ExpectedSHA256 = $artifactInfo.sha256

# ── 下载 ─────────────────────────────────────────────

Write-Host "Downloading $AppName $Version ..."

$ZipPath = "$env:TEMP\$AppName.zip"
curl.exe -sSfL $DownloadUrl -o $ZipPath
if ($LASTEXITCODE -ne 0) {
    Remove-Item $ZipPath -ErrorAction SilentlyContinue
    throw "Download failed"
}

# ── SHA256 校验 ──────────────────────────────────────

Write-Host "Verifying checksum..."

$ActualHash = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
if ($ActualHash -ne $ExpectedSHA256) {
    Remove-Item $ZipPath -Force -ErrorAction SilentlyContinue
    throw "Checksum mismatch: expected $ExpectedSHA256, got $ActualHash"
}
Write-Host "Checksum verified"

# ── 安装 ─────────────────────────────────────────────

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
}

$ExtractDir = Join-Path $env:TEMP "$AppName-install"
if (Test-Path $ExtractDir) {
    Remove-Item $ExtractDir -Recurse -Force
}
New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null

tar.exe xf $ZipPath -C $ExtractDir

Move-Item -Path (Join-Path $ExtractDir "$AppName.exe") -Destination (Join-Path $BinDir "$AppName.exe") -Force

if (Test-Path (Join-Path $ExtractDir "assets")) {
    $AssetsDir = Join-Path $Prefix "assets"
    if (Test-Path $AssetsDir) {
        Remove-Item $AssetsDir -Recurse -Force
    }
    Move-Item -Path (Join-Path $ExtractDir "assets") -Destination $AssetsDir -Force
}

Remove-Item $ExtractDir -Recurse -Force
Remove-Item $ZipPath -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "$AppName $Version installed to $BinDir\$AppName.exe"

# ── 验证 ─────────────────────────────────────────────

$exe = Join-Path $BinDir "$AppName.exe"
if (Test-Path $exe) {
    try {
        & $exe --version 2>&1 | Out-Null
        Write-Host "Verification: OK"
    } catch {
        Write-Warning "Please add $BinDir to your PATH: setx PATH `"%PATH%;$BinDir`""
    }
}

Write-Host ""
Write-Host "Run 'flowforge init' to get started in a project"
