package highlight

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

func GetSyntaxHighlightedContent(content []byte, filename string) string {
	var buf strings.Builder
	err := quick.Highlight(&buf, string(content), filepath.Ext(filename)[1:], "terminal", "monokai")
	if err != nil {
		return string(content)
	}
	return buf.String()
}
