package handlers

import (
	"github.com/go-kit/kit/log"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func NewAppService(conf contract.ConfigReader, log log.Logger, ur UserRepository, cr CodeRepository, fr FileRepository, sender contract.SmsSender, wechat wechat.Wechater) appService {
	return appService{conf: conf, logger: log, ur: ur, cr: cr, sender: sender, wechat: wechat, fr: fr}
}

type ServerMiddleware func(server pb.AppServer) pb.AppServer

func NewMonitorMiddleware(userBus UserBus, eventBus EventBus) ServerMiddleware {
	return func(server pb.AppServer) pb.AppServer {
		return &MonitoredAppService{userBus: userBus, eventBus: eventBus, AppServer: server}
	}
}

func NewInputEnrichMiddleware() ServerMiddleware {
	return func(server pb.AppServer) pb.AppServer {
		return &InputEnrichedAppService{AppServer: server}
	}
}

func Chain(outer ServerMiddleware, others ...ServerMiddleware) ServerMiddleware {
	return func(next pb.AppServer) pb.AppServer {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func ProvideAppServer(userBus UserBus, eventBus EventBus, service appService) pb.AppServer {
	return Chain(
		NewInputEnrichMiddleware(),
		NewMonitorMiddleware(userBus, eventBus),
	)(service)
}
