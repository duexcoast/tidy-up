package tidy

import (
	"io/fs"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

type sortScenario struct {
	testID              int
	initialDirsPresent  []string
	initialFilesPresent []string
	want                map[string][]string
}

func TestFiletypeSort(t *testing.T) {
	t.Log("Given the need to sort a directory by filetype.")
	tests := map[string]sortScenario{
		"No initial scaffolding. No files to sort.": {
			testID:              0,
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
		"No initial scaffolding. Small amount of files to sort.": {
			testID:             1,
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
				"Documents":   {"resume2023.docx", "story.txt"},
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
			t.Logf("\tTest %d:\t%s", tc.testID, name)

			Tidy := NewTidy(NewFiletypeSorter(), afero.NewMemMapFs())

			// Setup the initial state of the directory before testing.
			// Create the initial directories.
			for _, v := range tc.initialDirsPresent {
				err := Tidy.Fs.Mkdir(v, 0777)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to setup starting state of directories in the test filesystem, error: %v", failed, tc.testID, err)
				}
			}
			t.Logf("\t%s\tTest %d:\tShould be able to setup starting state of directories in the test filesystem.", success, tc.testID)

			// Create the initial files.
			for _, v := range tc.initialFilesPresent {
				file, err := Tidy.Fs.Create(v)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to setup starting state of files in the test filesystem: %v", failed, tc.testID, err)
				}
				defer file.Close()
			}
			t.Logf("\t%s\tTest %d:\tShould be able to setup starting state of files in the test filesystem.", success, tc.testID)
			t.Logf("\t%s\tTest %d:\tTest successfully setup mock MemMapFS.", success, tc.testID)

			// This is what we're testing
			if err := Tidy.Sort(); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to call Tidy.Sort() without error: %v", failed, tc.testID, err)
			}

			got, err := mapOfDirs(t, Tidy.Fs)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create map of final directory structure: %v", failed, tc.testID, err)
			}
			// t.Logf("\t%s\tTest %d: Should be able to create map of final directory structure: %v", success, tc.testID)

			if !cmp.Equal(got, tc.want) {
				t.Logf("\t\tTest %d:\texp: %v", tc.testID, tc.want)
				t.Logf("\t\tTest %d:\tgot: %v", tc.testID, got)
				t.Logf("\t\tTest %d:\tdiff: %v", tc.testID, cmp.Diff(got, tc.want))
				t.Fatalf("\t%s\tTest %d:\tShould have sorted files by their extension.", failed, tc.testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have sorted files by their extension", success, tc.testID)

		})
	}
}

type scaffoldScenario struct {
	testID              int
	initialDirsPresent  []string
	initialFilesPresent []string
	want                []string
}

func TestCreateScaffolding(t *testing.T) {
	t.Log("Given the need to scaffold the directory structure for sorting.")

	tests := map[string]scaffoldScenario{
		"empty dir": {
			testID:              0,
			initialDirsPresent:  []string{},
			initialFilesPresent: []string{},
			want:                []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
		},
		"scaffolding already present": {
			testID:              1,
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
			initialFilesPresent: []string{},
			want:                []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
		},
		"scaffolding partially present": {
			testID:              2,
			initialDirsPresent:  []string{"Compressed", "Directories", "Images"},
			initialFilesPresent: []string{},
			want:                []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
		},
		"scaffolding present alongside other files": {
			testID:              3,
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
			initialFilesPresent: []string{"story.txt", "intl-players-anthem.mp3", "programming-pearls.pdf", "kobe.iso", "config.lua"},
			want:                []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos"},
		},
		"scaffolding present alongside extra directories": {
			testID:              3,
			initialDirsPresent:  []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Images", "Other", "PDFs", "Videos", "Extra", "dotfiles"},
			initialFilesPresent: []string{},
			want:                []string{"Audio", "Code", "Compressed", "Directories", "Documents", "Extra", "Images", "Other", "PDFs", "Videos", "dotfiles"},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", tc.testID, name)

			ftSorter := NewFiletypeSorter()
			Tidy := NewTidy(ftSorter, afero.NewMemMapFs())

			// Setup the initial state of the directory before testing.
			// Create the initial directories.
			for _, v := range tc.initialDirsPresent {
				err := Tidy.Fs.Mkdir(v, 0777)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to setup starting state of directories in the test filesystem, error: %v", failed, tc.testID, err)
				}
			}
			t.Logf("\t%s\tTest %d:\tShould be able to setup starting state of directories in the test filesystem", success, tc.testID)

			// Create the initial files.
			for _, v := range tc.initialFilesPresent {
				file, err := Tidy.Fs.Create(v)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to setup starting state of files in the test filesystem: %v", failed, tc.testID, err)
				}
				defer file.Close()
			}
			t.Logf("\t%s\tTest %d:\tShould be able to setup starting state of files in the test filesystem.", success, tc.testID)
			t.Logf("\t%s\tTest %d:\tTest successfully setup mock MemMapFS.", success, tc.testID)

			// This is what we're testing
			err := Tidy.CreateScaffolding()
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould have called Tidy.CreateScaffolding() without error: %v", failed, tc.testID, err)
			}

			got, err := sliceOfDirs(t, Tidy.Fs)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a slice of final directory structure: %v", failed, tc.testID, err)
			}

			if !(cmp.Equal(tc.want, got)) {
				t.Logf("\t\tTest %d:\texp: %v", tc.testID, tc.want)
				t.Logf("\t\tTest %d:\tgot: %v", tc.testID, got)
				t.Logf("\t\tTest %d:\tdiff: %v", tc.testID, cmp.Diff(got, tc.want))
				t.Fatalf("\t%s\tTest %d:\tShould have scaffolded the correct directory structure.", failed, tc.testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have scaffolded the correct directory structure.", success, tc.testID)
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
		fileSlice := make([]string, 0)

		err := afero.Walk(fsys, dirName, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
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
