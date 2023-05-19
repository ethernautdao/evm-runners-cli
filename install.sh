#!/bin/bash

# Define variables
APP_NAME="evm-runners"
APP_NAME_ALT="evmr"
GITHUB_USER="ethernautdao"
GITHUB_REPO="evm-runners-cli"
VERSION="v0.2"

# Determine OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"
if [ "$OS" = "windows" ]; then
  echo "Error: Windows is not supported"
  exit 1
fi

# Determine architecture
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
else
  echo "Error: Unsupported architecture $ARCH"
  exit 1
fi

# Download release binary
DOWNLOAD_URL="https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/download/$VERSION/$APP_NAME-$OS-$ARCH"

echo "Downloading $DOWNLOAD_URL"

curl -L "$DOWNLOAD_URL" -o "$APP_NAME"

# Set executable permissions
chmod +x "$APP_NAME"

INSTALL_DIR="${HOME}/.${APP_NAME}"

echo "Moving the binary to $INSTALL_DIR"

mkdir -p "$INSTALL_DIR"
mv "$APP_NAME" "$INSTALL_DIR"
ln -sf "$INSTALL_DIR/$APP_NAME" "$INSTALL_DIR/$APP_NAME_ALT"

# Determine shell
SHELL_NAME="$(basename "$SHELL")"

# Update PATH environment variable in the appropriate shell configuration file
if [ "$SHELL_NAME" = "bash" ]; then
  CONFIG_FILE=~/.bashrc
elif [ "$SHELL_NAME" = "zsh" ]; then
  CONFIG_FILE=~/.zshrc
else
  echo "Error: Unsupported shell: $SHELL_NAME"
  exit 1
fi

echo ""
echo "Updating the PATH environment variable"

# Update PATH environment variable in the shell configuration file
if [ -z "$(grep "${INSTALL_DIR}" "$CONFIG_FILE")" ]; then
  echo "export PATH=${INSTALL_DIR}:\$PATH" >> "$CONFIG_FILE"
  echo "Updated $CONFIG_FILE with PATH modification"
else
  echo "$CONFIG_FILE already contains PATH modification"
fi

echo ""
echo "$APP_NAME version $VERSION installed successfully!"
echo ""
echo "Run 'source $CONFIG_FILE' or start a new terminal session to use evm-runners."
echo "Then run 'evm-runners help' to get started."
echo "Alternatively, you can also use 'evmr' instead of 'evm-runners'"
