package cmd

import (
	"fmt"
	"strconv"

	"github.com/speps/go-hashids"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(inviteCmd)
}

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Create an invite code",
	Long:  `Create an invite code that can be used in share module`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hd := hashids.NewData()
		hd.Salt = conf.String("app.salt")
		hd.MinLength = 10
		h, _ := hashids.NewWithData(hd)
		i, _ := strconv.Atoi(args[0])
		result, _ := h.Encode([]int{i})
		fmt.Println(result)
	},
}
