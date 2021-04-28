package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

type signParam struct {
	admin        bool
	id           uint64
	suuid        string
	openid       string
	channel      string
	versionCode  string
	mobile       string
	packageName  string
	thirdPartyId string
	ttl          time.Duration
	issuer       string
}

var s signParam

func init() {
	signCmd.Flags().BoolVar(&s.admin, "admin", false, "the admin token")
	signCmd.Flags().Uint64Var(&s.id, "id", 1, "the user id in the token")
	signCmd.Flags().StringVar(&s.suuid, "suuid", "", "the suuid in the token")
	signCmd.Flags().StringVar(&s.openid, "openid", "", "the wechat openid in the token")
	signCmd.Flags().StringVar(&s.channel, "channel", "", "the channel in the token")
	signCmd.Flags().StringVar(&s.versionCode, "versionCode", "", "the channel in the token")
	signCmd.Flags().StringVar(&s.mobile, "mobile", "", "the phone number in the token")
	signCmd.Flags().StringVar(&s.packageName, "packageName", "com.donews.www", "the package name of the token")
	signCmd.Flags().StringVar(&s.thirdPartyId, "thirdPartyId", "1", "the third party id of the token")
	signCmd.Flags().DurationVar(&s.ttl, "ttl", 12*30*24*time.Hour, "the ttl in the token")
	signCmd.Flags().StringVar(&s.issuer, "issuer", "signCmd", "the issuer in the token")
	rootCmd.AddCommand(signCmd)
}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign a jwt token",
	Long:  `Sign a valid jwt token for further use`,
	Run: func(cmd *cobra.Command, args []string) {
		claim := kittyjwt.NewClaim(
			s.id,
			s.issuer,
			s.suuid,
			s.channel,
			s.versionCode,
			s.openid,
			s.mobile,
			s.packageName,
			s.thirdPartyId,
			s.ttl,
		)
		if s.admin {
			claim = kittyjwt.NewAdminClaim(s.issuer, s.ttl)
		}
		tokenString, err := sign(claim)
		if err != nil {
			er(err)
			os.Exit(1)
		}
		fmt.Println(tokenString)
	},
}

func sign(claim *kittyjwt.Claim) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim,
	)
	token.Header["kid"] = coreModule.Conf.String("global.security.kid")
	tokenString, err := token.SignedString([]byte(coreModule.Conf.String("global.security.key")))
	return tokenString, err
}
