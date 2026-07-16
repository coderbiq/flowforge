#!/bin/sh
# FlowForge CLI — 一键安装脚本 (macOS / Linux)
# Usage: curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash
# Options: --version <ver> --prefix <dir>
set -eu

APP_NAME="flowforge"
APP_VERSION="latest"
INSTALL_PREFIX=""
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

    # 每个平台可能对应多个制品（.tar.gz / .zip 共享同一 platform 字段），
    # 仅匹配第一个块（制品按 tar.gz 先于 zip 的顺序写入，正是 Unix 需要的），
    # 并取解析结果第一行，避免多制品叠加导致 URL 包含多行而触发 curl 报错。
    local platform_entry
    platform_entry="$(echo "$manifest_json" | grep -A5 -m1 "\"platform\": *\"$platform\"")" || {
        return 1
    }

    local url sha256
    url="$(echo "$platform_entry" | sed -n 's/.*"url": *"\([^"]*\)".*/\1/p' | head -n1)"
    sha256="$(echo "$platform_entry" | sed -n 's/.*"sha256": *"\([^"]*\)".*/\1/p' | head -n1)"

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

# ── 选择安装目录 ─────────────────────────────────────

# 按优先级尝试已在 PATH 中的可写目录
find_install_prefix() {
    # 用户通过 --prefix 指定了目录，直接使用
    if [ "$INSTALL_PREFIX" != "" ]; then
        return
    fi

    # 候选目录：已在 PATH 中、常见、用户可写
    for dir in /usr/local /opt/homebrew "$HOME/.local"; do
        local bin_dir="$dir/bin"
        if mkdir -p "$bin_dir" 2>/dev/null && [ -w "$bin_dir" ]; then
            INSTALL_PREFIX="$dir"
            return
        fi
    done

    # 都不行就用用户目录
    INSTALL_PREFIX="$HOME/.flowforge"
}

# ── PATH 配置 ──────────────────────────────────────

configure_path() {
    local bin_dir="$1"
    local shell_name
    shell_name="$(basename "${SHELL:-unknown}")"

    local profile_file=""
    case "$shell_name" in
        zsh)  profile_file="$HOME/.zshrc" ;;
        bash) profile_file="$HOME/.bash_profile"
              [ -f "$HOME/.bashrc" ] && profile_file="$HOME/.bashrc" ;;
        fish) profile_file="$HOME/.config/fish/config.fish" ;;
        *)    profile_file="" ;;
    esac

    local export_line="export PATH=\"$bin_dir:\$PATH\""
    [ "$shell_name" = "fish" ] && export_line="set -gx PATH $bin_dir \$PATH"

    if [ -n "$profile_file" ] && ! grep -Fq "$bin_dir" "$profile_file" 2>/dev/null; then
        mkdir -p "$(dirname "$profile_file")"
        echo "$export_line" >> "$profile_file"
        info "Added $bin_dir to $profile_file"
        info "Run 'source $profile_file' or open a new terminal to use flowforge"
    elif [ -z "$profile_file" ]; then
        warn "Unknown shell. Add this to your shell profile:"
        warn "  export PATH=\"$bin_dir:\$PATH\""
    fi
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

    find_install_prefix

    local bin_dir="$INSTALL_PREFIX/bin"
    local need_path_config=false

    # 检查是否需要 PATH 配置（装到了非系统 PATH 目录）
    case "$bin_dir" in
        /usr/local/bin|/opt/homebrew/bin|"$HOME/.local/bin") ;;
        *) need_path_config=true ;;
    esac

    mkdir -p "$bin_dir" || error "Failed to create $bin_dir"

    mkdir -p "$bin_dir" || error "Failed to create $bin_dir"

    tar xzf "$archive" -C "$tmpdir"
    mv "$tmpdir/${APP_NAME}" "$bin_dir/"
    if [ -d "$tmpdir/assets" ]; then
        rm -rf "$INSTALL_PREFIX/assets"
        mv "$tmpdir/assets" "$INSTALL_PREFIX/assets"
    fi
    chmod +x "$bin_dir/${APP_NAME}"

    rm -rf "$tmpdir"

    info "${APP_NAME} ${version} installed to $bin_dir/${APP_NAME}"

    # 如果装到了非系统 PATH 目录，自动追加到 shell profile
    if [ "$need_path_config" = true ]; then
        configure_path "$bin_dir"
    fi

    # 验证
    if command -v "$bin_dir/${APP_NAME}" >/dev/null 2>&1 || PATH="$PATH:$bin_dir" "$bin_dir/${APP_NAME}" --version >/dev/null 2>&1; then
        info "Verification: OK"
    fi

    echo ""
    info "Run 'flowforge init' to get started in a project"
}

main "$@"
