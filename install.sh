#!/bin/sh
set -eu

if [ "$#" -ne 0 ]; then
    echo 'Usage: install.sh' >&2
    exit 1
fi

system=$(uname -s)
case "$system" in
    Linux) ;;
    *) echo "Unsupported system: $system. Download a binary from Releases." >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac

asset="badminton_linux_${arch}"
base="https://github.com/GGBondHere/BitBMTGrab/releases/latest/download"
dir="${XDG_CONFIG_HOME:-$HOME/.config}/BitBMTGrab"
case "$dir" in
    /*) ;;
    *) echo 'Config directory must be an absolute path.' >&2; exit 1 ;;
esac
program="$dir/badminton"
if [ "$(id -u)" = 0 ]; then
    bin_dir=/usr/local/bin
else
    bin_dir="$HOME/.local/bin"
fi
shortcut="$bin_dir/bmt"
if [ -e "$shortcut" ] || [ -L "$shortcut" ]; then
    if [ ! -L "$shortcut" ] || [ "$(readlink "$shortcut")" != "$program" ]; then
        echo "Command already exists, not overwritten: $shortcut" >&2
        exit 1
    fi
fi
if [ -d "$program" ]; then
    echo "Program path is a directory: $program" >&2
    exit 1
fi

mkdir -p "$dir" "$bin_dir"
chmod 700 "$dir"
work=$(mktemp -d "$dir/.download.XXXXXX")
trap 'rm -f "$work/program"; rmdir "$work"' 0

printf 'Downloading %s...\n' "$asset"
curl -fL --retry 3 --connect-timeout 15 -o "$work/program" "$base/$asset"
chmod 700 "$work/program"
mv -f "$work/program" "$program"
ln -sfn "$program" "$shortcut"
rmdir "$work"
trap - 0

printf 'Installed: %s\nCommand: %s\n' "$program" "$shortcut"
case ":$PATH:" in
    *":$bin_dir:"*) ;;
    *) printf 'To use bmt, add the command directory to PATH:\nexport PATH="%s:$PATH"\n' "$bin_dir" ;;
esac

# The installer may arrive through a pipe; the TUI must read from the terminal.
exec "$program" </dev/tty
