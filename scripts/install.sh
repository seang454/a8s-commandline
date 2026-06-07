#!/usr/bin/env sh
set -eu

BINARY_PATH="${1:-dist/a8s-linux-amd64}"
INSTALL_DIR="${2:-$HOME/.local/bin}"

if [ ! -f "$BINARY_PATH" ]; then
  echo "Binary not found: $BINARY_PATH" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
cp "$BINARY_PATH" "$INSTALL_DIR/a8s"
chmod +x "$INSTALL_DIR/a8s"

echo "Installed a8s to $INSTALL_DIR/a8s"
echo "Add this directory to PATH if it is not already present:"
echo "  $INSTALL_DIR"
