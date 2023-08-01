package tidy

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"testing"
)

func TestCreateScaffolding(t *testing.T) {
	// filesystem := fstest.MapFS{
	// 	"Images/":     {},
	// 	"Videos/":     {},
	// 	"Documents":   {},
	// 	"Audio/":      {},
	// 	"PDFs/":       {},
	// 	"Other/":      {},
	// 	"Compressed/": {},
	// 	"Code/":       {},
	// 	// "Directories/": {},
	// }
	//
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

	got := helperSliceOfDirectories(wd)
}

func helperSliceOfDirectories(cwd string) ([]string, error) {
	fsys := os.DirFS(cwd)
	dirsFound := make([]string, 0)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// CHANGE: make this a test helper func so I can use t.Fatal
			log.Fatal(err)
		}
		dirsFound = append(dirsFound, d.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirsFound, nil

}
