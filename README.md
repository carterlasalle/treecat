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
xattr -d com.apple.quarantine $(which treecat)
codesign --force --deep -s - $(which treecat)
```
Homebrew now needs to bypass macOS Gatekeeper

**macOS — direct binary download**
```bash
# Apple Silicon (M1/M2/M3/M4)
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_darwin_arm64.tar.gz | tar xz
sudo mv treecat /usr/local/bin/

# Intel Mac
curl -fsSL https://github.com/carterlasalle/treecat/releases/latest/download/treecat_darwin_amd64.tar.gz | tar xz
sudo mv treecat /usr/local/bin/
```

> **macOS Gatekeeper warning** — if you see *"cannot be opened because Apple cannot check it for malicious software"*, run:
> ```bash
> xattr -d com.apple.quarantine /usr/local/bin/treecat
> ```
> This removes the quarantine flag macOS applies to all internet-downloaded binaries. The binary is built from source in [public CI](https://github.com/carterlasalle/treecat/actions). 

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

# Include hidden files, ignore .gitignore
treecat . --hidden --no-ignore

# Interactive file selection TUI
treecat . -i
```

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
| `Tab` | Switch to preview panel |
| `x` | Toggle hex dump (binary files) |
| `H` | Toggle hidden files |
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

The extension filter bar shows all detected file types — it's display-only for now. Sorting by **size** (`s` twice) surfaces the largest files immediately — handy for catching accidentally included images or binaries.

## CI/CD

Releases are fully automated:

- **CI** runs on every push/PR: `go test`, `go vet`, `golangci-lint`
- **Release** triggers automatically when you push a `v*` tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GoReleaser builds binaries for Linux/macOS/Windows, creates `.deb`/`.rpm`/`.apk` packages, publishes a GitHub Release, and pushes the Homebrew formula automatically.

> **Note:** Add a `HOMEBREW_TAP_GITHUB_TOKEN` secret to your repo (Settings → Secrets → Actions) with a Personal Access Token that has `repo` scope on your `homebrew-treecat` repository.

## License

MIT
