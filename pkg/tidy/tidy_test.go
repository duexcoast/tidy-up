package tidy

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"testing"
)

func TestCreateScaffolding(t *testing.T) {
	ftSorter := NewFiletypeSort()

	wd, _ := os.Getwd()
	err := os.Chdir(path.Join(wd, "testFS"))
	if err != nil {
		t.Fatal(err)
	}
	wd, _ = os.Getwd()
	fmt.Println(wd)

	err = ftSorter.createScaffolding()
	if err != nil {
		t.Fatal(err)
	}

	got, err := helperSliceOfDirectories(wd)
	if err != nil {
		t.Fatal(err)
	}
	sliceOfDirs := ftSorter.sliceOfDirs()
	sliceOfDirs = append(sliceOfDirs, ".")
	sort.Slice(sliceOfDirs, func(i, j int) bool {
		return sliceOfDirs[i] < sliceOfDirs[j]
	})
	assertSlicesEqual(t, got, sliceOfDirs)

	err = helperCleanUpDirectories(wd)
	if err != nil {
		t.Fatal(err)
	}

}

func helperSliceOfDirectories(cwd string) ([]string, error) {
	fsys := os.DirFS(cwd)
	dirsFound := make([]string, 0)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dirsFound = append(dirsFound, d.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirsFound, nil
}

func helperCleanUpDirectories(cwd string) error {
	fsys := os.DirFS(cwd)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		_ = os.Remove(d.Name())
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func assertSlicesEqual(t testing.TB, got, want []string) {
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
		}

	}
}
