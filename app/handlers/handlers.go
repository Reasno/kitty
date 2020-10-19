package handlers

import (
	"context"
	pb "github.com/Reasno/kitty/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return &resp, status.Error(codes.Aborted, "test")
}
