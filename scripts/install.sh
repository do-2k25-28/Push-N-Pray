#!/usr/bin/env sh

set -euo pipefail

if [ "$(uname -ms)" != "Linux x86_64" ]; then
  echo "Install script only supports Linux x86_64"
  exit 1
fi

TARGET="$HOME/.local/bin/pushnpray"

if [ "$1" = "-local" ]; then
  echo "Installing from local file ($2)"
  mv "$2" "$TARGET"
else
  echo "Downloading from latest GitHub release"
  curl -fs -O "$TARGET" https://github.com/do-2k25-28/Push-N-Pray/releases/latest/download/pushnpray-linux-amd64
  chmod +x "$TARGET"
fi

case $(basename "$SHELL") in
bash)
  echo "Adding bash completions"
  mkdir -p "$HOME/.local/share/bash-completion/completions"
  pushnpray completion bash > "$HOME/.local/share/bash-completion/completions/pushnpray"
  ;;
*)
  echo "Install script doesn't support your shell completions"
  ;;
esac

echo "Successfuly installed Push'N'Pray command line interface $(pushnpray version)"
