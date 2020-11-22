package cmd

import (
	"github.com/go-kit/kit/log/level"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"

	app "glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/pkg/container"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/rule/client"
	rule "glab.tagtic.cn/ad_gains/kitty/rule/module"
	share "glab.tagtic.cn/ad_gains/kitty/share/module"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	ruleModule := rule.New(coreModule.Make("rule"))
	engine, _ := client.NewRuleEngine(client.WithRepository(ruleModule.GetRepository()))

	moduleContainer.Register(ruleModule)
	moduleContainer.Register(app.New(coreModule.MakeWithEngine("app", engine)))
	moduleContainer.Register(share.New(coreModule.MakeWithEngine("app", engine)))
	moduleContainer.Register(ots3.New(coreModule.Make("s3")))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Doc))
	moduleContainer.Register(container.HttpFunc(kittyhttp.HealthCheck))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Metrics))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Debug))

}

func shutdownModules() {
	for _, f := range moduleContainer.CloserProviders {
		f()
	}
}

func warn(msg string) {
	_ = level.Warn(coreModule.Logger).Log("msg", msg)
}

func er(err error) {
	_ = level.Error(coreModule.Logger).Log("err", err)
}

func debug(msg string) {
	_ = level.Debug(coreModule.Logger).Log("msg", msg)
}

func info(msg string) {
	_ = level.Info(coreModule.Logger).Log("msg", msg)
}

func conf() contract.ConfigReader {
	return coreModule.StaticConf
}
