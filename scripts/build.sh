#!/bin/bash
# FlowForge CLI — 跨平台编译脚本
set -eu

VERSION="${1:?Usage: build.sh <version> [native|all]}"
MODE="${2:-native}"

COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

LDFLAGS="-s -w"
LDFLAGS+=" -X flowforge/internal/version.injected=${VERSION}"

build_native() {
    echo "Building for current platform..."
    mkdir -p bin
    go build -ldflags="$LDFLAGS" -trimpath -o bin/flowforge ./cmd/flowforge
    echo "→ bin/flowforge ($(du -h bin/flowforge | cut -f1))"
}

build_all() {
    local out="dist/${VERSION}"
    rm -rf "$out"
    mkdir -p "$out"

    local targets=(
        "linux/amd64"
        "linux/arm64"
        "darwin/amd64"
        "darwin/arm64"
        "windows/amd64"
    )

    for target in "${targets[@]}"; do
        IFS='/' read -r goos goarch <<< "$target"

        local ext=""
        local archive_ext=".tar.gz"
        local platform="${goarch}-"

        case "$goos" in
            linux)   platform+="unknown-linux-gnu" ;;
            darwin)  platform+="apple-darwin" ;;
            windows) platform+="pc-windows-msvc"; ext=".exe"; archive_ext=".zip" ;;
        esac

        local bin_name="flowforge${ext}"
        local archive_name="flowforge-${platform}${archive_ext}"

        echo "Building ${goos}/${goarch} → ${archive_name}"

        GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 \
            go build -ldflags="$LDFLAGS" -trimpath \
            -o "${out}/${bin_name}" ./cmd/flowforge

        if [ "$goos" = "windows" ]; then
            (cd "$out" && zip "${archive_name}" "${bin_name}")
        else
            tar czf "${out}/${archive_name}" -C "$out" "${bin_name}"
        fi

        if command -v sha256sum >/dev/null 2>&1; then
            sha256sum "${out}/${archive_name}" | awk '{print $1}' > "${out}/${archive_name}.sha256"
        else
            shasum -a 256 "${out}/${archive_name}" | awk '{print $1}' > "${out}/${archive_name}.sha256"
        fi

        rm "${out}/${bin_name}"
    done

    cat "${out}"/*.sha256 > "${out}/checksums.txt"
    echo ""
    echo "Build complete: ${out}/"
    ls -lh "${out}/"
}

case "$MODE" in
    all)    build_all ;;
    native) build_native ;;
    *)      echo "Usage: $0 <version> [native|all]"; exit 1 ;;
esac
