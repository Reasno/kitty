package service

import (
	"github.com/go-kit/kit/log"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type DmpServers struct {
	DmpProd pb.DmpServer
	DmpDev  pb.DmpServer
}

func ProvideService(logger log.Logger, repo Repository, dmp DmpServers) Service {
	return &service{logger: logger, repo: repo, dmpServerProd: dmp.DmpProd, dmpServerDev: dmp.DmpDev}
}
