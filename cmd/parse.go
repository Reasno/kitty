package cmd

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

func init() {
	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a jwt token",
	Long:  `Parse a jwt token signed by kitty`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		claim := jwt2.Claim{}
		token, err := jwt.ParseWithClaims(args[0], &claim, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(conf.String("global.security.key")), nil
		})
		if err != nil {
			level.Error(logger).Log("err", err)
		}
		if !token.Valid {
			fmt.Println("token is NOT valid.")
		}
		fmt.Println("token is valid:")
		fmt.Printf("%+v\n", claim)
	},
}
