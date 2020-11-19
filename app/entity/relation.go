package entity

import (
	"errors"

	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	"gorm.io/gorm"
)

var ErrRewardClaimed = errors.New(msg.RewardClaimed)
var ErrOrientationHasNotBeenCompleted = errors.New(msg.OrientationHasNotBeenCompleted)

type Relation struct {
	gorm.Model
	MasterID             uint `gorm:"index"`
	ApprenticeID         uint `gorm:"index"`
	Master               User
	Apprentice           User
	Depth                int
	OrientationCompleted bool
	OrientationSteps     []OrientationStep
	RewardClaimed        bool
}

func NewRelation(apprentice *User, master *User, steps []OrientationStep) *Relation {
	return &Relation{
		MasterID:             master.ID,
		ApprenticeID:         apprentice.ID,
		Master:               *master,
		Apprentice:           *apprentice,
		Depth:                1,
		OrientationCompleted: len(steps) == 0,
		OrientationSteps:     steps,
		RewardClaimed:        false,
	}
}

func (r *Relation) CompleteStep(step OrientationStep) {
	var orientationCompleted = true
	for n := range r.OrientationSteps {
		// update step status
		if r.OrientationSteps[n].Name == step.Name {
			r.OrientationSteps[n].StepCompleted = true
		}
		// update orientationCompleted flag
		// 只要有一步没有完成，总的初始任务就没有完成
		if !r.OrientationSteps[n].StepCompleted {
			orientationCompleted = false
		}
	}
	r.OrientationCompleted = orientationCompleted
}

func (r *Relation) ClaimReward() error {
	if r.RewardClaimed {
		return ErrRewardClaimed
	}
	if !r.OrientationCompleted {
		return ErrOrientationHasNotBeenCompleted
	}
	r.RewardClaimed = true
	return nil
}

type OrientationStep struct {
	gorm.Model
	RelationID    uint `gorm:"index"`
	Name          string
	StepCompleted bool
}
