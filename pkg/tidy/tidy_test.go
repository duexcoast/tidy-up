package tidy

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

type fsState interface {
	dirsPresent() []string
	filesPresent() []string
}

type sortScenario struct {
	initialDirsPresent  []string
	initialFilesPresent []string
	want                map[string][]string
}

func (ss sortScenario) dirsPresent() []string {
	return ss.initialDirsPresent
}

func (ss sortScenario) filesPresent() []string {
	return ss.initialFilesPresent
}

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

			Tidy := NewTidy(NewFiletypeSorter())
			Tidy.Fs = afero.NewMemMapFs()

			// Setup the initial state of the directory before testing.
			if err := createInitialFsState(t, Tidy, tc); err != nil {
				t.Fatal(err)
			}

			if err := Tidy.Sort(); err != nil {
				t.Fatalf("could not sort the directory, error: %s", err)
			}

			// got, err := mapOfDirs(t, Tidy.Fs)
			// if err != nil {
			// 	t.Fatalf("could not assess the final state of fs, err: %s", err)
			// }
			//
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

func (ss scaffoldScenario) dirsPresent() []string {
	return ss.initialDirsPresent
}

func (ss scaffoldScenario) filesPresent() []string {
	return ss.initialFilesPresent
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
			Tidy := NewTidy(ftSorter)
			// By default the Fs is set to afero.NewOsFs but we want to change it
			// to a MemMapFs for efficient testing and painless clean up
			Tidy.Fs = afero.NewMemMapFs()

			// Setup the initial state of the directory before testing.
			if err := createInitialFsState(t, Tidy, tc); err != nil {
				t.Fatal(err)
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

// sliceOfDirContents function walks the current directory and returns a
// slice containing the name of every directory found.
func sliceOfDirs(t *testing.T, fsys afero.Fs) ([]string, error) {
	t.Helper()
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

func mapOfDirs(t *testing.T, fsys afero.Fs) (map[string][]string, error) {
	t.Helper()
	dirMap := make(map[string][]string)

	sliceOfDirs, err := sliceOfDirs(t, fsys)
	if err != nil {
		// t.Fatalf("failed constructing map of fs end state, error: %s", err)
		return nil, err
	}

	for _, dirName := range sliceOfDirs {
		fileSlice := make([]string, 0)

		err := afero.Walk(fsys, dirName, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fileSlice = append(fileSlice, info.Name())
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		dirMap[dirName] = fileSlice
	}

	return dirMap, nil
}

func createInitialFsState(t *testing.T, Tidy *Tidy, tc fsState) error {
	t.Helper()

	for _, v := range tc.dirsPresent() {
		err := Tidy.Fs.Mkdir(v, 0777)
		if err != nil {
			t.Fatalf("Could not create starting state of test dir, error: %s", err)
		}
	}

	for _, v := range tc.filesPresent() {
		file, err := Tidy.Fs.Create(v)
		if err != nil {
			t.Fatalf("Could not create the starting state of test dir, error: %s", err)
		}
		defer file.Close()
	}
	return nil

}
