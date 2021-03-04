package handlers

import (
	"context"

	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.DmpServer {
	return dmpService{}
}

type dmpService struct{}

func (s dmpService) UserMore(ctx context.Context, in *pb.DmpReq) (*pb.DmpResp, error) {
	var resp pb.DmpResp
	return &resp, nil
}
