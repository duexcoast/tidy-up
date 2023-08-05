package tidy

import (
	"path/filepath"
	"strings"
)

// getExtension returns the extension of a filename with no preceding dot.
// For example, passing in the string "telephone.txt" would return "txt".
// The extension is defined as the suffix after the final dot in the provided path.
// If the provided path does not contain a "." then getExtension will return an
// empty string "".
func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimPrefix(ext, ".")
}
