#!/bin/sh
# FlowForge CLI — 一键安装脚本 (macOS / Linux)
# Usage: curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash
# Options: --version <ver> --prefix <dir>
set -eu

APP_NAME="flowforge"
APP_VERSION="latest"
INSTALL_PREFIX="$HOME/.flowforge"
RELEASES_BASE="https://github.com/coderbiq/flowforge/releases"

# ── 参数解析 ─────────────────────────────────────────

while [ $# -gt 0 ]; do
    case "$1" in
        --version) APP_VERSION="$2"; shift 2 ;;
        --prefix)  INSTALL_PREFIX="$2"; shift 2 ;;
        *)         APP_VERSION="$1"; shift ;;
    esac
done

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
        Linux)  _ostype="linux" ;;
        Darwin) _ostype="darwin" ;;
        *) error "unsupported OS: $_ostype" ;;
    esac

    case "$_cputype" in
        x86_64|amd64)   _cputype="amd64" ;;
        aarch64|arm64)  _cputype="arm64" ;;
        *) error "unsupported CPU: $_cputype" ;;
    esac

    # Rosetta 2 检测
    if [ "$_ostype" = "darwin" ] && [ "$_cputype" = "amd64" ]; then
        if [ "$(sysctl -n sysctl.proc_translated 2>/dev/null)" = "1" ]; then
            _cputype="arm64"
            info "Detected Rosetta 2, downloading arm64 build"
        fi
    fi

    RETVAL="${_ostype}-${_cputype}"
}

# ── manifest 解析 ────────────────────────────────────

manifest_url() {
    if [ "$1" = "latest" ]; then
        echo "${RELEASES_BASE}/latest/download/manifest.json"
    else
        echo "${RELEASES_BASE}/download/$1/manifest.json"
    fi
}

# ── manifest 解析 ────────────────────────────────────

get_artifact_info() {
    local manifest_url="$1"
    local platform="$2"
    local manifest_json
    manifest_json="$(curl -sfL "$manifest_url")" || return 1

    local platform_entry
    platform_entry="$(echo "$manifest_json" | grep -A5 "\"platform\": *\"$platform\"")" || {
        return 1
    }

    local url sha256
    url="$(echo "$platform_entry" | sed -n 's/.*"url": *"\([^"]*\)".*/\1/p')"
    sha256="$(echo "$platform_entry" | sed -n 's/.*"sha256": *"\([^"]*\)".*/\1/p')"

    if [ -z "$url" ] || [ -z "$sha256" ]; then
        return 1
    fi

    RETVAL_URL="$url"
    RETVAL_SHA256="$sha256"
}

# ── 下载与校验 ───────────────────────────────────────

download_and_verify() {
    local version="$1"
    local arch="$2"
    local url sha256

    local murl
    murl="$(manifest_url "$version")"
    if get_artifact_info "$murl" "$arch"; then
        url="$RETVAL_URL"
        sha256="$RETVAL_SHA256"
    else
        error "Failed to find artifact for ${arch} version ${version}"
    fi

    local tmpdir
    tmpdir="$(mktemp -d)" || error "Failed to create temp directory"
    local archive="${tmpdir}/${APP_NAME}.tar.gz"

    info "Downloading ${APP_NAME} ${version}..."
    curl -sSfL "$url" -o "$archive" || {
        rm -rf "$tmpdir"
        error "Download failed"
    }

    info "Verifying checksum..."
    local actual_checksum
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum="$(sha256sum "$archive" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum="$(shasum -a 256 "$archive" | awk '{print $1}')"
    else
        rm -rf "$tmpdir"
        error "No sha256sum or shasum found"
    fi

    if [ "$actual_checksum" != "$sha256" ]; then
        rm -rf "$tmpdir"
        error "Checksum mismatch: expected ${sha256}, got ${actual_checksum}"
    fi
    info "Checksum verified"

    RETVAL_TMPDIR="$tmpdir"
    RETVAL_ARCHIVE="$archive"
}

# ── 主流程 ───────────────────────────────────────────

main() {
    get_architecture
    local arch="$RETVAL"
    info "Detected: $arch"

    local version
    if [ "$APP_VERSION" = "latest" ]; then
        local murl
        murl="$(manifest_url latest)"
        version="$(curl -sfL "$murl" 2>/dev/null | sed -n 's/.*"version": *"\([^"]*\)".*/\1/p')" || {
            error "Failed to fetch latest version"
        }
    else
        version="$APP_VERSION"
    fi

    download_and_verify "$version" "$arch"
    local tmpdir="$RETVAL_TMPDIR"
    local archive="$RETVAL_ARCHIVE"

    local bin_dir="$INSTALL_PREFIX/bin"
    mkdir -p "$bin_dir"

    tar xzf "$archive" -C "$tmpdir"
    mv "$tmpdir/${APP_NAME}" "$bin_dir/"
    if [ -d "$tmpdir/assets" ]; then
        rm -rf "$INSTALL_PREFIX/assets"
        mv "$tmpdir/assets" "$INSTALL_PREFIX/assets"
    fi
    chmod +x "$bin_dir/${APP_NAME}"

    rm -rf "$tmpdir"

    info "${APP_NAME} ${version} installed to $bin_dir/${APP_NAME}"

    if command -v "$bin_dir/${APP_NAME}" >/dev/null 2>&1 || PATH="$PATH:$bin_dir" "$bin_dir/${APP_NAME}" --version >/dev/null 2>&1; then
        info "Verification: OK"
    fi

    # PATH 配置提示
    if ! echo "$PATH" | tr ':' '\n' | grep -Fxq "$bin_dir"; then
        echo ""
        warn "══════════════════════════════════════════════"
        warn "  $bin_dir is not in your PATH."
        warn ""
        warn "  Add this to your shell profile:"
        echo ""
        case "$(basename "${SHELL:-unknown}")" in
            zsh)  echo '  echo '\''export PATH="$HOME/.flowforge/bin:$PATH"'\'' >> ~/.zshrc && source ~/.zshrc' ;;
            bash) echo '  echo '\''export PATH="$HOME/.flowforge/bin:$PATH"'\'' >> ~/.bash_profile && source ~/.bash_profile' ;;
            fish) echo '  echo '\''set -gx PATH $HOME/.flowforge/bin $PATH'\'' >> ~/.config/fish/config.fish' ;;
            *)    echo '  export PATH="$HOME/.flowforge/bin:$PATH"' ;;
        esac
        warn "══════════════════════════════════════════════"
        echo ""
    fi

    echo ""
    info "Run 'flowforge init' to get started in a project"
}

main "$@"
