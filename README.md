# today

A simple, command-line journal for tracking your daily wins and accomplishments.

![Demo](assets/demo.gif)

`today` helps you keep a record of what you've achieved, making it easier to write standup updates, performance reviews, or just reflect on your progress.

## Installation

### Using Nix

```bash
nix run github:mholtzscher/today
```

### Using Homebrew

```bash
brew tap mholtzscher/tap
brew install today
```

### From Source

```bash
git clone https://github.com/mholtzscher/today.git
cd today
nix build
```

## Usage

```bash
# Add an entry
today add "Fixed the login bug"

# Add an entry for a prior day
today add --date 2026-02-19 "Finished migration follow-up"

# Show today's entries
today show

# Show entries for the last 3 days
today show 3

# Archive an entry by id
today archive 12 --yes

# Show including archived entries
today show --all

# Restore an archived entry
today restore 12

# Show help
today --help
```

## Development

This project uses Nix for reproducible development environments.

```bash
# Enter development shell
nix develop

# Or use direnv
direnv allow

# Generate code (sqlc)
just generate

# Run checks
just check

# Build
just build

# Run tests
just test
```

## License

MIT
