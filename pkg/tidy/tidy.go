package tidy

import (
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

func createScaffolding(cfg TidyConfig) error {
	switch cfg.SortType {
	case "filetype":
		err := scaffold()
	}
}

type TidyConfig struct {	
	SortType string
}

// put this config inside of a function for now, 
func run() {
	cfg := struct {
			SortType string `default:"filetype"` 
	}

	// create scaffolding
	createScaffolding(cfg)
}
