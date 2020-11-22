package cmd

import (
	"fmt"
	"strconv"

	"github.com/speps/go-hashids"
	"github.com/spf13/cobra"
)

var decode bool

func init() {
	inviteCmd.Flags().BoolVarP(&decode, "decode", "d", false, "instead of encode an id, decode a token to id")
	rootCmd.AddCommand(inviteCmd)
}

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Create an invite code",
	Long:  `Create an invite code that can be used in share module`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hd := hashids.NewData()
		hd.Salt = coreModule.StaticConf.String("global.salt")
		hd.MinLength = 10
		h, _ := hashids.NewWithData(hd)
		if !decode {
			i, _ := strconv.Atoi(args[0])
			result, _ := h.Encode([]int{i})
			fmt.Println(result)
			return
		}
		result, err := h.DecodeWithError(args[0])
		if err != nil {
			er(err)
			return
		}
		fmt.Println(result[0])
	},
}
