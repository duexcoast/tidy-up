package tidy

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
)

type Tidy struct {
	Sorter Sorter

	// Using afero to interact with the filesystem which allows easier mocking of
	// filesystem in tests
	Fs afero.Fs

	// sortDir is the directory to be sorted. The current working directory will
	// be equal to the value of sortDir
	sortDir string
}

func NewTidy(sorter Sorter, fsys afero.Fs) (*Tidy, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Tidy{
		Sorter:  sorter,
		Fs:      fsys,
		sortDir: wd,
	}, nil
}

func (t *Tidy) ChangeSortDir(path string) error {
	// cleanPath := filepath.Clean(path)
	info, err := t.Fs.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		t.sortDir = path
		err := os.Chdir(path)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("the string passed is not a directory.")

}

// Method CreateScaffolding() creates the given scaffolding for the directory
// based upon the Sorter type.
func (t *Tidy) CreateScaffolding() error {
	if err := t.Sorter.createScaffolding(t.Fs); err != nil {
		return err
	}
	return nil
}

func (t *Tidy) Sort() error {
	if err := t.CreateScaffolding(); err != nil {
		return err
	}
	if err := t.Sorter.sort(t.Fs); err != nil {
		return err
	}
	return nil
}

// Sorter is an interface which allows different types of sorting. It requires a
// a CreateScaffolding() method which creates the directory structure for files to
// be sorted in, and a Sort() method which sorts a given directory
type Sorter interface {
	createScaffolding(fsys afero.Fs) error
	sort(fsys afero.Fs) error
	// undo() error
}

type FiletypeLookup map[string]*FiletypeSortingFolder

// FiletypeSorter implements the Sorter interface and is used for sorting a directory
// based on filetype.
type FiletypeSorter struct {
	// Dirs provides a slice of *FiletypeSortingFolders, representing the directories
	// in which files will be sorted.
	Dirs []*FiletypeSortingFolder

	// The Lookup map uses file extensions as keys, with *FiletypeSortingFolder as
	// values, this allows us to determine where a file should be sorted in constant
	// time.
	Lookup FiletypeLookup
}

// FiletypeSortingFolder represents an individual directory in which files will be sorted
// when using the FiletypeSorter.
type FiletypeSortingFolder struct {
	Name string

	// The Extensions field contains a slice of all file extensions that should be
	// sorted in this folder.
	Extensions []string
}

func (ftsf *FiletypeSortingFolder) String() string {
	return fmt.Sprintf("%s will store files with the following extensions: [ %s ]", ftsf.Name, strings.Join(ftsf.Extensions, ", "))
}

func NewFiletypeSorter() *FiletypeSorter {
	dirs := []*FiletypeSortingFolder{
		{
			Name:       "Audio",
			Extensions: []string{"aa", "aax", "act", "aiff", "alac", "au", "wav", "flac", "ra", "wma", "ac3", "m4b", "mp3", "aac", "ots"},
		},
		{
			Name:       "Code",
			Extensions: []string{"html", "js", "json", "ts", "tsx", "jsx", "go", "c", "cpp", "java", "awk", "sh", "zsh", "lua", "pl", "obj", "s", "sql", "py", "r", "rb", "rs", "cs", "kt", "php", "pm", "rkt", "rktl", "scm", "scala"},
		},
		{
			Name:       "Compressed",
			Extensions: []string{"a", "ar", "cpio", "shar", "lbr", "iso", "mar", "sbx", "tar", "br", "bz2", "f", "?xf", "genozip", "gz", "lz", "lz4", "lzma", "lzo", "rz", "sz", "sfark", "xz", "z", "zst", "7z", "s7z", "ace", "afa", "alz", "apk", "arc", "ark", "arc", "cdx", "arj", "b1", "b6z", "ba", "bh", "cab", "car", "cfs", "cpt", "dar", "dd", "dgc", "dmg", "ear", "gca", "genozip", "ha", "hki", "ice", "kgb", "lzh", "lha", "lzx", "pak", "partimg", "paq6", "paq7", "paq8", "pea", "phar", "pim", "pit", "qda", "rar", "rk", "sda", "sea", "sen", "sfx", "shk", "sit", "sitx", "sqx", "tar.gz", "tgz", "tar.z", "tar.bz2", "tbz2", "tar.lz", "tlz", "tar.xz", "txz", "tar.zst", "uc", "uc0", "uc2", "ucn", "ur2", "ue2", "uca", "uha", "war", "wim", "xar", "xp3", "yz1", "zip", "zipx", "zoo", "zpaq", "zz", "ecc", "ecsbx", "par", "par2", "rev"},
		},
		{
			Name:       "Directories",
			Extensions: []string{},
		},
		{
			Name:       "Documents",
			Extensions: []string{"doc", "docx", "odt", "msg", "rtf", "tex", "txt", "wks", "wps", "wpd", "md"},
		},
		{
			Name:       "Images",
			Extensions: []string{"jpeg", "jpg", "ai", "bmp", "gif", "heif", "heic", "ico", "max", "obj", "png", "ps", "psd", "svg", "tif", "tiff", "3ds", "3dm", "webp"},
		},
		{
			Name:       "Other",
			Extensions: []string{},
		},
		{
			Name:       "PDFs",
			Extensions: []string{"pdf"},
		},
		{
			Name:       "Videos",
			Extensions: []string{"avi", "flv", "h264", "m4v", "mkv", "mov", "mp4", "mpg", "mpeg", "mpeg-1", "mpeg-2", "mpeg-4", "", "rm", "swf", "vob", "wmv", "3g2", "3gp"},
		},
	}

	ftSorter := &FiletypeSorter{Dirs: dirs}
	ftSorter.Lookup = ftSorter.newLookup()

	return ftSorter
}

// The newLookup() method returns a FiletypeLookup map for the FiletypeSorter. The keys are
// retrieved by looping through the Extensions field of every FiletypeSortFolder
// in the FiletypeSorter.Dirs slice. The main purpose of this function is for use when
// initializing a new FiletypeSorter.
//
// newLookup() is guaranteed to be correct upon creation, but care must be taken to
// keep entries consistent across all data structures if the mappings are to be
// changed.
func (fts *FiletypeSorter) newLookup() FiletypeLookup {
	lookup := make(map[string]*FiletypeSortingFolder)

	for _, sortingFolder := range fts.Dirs {
		for _, extension := range sortingFolder.Extensions {
			lookup[extension] = sortingFolder
		}
	}
	return lookup
}

// DirsSlice returns a slice containing the names of all FiletypeSortingFolders in
// the Dirs slice.
func (fts *FiletypeSorter) dirsSlice() []string {
	dirs := make([]string, 0, len(fts.Dirs))

	for _, v := range fts.Dirs {
		dirs = append(dirs, v.Name)
	}
	return dirs
}

// createScaffolding reads the names of the elements in fts.Dirs and creates
// directories of the same names in the current working directory.
//
// If there is already a folder with the same name then createScaffolding
// will refrain from creating that directory. If there is a file with the same
// name, however, then an error will be returned.
func (fts *FiletypeSorter) createScaffolding(fsys afero.Fs) error {
	for _, v := range fts.Dirs {
		err := idempotentMkdir(v.Name, fs.ModePerm, fsys)
		if err != nil {
			return err
		}
	}
	return nil
}

// idempotentMkdir will create a directory with the given name if it does not exist
// if the directory already exists, idempotentMkdir will return without an error.
// This function is safe for concurrent execution.
//
// Taken from stackoverflow user @pr-pal: https://stackoverflow.com/a/56600630/18245016
func idempotentMkdir(name string, perm fs.FileMode, fsys afero.Fs) error {
	// We "do then check" to avoid race conditions, as opposed
	// to "check then do".
	err := fsys.Mkdir(name, perm)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := fsys.Stat(name)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// TODO: Need to deal with the case in which there is an existing file, that is not
			// a directory, but that has the same name as one of the needed directories. Currently
			// the behavior is to return an error. There should probably be a better solution.
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func (fts *FiletypeSorter) sort(fsys afero.Fs) error {
	err := afero.Walk(fsys, ".", func(path string, f fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if this is a directory, if it is check to see if it is part of the
		// scaffolding. If it is part of the scaffolding then return, if not - then
		// move it to the 'Directories' folder.
		if f.IsDir() {
			if slices.Contains(fts.dirsSlice(), f.Name()) {
				return filepath.SkipDir
			}
			if path == "." {
				return nil
			}

			dest := filepath.Join("Directories", f.Name())

			err := fsys.Rename(f.Name(), dest)
			if err != nil {
				return err
			}
			return filepath.SkipDir
		}
		ext := getExtension(f.Name())

		val, ok := fts.Lookup[ext]
		if !ok || ext == "" {
			dest := filepath.Join("Other", f.Name())
			err := fsys.Rename(f.Name(), dest)
			if err != nil {
				return err
			}
			return nil
		}
		dest := filepath.Join(val.Name, f.Name())
		err = fsys.Rename(f.Name(), dest)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

type CreatedAtSorter struct {
}
