//go:generate mockery --name=RelationRepository

package internal

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
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

type InvitationManager struct {
	conf        contract.ConfigReader
	rr          RelationRepository
	tokenizer   EncodeDecoder
	xtaskClient XTaskRequester
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
	steps := getSteps(im.conf.Strings("orientation_events"))
	relation := entity.NewRelation(&apprentice, &master, steps)

	return im.rr.AddRelations(ctx, relation)
}

func (im *InvitationManager) ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error {

	apprentice := user(uint(apprenticeId))

	return im.rr.UpdateRelations(ctx, &apprentice, func(relations []entity.Relation) error {
		for _, rel := range relations {
			if rel.MasterID == uint(masterId) {
				// TODO： 真正发奖
				if err := rel.ClaimReward(); err != nil {
					return err
				} else {
					amount := im.conf.Int("reward.level" + strconv.Itoa(rel.Depth))
					resp, err := im.xtaskClient.Request(ctx, &XTaskRequest{
						ScoreDesc:  "邀请好友获得奖励",
						ScoreValue: amount,
						TaskId:     "666666",
						UniqueId:   xid.New().String(),
					})
					if err != nil {
						return err
					}
					if resp.Code != 0 {
						return errors.New(resp.Msg)
					}
					return nil
				}
			}
		}
		return ErrNoRewardAvailable
	})
}

func (im *InvitationManager) AdvanceStep(ctx context.Context, apprenticeId uint64, eventName string) error {

	apprentice := user(uint(apprenticeId))

	return im.rr.UpdateRelations(ctx, &apprentice, func(relations []entity.Relation) error {
		for _, rel := range relations {
			rel.CompleteStep(entity.OrientationStep{Name: eventName})
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
	amount := im.conf.Int("reward.level" + strconv.Itoa(depth))
	for _, rel := range rels {
		out = append(out, RelationWithRewardAmount{
			Relation: &rel,
			Amount:   amount,
		})
	}
	return out, nil
}

func (im *InvitationManager) GetToken(ctx context.Context) string {
	claim := jwt2.GetClaim(ctx)
	value, _ := im.tokenizer.Encode(uint(claim.UserId))
	return value
}

func (im *InvitationManager) GetUrl(ctx context.Context) string {
	claim := jwt2.GetClaim(ctx)
	args := url.Values{}
	args.Add("user_id", strconv.FormatUint(claim.UserId, 10))
	args.Add("channel", claim.Channel)
	args.Add("version_code", claim.VersionCode)
	args.Add("package_name", claim.PackageName)
	args.Add("invite_code", im.GetToken(ctx))
	argsStr := args.Encode()

	format := im.conf.String("url")
	return fmt.Sprintf(format, argsStr)
}

func getSteps(names []string) []entity.OrientationStep {
	var steps []entity.OrientationStep
	for _, s := range names {
		steps = append(steps, entity.OrientationStep{Name: s})
	}
	return steps
}
