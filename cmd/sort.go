/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/duexcoast/tidy-up/pkg/tidy"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	sortType string
)

// sortCmd represents the clean command
var sortCmd = &cobra.Command{
	Use:     "sort",
	Aliases: []string{"s"},
	Short:   "This command will sort the specified directory.",
	Long:    ``,
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		Tidy, err := tidy.NewTidy(tidy.NewFiletypeSorter(), afero.NewOsFs())
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		// arg is path of directory to be sorted
		if len(args) == 1 {
			err := Tidy.ChangeSortDir(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		}
		err = Tidy.Sort()
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	},
}

func init() {
	sortCmd.Flags().StringVarP(&sortType, "type", "t", "", "The sort type to be used")
	rootCmd.AddCommand(sortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sortCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sortCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
