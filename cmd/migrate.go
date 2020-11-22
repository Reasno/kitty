package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"os"
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

		env := config.ProvideEnv(coreModule.StaticConf)
		if env.IsProd() {
			er(fmt.Errorf("migrations and rollback in production requires force flag to be set"))
			os.Exit(1)
			return
		}

		if rollbackId != "" {
			for _, f := range moduleContainer.MigrationProvider {
				if err := f.Rollback(rollbackId); err != nil {
					er(fmt.Errorf("unable to rollback: %w", err))
					os.Exit(1)
				}
			}

			info("rollback successfully completed")
			return
		}

		for _, f := range moduleContainer.MigrationProvider {
			if err := f.Migrate(); err != nil {
				er(fmt.Errorf("unable to migrate: %w", err))
				os.Exit(1)
			}
		}

		info("migration successfully completed")
	},
}
