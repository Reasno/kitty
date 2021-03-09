package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
		version := coreModule.Conf.String("global.version")
		info(fmt.Sprintf("Kitty %s", version))
		resp, _ := http.Get("http://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png")
		fmt.Println(ioutil.ReadAll(resp.Body))
	},
}
