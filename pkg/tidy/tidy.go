package tidy

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func ChangeWorkingDir(dirName string) error {
	// TODO: I should do a better job of sanitizing/validating the dirName arg
	if dirName == "" {
		// For now, this is the default for MacOS, I'll have to add a better
		// default solution for cross-platform later.
		dirName = "Downloads"
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	downloadsDir := filepath.Join(usr.HomeDir, dirName)
	err = os.Chdir(downloadsDir)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return err
	}

	cwd, _ := os.Getwd()
	log.Info().Str("cwd", cwd).Send()

	return nil
}

// func createScaffolding(cfg TidyConfig) error {
// 	switch cfg.SortType {
// 	case "filetype":
// 		err := scaffold()
// 	}
// }

// type TidyConfig struct {
// 	// embedded type sortType provides information on how the directory should
// 	// be sorted.
// 	sortType
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

// NOTE: the logic for how to apply the sorting rules for each sort type will be
// in the sort method on the sortable interface. So for example, I will have the
// struct filetypeSort, which is of type sortType. It implements the sortable
// interface. It has the sort() method, which defines the rules on which files go
// in which directories.
type sortable interface {
	createScaffolding() error
	sort()
}

type filetypeSort struct {
	// dirsToExtension is a map where the keys are the directories in which files
	// will be sorted, and the values are a slice containing the file extensions
	// which belong in that directory.
	dirsToExtension map[string][]string
	// extensionToDir is an inverse map of dirsToExtension, where the key is a
	// file extension and the value is the directory to which it should be sorted.
	extensionToDir map[string]string
}

func NewFiletypeSort() *filetypeSort {
	dirs := map[string][]string{
		"Images":      {"jpeg", "jpg", "ai", "bmp", "gif", "heif", "heic", "ico", "max", "obj", "png", "ps", "psd", "svg", "tif", "tiff", "3ds", "3dm"},
		"Videos":      {"avi", "flv", "h264", "m4v", "mkv", "mov", "mp4", "mpg", "mpeg", "mpeg-1", "mpeg-2", "mpeg-4", "", "rm", "swf", "vob", "wmv", "3g2", "3gp"},
		"Documents":   {"doc", "docx", "odt", "msg", "rtf", "tex", "txt", "wks", "wps", "wpd", "md"},
		"Code":        {"html", "js", "json", "ts", "tsx", "jsx", "go", "c", "cpp", "java", "awk", "sh", "zsh", "lua", "pl", "obj", "s", "sql", "py", "r", "rb", "rs", "cs", "kt", "php", "pm", "rkt", "rktl", "scm", "scala"},
		"Audio":       {"aa", "aax", "act", "aiff", "alac", "au", "wav", "flac", "ra", "wma", "ac3", "m4b", "mp3", "aac", "ots"},
		"PDFs":        {"pdf", "epub"},
		"Compressed":  {"a", "ar", "cpio", "shar", "lbr", "iso", "mar", "sbx", "tar", "br", "bz2", "f", "?xf", "genozip", "gz", "lz", "lz4", "lzma", "lzo", "rz", "sz", "sfark", "xz", "z", "zst", "7z", "s7z", "ace", "afa", "alz", "apk", "arc", "ark", "arc", "cdx", "arj", "b1", "b6z", "ba", "bh", "cab", "car", "cfs", "cpt", "dar", "dd", "dgc", "dmg", "ear", "gca", "genozip", "ha", "hki", "ice", "kgb", "lzh", "lha", "lzx", "pak", "partimg", "paq6", "paq7", "paq8", "pea", "phar", "pim", "pit", "qda", "rar", "rk", "sda", "sea", "sen", "sfx", "shk", "sit", "sitx", "sqx", "tar.gz", "tgz", "tar.z", "tar.bz2", "tbz2", "tar.lz", "tlz", "tar.xz", "txz", "tar.zst", "uc", "uc0", "uc2", "ucn", "ur2", "ue2", "uca", "uha", "war", "wim", "xar", "xp3", "yz1", "zip", "zipx", "zoo", "zpaq", "zz", "ecc", "ecsbx", "par", "par2", "rev"},
		"Other":       {},
		"Directories": {},
	}

	extensionToDirMap := invertMap(dirs)

	return &filetypeSort{dirsToExtension: dirs, extensionToDir: extensionToDirMap}
}

// invertMap function takes an argument myMap of type map[string][]string, returning a new
// map where the keys correspond to the individual string values in each slice of myMap.
func invertMap(myMap map[string][]string) map[string]string {
	invertedMap := make(map[string]string)

	for k, val := range myMap {

		for _, v := range val {
			invertedMap[v] = k
		}
	}
	return invertedMap
}

// createScaffolding method will create the correct directory structure for the
// filetypeSort. The list of directories are taken from the keys in fs.dirsToExtension
// It will first check if the directories already exist, if they do not, then it will
// proceed to create them. Any failure is returned as an error.
//
// createScaffolding method will run in whatever the current working directory is.
func (fs *filetypeSort) createScaffolding() error {
	directories := make([]string, 0, len(fs.dirsToExtension))
	for k := range fs.dirsToExtension {
		directories = append(directories, k)
	}

	for _, v := range directories {

	}

	return nil
}

func (fs *filetypeSort) sort() {
	err := fs.createScaffolding()
	if err != nil {
		fmt.Println(err)
	}

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
