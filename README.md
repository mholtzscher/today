# today

A Go CLI tool built with Nix

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
# Show help
today --help

# Run example command
today example

# Run with verbose output
today --verbose example
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
