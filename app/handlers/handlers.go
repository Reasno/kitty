package handlers

import (
	"context"
	"github.com/Reasno/kitty/app/msg"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.AppServer {
	return injectAppServer()
}

type appService struct {
	log log.Logger
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	if len(in.Wechat) == 0 && len(in.Mobile) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}
	if len(in.Mobile) != 0 && len(in.Code) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}

	var resp pb.UserLoginReply
	return &resp, nil
}
