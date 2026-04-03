# treecat

Print a recursive directory tree followed by syntax-highlighted contents of every selected file. Built for sharing code context with LLMs, documentation, and code review.

```
Directory Structure:

myproject/
├── main.go
├── README.md
└── src
    ├── handler.go
    └── handler_test.go

---
File: /myproject/main.go
---

package main
...
```

## Install

**Homebrew** (macOS / Linux — recommended)
```bash
brew tap carterlasalle/treecat
brew install treecat
```
`treecat` is now distributed as a Homebrew **formula** (not a cask), so Homebrew builds it from source and avoids the macOS quarantine/codesign workaround.

**macOS — direct binary download**
```bash
# Apple Silicon (M1/M2/M3/M4)
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_darwin_arm64.tar.gz | tar xz
sudo mv treecat /usr/local/bin/

# Intel Mac
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_darwin_amd64.tar.gz | tar xz
sudo mv treecat /usr/local/bin/
```

> **macOS Gatekeeper warning (direct download only)** — if you see *"cannot be opened because Apple cannot check it for malicious software"*, run:
> ```bash
> xattr -d com.apple.quarantine /usr/local/bin/treecat
> ```
> This removes the quarantine flag macOS applies to internet-downloaded binaries. Homebrew formula installs should not need this.

**Debian/Ubuntu**
```bash
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_linux_amd64.deb -o treecat.deb
sudo dpkg -i treecat.deb
```

**RPM (Fedora/RHEL)**
```bash
sudo rpm -i https://github.com/carterlasalle/treecat/releases/latest/download/treecat_linux_amd64.rpm
```

**Alpine**
```bash
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_linux_amd64.apk -o treecat.apk
sudo apk add --allow-untrusted treecat.apk
```

**Go install**
```bash
go install github.com/carterlasalle/treecat/cmd/treecat@latest
```

## Usage

```
treecat [path] [flags]

Flags:
  -o, --output string    terminal | md | txt  (default "terminal")
  -f, --file string      write output to file
  -e, --ext strings      include extensions: .go,.ts
      --no-ignore        disable .gitignore
      --hidden           include hidden files
      --max-depth int    max depth (0=unlimited)
      --max-size string  skip files larger than size (e.g. 1MB)
      --hex              hex dump binary files
  -s, --sort string      name|size|lines|ext (default "name")
      --tree-only        print tree only, skip file contents
      --no-tree          skip tree header
      --no-color         disable ANSI colors
      --no-syntax        disable syntax highlighting
  -i, --interactive      launch interactive TUI
      --relative         render file headers relative to the scanned root

Commands:
  completion [shell]     generate shell completions
  man                    generate a roff man page
  version                print version information
```

## Examples

```bash
# Current directory
treecat

# Export as markdown (great for LLM context)
treecat . -o md -f context.md

# Only Go files
treecat . -e .go

# Show biggest files first (useful for spotting accidentally included large files)
treecat . -s size

# Use relative file headers for shareable output
treecat . --relative --no-color --no-syntax

# Generate shell completions or a man page
treecat completion zsh > ~/.zsh/completions/_treecat
treecat man > treecat.1

# Interactive file selection TUI
treecat . -i
```

## Compatibility and stability

`treecat` follows semantic versioning for user-facing CLI behavior. Patch releases should only fix bugs. Minor releases may add flags, commands, or output improvements without breaking existing invocations. Breaking changes to flags, command names, or output defaults should be reserved for major releases and called out explicitly in release notes.

For script-friendly usage, prefer these stable entry points:

- `treecat [path] [flags]`
- `treecat completion <shell>`
- `treecat man`
- `treecat version`

## Shell completion and man pages

```bash
# bash
treecat completion bash > /etc/bash_completion.d/treecat

# zsh
mkdir -p ~/.zsh/completions
treecat completion zsh > ~/.zsh/completions/_treecat

# fish
treecat completion fish > ~/.config/fish/completions/treecat.fish

# PowerShell
treecat completion powershell > treecat.ps1

# man page
treecat man > treecat.1
sudo install -m 0644 treecat.1 /usr/local/share/man/man1/treecat.1
```

## Install verification matrix

After installing from any release channel, verify the package the same way:

```bash
treecat version
treecat --help
treecat completion bash >/dev/null
treecat man >/dev/null
treecat . --tree-only --no-color
```

Recommended spot checks by channel:

- **Homebrew:** `brew info carterlasalle/treecat/treecat && treecat version`
- **Direct tarball:** confirm the extracted `treecat` binary is on your `PATH`
- **`.deb` / `.rpm` / `.apk`:** verify the package manager install succeeded, then run the common checks above
- **`go install`:** confirm the binary under `$(go env GOPATH)/bin` or `$(go env GOBIN)` is the expected version

## Interactive TUI

Run `treecat -i` to open the interactive file selector:

```
┌─ treecat ./myproject ───────────────────────────────────────────┐
│ Tree                        │ Preview: main.go                   │
│─────────────────────────────│────────────────────────────────────│
│ ▼ [✓] ./               dir  │  1  package main                   │
│   [✓] main.go      2.1KB 7L │  2                                 │
│ → [✓] README.md      512B   │  3  import "fmt"                   │
│   [✗] logo.png     94.2KB   │                                    │
│        ↑ binary             │  [binary — use H to toggle hex]    │
│─────────────────────────────│────────────────────────────────────│
│ Extensions: [.go ✓] [.md ✓] [.png ✗]                           │
│─────────────────────────────────────────────────────────────────│
│ 2 files · 2.6KB · sort:name  ↑↓ move  spc toggle  ctrl+g gen   │
└─────────────────────────────────────────────────────────────────┘
```

**Tree panel**

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Move cursor |
| `pgup/pgdn` or `ctrl+u/ctrl+d` | Jump half a page |
| `←/→` or `h/l` | Collapse / expand directory |
| `Enter` | Toggle collapse/expand |
| `Space` | Toggle file selection |
| `a` | Select/deselect all direct children of dir |
| `A` | Select/deselect entire tree (recursive) |
| `s` | Cycle sort: name → size → lines → ext |
| `e` | Toggle the current file's extension filter |
| `E` | Reset extension filters |
| `g` | Toggle `.gitignore` filtering |
| `Tab` | Switch to preview panel |
| `x` | Toggle hex dump (binary files) |
| `H` | Toggle hidden files |
| `?` | Toggle help summary |
| `ctrl+g` | Open save dialog |
| `q` / `Esc` | Quit |

**Preview panel** (after pressing `Tab`)

| Key | Action |
|-----|--------|
| `↑/↓` | Scroll content |
| `pgup/pgdn` or `ctrl+u/ctrl+d` | Jump half a page |
| `Tab` | Switch back to tree |

**Save dialog** (after pressing `ctrl+g`)

| Key | Action |
|-----|--------|
| `Tab` | Cycle: Terminal → File → Both |
| Type / backspace | Edit the output file path |
| `Enter` | Confirm and generate output |
| `Esc` | Cancel, return to tree |

The extension filter bar shows all detected file types and updates live as you toggle file-type filters with `e`/`E`. Sorting by **size** (`s` twice) surfaces the largest files immediately — handy for catching accidentally included images or binaries.

### Accessibility and terminal notes

- `?` toggles an in-app help summary.
- Narrow terminals automatically switch to a single-panel layout; use `Tab` to swap between tree and preview.
- `--no-color` disables ANSI colors in exported output for accessibility, copy/paste, and CI logs.
- `NO_COLOR=1 treecat ...` is also respected in non-interactive mode.
- For the most predictable screen-reader and plain-text behavior, use `--no-color --no-syntax --output txt`.

### Diagnostics and bug reports

When reporting an issue, include:

- `treecat version`
- your install source (Homebrew, tarball, package, or `go install`)
- your operating system and terminal emulator
- the exact command you ran
- whether the issue reproduces with `--no-color --no-syntax`

## CI/CD

Releases are fully automated:

- **CI** runs on every push/PR: `go test`, `go vet`, `golangci-lint`
- **Release** triggers automatically when you push a `v*` tag:

```bash
# replace X.Y.Z with the next semver release (for example: 1.2.3)
git tag vX.Y.Z
git push origin vX.Y.Z
```

GoReleaser builds binaries for Linux/macOS/Windows, creates `.deb`/`.rpm`/`.apk` packages, publishes a GitHub Release, and pushes the Homebrew formula automatically.

> **Note:** Add a `HOMEBREW_TAP_GITHUB_TOKEN` secret to your repo (Settings → Secrets → Actions) with a Personal Access Token that has `repo` scope on your `homebrew-treecat` repository.

### Homebrew formula publishing checklist

If your `homebrew-treecat` repo used to publish a cask, migrate it once to formula layout:

```bash
# in carterlasalle/homebrew-treecat
mkdir -p Formula
if [ -f Casks/treecat.rb ]; then
  echo "Removing legacy Casks/treecat.rb"
  git rm -f Casks/treecat.rb
fi
if [ -f Formula/treecat.rb ] && grep -Eq "^cask [\"']" Formula/treecat.rb; then
  echo "Removing cask-based Formula/treecat.rb"
  git rm -f Formula/treecat.rb
fi
if [ -d Casks ]; then
  if [ -z "$(ls -A Casks)" ]; then
    echo "Removing empty Casks/ directory"
    rmdir Casks
  else
    echo "Leaving non-empty Casks/ directory for manual review"
  fi
fi
```

Do **not** move the old cask file into `Formula/`. A cask file starts with `cask "treecat" do` and cannot be loaded as a formula.
If `Formula/treecat.rb` currently starts with `cask`, delete it and let the next release regenerate it as a formula.

Your tap should end up like:

```text
homebrew-treecat/
├── Formula/
│   └── treecat.rb
└── README.md
```

Make sure the default branch is writable by the token in `HOMEBREW_TAP_GITHUB_TOKEN`.

Then verify each release updates the formula:

```bash
# in this repo
# replace X.Y.Z with the next semver release (for example: 1.2.3)
git tag vX.Y.Z
git push origin vX.Y.Z

# after the Release workflow finishes
brew update
brew tap carterlasalle/treecat
brew info carterlasalle/treecat/treecat
brew reinstall carterlasalle/treecat/treecat
treecat --version
```

Expected result: the release workflow creates/updates `Formula/treecat.rb` in `homebrew-treecat`, and users install from the formula with:

```bash
brew install carterlasalle/treecat/treecat
```

If you previously installed the cask or have a stale local tap checkout, reset locally:

```bash
brew uninstall --cask treecat 2>/dev/null || true
brew untap carterlasalle/treecat
brew tap carterlasalle/treecat
brew update
brew install carterlasalle/treecat/treecat
```

## License

MIT
