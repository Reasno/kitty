package service

import (
	"github.com/go-kit/kit/log"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func ProvideService(logger log.Logger, repo Repository, dmpServer pb.DmpServer) Service {
	return &service{logger: logger, repo: repo, dmpServer: dmpServer}
}
