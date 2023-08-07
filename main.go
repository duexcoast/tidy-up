/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/duexcoast/tidy-up/cmd"
	"github.com/duexcoast/tidy-up/pkg/logger"
)

func init() {
	l := logger.Get()
	l.Debug().Msg("Logger initialized")
}

func main() {
	l := logger.Get()
	l.Debug().Msg("Program starting")
	cmd.Execute()
}
