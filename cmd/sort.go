/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dirPath string
)

// sortCmd represents the clean command
var sortCmd = &cobra.Command{
	Use:     "sort",
	Aliases: []string{"c"},
	Short:   "This command will sort up the specified directory.",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sort called")
	},
}

func init() {
	sortCmd.Flags().StringVarP(&dirPath, "dir", "d", "", "The directory to be sorted.")
	rootCmd.AddCommand(sortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sortCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sortCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
