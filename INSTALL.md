# Installation Guide

Quick step-by-step guide to install and set up logseq-claude-indexer.

## Step 1: Install the Binary

### Option A: From Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/dyluth/logseq-claude-indexer.git
cd logseq-claude-indexer

# Build and install to ~/.local/bin
make build
make test
make install-user
```

### Option B: Using Go

```bash
go install github.com/dyluth/logseq-claude-indexer/cmd/logseq-claude-indexer@latest
```

## Step 2: Add to PATH (if using install-user)

Add this to your `~/.bashrc` or `~/.zshrc`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then reload your shell:

```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Step 3: Verify Installation

```bash
logseq-claude-indexer version
```

You should see version information.

## Step 4: Generate Initial Indexes

Navigate to your Logseq repository and generate indexes:

```bash
cd /path/to/your/logseq
logseq-claude-indexer generate --repo . --output .claude/indexes
```

You should see:
```
Scanning Logseq repository: .
Found X markdown files
Extracted Y tasks and Z references
✓ Created .claude/indexes/tasks-by-status.md
✓ Created .claude/indexes/tasks-by-priority.md
✓ Created .claude/indexes/timeline-recent.md
✓ Created .claude/indexes/timeline-full.md
✓ Created .claude/indexes/missing-pages.md
✓ Created .claude/indexes/time-tracking.md
✓ Created .claude/indexes/reference-graph.md
✓ Created .claude/indexes/dashboard.md
Index generation complete!
```

## Step 5: View the Dashboard

```bash
cat .claude/indexes/dashboard.md
```

This is the main overview file to share with Claude!

## Step 6: Set Up Git Hook (Automatic Updates)

While still in your Logseq repository:

```bash
# Create the post-commit hook
cat > .git/hooks/post-commit << 'EOF'
#!/bin/bash
# Auto-generate Claude indexes after each commit
logseq-claude-indexer generate --repo . --output .claude/indexes --quiet
EOF

# Make it executable
chmod +x .git/hooks/post-commit
```

## Step 7: Test the Git Hook

```bash
# Make a test commit
echo "test" >> pages/test-hook.md
git add pages/test-hook.md
git commit -m "test git hook"

# Verify indexes were updated (check timestamps)
ls -lt .claude/indexes/ | head
```

If the timestamps are recent, it worked! 🎉

## Optional: Add to .gitignore

If you don't want to commit the indexes:

```bash
echo ".claude/" >> .gitignore
git add .gitignore
git commit -m "gitignore Claude indexes"
```

Or if you want to commit them for team sharing:

```bash
git add .claude/indexes/
git commit -m "add Claude indexes"
```

## Troubleshooting

### "command not found: logseq-claude-indexer"

- Check if `~/.local/bin` is in your PATH: `echo $PATH`
- Verify the binary exists: `ls -l ~/.local/bin/logseq-claude-indexer`
- Reload your shell: `source ~/.bashrc` or `source ~/.zshrc`

### "No markdown files found"

- Make sure you're in the root of your Logseq repository
- Verify you have `pages/` or `journals/` directories
- Try with `--verbose` flag: `logseq-claude-indexer generate --repo . --verbose`

### Git hook not running

- Check hook is executable: `ls -l .git/hooks/post-commit`
- Make it executable: `chmod +x .git/hooks/post-commit`
- Verify hook content: `cat .git/hooks/post-commit`
- Check indexer is in PATH: `which logseq-claude-indexer`

### Need help?

- Check the [README.md](README.md) for detailed documentation
- File an issue: https://github.com/dyluth/logseq-claude-indexer/issues

## Next Steps

1. Share `dashboard.md` with Claude to get context on your knowledge base
2. Read `.claude/indexes/README.md` to understand each index file
3. Explore different indexes based on your needs (see README.md)

---

**You're all set!** Indexes will now auto-generate after every commit. 🚀
