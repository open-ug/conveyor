#!/bin/bash

set -e  # Exit immediately if a command fails

# Ensure LIBRARY_VERSION is set


# Set variables
MODULE_NAME="conveyor.cloud.cranom.tech"
VERSION="v0.0.1"
OUTPUT_DIR="./build"
ARCHIVE_NAME="driver-runtime.a"
ATHENS_PROXY="http://localhost:3000"
# Ensure the output directory exists
mkdir -p "$OUTPUT_DIR"

echo "🔧 Building precompiled Go archive (.a) for $MODULE_NAME@$VERSION..."

# Build the package and generate the .a file
go build -pkgdir "$OUTPUT_DIR" -o "$OUTPUT_DIR/$ARCHIVE_NAME" "$MODULE_NAME"

if [ $? -eq 0 ]; then
    echo "✅ Build successful! Archive saved at: $OUTPUT_DIR/$ARCHIVE_NAME"
else
    echo "❌ Build failed!"
    exit 1
fi

# Display the generated file
ls -lh "$OUTPUT_DIR/$ARCHIVE_NAME"

echo "🚀 Publishing $MODULE_NAME@$VERSION to Local Athens Proxy at $ATHENS_PROXY..."

# Set Go proxy to use local Athens
export GOPROXY="$ATHENS_PROXY"

# Initialize the module (required for new versions)
go get "$MODULE_NAME@$VERSION" || {
    echo "❌ Failed to initialize module $MODULE_NAME@$VERSION!"
    exit 1
}

# Push the module version to Athens
go list -m "$MODULE_NAME@$VERSION" || {
    echo "❌ Failed to publish $MODULE_NAME@$VERSION to Athens!"
    exit 1
}

echo "✅ Successfully published $MODULE_NAME@$VERSION to Athens Proxy at $ATHENS_PROXY"