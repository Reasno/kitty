//go:generate mockery --name=InvitationManager
//go:generate mockery --name=UserRepository
package handlers

import (
	"context"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	shareevent "glab.tagtic.cn/ad_gains/kitty/share/event"
	"time"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	token "glab.tagtic.cn/ad_gains/kitty/pkg/invitecode"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/internal"
)

var ErrReenteringInviteCode = errors.New("不能重复填写邀请码")

type shareService struct {
	manager    InvitationManager
	ur         UserRepository
	dispatcher contract.Dispatcher
	tokenizer  internal.EncodeDecoder
}

type InvitationManager interface {
	AddToken(ctx context.Context, userId uint64, token string) error
	ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error
	CompleteStep(ctx context.Context, apprenticeId uint64, event internal.ReceivedEvent) error
	ListApprentices(ctx context.Context, masterId uint64, depth int) ([]internal.RelationWithRewardAmount, error)
	GetToken(ctx context.Context, id uint) string
	GetUrl(ctx context.Context, claim *kittyjwt.Claim) string
}

type UserRepository interface {
	UpdateCallback(ctx context.Context, id uint, f func(user *entity.User) error) (err error)
	Exists(ctx context.Context, id uint) bool
}

func (s shareService) AddInvitationCode(ctx context.Context, in *pb.ShareAddInvitationRequest) (*pb.ShareGenericReply, error) {
	var err error

	claim := kittyjwt.ClaimFromContext(ctx)

	inviterId, err := s.tokenizer.Decode(in.GetInviteCode())

	if err != nil || !s.ur.Exists(ctx, inviterId) {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}

	err = s.ur.UpdateCallback(ctx, uint(claim.UserId), func(user *entity.User) error {
		if user.InviteCode != "" {
			return ErrReenteringInviteCode
		}

		err := s.manager.AddToken(ctx, claim.UserId, in.InviteCode)
		if err != nil {
			return errors.Wrap(err, msg.InvalidInviteCode)
		}
		user.InviteCode = in.InviteCode
		return nil
	})

	if errors.Is(err, ErrReenteringInviteCode) {
		return nil, kerr.InvalidArgumentErr(err, msg.ReenteringCode)
	}
	if errors.Is(err, entity.ErrRelationArgument) {
		return nil, kerr.InvalidArgumentErr(err, msg.InvalidInviteTarget)
	}
	if errors.Is(err, entity.ErrRelationCircled) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorCircledInvitation)
	}
	if errors.Is(err, entity.ErrRelationExist) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorRelationAlreadyExists)
	}
	if errors.Is(err, token.ErrFailedToDecodeToken) {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}

	// 触发事件
	e := shareevent.InvitationCodeAdded{
		InviteeId:   claim.UserId,
		InviterId:   uint64(inviterId),
		PackageName: claim.PackageName,
		InviteCode:  in.GetInviteCode(),
		DateTime:    time.Now(),
		Channel:     claim.Channel,
	}

	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, e))

	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) ClaimReward(ctx context.Context, in *pb.ShareClaimRewardRequest) (*pb.ShareGenericReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
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
	claim := kittyjwt.ClaimFromContext(ctx)
	rels, err := s.manager.ListApprentices(ctx, claim.UserId, int(in.Depth))
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var (
		resp          pb.ShareListFriendReply
		countNotReady int32
	)
	resp.Data = new(pb.ShareListFriendData)
	for _, rel := range rels {
		item := &pb.ShareListFriendDataItem{
			Id:       uint64(rel.ApprenticeID),
			UserName: rel.Apprentice.UserName,
			HeadImg:  rel.Apprentice.HeadImg,
			Gender:   pb.Gender(rel.Apprentice.Gender),
			Coin:     int32(rel.Amount),
			Steps:    make(map[string]bool),
			CreateAt: rel.CreatedAt.Unix(),
		}
		item.ClaimStatus = status(&rel, &countNotReady)

		for _, step := range rel.OrientationSteps {
			item.Steps[step.ChineseName] = step.StepCompleted
		}
		resp.Data.Items = append(resp.Data.Items, item)
	}
	resp.Data.CountNotReady = countNotReady
	resp.Data.CountAll = int32(len(rels))
	return &resp, nil
}

func status(item *internal.RelationWithRewardAmount, countNotReady *int32) pb.ClaimStatus {
	if item.RewardClaimed {
		return pb.ClaimStatus_DONE
	}
	if item.OrientationCompleted {
		return pb.ClaimStatus_READY
	}
	*countNotReady++
	return pb.ClaimStatus_NOT_READY
}

func (s shareService) InviteByUrl(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataUrlReply, error) {
	url := s.manager.GetUrl(ctx, kittyjwt.ClaimFromContext(ctx))
	var resp = pb.ShareDataUrlReply{
		Code: 0,
		Data: &pb.ShareDataUrlReply_Url{
			Url: url,
		},
	}
	return &resp, nil
}

func (s shareService) InviteByToken(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataTokenReply, error) {
	id := uint(kittyjwt.ClaimFromContext(ctx).UserId)
	code := s.manager.GetToken(ctx, id)
	var resp = pb.ShareDataTokenReply{
		Code: 0,
		Data: &pb.ShareDataTokenReply_Code{
			Code: code,
		},
	}
	return &resp, nil
}

func (s shareService) PushSignEvent(ctx context.Context, in *pb.SignEvent) (*pb.ShareGenericReply, error) {
	var event internal.ReceivedEvent
	event.Id = int(in.Id)
	event.Type = "sign"

	err := s.manager.CompleteStep(ctx, in.UserId, event)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) PushTaskEvent(ctx context.Context, in *pb.TaskEvent) (*pb.ShareGenericReply, error) {

	var event internal.ReceivedEvent
	event.Id = int(in.Id)
	event.Type = "task"

	err := s.manager.CompleteStep(ctx, in.UserId, event)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}
