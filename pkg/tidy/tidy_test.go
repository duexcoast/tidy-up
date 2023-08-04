package tidy

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

// type fsState interface {
// 	dirsPresent() []string
// 	filesPresent() []string
// }

type sortScenario struct {
	initialDirsPresent  []string
	initialFilesPresent []string
	want                map[string][]string
}

// func (ss sortScenario) dirsPresent() []string {
// 	return ss.initialDirsPresent
// }
//
// func (ss sortScenario) filesPresent() []string {
// 	return ss.initialFilesPresent
// }

func TestSort(t *testing.T) {
	tests := map[string]sortScenario{
		"no initial scaffolding. no files to sort": {
			initialDirsPresent:  []string{},
			initialFilesPresent: []string{},
			want: map[string][]string{
				"Audio":       {},
				"Code":        {},
				"Compressed":  {},
				"Directories": {},
				"Documents":   {},
				"Images":      {},
				"Other":       {},
				"PDFs":        {},
				"Videos":      {},
			},
		},
		"no initial scaffolding. small amount of files to sort": {
			initialDirsPresent: []string{},
			initialFilesPresent: []string{
				"story.txt",
				"intl-players-anthem.mp3",
				"programming-pearls.pdf",
				"kobe.iso",
				"config.lua",
				"home-video.mp4",
				"resume2023.docx",
				"random.xxx",
			},
			want: map[string][]string{
				"Audio":       {"intl-players-anthem.mp3"},
				"Code":        {"config.lua"},
				"Compressed":  {"kobe.iso"},
				"Directories": {},
				"Documents":   {"story.txt", "resume2023.docx"},
				"Images":      {},
				"Other":       {"random.xxx"},
				"PDFs":        {"programming-pearls.pdf"},
				"Videos":      {"home-video.mp4"},
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {

			Tidy := NewTidy(NewFiletypeSorter(), afero.NewMemMapFs())

			// Setup the initial state of the directory before testing.
			// Create the initial directories.
			for _, v := range tc.initialDirsPresent {
				err := Tidy.Fs.Mkdir(v, 0777)
				if err != nil {
					t.Fatalf("Could not create starting state of test dir, error: %s", err)
				}
			}

			// Create the initial files.
			for _, v := range tc.initialFilesPresent {
				file, err := Tidy.Fs.Create(v)
				if err != nil {
					t.Fatalf("Could not create the starting state of test dir, error: %s", err)
				}
				defer file.Close()
			}

			// This is what we're testing
			if err := Tidy.Sort(); err != nil {
				t.Fatalf("could not sort the directory, error: %s", err)
			}

			slice, err := sliceOfDirs(t, Tidy.Fs)
			fmt.Printf("[SLICE] %#v", slice)
			if err != nil {
				t.Fatal(err)
			}

			got, err := mapOfDirs(t, Tidy.Fs)
			if err != nil {
				t.Fatalf("could not assess the final state of fs, err: %s", err)
			}
			fmt.Printf("[GOT] %v\n", got)

			// if !cmp.Equal(tc.want, got) {
			// 	t.Fatalf("\ngot:\n\t%#v, \nwant:\n\t%#v", got, tc.want)
			// }

		})
	}
}

type scaffoldScenario struct {
	initialDirsPresent  []string
	initialFilesPresent []string
}

func TestCreateScaffolding(t *testing.T) {

	tests := map[string]scaffoldScenario{
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
			// Create the initial directories
			for _, v := range tc.initialDirsPresent {
				err := Tidy.Fs.Mkdir(v, 0777)
				if err != nil {
					t.Fatalf("Could not create starting state of test dir, error: %s", err)
				}
			}

			// Create the initial files
			for _, v := range tc.initialFilesPresent {
				file, err := Tidy.Fs.Create(v)
				if err != nil {
					t.Fatalf("Could not create the starting state of test dir, error: %s", err)
				}
				defer file.Close()
			}

			// This is what we're testing
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

// sliceOfDirContents function walks the current directory and returns a
// slice containing the name of every directory found.
func sliceOfDirs(t *testing.T, fsys afero.Fs) ([]string, error) {
	t.Helper()
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
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirsFound, nil
}

// mapOfDirs is a helper function that returns a map where the keys are all the
// directories present at the root, and the values are the files/dirs within them.
// The map is not recursive, it only displays one level of depth.
// This function is used for comparing the final sorted state against the desired
// state.
func mapOfDirs(t *testing.T, fsys afero.Fs) (map[string][]string, error) {
	t.Helper()
	dirMap := make(map[string][]string)

	sliceOfDirs, err := sliceOfDirs(t, fsys)
	if err != nil {
		return nil, err
	}

	for _, dirName := range sliceOfDirs {
		fmt.Printf("[loop] %s\n", dirName)
		fileSlice := make([]string, 0)

		err := afero.Walk(fsys, dirName, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fmt.Printf("[skip]\t%s\n", info.Name())
				return nil
			}
			fileSlice = append(fileSlice, info.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
		dirMap[dirName] = fileSlice
	}

	return dirMap, nil
}
