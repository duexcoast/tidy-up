package tidy

import (
	"errors"
	"io/fs"
	"os"

	"github.com/spf13/afero"
)

// func ChangeWorkingDir(dirName string) error {
// 	// TODO: I should do a better job of sanitizing/validating the dirName arg
// 	if dirName == "" {
// 		// For now, this is the default for MacOS, I'll have to add a better
// 		// default solution for cross-platform later.
// 		dirName = "Downloads"
// 	}
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatal().Msg(err.Error())
// 	}
// 	downloadsDir := filepath.Join(usr.HomeDir, dirName)
// 	err = os.Chdir(downloadsDir)
// 	if err != nil {
// 		log.Fatal().Msg(err.Error())
// 		return err
// 	}
//
// 	cwd, _ := os.Getwd()
// 	log.Info().Str("cwd", cwd).Send()
//
// 	return nil
// }

// func createScaffolding(cfg TidyConfig) error {
// 	switch cfg.SortType {
// 	case "filetype":
// 		err := scaffold()
// 	}
// }

//
// // sortType struct contains configuration on how the directory should be sorted
// type sortType struct {
// 	sortMethod string
// 	sortDirs   []string
// }

// func (s *sortType) setDirs() {
// 	// TODO: Need to associate rules with each directory name. This needs to be
// 	// done in a way in which different types of rules can be applied to
// 	// different sorting methods
// 	tidyDirs := map[string]s.sortDirs{ // is s.sortDirs a valid type here? I
// 		// simply want to make the type []string,
// 		// but am aiming for clarity here
// 		"filetype":  {"Images", "Video", "PDFs", "Audio", "Applications", "Other"},
// 		"createdAt": {"Today", "This Week", "This Month", "This Year", "Older"},
// 	}
//
// 	s.sortDirs = tidyDirs[s.sortMethod]
// }

type Common struct {
	Fs afero.Fs
}

type Tidy struct {
	Sorter Sorter

	// Using afero to interact with the filesystem which allows easier mocking of
	// filesystem in tests
	Fs afero.Fs
}

func NewTidy(sorter Sorter, Fs afero.Fs) *Tidy {
	return &Tidy{
		Sorter: sorter,
		Fs:     Fs,
	}

}

func (t *Tidy) CreateScaffolding() {
	t.Sorter.CreateScaffolding(t.Fs)
}

func (t *Tidy) Sort() {
	t.Sorter.Sort()
}

// Sorter is an interface which allows different types of sorting. It requires a
// Sort() method which sorts a given directory.
type Sorter interface {
	CreateScaffolding(fsys afero.Fs) error
	Sort() error
}

// FiletypeSorter implements the Sorter interface and is used for sorting a directory
// based on filetype.
//
// It contains two maps which are used internally for determining where a file
// should be sorted based on its extension. These mappings can be updated through
// the UpdateMap() method.
//
//	type FiletypeSorter struct {
//		// dirsToExtension is a map where the keys are the directories in which files
//		// will be sorted, and the values are a slice containing the file extensions
//		// which belong in that directory.
//		dirsToExtension map[string][]string
//		// extensionToDir is an inverse map of dirsToExtension, where the key is a
//		// file extension and the value is the directory to which it should be sorted.
//		extensionToDir map[string]string
//	}
//
//	type SortingFolder interface {
//		Update(string) error
//		Rename(newname string) error
//	}
type FiletypeLookup map[string]*FiletypeSortingFolder

type FiletypeSorter struct {
	// Dirs provides a slice of *FiletypeSortingFolders, representing the directories
	// in which files will be sorted.
	Dirs []*FiletypeSortingFolder

	// The Lookup map uses file extensions as keys, with *FiletypeSortingFolder as
	// values, this allows us to determine where a file should be sorted in constant
	// time.
	Lookup FiletypeLookup
}

type FiletypeSortingFolder struct {
	Name       string
	Extensions []string
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
			Extensions: []string{"jpeg", "jpg", "ai", "bmp", "gif", "heif", "heic", "ico", "max", "obj", "png", "ps", "psd", "svg", "tif", "tiff", "3ds", "3dm"},
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
	ftSorter.Lookup = ftSorter.NewLookup()

	return ftSorter
}

func (fts *FiletypeSorter) NewLookup() FiletypeLookup {
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
func (fts *FiletypeSorter) DirsSlice() []string {
	dirs := make([]string, 0, len(fts.Dirs))

	for _, v := range fts.Dirs {
		dirs = append(dirs, v.Name)
	}
	return dirs
}

func (fts *FiletypeSorter) CreateScaffolding(fsys afero.Fs) error {
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
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

// func initFiletypeLookupMap(dirs []*FiletypeSortingFolder) map[string]*FiletypeSortingFolder {
// 	lookup := make(map[string]*FiletypeSortingFolder)
//
// 	for _, v := range dirs {
//
// 	}
//
// }

// func newFiletypeSorter() *FiletypeSorter {
// 	// The dirs map keys contain the scaffolding structure for sorting. The values
// 	// represent th
// 	dirs := map[string][]string{
// 		"Images":      {"jpeg", "jpg", "ai", "bmp", "gif", "heif", "heic", "ico", "max", "obj", "png", "ps", "psd", "svg", "tif", "tiff", "3ds", "3dm"},
// 		"Videos":      {"avi", "flv", "h264", "m4v", "mkv", "mov", "mp4", "mpg", "mpeg", "mpeg-1", "mpeg-2", "mpeg-4", "", "rm", "swf", "vob", "wmv", "3g2", "3gp"},
// 		"Documents":   {"doc", "docx", "odt", "msg", "rtf", "tex", "txt", "wks", "wps", "wpd", "md"},
// 		"Code":        {"html", "js", "json", "ts", "tsx", "jsx", "go", "c", "cpp", "java", "awk", "sh", "zsh", "lua", "pl", "obj", "s", "sql", "py", "r", "rb", "rs", "cs", "kt", "php", "pm", "rkt", "rktl", "scm", "scala"},
// 		"Audio":       {"aa", "aax", "act", "aiff", "alac", "au", "wav", "flac", "ra", "wma", "ac3", "m4b", "mp3", "aac", "ots"},
// 		"PDFs":        {"pdf", "epub"},
// 		"Compressed":  {"a", "ar", "cpio", "shar", "lbr", "iso", "mar", "sbx", "tar", "br", "bz2", "f", "?xf", "genozip", "gz", "lz", "lz4", "lzma", "lzo", "rz", "sz", "sfark", "xz", "z", "zst", "7z", "s7z", "ace", "afa", "alz", "apk", "arc", "ark", "arc", "cdx", "arj", "b1", "b6z", "ba", "bh", "cab", "car", "cfs", "cpt", "dar", "dd", "dgc", "dmg", "ear", "gca", "genozip", "ha", "hki", "ice", "kgb", "lzh", "lha", "lzx", "pak", "partimg", "paq6", "paq7", "paq8", "pea", "phar", "pim", "pit", "qda", "rar", "rk", "sda", "sea", "sen", "sfx", "shk", "sit", "sitx", "sqx", "tar.gz", "tgz", "tar.z", "tar.bz2", "tbz2", "tar.lz", "tlz", "tar.xz", "txz", "tar.zst", "uc", "uc0", "uc2", "ucn", "ur2", "ue2", "uca", "uha", "war", "wim", "xar", "xp3", "yz1", "zip", "zipx", "zoo", "zpaq", "zz", "ecc", "ecsbx", "par", "par2", "rev"},
// 		"Other":       {},
// 		"Directories": {},
// 	}
//
// 	extensionToDirMap := invertMap(dirs)
//
// 	return &FiletypeSorter{dirsToExtension: dirs, extensionToDir: extensionToDirMap}
// }

// invertMap function takes an argument myMap of type map[string][]string, returning a new
// map where the keys correspond to the individual string values in each slice of myMap.
// func invertMap(myMap map[string][]string) map[string]string {
// 	invertedMap := make(map[string]string)
//
// 	for k, val := range myMap {
//
// 		for _, v := range val {
// 			invertedMap[v] = k
// 		}
// 	}
// 	return invertedMap
// }

// createScaffolding method creates the correct directory structure for the
// FiletypeSorter. The list of directories are taken from the keys in fs.dirsToExtension
// It will first check if the directories already exist, if they do not, then it will
// proceed to create them. Any failure is returned as an error.
// func (fts *FiletypeSorter) createScaffolding() error {
//
// 	scaffoldFolderNames := fts.sliceOfDirs()
// 	for _, v := range scaffoldFolderNames {
// 		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
//
// 			err := os.Mkdir(v, fs.ModePerm)
// 			if err != nil {
// 				log.Print(err)
// 			}
// 		}
// 	}
//
// 	return nil
// }

// sliceOfDirs method takes the keys in the dirsToExtension map, and creates a
// slice from them. This method allows any changes in the dirsToExtension map to
// be reflected every time the sliceOfDirs function is called.
// func (fts *FiletypeSorter) sliceOfDirs() []string {
// 	scaffoldFolderNames := make([]string, 0, len(fts.dirsToExtension))
// 	for k := range fts.dirsToExtension {
// 		scaffoldFolderNames = append(scaffoldFolderNames, k)
// 	}
// 	return scaffoldFolderNames
// }

func (fts *FiletypeSorter) Sort() error {
	// // TODO: use sortDir field in TidyConfig struct to set the directory to be
	// // sorted. For now I will take an argument into sort()
	// wd := "/"
	// err := os.Chdir(wd)
	// if err != nil {
	// 	log.Print(err)
	// }
	//
	// err = fts.createScaffolding()
	// if err != nil {
	// 	log.Print(err)
	// }
	// scaffoldFolderNames := fts.sliceOfDirs()
	//
	// fileSystem := os.DirFS(wd)
	//
	// err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
	// 	// check if d is a directory, if it is check to see if it is part of the
	// 	// scaffolding. If d is part of the scaffolding then return, if not, then
	// 	// move it to the 'Directories' folder.
	// 	if d.IsDir() {
	// 		if contains(scaffoldFolderNames, d.Name()) || d.Name() == "." {
	// 			return nil
	// 		}
	// 		destinationDir := filepath.Join(wd, "Directories")
	// 		destinationPath := filepath.Join(destinationDir, d.Name())
	// 		err := os.Rename(path, destinationPath)
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 		return nil
	//
	// 	}
	// 	ext := getExtension(d.Name())
	//
	// 	val, ok := fts.extensionToDir[ext]
	// 	newLocation := filepath.Join(wd, val, d.Name())
	// 	if !ok {
	// 		newLocation = filepath.Join(wd, "Other", d.Name())
	// 		err := os.Rename(path, newLocation)
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 	}
	// 	err = os.Rename(path, newLocation)
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }
	//
	return nil

}

// func (sortType string) TidyConfig {
// 	switch sortType {
// 	case "filetype":
//
//
// 	}
// }
//
// // put this config inside of a function for now,
// func run() {
// 	cfg := struct {
// 			SortType string `default:"filetype"`
// 	}
//
// 	// create scaffolding
// 	createScaffolding(cfg)
// }
