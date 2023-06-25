#!/bin/bash

# Define variables
APP_NAME="evm-runners"
APP_NAME_ALT="evmr"
GITHUB_USER="ethernautdao"
GITHUB_REPO="evm-runners-cli"
VERSION="v0.4.0"

main() {
  say "Installing evm-runners..."
  
  PLATFORM="$(uname -s)"
  case $PLATFORM in
    Linux)
      PLATFORM="linux"
      ;;
    Darwin)
      PLATFORM="darwin"
      ;;
    *)
      err "unsupported platform: $PLATFORM"
      ;;
  esac

  # Determine architecture
  ARCH="$(uname -m)"

  if [ "$ARCH" = "x86_64" ]; then
    # Redirect stderr to /dev/null to avoid printing errors if non Rosetta.
    if [ "$(sysctl -n sysctl.proc_translated 2>/dev/null)" = "1" ]; then
      ARCH="arm64" # Rosetta
    else
      ARCH="amd64" # Intel
    fi
  elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64" # Arm
  else
    ARCH="amd64" # Amd
  fi

  # Download release binary
  DOWNLOAD_URL="https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/download/$VERSION/$APP_NAME-$PLATFORM-$ARCH"

  echo
  say "Downloading the binary from '$DOWNLOAD_URL'"

  curl -L "$DOWNLOAD_URL" -o "$APP_NAME"

  # Set executable permissions
  chmod +x "$APP_NAME"

  INSTALL_DIR="${HOME}/.${APP_NAME}"

  echo
  say "Moving the binary to $INSTALL_DIR"
  echo 
  
  mkdir -p "$INSTALL_DIR"
  mv "$APP_NAME" "$INSTALL_DIR"
  ln -sf "$INSTALL_DIR/$APP_NAME" "$INSTALL_DIR/$APP_NAME_ALT"

  # Determine shell
  SHELL_NAME="$(basename "$SHELL")"

  # Update PATH environment variable in the appropriate shell configuration file
  if [ "$SHELL_NAME" = "bash" ]; then
    CONFIG=$HOME/.bashrc
  elif [ "$SHELL_NAME" = "zsh" ]; then
    CONFIG=$HOME/.zshrc
  elif [ "$SHELL_NAME" = "fish" ]; then
    CONFIG=$HOME/.config/fish/config.fish
  else
    err "unsupported shell: $SHELL_NAME, manually add ${INSTALL_DIR} to your PATH."
  fi

  say "Updating PATH..."

  # Update PATH environment variable in the shell configuration file
  if [ -z "$(grep "${INSTALL_DIR}" "$CONFIG")" ]; then
    echo >> $CONFIG && echo "export PATH=${INSTALL_DIR}:\$PATH" >> $CONFIG
    say "Updated $CONFIG with PATH modification"
  else
    say "$CONFIG already contains PATH modification"
  fi


  echo
  say "$APP_NAME version $VERSION installed successfully!"
  echo
  say "Run 'source $CONFIG' or start a new terminal session to use evm-runners."
  say "Then run 'evm-runners help', or alternatively 'evmr help', to get started."
}

say() {
  printf '%s\n' "$1"
}

warn() {
  say "warning: ${1}" >&2
}

err() {
  say "$1" >&2
  exit 1
}

main "$@" || exit 1
