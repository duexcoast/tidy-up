package tidy

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
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

// dirsInCwd walks the current directory and returns a slice containing the name of
// every directory found. The returned slice will be lexicographically sorted.
func dirsInCwd(fsys afero.Fs) ([]string, error) {
	dirsFound := make([]string, 0)

	err := afero.Walk(fsys, ".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == "." {
				return nil
			}
			dirsFound = append(dirsFound, info.Name())
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(dirsFound, func(i, j int) bool {
		return dirsFound[i] < dirsFound[j]
	})
	return dirsFound, nil
}

func moveToParentDir(fsys afero.Fs, path string) error {
	fileName := filepath.Base(path)
	wd, _ := os.Getwd()

	dest := filepath.Join(wd, fileName)

	err := fsys.Rename(path, dest)
	if err != nil {
		return err
	}
	return nil
}
