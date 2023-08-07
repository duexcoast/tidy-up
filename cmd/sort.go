/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/duexcoast/tidy-up/pkg/logger"
	"github.com/duexcoast/tidy-up/pkg/tidy"
	"github.com/joho/godotenv"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type sortCmdOptions struct {
	sortType string
	verbose  bool
	envFiles []string
}

// sortCmd represents the clean command

func init() {
	opts := &sortCmdOptions{}
	cmd := newSortCommand(opts)
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&opts.sortType, "type", "t", "filetypeSorter", "The sort type to be used")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")

	cmd.PersistentFlags().StringSliceVar(&opts.envFiles, "env-file", []string{}, "Env files to parse environment variables (looks for .env by default).")
}

func newSortCommand(opts *sortCmdOptions) *cobra.Command {
	return &cobra.Command{

		Use:     "sort <path> [--type <sort type>]",
		Aliases: []string{"s"},
		Short:   "This command will sort the specified directory.",
		Long:    ``,
		Args:    cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			l := logger.Get()
			err := godotenv.Load(opts.envFiles...)
			if err != nil {
				l.Error().Err(err).Msg("error loading env files.")
			}
			runSort(opts, args)
		},
	}
}

func runSort(opts *sortCmdOptions, args []string) {
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
}
