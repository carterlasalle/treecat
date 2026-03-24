package highlight

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Options controls highlighting output.
type Options struct {
	Color bool // ANSI color output; false = plain text
}

// File syntax-highlights src as if it were filename.
// Returns plain src when Color=false or no lexer found.
func File(filename, src string, opts Options) (string, error) {
	if !opts.Color {
		return src, nil
	}

	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(src)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, src)
	if err != nil {
		return src, nil // fallback to plain
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return src, nil
	}
	return buf.String(), nil
}

// HexDump returns an xxd-style hex + ASCII dump of data.
func HexDump(data []byte) string {
	const cols = 16
	var sb strings.Builder
	for i := 0; i < len(data); i += cols {
		end := i + cols
		if end > len(data) {
			end = len(data)
		}
		row := data[i:end]

		fmt.Fprintf(&sb, "%08x  ", i)

		for j, b := range row {
			fmt.Fprintf(&sb, "%02x ", b)
			if j == 7 {
				sb.WriteString(" ")
			}
		}
		// Pad short rows
		for j := len(row); j < cols; j++ {
			sb.WriteString("   ")
			if j == 7 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(" |")

		for _, b := range row {
			if b >= 32 && b < 127 && unicode.IsPrint(rune(b)) {
				sb.WriteByte(b)
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteString("|\n")
	}
	return sb.String()
}
