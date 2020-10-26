package cmd

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/config"
	"github.com/Reasno/kitty/pkg/container"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
	kitty_log "github.com/Reasno/kitty/pkg/log"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	cfgCheck bool
	serviceContainer container.ModuleContainer
	logger log.Logger

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

func initConfig(_ *cobra.Command, _ []string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name "kitty" (without extension).
		viper.AddConfigPath("./config/")
		viper.AddConfigPath(home)
		viper.SetConfigName("kitty")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if cfgCheck {
		var validate config.AutoGenerated
		if err := viper.UnmarshalExact(&validate); err != nil {
			return err
		}
	}
	return nil
}

func initLogger(cmd *cobra.Command, _ []string) error {
	logger = kitty_log.NewLogger(viper.GetString("global.env"))
	logger = log.With(logger, "subcommand", cmd.Use)
	logger = level.Info(logger)
	logger.Log("config", viper.ConfigFileUsed())

	return nil
}

func initServiceContainer(_ *cobra.Command, _ []string) error {
	serviceContainer = container.NewModuleContainer()
	serviceContainer.Register(handlers.New())
	serviceContainer.Register(container.HttpFunc(kittyhttp.Doc))
	serviceContainer.Register(container.HttpFunc(kittyhttp.HealthCheck))
	serviceContainer.Register(container.HttpFunc(kittyhttp.Metrics))
	return nil
}
