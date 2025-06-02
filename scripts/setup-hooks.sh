#!/bin/bash

# Setup git hooks for documentation reminders

echo "🔧 Setting up git hooks for documentation reminders..."

# Set git hooks directory
git config core.hooksPath .githooks

echo "✅ Git hooks configured successfully!"
echo ""
echo "The pre-commit hook will now remind you to update documentation when:"
echo "  • package.json is modified"
echo "  • API endpoints are changed"
echo "  • New components are added"
echo "  • New backend handlers are added"
echo ""
echo "To disable hooks temporarily: git commit --no-verify"