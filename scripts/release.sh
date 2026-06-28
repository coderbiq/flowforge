#!/bin/bash
# FlowForge CLI — 发布打包脚本
# 生成 manifest.json 并为每个制品生成 Ed25519 签名
# 发布由 GoReleaser + GitHub Actions 自动完成
set -eu

VERSION="${1:?Usage: release.sh <version>}"
DIST_DIR="dist"
RELEASES_BASE_URL="https://github.com/coderbiq/flowforge/releases/download/${VERSION}"

if [ ! -d "$DIST_DIR" ]; then
    echo "Error: ${DIST_DIR} not found. Run build.sh first."
    exit 1
fi

SIGNING_KEY="${FLOWFORGE_SIGNING_KEY_FILE:-${FLOWFORGE_SIGNING_KEY:-}}"

# ── Artifact 签名 ─────────────────────────────────────

sign_artifact() {
    local archive="$1"
    if [ ! -f "$SIGNING_KEY" ]; then
        echo "  (no signing key, skipping)"
        return 0
    fi

    local sig_file="${archive}.sig"
    if openssl pkeyutl -sign \
        -inkey "$SIGNING_KEY" \
        -rawin -in "$archive" \
        -out "$sig_file" 2>/dev/null; then
        echo "  → $(basename "$sig_file")"
    else
        echo "  Warning: signing failed"
        return 0
    fi
}

# ── Manifest 生成 ────────────────────────────────────

generate_manifest() {
    echo "Generating manifest.json..."

    printf '{\n  "version": "%s",\n' "$VERSION" > "${DIST_DIR}/manifest.json"
    printf '  "published_at": "%s",\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "${DIST_DIR}/manifest.json"
    printf '  "release_notes": "Release %s",\n' "$VERSION" >> "${DIST_DIR}/manifest.json"
    printf '  "artifacts": [\n' >> "${DIST_DIR}/manifest.json"

    local first=true
    for archive in "${DIST_DIR}"/flowforge-*.tar.gz "${DIST_DIR}"/flowforge-*.zip; do
        [ -f "$archive" ] || continue
        local filename
        filename=$(basename "$archive")
        local platform
        platform=$(echo "$filename" | sed 's/flowforge-//;s/\.\(tar\.gz\|zip\)$//')
        local sha256
        sha256=$(cat "${archive}.sha256" | awk '{print $1}')
        local size
        size=$(stat -f%z "$archive" 2>/dev/null || stat -c%s "$archive")

        if [ "$first" = true ]; then first=false; else printf ',\n' >> "${DIST_DIR}/manifest.json"; fi

        printf '    {\n' >> "${DIST_DIR}/manifest.json"
        printf '      "platform": "%s",\n' "$platform" >> "${DIST_DIR}/manifest.json"
        printf '      "url": "%s/%s",\n' "$RELEASES_BASE_URL" "$filename" >> "${DIST_DIR}/manifest.json"
        printf '      "sha256": "%s",\n' "$sha256" >> "${DIST_DIR}/manifest.json"
        printf '      "size_bytes": %s' "$size" >> "${DIST_DIR}/manifest.json"

        if [ -f "${archive}.sig" ]; then
            printf ',\n      "signature_url": "%s/%s.sig"' "$RELEASES_BASE_URL" "$filename" >> "${DIST_DIR}/manifest.json"
        fi

        printf '\n    }' >> "${DIST_DIR}/manifest.json"
    done

    printf '\n  ]\n}\n' >> "${DIST_DIR}/manifest.json"
}

# ── 签名所有制品 ─────────────────────────────────────

echo "Signing artifacts..."
for archive in "${DIST_DIR}"/flowforge-*.tar.gz "${DIST_DIR}"/flowforge-*.zip; do
    [ -f "$archive" ] || continue
    printf "  %s\n" "$(basename "$archive")"
    sign_artifact "$archive"
done

# ── 主流程 ───────────────────────────────────────────

generate_manifest

echo ""
echo "Release ${VERSION} artifacts ready."
echo "  Directory: ${DIST_DIR}/"
echo "  Manifest:  ${DIST_DIR}/manifest.json"
echo ""
echo "Signing key: ${SIGNING_KEY:-not configured}"
echo "Upload to GitHub Releases is handled by GoReleaser in CI."
