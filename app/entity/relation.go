package entity

import (
	"errors"
	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	"gorm.io/gorm"
)

var ErrRewardClaimed = errors.New(msg.RewardClaimed)
var ErrRelationCircled = errors.New("关系中不能有环")
var ErrRelationArgument = errors.New("错误的关系参数")
var ErrRelationExist = errors.New("关系已经存在")
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

func NewIndirectRelation(apprentice *User, master *User, steps []OrientationStep) *Relation {
	return &Relation{
		MasterID:             master.ID,
		ApprenticeID:         apprentice.ID,
		Master:               *master,
		Apprentice:           *apprentice,
		Depth:                2,
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

func (r *Relation) Validate() error {
	if r.MasterID == 0 {
		return ErrRelationArgument
	}
	if r.ApprenticeID == 0 {
		return ErrRelationArgument
	}
	if r.ApprenticeID == r.MasterID {
		return ErrRelationArgument
	}
	return nil
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

func (r *Relation) Connect(grandMaster *User, descendants []Relation) (addition []Relation, err error) {
	newRelations := []Relation{*r}
	if grandMaster != nil && grandMaster.ID != 0 {
		newRelations = append(newRelations, *NewIndirectRelation(&r.Apprentice, grandMaster, r.OrientationSteps))
	}

	// 检测四阶环
	if circleDetected(&r.Master, grandMaster, descendants) {
		return nil, ErrRelationCircled
	}

	for _, descendant := range descendants {
		if descendant.Depth == 2 {
			continue
		}
		apprentice := User{Model: gorm.Model{ID: descendant.ApprenticeID}}
		newRelations = append(newRelations, *NewIndirectRelation(&apprentice, &r.Master, r.OrientationSteps))
	}
	return newRelations, nil
}

func circleDetected(master, grandMaster *User, descendants []Relation) bool {
	if grandMaster != nil && grandMaster.ID != 0 {
		return in(grandMaster, descendants) || in(master, descendants)
	}
	return in(master, descendants)
}

func in(user *User, descendants []Relation) bool {
	for _, v := range descendants {
		if user.ID == v.ApprenticeID {
			return true
		}
	}
	return false
}

type OrientationStep struct {
	gorm.Model
	RelationID    uint `gorm:"index"`
	Name          string
	StepCompleted bool
}
