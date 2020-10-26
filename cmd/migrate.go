package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate gorm tables",
	Long:  `Run all gorm table migrations.`,
	PreRunE: initServiceContainer,
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range serviceContainer.MigrationProvider {
			if err := f(); err != nil {
				level.Error(logger).Log("err", fmt.Sprintf("Unable to migrate: %s", err.Error()))
				os.Exit(1)
			}
		}
		level.Info(logger).Log("msg", "migration successfully completed")
	},
}
