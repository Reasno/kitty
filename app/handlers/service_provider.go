package handlers

import (
	"github.com/go-kit/kit/log"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/sms"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
)

func NewAppService(conf contract.ConfigReader, log log.Logger, ur UserRepository, cr CodeRepository, er ExtraRepository, sender *sms.SenderFactory, wechat *wechat.Transport, uploader contract.Uploader, fr FileRepository) *appService {
	return &appService{conf: conf, log: log, ur: ur, cr: cr, er: er, sender: sender, wechat: wechat, uploader: uploader, fr: fr}
}

func NewMonitoredAppService(userBus UserBus, eventBus EventBus, appServer *appService) *MonitoredAppService {
	return &MonitoredAppService{userBus: userBus, eventBus: eventBus, AppServer: appServer}
}
