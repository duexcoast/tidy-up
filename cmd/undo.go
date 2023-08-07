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

// undoCmd represents the undo command
var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last sort",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Tidy, err := tidy.NewTidy(tidy.NewFiletypeSorter(), afero.NewOsFs())
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}

		// arg is path of directory to be unsorted
		if len(args) == 1 {
			err := Tidy.ChangeSortDir(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		}
		err = Tidy.Undo()
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// undoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// undoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
