#!/usr/bin/env bash
set -e

echo "Installing evmrup..."

# Define variables
EVMR_DIR="${HOME}/.evm-runners"
EVMR_BIN_DIR="${EVMR_DIR}/bin"

EVMRUP_URL="https://raw.githubusercontent.com/ethernautdao/evm-runners-cli/main/evmrup"
EVMRUP_PATH="$EVMR_BIN_DIR/evmrup"

# Create bin directory if it doesn't exist and download evmrup 
mkdir -p "$EVMR_BIN_DIR"
curl -# -L $EVMRUP_URL -o $EVMRUP_PATH
chmod +x $EVMRUP_PATH

# Determine shell
SHELL_NAME="$(basename "$SHELL")"

# Determine appropriate shell configuration file
if [ "$SHELL_NAME" = "bash" ]; then
    CONFIG=$HOME/.bashrc
elif [ "$SHELL_NAME" = "zsh" ]; then
    CONFIG=$HOME/.zshrc
elif [ "$SHELL_NAME" = "fish" ]; then
    CONFIG=$HOME/.config/fish/config.fish
else
    echo "unsupported shell: $SHELL_NAME, manually add ${EVMR_BIN_DIR} to your PATH."
    exit 1
fi

# Update PATH environment variable in the shell configuration file if it doesn't exist
if [[ ":$PATH:" != *":${EVMR_BIN_DIR}:"* ]]; then
  echo >> $CONFIG && echo "export PATH=\"\$PATH:$EVMR_BIN_DIR\"" >> $CONFIG
  echo "Updated $CONFIG with PATH modification"
else
  echo "$CONFIG already contains PATH modification"
fi

echo && echo "Detected your preferred shell is $SHELL_NAME and added evmrup to PATH."
echo "Installing evm-runners by running 'evmrup'..."

# Run evmrup
$EVMRUP_PATH