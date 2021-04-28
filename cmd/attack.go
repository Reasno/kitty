package cmd

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	"gorm.io/gorm/clause"
	"os"
	"time"
)

func init() {
	rootCmd.AddCommand(attackCmd)
}

var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "Load testing",
	Long:  `Load testing target module with real users`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dialector, err := module.ProvideDialector(coreModule.Conf.Cut("app"))
		if err != nil {
			return err
		}
		gormConfig := module.ProvideGormConfig(coreModule.Logger, coreModule.Conf.Cut("app"))
		db, closer, err := module.ProvideGormDB(dialector, gormConfig, opentracing.NoopTracer{})
		if err != nil {
			return err
		}
		defer closer()
		repo := repository.NewUserRepo(db, nil, nil)

		users, err := repo.GetAll(
			cmd.Context(),
			clause.Where{Exprs: []clause.Expression{clause.Eq{
				Column: "package_name",
				Value:  "com.skin.v10mogul",
			}}},
			clause.Limit{
				Limit:  500,
				Offset: 0,
			},
			clause.OrderBy{
				Expression: clause.Expr{
					SQL:                "updated_at desc",
					Vars:               nil,
					WithoutParentheses: false,
				},
			},
		)
		if err != nil {
			return err
		}
		var targets []vegeta.Target

		rate := vegeta.Rate{Freq: 150, Per: time.Second}
		duration := 5 * time.Second
		for i := range users {
			token, _ := sign(kittyjwt.NewClaim(
				uint64(users[i].ID),
				"bench",
				users[i].CommonSUUID,
				users[i].Channel,
				users[i].VersionCode,
				users[i].WechatOpenId.String,
				users[i].Mobile.String,
				users[i].PackageName,
				users[i].ThirdPartyId,
				time.Hour,
			))
			targets = append(targets, vegeta.Target{
				Method: "POST",
				URL:    "https://commercial-products-b-dev.xg.tagtic.cn/v10mogul/addAction",
				Body: []byte(`{
  					"type": 6,
  					"source": 1,
  					"id": 351
				}`),
				Header: map[string][]string{
					"Authorization": {"bearer " + token},
					"Content-Type":  {"application/json; charset=utf-8"},
				},
			})
			targets = append(targets, vegeta.Target{
				Method: "POST",
				URL:    "https://commercial-products-b-dev.xg.tagtic.cn/v10mogul/panicbuy",
				Body: []byte(`{
					"id": 1,
					"skinId": 351
				}`),
				Header: map[string][]string{
					"Authorization": {"bearer " + token},
					"Content-Type":  {"application/json; charset=utf-8"},
				},
			})
		}
		fmt.Println("attacking")
		targeter := vegeta.NewStaticTargeter(targets...)
		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
			fmt.Println(string(res.Body))
			metrics.Add(res)
		}
		metrics.Close()

		vegeta.NewTextReporter(&metrics).Report(os.Stdout)
		return nil
	},
}
