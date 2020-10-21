package domain

import (
	"context"
	pb "github.com/Reasno/kitty/proto"
)

type UserService struct {
	repo Repository
}

func NewUserService() *UserService {
	return &UserService{}
}

type Repository interface {
	GetFromWechat(ctx context.Context, wechat string, device *Device) (*entity.User, error)
	GetFromMobile(ctx context.Context, mobile, code string, device *Device) (sequence string, err error)
}

func (us *UserService) login(ctx context.Context, wechat, mobile, code string, device *Device) (*pb.User, error) {
	var (
		u pb.User
		err error
	)
	if len(wechat) != 0 {
		u, err = us.repo.GetFromWechat(ctx, wechat, device)
	} else {
		u, err = us.repo.GetFromMobile(ctx, mobile, code, device)
	}
	if err != nil {
		return nil, err
	}

	if len(mobile) != 0 {

	}

	// Sign JWT

	// Construct DTO

}
