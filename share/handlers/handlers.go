package handlers

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/internal"
)

type shareService struct {
	manager InvitationManager
	ur      UserRepository
}

type InvitationManager interface {
	AddToken(ctx context.Context, userId uint64, token string) error
	ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error
	AdvanceStep(ctx context.Context, apprenticeId uint64, eventName string) error
	ListApprentices(ctx context.Context, masterId uint64, depth int) ([]internal.RelationWithRewardAmount, error)
	GetToken(ctx context.Context) string
	GetUrl(ctx context.Context) string
}

type UserRepository interface {
	UpdateCallback(ctx context.Context, id uint, f func(user *entity.User) error) (err error)
}

func (s shareService) AddInvitationCode(ctx context.Context, in *pb.ShareAddInvitationRequest) (*pb.ShareGenericReply, error) {

	claim := kittyjwt.GetClaim(ctx)

	err := s.ur.UpdateCallback(ctx, uint(claim.UserId), func(user *entity.User) error {
		if user.InviteCode == "" {
			user.InviteCode = in.InviteCode
		}
		err := s.manager.AddToken(ctx, claim.UserId, in.InviteCode)
		if err != nil {
			return errors.Wrap(err, msg.InvalidInviteCode)
		}
		return nil
	})

	if errors.Is(err, repository.ErrRelationArgument) {
		return nil, kerr.InvalidArgumentErr(err, msg.InvalidInviteTarget)
	}
	if errors.Is(err, repository.ErrRelationCircled) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorCircledInvitation)
	}
	if errors.Is(err, repository.ErrRelationExist) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorRelationAlreadyExists)
	}
	if errors.Is(err, internal.ErrFailedToDecodeToken) {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}

	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) ClaimReward(ctx context.Context, in *pb.ShareClaimRewardRequest) (*pb.ShareGenericReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	err := s.manager.ClaimReward(ctx, claim.UserId, in.ApprenticeId)
	if err != nil {
		if errors.Is(err, entity.ErrOrientationHasNotBeenCompleted) {
			return nil, kerr.FailedPreconditionErr(err, msg.OrientationHasNotBeenCompleted)
		}
		if errors.Is(err, entity.ErrRewardClaimed) {
			return nil, kerr.FailedPreconditionErr(err, msg.RewardClaimed)
		}
		if errors.Is(err, internal.ErrFailedXtaskRequest) {
			return nil, kerr.FailedPreconditionErr(err, msg.XTastAbnormally)
		}
		return nil, kerr.InternalErr(err, msg.NoRewardAvailable)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) ListFriend(ctx context.Context, in *pb.ShareListFriendRequest) (*pb.ShareListFriendReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	lists, err := s.manager.ListApprentices(ctx, claim.UserId, int(in.Depth))
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var resp pb.ShareListFriendReply
	resp.Data = new(pb.ShareListFriendData)
	for i := range lists {
		item := &pb.ShareListFriendDataItem{
			Id:          uint64(lists[i].ApprenticeID),
			UserName:    lists[i].Apprentice.UserName,
			HeadImg:     lists[i].Apprentice.HeadImg,
			Gender:      pb.Gender(lists[i].Apprentice.Gender),
			ClaimStatus: pb.ClaimStatus_NOT_READY,
			Coin:        int32(lists[i].Amount),
			Steps:       make(map[string]bool),
			CreateAt:    lists[i].CreatedAt.Unix(),
		}
		if lists[i].RewardClaimed {
			item.ClaimStatus = pb.ClaimStatus_DONE
		}
		if lists[i].OrientationCompleted {
			item.ClaimStatus = pb.ClaimStatus_NOT_READY
		}
		for _, step := range lists[i].OrientationSteps {
			item.Steps[step.Name] = step.StepCompleted
		}
		resp.Data.Items = append(resp.Data.Items, item)
	}
	return &resp, nil
}

func (s shareService) InviteByUrl(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	url := s.manager.GetUrl(ctx)
	var resp = pb.ShareDataReply{
		Code: 0,
		Data: map[string]string{"invite_url": url},
	}
	return &resp, nil
}

func (s shareService) InviteByToken(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	code := s.manager.GetToken(ctx)
	var resp = pb.ShareDataReply{
		Code: 0,
		Data: map[string]string{"invite_code": code},
	}
	return &resp, nil
}
