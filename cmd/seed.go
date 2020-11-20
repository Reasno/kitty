package cmd

import (
	"fmt"

	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed the database",
	Long:  `use the provided seeds to bootstrap fake data in database`,
	Run: func(cmd *cobra.Command, args []string) {
		initModules()
		defer shutdownModules()

		for _, f := range moduleContainer.SeedProvider {
			if err := f(); err != nil {
				level.Error(logger).Log("err", fmt.Sprintf("unable to seed %s", err.Error()))
				return
			}
		}

		level.Info(logger).Log("msg", "seeding successfully completed")
	},
}
