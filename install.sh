#!/bin/bash
set -eu

# TODO: test both dynamic and static binaries
# TODO: detect architecture and os

bin_dir="$HOME/.local/bin"
mkdir -p "$bin_dir"

function download {
  if which wget >/dev/null 2>/dev/null; then
    wget -O "$1" "$2"
  else
    curl --output "$1" "$2"
  fi
}

download "$bin_dir/tunme" https://github.com/azert9/tunme/releases/download/latest/tunme_0.1.0_linux_amd64

chmod +x "$bin_dir/tunme"

if [ -d "$HOME/.oh-my-zsh/completions/" ]; then
  echo "Installing completions for ZSH."
  (exec -a tunme "$bin_dir/tunme" completion zsh) >"$HOME/.oh-my-zsh/completions/_tunme"
fi
