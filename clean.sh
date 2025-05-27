#!/bin/bash
# NeruBot Cleanup - Remove temporary files and cache
echo "🧹 Cleaning NeruBot..."

# Remove Python cache and compiled files
find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
find . -name "*.pyc" -delete 2>/dev/null || true
find . -name "*.pyo" -delete 2>/dev/null || true

# Remove old logs (>7 days)
find . -name "*.log" -mtime +7 -delete 2>/dev/null || true

# Remove temporary and backup files
find . -name "*.tmp" -delete 2>/dev/null || true
find . -name "*.bak" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true
find . -name ".DS_Store" -delete 2>/dev/null || true

echo "✅ Cleanup complete! Project size: $(du -sh . | cut -f1)"
