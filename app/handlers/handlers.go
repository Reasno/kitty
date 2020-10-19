package handlers

import (
	"context"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/log"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.AppServer {
	return injectAppServer()
}

type appService struct {
	log log.Logger
}

func (s appService) Create(ctx context.Context, in *pb.UserRequest) (*pb.GenericReply, error) {
	s.log.Log("fff", "bbb")
	var resp pb.GenericReply
	return &resp, nil
}

func (s appService) Code(ctx context.Context, in *pb.EmptyRequest) (*pb.GenericReply, error) {
	s.log.Log("foo", "bar")
	var resp pb.GenericReply
	return &resp, nil
}
