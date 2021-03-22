package handlers

import (
	"context"
	"errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	"glab.tagtic.cn/ad_gains/kitty/pkg/invitecode"
	"testing"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/handlers/mocks"
	"glab.tagtic.cn/ad_gains/kitty/share/internal"
)

func ctx(id int) context.Context {
	return context.WithValue(context.Background(), jwt.JWTClaimsContextKey, &kjwt.Claim{
		UserId: uint64(id),
	})
}

func TestShareService_AddInvitationCode(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		service shareService
		ctx     context.Context
		req     pb.ShareAddInvitationRequest
		err     error
	}{
		{
			name: "正常请求",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return &m
				}(),
				ur: func() UserRepository {
					m := mocks.UserRepository{}
					m.On("UpdateCallback", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, f func(user *entity.User) error) error {
						user := &entity.User{}
						user.ID = id
						return f(user)
					})
					m.On("Exists", mock.Anything, mock.Anything).Return(true)
					return &m
				}(),
				dispatcher: func() contract.Dispatcher {
					return &event.Dispatcher{}
				}(),
				tokenizer: func() internal.EncodeDecoder {
					return invitecode.NewTokenizer("DonewsTeaParty")
				}(),
			},
			ctx: ctx(2),
			req: pb.ShareAddInvitationRequest{
				InviteCode: "0w427W2zBG",
			},
			err: nil,
		},
		{
			name: "异常请求",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(ErrReenteringInviteCode)
					return &m
				}(),
				ur: func() UserRepository {
					m := mocks.UserRepository{}
					m.On("UpdateCallback", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, f func(user *entity.User) error) error {
						user := &entity.User{}
						user.ID = id
						return f(user)
					})
					m.On("Exists", mock.Anything, mock.Anything).Return(true)
					return &m
				}(),
				dispatcher: func() contract.Dispatcher {
					return &event.Dispatcher{}
				}(),
				tokenizer: func() internal.EncodeDecoder {
					return invitecode.NewTokenizer("DonewsTeaParty")
				}(),
			},
			ctx: ctx(2),
			req: pb.ShareAddInvitationRequest{
				InviteCode: "0w427W2zBG",
			},
			err: ErrReenteringInviteCode,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			_, err := cc.service.AddInvitationCode(cc.ctx, &cc.req)
			assert.True(t, errors.Is(err, cc.err))
		})
	}
}

func TestShareService_ClaimReward(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		service shareService
		ctx     context.Context
		req     pb.ShareClaimRewardRequest
		err     error
	}{
		{
			name: "正常请求",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ClaimReward", mock.Anything, uint64(1), uint64(2)).Return(nil)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareClaimRewardRequest{
				ApprenticeId: 2,
			},
			err: nil,
		},
		{
			name: "异常请求",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ClaimReward", mock.Anything, mock.Anything, mock.Anything).Return(internal.ErrFailedXtaskRequest)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareClaimRewardRequest{
				ApprenticeId: 2,
			},
			err: internal.ErrFailedXtaskRequest,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			_, err := cc.service.ClaimReward(cc.ctx, &cc.req)
			assert.True(t, errors.Is(err, cc.err))
		})
	}
}

func TestShareService_ListFriend(t *testing.T) {
	//t.Parallel()
	cases := []struct {
		name          string
		service       shareService
		ctx           context.Context
		req           pb.ShareListFriendRequest
		status        pb.ClaimStatus
		countAll      int32
		countNotReady int32
		err           error
	}{
		{
			name: "正常请求ready",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ListApprentices", mock.Anything, uint64(1), 1).Return([]internal.RelationWithRewardAmount{
						{
							Relation: &entity.Relation{
								MasterID:             1,
								ApprenticeID:         2,
								Depth:                1,
								OrientationCompleted: true,
								OrientationSteps:     nil,
								RewardClaimed:        false,
							},
							Amount: 0,
						},
					}, nil)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareListFriendRequest{
				Depth: 1,
			},
			status:        pb.ClaimStatus_READY,
			countAll:      1,
			countNotReady: 0,
			err:           nil,
		},
		{
			name: "正常请求not ready",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ListApprentices", mock.Anything, uint64(1), 1).Return([]internal.RelationWithRewardAmount{
						{
							Relation: &entity.Relation{
								MasterID:             1,
								ApprenticeID:         2,
								Depth:                1,
								OrientationCompleted: false,
								OrientationSteps:     nil,
								RewardClaimed:        false,
							},
							Amount: 0,
						},
					}, nil)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareListFriendRequest{
				Depth: 1,
			},
			status:        pb.ClaimStatus_NOT_READY,
			countAll:      1,
			countNotReady: 1,
			err:           nil,
		},
		{
			name: "正常请求claimed",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ListApprentices", mock.Anything, uint64(1), 1).Return([]internal.RelationWithRewardAmount{
						{
							Relation: &entity.Relation{
								MasterID:             1,
								ApprenticeID:         2,
								Depth:                1,
								OrientationCompleted: true,
								OrientationSteps:     nil,
								RewardClaimed:        true,
							},
							Amount: 0,
						},
					}, nil)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareListFriendRequest{
				Depth: 1,
			},
			countAll:      1,
			countNotReady: 0,
			status:        pb.ClaimStatus_DONE,
			err:           nil,
		},
		{
			name: "异常请求",
			service: shareService{
				manager: func() InvitationManager {
					m := mocks.InvitationManager{}
					m.On("ListApprentices", mock.Anything, mock.Anything, mock.Anything).Return([]internal.RelationWithRewardAmount{
						{
							Relation: &entity.Relation{
								MasterID:             1,
								ApprenticeID:         2,
								Depth:                1,
								OrientationCompleted: true,
								OrientationSteps:     nil,
								RewardClaimed:        false,
							},
							Amount: 0,
						},
					}, internal.ErrFailedXtaskRequest)
					return &m
				}(),
			},
			ctx: ctx(1),
			req: pb.ShareListFriendRequest{
				Depth: 1,
			},
			err: internal.ErrFailedXtaskRequest,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			resp, err := cc.service.ListFriend(cc.ctx, &cc.req)
			assert.True(t, errors.Is(err, cc.err))
			if err == nil {
				assert.Equal(t, cc.status, resp.Data.Items[0].ClaimStatus)
				assert.Equal(t, cc.countAll, resp.Data.CountAll)
				assert.Equal(t, cc.countNotReady, resp.Data.CountNotReady)
			}
		})
	}
}
