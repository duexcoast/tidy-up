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
// TODO: I can refactor this and many other pieces of code in this package using
// the os.ReadDir function, which returns all directory entries sorted by filename
// This would greatly simplify the code - the Walk function is needlessly complex
// for its simple purpose here.
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

// moveToParentDir will take the file located at the path arguement, and move it
// into the parent directory. An error will be returned if the renaming was unsuccesful.
// If the renaming was succesfull the newPath return value will be an absolute path to
// the file.
//
// For example: if you pass in a file located at "/Users/duexcoast/Downloads/test/myfile.txt"
// it will be moved to "/Users/duexcoast/Downloads/myfile.txt"
func moveToParentDir(fsys afero.Fs, path string) (newPath string, err error) {
	fileName := filepath.Base(path)
	wd, _ := os.Getwd()

	dest := filepath.Join(wd, fileName)

	err = fsys.Rename(path, dest)
	if err != nil {
		return "", &SortingError{Filename: fileName, AbsPath: dest, Sort: false, Err: err}
	}
	return dest, nil
}

// sliceIsSubset will return true if s1 is a subset of s2. Otherwise it will return false.
// This function requires that both slices are **sorted**, and will return incorrect values
// if they are not sorted.
func sliceIsSubset(s1, s2 []string) bool {
	for len(s1) > 0 {
		switch {
		case len(s2) == 0:
			return false
		case s1[0] == s2[0]:
			s1 = s1[1:]
			s2 = s2[1:]
		case s1[0] < s2[0]:
			return false
		case s1[0] > s2[0]:
			s2 = s2[1:]
		}
	}

	return true
}
