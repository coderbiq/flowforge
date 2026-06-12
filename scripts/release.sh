#!/bin/bash
# FlowForge CLI — 发布打包脚本
set -eu

VERSION="${1:?Usage: release.sh <version>}"
DIST_DIR="dist/${VERSION}"

if [ ! -d "$DIST_DIR" ]; then
    echo "Error: ${DIST_DIR} not found. Run build.sh first."
    exit 1
fi

generate_manifest() {
    echo "Generating manifest.json..."
    cat > "${DIST_DIR}/manifest.json" << EOF
{
  "version": "${VERSION}",
  "published_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "release_notes": "Release ${VERSION}",
  "artifacts": [
EOF

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

        if [ "$first" = true ]; then first=false; else echo "," >> "${DIST_DIR}/manifest.json"; fi

        cat >> "${DIST_DIR}/manifest.json" << EOF
    {
      "platform": "${platform}",
      "url": "https://cdn.flowforge.dev/release/${VERSION}/${filename}",
      "sha256": "${sha256}",
      "size_bytes": ${size}
    }
EOF
    done

    echo "  ]" >> "${DIST_DIR}/manifest.json"
    echo "}" >> "${DIST_DIR}/manifest.json"
}

sign_manifest() {
    if [ -f "$HOME/.flowforge-signing-key.pem" ]; then
        echo "Signing manifest..."
        openssl pkeyutl -sign \
            -inkey "$HOME/.flowforge-signing-key.pem" \
            -rawin -in "${DIST_DIR}/manifest.json" \
            -out "${DIST_DIR}/manifest.json.sig"
        echo "→ manifest.json.sig"
    else
        echo "Warning: No signing key found at ~/.flowforge-signing-key.pem"
    fi
}

upload_to_cdn() {
    if command -v ossutil >/dev/null 2>&1; then
        echo "Uploading to OSS..."
        ossutil cp -r "${DIST_DIR}/" "oss://flowforge-releases/release/${VERSION}/" --update
        echo "${VERSION}" > /tmp/release-latest.txt
        ossutil cp /tmp/release-latest.txt oss://flowforge-releases/release-latest.txt
        echo "Done."
    else
        echo "Warning: ossutil not found. Files ready in ${DIST_DIR}/"
        echo "Upload manually to your CDN."
        echo ""
        echo "Files to upload:"
        ls -1 "${DIST_DIR}/"
        echo ""
        echo "Then update release-latest.txt to: ${VERSION}"
    fi
}

generate_manifest
sign_manifest
upload_to_cdn
