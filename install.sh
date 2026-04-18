#!/bin/bash

set -e

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

REPO="DuckInAShert/gocommit"
BINARY="gocommit"
INSTALL_DIR="/usr/local/bin"

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
    echo "Error: Could not determine latest version"
    exit 1
fi

echo "Installing gocommit ${LATEST} for ${OS}-${ARCH}..."

URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}_${LATEST:1}_${OS}_${ARCH}.tar.gz"

TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

curl -fsSL "$URL" | tar -xz -C "$TMPDIR"

mkdir -p "$INSTALL_DIR"
mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "gocommit installed successfully!"
echo ""
echo "Next step: configure your AI provider"
echo "  gocommit setup"
echo ""
echo "Or manually:"
echo "  gocommit config api_key=YOUR_KEY base_url=https://opencode.ai/zen/go/v1 model=kimi-k2.5"
