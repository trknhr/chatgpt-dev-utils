#!/bin/bash

set -e

REPO="your-org-or-username/chatgpt-dev-utils"
BINARY_NAME="cdev"
INSTALL_DIR="/usr/local/bin"

# Detect OS and ARCH
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="x86_64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Fetch latest release tag
TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Build download URL
FILENAME="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${FILENAME}"

# Download and extract
TMPDIR=$(mktemp -d)
echo "⬇️ Downloading $URL..."
curl -L "$URL" -o "$TMPDIR/$FILENAME"

echo "📦 Extracting..."
tar -xf "$TMPDIR/$FILENAME" -C "$TMPDIR"

echo "🚀 Installing to $INSTALL_DIR..."
chmod +x "$TMPDIR/$BINARY_NAME"
sudo mv "$TMPDIR/$BINARY_NAME" "$INSTALL_DIR/"

echo "✅ Installed: $INSTALL_DIR/$BINARY_NAME"
