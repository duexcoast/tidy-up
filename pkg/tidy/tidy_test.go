package tidy

import (
	"io/fs"
	"os"
	"sort"
	"testing"
)

// func TestCreateScaffolding(t *testing.T) {
// 	ftSorter := NewFiletypeSorter()
// 	want := helperSortedSliceOfDirs(ftSorter)
//
// 	wd, _ := os.Getwd()
// 	// TODO: This is brittle. need the path to be more deterministic.
// 	err := os.Chdir(path.Join(wd, "testFS"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wd, _ = os.Getwd()
//
// 	err = helperCleanUpDirectories(wd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	t.Run("generate correct list of directories.", func(t *testing.T) {
// 		err = ftSorter.createScaffolding()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		got, err := helperSliceOfDirectories(wd)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assertSlicesEqual(t, got, want)
//
// 		err = helperCleanUpDirectories(wd)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})
//
// 	t.Run("don't create duplicate directories.", func(t *testing.T) {
// 		err := os.Mkdir("Audio", fs.ModePerm)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		err = os.Mkdir("Compressed", fs.ModePerm)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		err = ftSorter.createScaffolding()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		got, err := helperSliceOfDirectories(wd)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assertSlicesEqual(t, got, want)
//
// 		err = helperCleanUpDirectories(wd)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 	})
//
// }

func TestSort(t *testing.T) {
	// TODO: Also brittle. need to make this more deterministic. Full path name
	// maybe
	wd, _ := os.Getwd()
	err := helperCleanUpDirectories(wd)
	if err != nil {
		t.Fatal(err)
	}
	ftSorter := NewFiletypeSorter()

	// wd, _ := os.Getwd()
	// err := os.Chdir(path.Join(wd, "testFS"))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	fileNames := []string{
		"story.txt",
		"intl-players-anthem.mp3",
		"programming-pearls.pdf",
		"kobe.iso",
		"config.lua",
		"home-video.mp4",
		"resume2023.docx",
		"random.xxx",
	}

	for _, v := range fileNames {
		file, err := os.Create(v)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}
	_ = os.Mkdir("Inception", fs.ModePerm)
	fsys := os.DirFS(wd)
	_ = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		return nil
	})

	// wd, _ := os.Getwd()
	err = ftSorter.Sort(wd)
	if err != nil {
		t.Fatal(err)
	}
	err = helperCleanUpDirectories(wd)
	if err != nil {
		t.Fatal(err)
	}
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

func helperSortedSliceOfDirs(fts *FiletypeSorter) []string {
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
		_ = os.RemoveAll(d.Name())
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

// setWD function sets the current working directory to the subdirectory /testFS
// It returns the directory name as a string and an error. Calling this function
// at the beginning of each test func ensures that we are performing all work in
// the correct location, no matter where we were when we started.
// func setWD() error {
// 	filepath.
// }

func chtmpdir(t *testing.T) func() {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := os.Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		os.RemoveAll(d)
	}
}

func TestCreateScaffolding(t *testing.T) {

	tests := map[string]struct {
		initialDirsPresent  []string
		initialFilesPresent []string
	}{
		"empty dir": {
			initialDirsPresent:  []string{},
			initialFilesPresent: []string{},
		},
		"scaffolding already present": {
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
			initialFilesPresent: []string{},
		},
		"scaffolding partially present": {
			initialDirsPresent:  []string{"Code", "Directories", "Images"},
			initialFilesPresent: []string{},
		},
		"scaffolding present alongside other files": {
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
			initialFilesPresent: []string{"story.txt", "intl-players-anthem.mp3", "programming-pearls.pdf", "kobe.iso", "config.lua"},
		},
	}

	ftSorter := NewFiletypeSorter()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

		})
	}

}
