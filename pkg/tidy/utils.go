package tidy

import (
	"path/filepath"
	"strings"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// getExtension returns the extension of a filename with no preceding ".".
// For example, passing in the string "telephone.txt" would return "txt".
// The extension is defined as the suffix after the final "." in the provided path.
// If the provided path does not contain a "." then getExtension will return an
// empty string "".
func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimPrefix(ext, ".")
}
