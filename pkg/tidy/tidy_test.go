package tidy

import (
	"io/fs"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
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

// func TestSort(t *testing.T) {
// 	// TODO: Also brittle. need to make this more deterministic. Full path name
// 	// maybe
// 	wd, _ := os.Getwd()
// 	err := helperCleanUpDirectories(wd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ftSorter := NewFiletypeSorter()
//
// 	// wd, _ := os.Getwd()
// 	// err := os.Chdir(path.Join(wd, "testFS"))
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
//
// 	fileNames := []string{
// 		"story.txt",
// 		"intl-players-anthem.mp3",
// 		"programming-pearls.pdf",
// 		"kobe.iso",
// 		"config.lua",
// 		"home-video.mp4",
// 		"resume2023.docx",
// 		"random.xxx",
// 	}
//
// 	for _, v := range fileNames {
// 		file, err := os.Create(v)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		file.Close()
// 	}
// 	_ = os.Mkdir("Inception", fs.ModePerm)
// 	fsys := os.DirFS(wd)
// 	_ = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 		return nil
// 	})
//
// 	// wd, _ := os.Getwd()
// 	err = ftSorter.Sort(wd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = helperCleanUpDirectories(wd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// sliceOfDirContents function walks the current directory and returns a
// slice containing the name of every directory found.
func sliceOfDirs(t *testing.T, fsys afero.Fs) ([]string, error) {
	dirsFound := make([]string, 0)
	err := afero.Walk(fsys, ".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != "." {
			dirsFound = append(dirsFound, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirsFound, nil
}

// func helperSortedSliceOfDirs(fts *FiletypeSorter) []string {
// 	want := fts.sliceOfDirs()
// 	// want = append(want, ".")
// 	sort.Slice(want, func(i, j int) bool {
// 		return want[i] < want[j]
// 	})
// 	return want
//
// }
//
// // helperCleanUpDirectories function will erase every file or empty directory in
// // the given working directory.
// func helperCleanUpDirectories(cwd string) error {
// 	fsys := os.DirFS(cwd)
// 	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 		_ = os.RemoveAll(d.Name())
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func assertSlicesEqual(t testing.TB, got, want []string) {
// 	if len(got) != len(want) {
// 		t.Errorf("got %v, want %v", got, want)
// 	}
// 	for i := range got {
// 		if got[i] != want[i] {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
//
// 	}
// }
//
// setWD function sets the current working directory to the subdirectory /testFS
// It returns the directory name as a string and an error. Calling this function
// at the beginning of each test func ensures that we are performing all work in
// the correct location, no matter where we were when we started.
// func setWD() error {
// 	filepath.
// }

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
			initialDirsPresent:  []string{"Compressed", "Directories", "Images"},
			initialFilesPresent: []string{},
		},
		"scaffolding present alongside other files": {
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
			initialFilesPresent: []string{"story.txt", "intl-players-anthem.mp3", "programming-pearls.pdf", "kobe.iso", "config.lua"},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {

			ftSorter := NewFiletypeSorter()
			Tidy := NewTidy(ftSorter, afero.NewMemMapFs())

			// Setup the initial state of the directory before testing.
			for _, v := range tc.initialDirsPresent {
				err := Tidy.Fs.Mkdir(v, 0777)
				if err != nil {
					t.Fatalf("Could not create starting state of test dir, error: %s", err)
				}
			}

			err := Tidy.Sorter.CreateScaffolding(Tidy.Fs)
			if err != nil {
				t.Fatalf("Couldn't create scaffolding, error: %s", err)
			}

			got, err := sliceOfDirs(t, Tidy.Fs)
			if err != nil {
				t.Fatalf("Couldn't read the dirs in the test filesystem, error: %s", err)
			}
			want := ftSorter.DirsSlice()

			// diff := cmp.Diff(got, want)
			// if diff != "" {
			// 	t.Fatal(diff)
			// }
			if !(cmp.Equal(want, got)) {
				t.Fatalf("\ngot:\n\t%#v, \nwant:\n\t%#v", got, want)
			}

		})
	}
}
