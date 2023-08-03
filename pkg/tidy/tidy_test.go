package tidy

import (
	"io/fs"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

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

			err := Tidy.CreateScaffolding()
			if err != nil {
				t.Fatalf("Couldn't create scaffolding, error: %s", err)
			}

			got, err := sliceOfDirs(t, Tidy.Fs)
			if err != nil {
				t.Fatalf("Couldn't read the dirs in the test filesystem, error: %s", err)
			}
			want := ftSorter.dirsSlice()

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
