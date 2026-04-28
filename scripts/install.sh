#!/usr/bin/env sh
set -eu

repo="Ebube-Ugwu/snape"
install_dir="${SNAPE_INSTALL_DIR:-$HOME/.local/bin}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
ext=""

case "$os" in
  linux) platform="linux" ;;
  darwin) platform="macos" ;;
  mingw*|msys*|cygwin*) platform="windows"; ext=".exe" ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

asset="snape-${platform}-${arch}${ext}"
url="https://github.com/${repo}/releases/latest/download/${asset}"
tmp="$(mktemp)"

echo "downloading $url"
curl -fL "$url" -o "$tmp"
mkdir -p "$install_dir"
install -m 755 "$tmp" "$install_dir/snape${ext}"
rm -f "$tmp"

echo "installed $install_dir/snape${ext}"
echo "make sure $install_dir is on your PATH"
