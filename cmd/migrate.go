package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var force bool
var rollbackId string

func init() {
	migrateCommand.Flags().StringVarP(&rollbackId, "rollback", "r", "", "rollback to the given migration id")
	migrateCommand.Flag("rollback").NoOptDefVal = "-1"
	migrateCommand.Flags().BoolVarP(&force, "force", "f", false, "migrations and rollback in production requires force flag to be set")
	rootCmd.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate gorm tables",
	Long:  `Run all gorm table migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		initModules()
		defer shutdownModules()

		env := strings.ToLower(viper.GetString("global.env"))
		if env == "prod" || env == "production" {
			level.Error(logger).Log("err", "migrations and rollback in production requires force flag to be set")
			os.Exit(1)
			return
		}

		if rollbackId != "" {
			for _, f := range moduleContainer.MigrationProvider {
				if err := f.Rollback(rollbackId); err != nil {
					level.Error(logger).Log("err", fmt.Sprintf("Unable to rollback: %s", err.Error()))
					os.Exit(1)
				}
			}

			level.Info(logger).Log("msg", "rollback successfully completed")
			return
		}

		for _, f := range moduleContainer.MigrationProvider {
			if err := f.Migrate(); err != nil {
				level.Error(logger).Log("err", fmt.Sprintf("Unable to migrate: %s", err.Error()))
				os.Exit(1)
			}
		}

		level.Info(logger).Log("msg", "migration successfully completed")
	},
}
