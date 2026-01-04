# ProtonDB Checker - pdbc

CLI tool to search for games on Steam and display their ProtonDB compatibility tier.

## Installation

### From GitHub Releases (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/R4yL-dev/pdbc/releases).

#### Linux (amd64)

```bash
# Download the archive
wget https://github.com/R4yL-dev/pdbc/releases/latest/download/pdbc_<version>_linux_x86_64.tar.gz

# Extract
tar xzf pdbc_<version>_linux_x86_64.tar.gz

# Move to PATH
sudo mv pdbc /usr/local/bin/

# Verify installation
pdbc --version
```

#### macOS (Apple Silicon)

```bash
# Download and extract
curl -L https://github.com/R4yL-dev/pdbc/releases/latest/download/pdbc_<version>_darwin_arm64.tar.gz | tar xz

# Move to PATH
sudo mv pdbc /usr/local/bin/

# Verify installation
pdbc --version
```

#### macOS (Intel)

```bash
# Download and extract
curl -L https://github.com/R4yL-dev/pdbc/releases/latest/download/pdbc_<version>_darwin_x86_64.tar.gz | tar xz

# Move to PATH
sudo mv pdbc /usr/local/bin/

# Verify installation
pdbc --version
```

#### Windows

Download the appropriate `.zip` file from the [releases page](https://github.com/R4yL-dev/pdbc/releases), extract it, and add the directory to your PATH.

#### Available Platforms

- **Linux**: x86_64, i386, arm64, armv6 (Raspberry Pi 1/Zero), armv7 (Raspberry Pi 2/3)
- **macOS**: x86_64 (Intel), arm64 (Apple Silicon M1/M2/M3)
- **Windows**: x86_64, i386

#### Verify Checksums

```bash
# Download checksums file
wget https://github.com/R4yL-dev/pdbc/releases/latest/download/checksums.txt

# Verify (Linux/macOS)
sha256sum -c checksums.txt --ignore-missing
```

### From Source

#### Build

```bash
make
```

This compiles the project and creates the `pdbc` binary.

#### Install

Install to `/usr/local/bin`:

```bash
make install
```

Uninstall:

```bash
make uninstall
```

#### Clean

Remove binary and build artifacts:

```bash
make clean
```

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
