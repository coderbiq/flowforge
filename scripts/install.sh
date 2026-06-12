#!/bin/sh
# FlowForge CLI — 一键安装脚本 (macOS / Linux)
# Usage: curl -fsSL https://get.flowforge.dev | sh
# 或指定版本: curl -fsSL https://get.flowforge.dev | sh -s v0.1.0
set -eu

APP_NAME="flowforge"
APP_VERSION="${1:-latest}"
CDN_BASE="${FLOWFORGE_CDN:-https://cdn.flowforge.dev}"

# ── 工具函数 ─────────────────────────────────────────

info()  { printf "\033[0;32m%s\033[0m\n" "$1"; }
warn()  { printf "\033[0;33m%s\033[0m\n" "$1"; }
error() { printf "\033[0;31merror: %s\033[0m\n" "$1"; exit 1; }

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        error "need '$1' (command not found)"
    fi
}

need_cmd uname
need_cmd curl
need_cmd mktemp
need_cmd chmod
need_cmd mkdir

# ── 架构检测 ─────────────────────────────────────────

get_architecture() {
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"

    case "$_ostype" in
        Linux)
            if ldd /bin/sh 2>/dev/null | grep -qi musl; then
                _ostype="unknown-linux-musl"
            else
                _ostype="unknown-linux-gnu"
            fi
            ;;
        Darwin)  _ostype="apple-darwin" ;;
        MINGW*|MSYS*|CYGWIN*) _ostype="pc-windows-msvc" ;;
        *) error "unsupported OS: $_ostype" ;;
    esac

    case "$_cputype" in
        x86_64|amd64)   _cputype="x86_64" ;;
        aarch64|arm64)  _cputype="aarch64" ;;
        *) error "unsupported CPU: $_cputype" ;;
    esac

    # Rosetta 2 检测
    if [ "$_ostype" = "apple-darwin" ] && [ "$_cputype" = "x86_64" ]; then
        if [ "$(sysctl -n sysctl.proc_translated 2>/dev/null)" = "1" ]; then
            _cputype="aarch64"
            info "Detected Rosetta 2, downloading arm64 build"
        fi
    fi

    RETVAL="${_cputype}-${_ostype}"
}

# ── 主流程 ───────────────────────────────────────────

main() {
    get_architecture
    local arch="$RETVAL"

    info "Detected: $arch"

    # 版本解析
    local version
    if [ "$APP_VERSION" = "latest" ]; then
        version=$(curl -sfL "${CDN_BASE}/release-latest.txt") || {
            error "Failed to fetch latest version"
        }
    else
        version="$APP_VERSION"
    fi

    # 下载 URL
    local url="${CDN_BASE}/release/${version}/${APP_NAME}-${arch}.tar.gz"
    local tmpdir
    tmpdir="$(mktemp -d)" || error "Failed to create temp directory"
    local archive="${tmpdir}/${APP_NAME}.tar.gz"

    # 下载
    info "Downloading ${APP_NAME} ${version}..."
    curl -sSfL "$url" -o "$archive" || error "Download failed"

    # SHA256 校验
    info "Verifying checksum..."
    local expected_checksum
    expected_checksum="$(curl -sfL "${url}.sha256" 2>/dev/null | awk '{print $1}')" || true
    if [ -n "$expected_checksum" ]; then
        local actual_checksum
        if command -v sha256sum >/dev/null 2>&1; then
            actual_checksum="$(sha256sum "$archive" | awk '{print $1}')"
        elif command -v shasum >/dev/null 2>&1; then
            actual_checksum="$(shasum -a 256 "$archive" | awk '{print $1}')"
        fi
        if [ "$actual_checksum" != "$expected_checksum" ]; then
            error "Checksum mismatch"
        fi
        info "Checksum verified"
    else
        warn "Skipping checksum verification"
    fi

    # 安装
    local install_dir="${FLOWFORGE_INSTALL:-$HOME/.flowforge}"
    local bin_dir="$install_dir/bin"
    mkdir -p "$bin_dir"

    tar xzf "$archive" -C "$tmpdir"
    mv "$tmpdir/${APP_NAME}" "$bin_dir/"
    chmod +x "$bin_dir/${APP_NAME}"

    rm -rf "$tmpdir"

    info "${APP_NAME} ${version} installed to $bin_dir/${APP_NAME}"

    # PATH 配置
    case "$(basename "${SHELL:-unknown}")" in
        zsh)
            echo "export PATH=\"\$PATH:$bin_dir\"" >> "$HOME/.zshrc"
            info "Added to PATH in ~/.zshrc"
            ;;
        bash)
            echo "export PATH=\"\$PATH:$bin_dir\"" >> "$HOME/.bash_profile"
            info "Added to PATH in ~/.bash_profile"
            ;;
        fish)
            echo "set -gx PATH \$PATH $bin_dir" >> "$HOME/.config/fish/config.fish"
            info "Added to PATH in ~/.config/fish/config.fish"
            ;;
        *)
            warn "Please add $bin_dir to your PATH manually"
            ;;
    esac

    echo ""
    info "Run 'flowforge --help' to get started"
}

main "$@"
