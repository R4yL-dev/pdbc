# ProtonDB Checker - pdbc

CLI tool to search for games on Steam and display their ProtonDB compatibility tier.

## Build

```bash
make
```

This compiles the project and creates the `pdbc` binary.

## Installation

Install to `/usr/local/bin`:

```bash
make install
```

Uninstall:

```bash
make uninstall
```

## Other commands

- `make run` - Build and run the program
- `make clean` - Remove binary and build artifacts

## Usage

```bash
pdbc <search_term> [search_term2] [search_term3] ...
```

### Examples

Search for a single game:
```bash
pdbc "Half-Life"
```

Search for multiple games:
```bash
pdbc "Anno" "Cyberpunk" "Half-Life"
```

### Output

Results are displayed in a table with the following columns:
- **Game Name**: Title of the game from Steam
- **Tier**: ProtonDB compatibility rating (Platinum, Gold, Silver, Bronze, Borked, Pending, Unknown)
- **Confidence**: Rating confidence level (Strong, Good, Moderate, Weak, Low, Inadequate, Unknown)
