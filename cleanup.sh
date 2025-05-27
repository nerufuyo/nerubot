#!/bin/bash
# NeruBot Cleanup Script
# Removes temporary files, cache, and other junk

echo "🧹 NeruBot Cleanup Script"
echo "========================="

# Remove Python cache directories
echo "🗑️  Removing Python cache directories..."
find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true

# Remove Python compiled files
echo "🗑️  Removing Python compiled files..."
find . -name "*.pyc" -delete 2>/dev/null || true
find . -name "*.pyo" -delete 2>/dev/null || true

# Remove log files older than 7 days
echo "🗑️  Cleaning old log files..."
find . -name "*.log" -mtime +7 -delete 2>/dev/null || true

# Remove temporary files
echo "🗑️  Removing temporary files..."
find . -name "*.tmp" -delete 2>/dev/null || true
find . -name "*.temp" -delete 2>/dev/null || true
find . -name ".DS_Store" -delete 2>/dev/null || true

# Remove backup files
echo "🗑️  Removing backup files..."
find . -name "*.bak" -delete 2>/dev/null || true
find . -name "*.backup" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true

# Remove IDE specific files
echo "🗑️  Removing IDE files..."
find . -name ".vscode" -type d -exec rm -rf {} + 2>/dev/null || true
find . -name ".idea" -type d -exec rm -rf {} + 2>/dev/null || true

# Display summary
echo ""
echo "✅ Cleanup completed!"
echo "📊 Current project size:"
du -sh . | cut -f1
echo ""
echo "💡 Run this script regularly to keep your project clean!"
