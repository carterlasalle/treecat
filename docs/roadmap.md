# Roadmap and known limitations

## Near-term roadmap

- Keep release/install paths reproducible and verifiable across macOS, Linux, and Windows.
- Expand scanner and renderer regression coverage as new output modes are added.
- Improve TUI keyboard ergonomics and behavior in very narrow terminals.
- Keep large-tree performance predictable with bounded memory use and benchmark-driven changes.

## Known limitations

- Treecat deliberately skips symbolic links instead of following them, so linked trees are not expanded.
- Interactive mode depends on terminal capabilities and is less suitable for screen readers than the non-interactive `--no-color --no-syntax --output txt` path.
- Windows releases are ZIP archives; native MSI/WinGet publishing is not currently part of the release pipeline.
- Homebrew publication depends on the separate `carterlasalle/homebrew-treecat` tap and its release token remaining writable.
- Syntax highlighting is best-effort; unknown file types still render as plain text.

These are intentional current boundaries rather than hidden behavior. Changes that alter the stable CLI surface follow semantic versioning and should be called out in release notes.
