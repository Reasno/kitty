package cmd

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Kitty",
	Long:  `All software has versions. This is Kitty's`,
	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("global.version")
		logger.Log("msg", fmt.Sprintf("Kitty %s", version))
	},
}
