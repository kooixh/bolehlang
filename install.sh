#!/bin/bash
set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="boleh"

echo "Building bolehlang..."
go build -o "$BINARY_NAME" .

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "Done! bolehlang installed as '$BINARY_NAME'."
echo "Try it out: boleh examples/hello.myl"
