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
	Run: func(cmd *cobra.Command, args []string) {
		initModules()
		defer shutdownModules()

		for _, f := range moduleContainer.MigrationProvider {
			if err := f(); err != nil {
				level.Error(logger).Log("err", fmt.Sprintf("Unable to migrate: %s", err.Error()))
				os.Exit(1)
			}
		}
		level.Info(logger).Log("msg", "migration successfully completed")
	},
}
