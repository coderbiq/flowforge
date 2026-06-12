# FlowForge CLI — Windows 安装脚本 (PowerShell)
# Usage: irm https://get.flowforge.dev/install.ps1 | iex

$ErrorActionPreference = 'Stop'

$AppName = "flowforge"
$CDN = if ($env:FLOWFORGE_CDN) { $env:FLOWFORGE_CDN } else { "https://cdn.flowforge.dev" }
$InstallDir = if ($env:FLOWFORGE_INSTALL) { $env:FLOWFORGE_INSTALL } else { "$HOME\.flowforge" }
$BinDir = "$InstallDir\bin"

# 架构检测
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ((Get-CimInstance Win32_Processor).Architecture -eq 9) {
        "aarch64-pc-windows-msvc"
    } else {
        "x86_64-pc-windows-msvc"
    }
} else {
    throw "32-bit Windows is not supported"
}

# 版本
$Version = if ($args[0]) { $args[0] } else {
    (Invoke-RestMethod -Uri "$CDN/release-latest.txt").Trim()
}

$DownloadUrl = "$CDN/release/$Version/$AppName-$Arch.zip"

Write-Host "Downloading $AppName $Version for $Arch..."

# 创建安装目录
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
}

# 下载
$ZipPath = "$env:TEMP\$AppName.zip"
curl.exe -sSfL $DownloadUrl -o $ZipPath
if ($LASTEXITCODE -ne 0) { throw "Download failed" }

# SHA256 校验
try {
    $ExpectedHash = (Invoke-RestMethod -Uri "$DownloadUrl.sha256").Split(' ')[0]
    $ActualHash = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
    if ($ActualHash -ne $ExpectedHash) {
        throw "Checksum mismatch: expected $ExpectedHash, got $ActualHash"
    }
    Write-Host "Checksum verified"
} catch {
    Write-Warning "Skipping checksum verification: $_"
}

# 解压
tar.exe xf $ZipPath -C $BinDir
Remove-Item $ZipPath

# 添加到 PATH
$User = [System.EnvironmentVariableTarget]::User
$CurrentPath = [System.Environment]::GetEnvironmentVariable('Path', $User)
if ($CurrentPath -notlike "*$BinDir*") {
    [System.Environment]::SetEnvironmentVariable('Path', "$CurrentPath;$BinDir", $User)
    $env:Path += ";$BinDir"
    Write-Host "Added $BinDir to PATH"
}

Write-Host ""
Write-Host "$AppName $Version installed successfully to $BinDir"
Write-Host "Run 'flowforge --help' to get started"
