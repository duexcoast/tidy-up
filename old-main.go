package main

//
// import (
// 	"os"
// 	"os/user"
// 	"path/filepath"
//
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )
//
// func oldMain() {
// 	// initialize zerolog logger with pretty print for CLI
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
//
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatal().Msg(err.Error())
// 	}
//
// 	downloadsDir := filepath.Join(usr.HomeDir, "Downloads")
//
// 	err = os.Chdir(downloadsDir)
// 	if err != nil {
// 		log.Fatal().Msg(err.Error())
// 	}
//
// 	cwd, _ := os.Getwd()
// 	log.Info().Str("cwd", cwd).Send()
// }
