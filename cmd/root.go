package cmd

import (
	"github.com/spf13/cobra"
	core "glab.tagtic.cn/ad_gains/kitty/core/module"
)

var (
	// Used for flags.
	cfgFile string

	coreModule *core.Module

	rootCmd = &cobra.Command{
		Use:   "kitty",
		Short: "A Pragmatic and Opinionated Go Application",
		Long:  `Kitty is a starting point to write 12-factor Go Applications.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			coreModule = core.New(cfgFile)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .config/kitty.yaml)")
}
