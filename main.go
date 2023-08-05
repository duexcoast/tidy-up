/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/duexcoast/tidy-up/cmd"
	"github.com/duexcoast/tidy-up/pkg/tidy/logger"
)

func init() {
	l := logger.Get()
}

func main() {
	l := logger.Get()
	cmd.Execute()
}
