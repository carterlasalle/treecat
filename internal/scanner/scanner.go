package scanner

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileNode represents a file or directory in the scanned tree.
type FileNode struct {
	Path      string
	Name      string
	IsDir     bool
	Size      int64
	Lines     int64
	Chars     int64
	Ext       string
	IsBinary  bool
	Children  []*FileNode
	Selected  bool
	Collapsed bool
	Depth     int
}

// Options controls scanner behaviour.
type Options struct {
	MaxDepth int  // 0 = unlimited
	Hidden   bool // include dot-prefixed files/dirs
}

// Scan walks root and returns a FileNode tree.
func Scan(root string, opts Options) (*FileNode, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return scan(abs, 0, opts)
}

func scan(path string, depth int, opts Options) (*FileNode, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Path:  path,
		Name:  info.Name(),
		IsDir: info.IsDir(),
		Ext:   strings.ToLower(filepath.Ext(info.Name())),
		Depth: depth,
	}

	if info.IsDir() {
		if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
			return node, nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return node, nil // skip unreadable dirs
		}
		for _, e := range entries {
			if !opts.Hidden && strings.HasPrefix(e.Name(), ".") {
				continue
			}
			child, err := scan(filepath.Join(path, e.Name()), depth+1, opts)
			if err != nil {
				continue
			}
			node.Children = append(node.Children, child)
		}
		return node, nil
	}

	// File: read metadata
	node.Size = info.Size()
	node.Lines, node.Chars, node.IsBinary = readMeta(path)
	return node, nil
}

const binaryCheckBytes = 8192

// readMeta counts lines/chars and detects binary files.
func readMeta(path string) (lines, chars int64, binary bool) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, false
	}
	defer f.Close()

	buf := make([]byte, binaryCheckBytes)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return 0, 0, false
	}
	buf = buf[:n]

	// Binary detection: null bytes = binary
	for _, b := range buf {
		if b == 0 {
			return 0, 0, true
		}
	}

	// Count lines/chars for whole file by re-reading from start
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return 0, 0, false
	}
	all, err := io.ReadAll(f)
	if err != nil {
		return 0, 0, false
	}
	chars = int64(len(all))
	for _, b := range all {
		if b == '\n' {
			lines++
		}
	}
	if chars > 0 && all[len(all)-1] != '\n' {
		lines++ // count final line without trailing newline
	}
	return lines, chars, false
}
