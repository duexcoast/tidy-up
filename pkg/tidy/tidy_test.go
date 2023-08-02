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
	want := helperSortedSliceOfDirs(ftSorter)

	wd, _ := os.Getwd()
	err := os.Chdir(path.Join(wd, "testFS"))
	if err != nil {
		t.Fatal(err)
	}
	wd, _ = os.Getwd()
	fmt.Println(wd)

	t.Run("generate correct list of directories.", func(t *testing.T) {
		err = ftSorter.createScaffolding()
		if err != nil {
			t.Fatal(err)
		}

		got, err := helperSliceOfDirectories(wd)
		if err != nil {
			t.Fatal(err)
		}
		assertSlicesEqual(t, got, want)

		err = helperCleanUpDirectories(wd)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("don't create duplicate directories.", func(t *testing.T) {
		err := os.Mkdir("Audio", fs.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		err = os.Mkdir("Compressed", fs.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		err = ftSorter.createScaffolding()
		if err != nil {
			t.Fatal(err)
		}

		got, err := helperSliceOfDirectories(wd)
		if err != nil {
			t.Fatal(err)
		}
		assertSlicesEqual(t, got, want)

		err = helperCleanUpDirectories(wd)
		if err != nil {
			t.Fatal(err)
		}

	})

}

// helperSliceOfDirectories function walks the current directory and returns a
// slice containing the name of every file/directory found.
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

func helperSortedSliceOfDirs(fts *filetypeSort) []string {
	want := fts.sliceOfDirs()
	want = append(want, ".")
	sort.Slice(want, func(i, j int) bool {
		return want[i] < want[j]
	})
	return want

}

// helperCleanUpDirectories function will erase every file or empty directory in
// the given working directory.
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
