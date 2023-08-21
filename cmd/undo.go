/*
Copyright © 2023 DUEX COAST duexcoast@gmail.com
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

type undoCmdOptions struct {
	sortType string
	verbose  bool
	watch    bool
	envFiles []string
}

func init() {
	opts := &undoCmdOptions{}
	cmd := newUndoCommand(opts)
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&opts.sortType, "type", "t", "filetypeSorter", "The sort type to be used")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&opts.watch, "watch", "w", false, "Watch the filesystem and continuously sort as new files are created.")

	cmd.PersistentFlags().StringSliceVar(&opts.envFiles, "env-file", []string{}, "Env files to parse environment variables (looks for .env by default).")
}

func newUndoCommand(opts *undoCmdOptions) *cobra.Command {
	return &cobra.Command{

		Use:     "undo <path> [--type <sort type>]",
		Aliases: []string{"s"},
		Short:   "This command will unsort the specified directory.",
		Long:    ``,
		Args:    cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			l := logger.Get()
			err := godotenv.Load(opts.envFiles...)
			if err != nil {
				l.Error().Err(err).Msg("error loading env files.")
			}
			runUndo(opts, args)
		},
	}
}

func runUndo(opts *undoCmdOptions, args []string) {
	flags := &tidy.TidyFlags{Verbose: opts.verbose}
	Tidy, err := tidy.NewTidy(tidy.NewFiletypeSorter(), flags, afero.NewOsFs())
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

}
