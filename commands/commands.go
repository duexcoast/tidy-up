package cmd

//
// import (
// 	"fmt"
// 	"path/filepath"
//
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )
//
// var (
// 	verbose bool
//
// 	config   string // config file location
// 	showVers bool   // whether to print the version info
//
// 	// populated by linker
// 	version string
// 	commit  string
//
// 	TidyCmd = &cobra.Command{
// 		Use:           "tidy",
// 		Short:         "tidy - clean up your directories",
// 		Long:          ``,
// 		SilenceErrors: true,
// 		SilenceUsage:  true,
//
// 		// parse the config if one is provided, or use the defaults
// 		PersistentPreRunE: readConfig,
//
// 		// print version or help, or continue, depending on flag settings
// 		PreRunE: preFlight,
//
// 		// either run tidy as a server, or run it as a CLI depending on what flags
// 		// are provided
// 		RunE: startTidy,
// 	}
// )
//
// func readConfig(cmd *cobra.Command, args []string) error {
// 	// if --config is passed, attempt to read the config file
// 	if config != "" {
// 		filename := filepath.Base(config)
// 		viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
// 		viper.AddConfigPath(filepath.Dir(config))
//
// 		err := viper.ReadInConfig()
// 		if err != nil {
// 			return fmt.Errorf("Failed to read config file - %s", err)
// 		}
// 	}
// 	return nil
// }
//
// func preFlight(cmd *cobra.Command, args []string) error {
// 	// if --version is passed, print the version info
// 	if showVers {
// 		fmt.Printf("tidy %s (%s)\n", version, commit)
// 		return fmt.Errorf("")
//
// 	}
// 	return nil
// }
//
// func startTidy(cmd *cobra.Command, args []string) error {
// 	return nil
// }
//
// func init() {
// 	// set config defaults
// 	logLevel := "INFO"
//
// 	// cli flags
// 	TidyCmd.PersistentFlags().String("log-level", logLevel, "Output level of logs (DEBUG, INFO, WARN, ERROR, FATAL)")
//
// 	// bind config to cli flags
// 	viper.BindPFlag("log-level", TidyCmd.PersistentFlags().Lookup("log-level"))
//
// 	// cli-only flags
// 	TidyCmd.Flags().StringVarP(&config, "config", "c", "", "Path to config file (with extension)")
// 	TidyCmd.Flags().BoolVarP(&showVers, "version", "v", false, "Display the current version of this CLI")
//
// 	// add commands
// 	TidyCmd.AddCommand(sortCmd)
// 	TidyCmd.AddCommand(undoCmd)
//
// }
