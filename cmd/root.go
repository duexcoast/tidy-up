/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type rootCmdOptions struct {
	toggle bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Tidys up a messy directory in a configurable fashion.",
	Long: `Have your directories gotten out of control? Do you need help?
tidy-up is here to help you gain back control. Provide a chosen
directory, by default tidy-up will sort the directory into sub-
folders based on filetype.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// func addSubcommandPalettes() {
// 	rootCmd.AddCommand(cleanCmd)
// }

func init() {

	opts := &rootCmdOptions{}
	rootCmd.Flags().BoolVarP(&opts.toggle, "toggle", "t", false, "Help message for toggle")
	// rootCmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
}
