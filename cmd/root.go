package cmd

import (
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile  string
	cfgCheck bool

	logger log.Logger
	conf   contract.ConfigReader

	rootCmd = &cobra.Command{
		Use:   "kitty",
		Short: "A Pragmatic and Opinionated Go Application",
		Long:  `Kitty is a starting point to write 12-factor Go Applications.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initConfig(cmd, args); err != nil {
				return err
			}
			if err := initLogger(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/kitty.yaml)")
	rootCmd.PersistentFlags().BoolVar(&cfgCheck, "check", false, "check config file integrity during boot up")
}
