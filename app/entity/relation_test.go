package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRelation_CompleteStep(t *testing.T) {
	var u1 = User{Model: gorm.Model{ID: 1}}
	var u2 = User{Model: gorm.Model{ID: 2}}
	rel := NewRelation(&u1, &u2, []OrientationStep{
		{
			EventId:       1,
			StepCompleted: false,
		},
		{
			EventId:       1,
			StepCompleted: false,
		},
	})
	rel.CompleteStep(OrientationStep{
		EventId: 1,
	})
	rel.CompleteStep(OrientationStep{
		EventId: 1,
	})
	if !rel.OrientationCompleted {
		t.Fatal("orientation should be completed by now")
	}
}

func TestRelation_ClaimReward(t *testing.T) {
	var u1 = User{Model: gorm.Model{ID: 1}}
	var u2 = User{Model: gorm.Model{ID: 2}}
	rel := NewRelation(&u1, &u2, []OrientationStep{
		{
			EventId:       1,
			StepCompleted: false,
		},
	})
	rel.CompleteStep(OrientationStep{
		EventId: 1,
	})
	assert.False(t, rel.RewardClaimed)
	err := rel.ClaimReward()
	assert.NoError(t, err)
	assert.True(t, rel.RewardClaimed)
}
