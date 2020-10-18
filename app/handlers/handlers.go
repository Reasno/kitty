package handlers

import (
	"context"
	"fmt"
	pb "github.com/Reasno/kitty/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.AppServer {
	return appService{}
}

type appService struct {
}

func (s appService) Create(ctx context.Context, in *pb.UserRequest) (*pb.GenericReply, error) {
	var resp pb.GenericReply
	return &resp, nil
}

func (s appService) Code(ctx context.Context, in *pb.EmptyRequest) (*pb.GenericReply, error) {
	var resp pb.GenericReply
	return &resp, fmt.Errorf("foo")
	//resp.Message = "fff"
	//return &resp, nil
}
