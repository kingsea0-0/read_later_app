#!/bin/bash
# Build Chrome Extension zip for store submission
set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

ZIP_NAME="hoard-extension.zip"

rm -f "$ZIP_NAME"
zip -r "$ZIP_NAME" \
  manifest.json \
  background.js \
  content-hoard.js \
  icons/ \
  -x "*.DS_Store"

echo "Created $ZIP_NAME"
