#!/usr/bin/env bash
# ============================================================
# 🚀 AZKIA — INSTALL SEKALI-RUN (build dari source)
# ------------------------------------------------------------
# Jalankan:  bash install.sh
# (aman dijalankan ulang — idempotent)
#
# Yang dilakukan:
#   1. Cek Go toolchain, install otomatis kalau belum ada
#   2. go build binary 'azkia'
#   3. Install ke /usr/local/bin (root) atau ~/.local/bin
# ============================================================
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

GO_VERSION_NEEDED="$(grep '^go ' go.mod | awk '{print $2}')"
GO_VERSION_NEEDED="${GO_VERSION_NEEDED:-1.24}"
GO_BIN="$(command -v go || true)"
# Go sering terinstall di /usr/local/go/bin tapi belum ada di PATH sesi ini
if [ -z "$GO_BIN" ] && [ -x /usr/local/go/bin/go ]; then
    GO_BIN=/usr/local/go/bin/go
    export PATH="$PATH:/usr/local/go/bin"
fi

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

echo ""
echo "=============================================="
echo "  🚀 AZKIA — install sekali-run"
echo "=============================================="

# ---------- 1. Cek / install Go ----------
ensure_go() {
    if [ -n "$GO_BIN" ]; then
        current="$(go version 2>/dev/null | grep -oE 'go[0-9]+\.[0-9]+' | head -1 | tr -d 'go')"
        echo "   ✓ Go ditemukan: $(go version | awk '{print $3}') (butuh >= $GO_VERSION_NEEDED)"
        # Bandingkan versi: butuh versi utama.minor >= kebutuhan
        if [ "$(printf '%s\n%s\n' "$current" "$GO_VERSION_NEEDED" | sort -V | head -1)" = "$current" ] \
           && [ "$current" != "$GO_VERSION_NEEDED" ]; then
            echo "   ⚠️  Versi Go ($current) lebih rendah dari kebutuhan ($GO_VERSION_NEEDED). Install ulang..."
            install_go
        fi
        return
    fi
    echo "   ⚠️  Go belum terinstall. Install otomatis dari go.dev..."
    install_go
}

install_go() {
    local arch="$(uname -m)"
    [ "$arch" = "x86_64" ] && arch="amd64"
    [ "$arch" = "aarch64" ] && arch="arm64"
    local os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    local version
    version="$(curl -fsSL "https://go.dev/VERSION?m=text" | head -1)"
    echo "   ⏬ Download $version ($os/$arch)..."
    curl -fsSL -o /tmp/go-install.tgz "https://go.dev/dl/${version}.${os}-${arch}.tar.gz"
    echo "   📦 Ekstrak ke /usr/local/go..."
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go-install.tgz
    rm -f /tmp/go-install.tgz
    # Tambahkan ke PATH untuk sesi login berikutnya
    if [ -d /etc/profile.d ] && ! grep -q "/usr/local/go/bin" /etc/profile.d/golang.sh 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/golang.sh
    fi
    export PATH="$PATH:/usr/local/go/bin"
    echo "   ✓ Go terinstall: $(go version | awk '{print $3}')"
}

echo ""
echo "[1/4] Cek Go toolchain..."
ensure_go
export PATH="$PATH:/usr/local/go/bin"
export GOPATH="${GOPATH:-$HOME/go}"

# ---------- 2. Build ----------
echo ""
echo "[2/4] Build binary azkia (ini butuh beberapa menit pertama kali)..."
go build -o azkia .
echo "   ✓ Build selesai"

# ---------- 3. Install ----------
echo ""
echo "[3/4] Install binary..."
INSTALL_DIR="/usr/local/bin"
if [ "$(id -u)" != "0" ] && [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        echo "   ⚠️  Tambahkan ke PATH:  export PATH=\$PATH:$INSTALL_DIR"
    fi
fi
install -m 0755 azkia "$INSTALL_DIR/azkia"
rm -f azkia   # hapus binary build di folder repo (biar tidak ke-commit)
echo "   ✓ Terinstall di: $INSTALL_DIR/azkia"

# ---------- 4. Verifikasi ----------
echo ""
echo "[4/4] Verifikasi..."
azkia --version 2>&1 | tail -1 || "$INSTALL_DIR/azkia" --version 2>&1 | tail -1

echo ""
echo "=============================================="
echo "  ✅ AZKIA siap dipakai!"
echo "----------------------------------------------"
echo "  Mulai TUI     : azkia"
echo "  Setup provider : azkia providers     (API key: OpenAI, Claude, Gemini, dll.)"
echo "  Prompt langsung: azkia -p \"pesan\""
echo "  Bantuan        : azkia --help"
echo "=============================================="
