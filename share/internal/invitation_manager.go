//go:generate mockery --name=RelationRepository

package internal

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	"gorm.io/gorm"
)

type RelationRepository interface {
	QueryRelations(ctx context.Context, condition entity.Relation) ([]entity.Relation, error)
	AddRelations(
		ctx context.Context,
		candidate *entity.Relation,
	) error
	UpdateRelations(
		ctx context.Context,
		apprentice *entity.User,
		existingRelationCallback func(relations []entity.Relation) error,
	) error
}

type OrientationEvent struct {
	Id          int    `yaml:"id"`
	Type        string `yaml:"type"`
	ChineseName string `yaml:"chineseName"`
}

type ReceivedEvent struct {
	Id   int    `yaml:"id"`
	Type string `yaml:"type"`
}

type ShareConfig struct {
	OrientationEvents []OrientationEvent `yaml:"orientation_events"`
	Url               string             `yaml:"url"`
	Reward            struct {
		Level1 int `yaml:"level1"`
		Level2 int `yaml:"level2"`
	} `yaml:"reward"`
	TaskId string `yaml:"task_id"`
}

func (s *ShareConfig) reward(depth int) int {
	if depth == 1 {
		return s.Reward.Level1
	}
	if depth == 2 {
		return s.Reward.Level2
	}
	return 0
}

type InvitationManager struct {
	conf        *ShareConfig
	rr          RelationRepository
	tokenizer   EncodeDecoder
	xtaskClient XTaskRequester
	logger      log.Logger
}

type EncodeDecoder interface {
	Decode(encoded string) (uint, error)
	Encode(id uint) (string, error)
}

type XTaskRequester interface {
	Request(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error)
}

type RelationWithRewardAmount struct {
	*entity.Relation
	Amount int
}

func user(id uint) entity.User {
	return entity.User{
		Model: gorm.Model{
			ID: id,
		},
	}
}

func (im *InvitationManager) AddToken(ctx context.Context, userId uint64, token string) error {
	masterId, err := im.tokenizer.Decode(token)
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}

	master := user(masterId)
	apprentice := user(uint(userId))
	steps := getSteps(im.conf.OrientationEvents)
	relation := entity.NewRelation(&apprentice, &master, steps)

	return im.rr.AddRelations(ctx, relation)
}

func (im *InvitationManager) ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error {

	apprentice := user(uint(apprenticeId))

	return im.rr.UpdateRelations(ctx, &apprentice, func(relations []entity.Relation) error {
		for _, rel := range relations {
			if rel.MasterID == uint(masterId) {

				if err := rel.ClaimReward(); err != nil {
					return err
				} else {
					resp, err := im.xtaskClient.Request(ctx, &XTaskRequest{
						ScoreDesc:  "邀请好友获得奖励",
						ScoreValue: im.conf.reward(rel.Depth),
						TaskId:     im.conf.TaskId,
						UniqueId:   xid.New().String(),
					})
					if err != nil {
						return errors.Wrap(err, "xtask request failed")
					}
					if resp.Code != 0 {
						return kerr.CustomErr(uint32(resp.Code), ErrFailedXtaskRequest, resp.Msg)
					}
					return nil
				}
			}
		}
		return ErrNoRewardAvailable
	})
}

func (im *InvitationManager) CompleteStep(ctx context.Context, apprenticeId uint64, event ReceivedEvent) error {

	if !in(event, im.conf.OrientationEvents) {
		level.Info(im.logger).Log("msg", fmt.Sprintf("invalid event %+v, want %+v", event, im.conf.OrientationEvents))
		return nil
	}

	apprentice := user(uint(apprenticeId))

	return im.rr.UpdateRelations(ctx, &apprentice, func(relations []entity.Relation) error {
		for i := range relations {
			relations[i].CompleteStep(entity.OrientationStep{EventType: event.Type, EventId: event.Id})
		}
		return nil
	})
}

func (im *InvitationManager) ListApprentices(ctx context.Context, masterId uint64, depth int) ([]RelationWithRewardAmount, error) {
	var out []RelationWithRewardAmount
	rels, err := im.rr.QueryRelations(ctx, entity.Relation{
		MasterID: uint(masterId),
		Depth:    depth,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error querying relations")
	}
	amount := im.conf.reward(depth)
	for i := range rels {
		out = append(out, RelationWithRewardAmount{
			Relation: &rels[i],
			Amount:   amount,
		})
	}
	return out, nil
}

func (im *InvitationManager) GetToken(_ context.Context, id uint) string {
	value, err := im.tokenizer.Encode(id)
	_ = level.Warn(im.logger).Log("err", err)
	return value
}

func (im *InvitationManager) GetUrl(ctx context.Context, claim *jwt2.Claim) string {
	args := url.Values{}
	args.Add("user_id", strconv.FormatUint(claim.UserId, 10))
	args.Add("channel", claim.Channel)
	args.Add("version_code", claim.VersionCode)
	args.Add("package_name", claim.PackageName)
	args.Add("invite_code", im.GetToken(ctx, uint(claim.UserId)))
	argsStr := args.Encode()

	return fmt.Sprintf(im.conf.Url, argsStr)
}

func getSteps(names []OrientationEvent) []entity.OrientationStep {
	var steps []entity.OrientationStep
	for _, s := range names {
		steps = append(steps, entity.OrientationStep{EventId: s.Id, EventType: s.Type})
	}
	return steps
}

func in(event ReceivedEvent, events []OrientationEvent) bool {
	for _, e := range events {
		if e.Id == event.Id && e.Type == event.Type {
			return true
		}
	}
	return false
}
