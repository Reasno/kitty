package cmd

import (
	"fmt"
	kitty_log "github.com/Reasno/kitty/pkg/log"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	logger log.Logger

	rootCmd = &cobra.Command{
		Use:   "kitty",
		Short: "A Pragmatic and Opinionated Go Application",
		Long:  `Kitty is a starting point to write 12-factor Go Applications.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
			logger = kitty_log.NewLogger(viper.GetString("app_env"))
			logger = log.With(logger, "subcommand", cmd.Use)
			_ = level.Debug(logger).Log("config", viper.ConfigFileUsed())
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/kitty.yaml)")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./config/")
		viper.AddConfigPath("../config/")
		viper.AddConfigPath(home)
		viper.SetConfigName("kitty")
	}

	viper.AutomaticEnv()
}
