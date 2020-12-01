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
			EventName:     "hello",
			StepCompleted: false,
		},
		{
			EventName:     "world",
			StepCompleted: false,
		},
	})
	rel.CompleteStep(OrientationStep{
		EventName: "hello",
	})
	rel.CompleteStep(OrientationStep{
		EventName: "world",
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
			EventName:     "hello",
			StepCompleted: false,
		},
	})
	rel.CompleteStep(OrientationStep{
		EventName: "hello",
	})
	assert.False(t, rel.RewardClaimed)
	err := rel.ClaimReward()
	assert.NoError(t, err)
	assert.True(t, rel.RewardClaimed)
}
