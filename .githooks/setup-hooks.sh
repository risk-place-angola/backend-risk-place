#!/bin/sh

echo "🔗 Configuring Git hooks path..."

git config core.hooksPath .githooks

echo "✅ Hooks installed successfully (using .githooks as hooks directory)."
