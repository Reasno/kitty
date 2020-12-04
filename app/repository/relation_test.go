package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
)

func TestRelationRepo_QueryRelations(t *testing.T) {
	setUp(t)
	defer tearDown()

	repo := RelationRepo{db}
	ctx := context.Background()

	data := []struct {
		apprentice entity.User
		master     entity.User
	}{
		{
			user(1),
			user(2),
		},
		{
			user(3),
			user(5),
		},
		{
			user(4),
			user(5),
		},
	}

	for _, d := range data {
		repo.AddRelations(ctx, entity.NewRelation(&d.apprentice, &d.master, nil))
	}

	cases := []struct {
		name   string
		master entity.User
		len    int
	}{
		{
			"0apprentice",
			user(14),
			0,
		},
		{
			"1apprentice",
			user(2),
			1,
		},
		{
			"2apprentice",
			user(5),
			2,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			relations, _ := repo.QueryRelations(ctx, entity.Relation{
				MasterID: cc.master.ID,
			})
			assert.Equal(t, cc.len, len(relations))
		})
	}
}

func TestRelationRepo_UpdateRelations(t *testing.T) {
	setUp(t)
	defer tearDown()

	repo := RelationRepo{db}
	ctx := context.Background()

	cases := []struct {
		name       string
		apprentice entity.User
		master     entity.User
		claimed    bool
	}{
		{
			"2to1",
			user(2),
			user(1),
			false,
		},
		{
			"4to1",
			user(4),
			user(1),
			true,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			repo.AddRelations(ctx, entity.NewRelation(&cc.apprentice, &cc.master, nil))
			repo.UpdateRelations(ctx, &cc.apprentice, func(relations []entity.Relation) error {
				for i := range relations {
					relations[i].RewardClaimed = cc.claimed
				}
				return nil
			})
			var relation entity.Relation
			db.First(&relation, "apprentice_id = ? and master_id = ?", cc.apprentice.ID, cc.master.ID)
			assert.Equal(t, cc.claimed, relation.RewardClaimed)
		})
	}
}

func TestRelationRepo_AddRelations(t *testing.T) {
	setUp(t)
	defer tearDown()

	repo := RelationRepo{db}
	ctx := context.Background()

	cases := []struct {
		name       string
		apprentice entity.User
		master     entity.User
		ok         bool
	}{
		{
			"0to1",
			user(0),
			user(1),
			false,
		},
		{
			"1to0",
			user(1),
			user(0),
			false,
		},
		{
			"1to1",
			user(1),
			user(1),
			false,
		},
		{
			"1to2",
			user(1),
			user(2),
			true,
		},
		{
			"1to2again",
			user(1),
			user(2),
			false,
		},
		{
			"2to3",
			user(2),
			user(3),
			true,
		},
		{
			"3to1",
			user(3),
			user(1),
			false,
		},
		{
			"2to1",
			user(2),
			user(1),
			false,
		},
		{
			"3to4",
			user(3),
			user(4),
			true,
		},
		{
			"4to2",
			user(4),
			user(2),
			false,
		},
		{
			"4to1",
			user(4),
			user(1),
			false,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			err := repo.AddRelations(ctx, entity.NewRelation(&cc.apprentice, &cc.master, nil))
			fmt.Println(err)
			assert.Equal(t, cc.ok, err == nil)
		})
	}
}

func TestRelationRepo_AddRelationsWithNumApprentice(t *testing.T) {
	setUp(t)
	defer tearDown()

	repo := RelationRepo{db}
	ctx := context.Background()

	cases := []struct {
		name          string
		apprentice    entity.User
		master        entity.User
		numApprentice int
	}{
		{
			"1to2",
			user(1),
			user(2),
			1,
		},
		{
			"2to3",
			user(2),
			user(3),
			2,
		},
		{
			"4to5",
			user(4),
			user(5),
			1,
		},
		{
			"3to4",
			user(3),
			user(4),
			2,
		},
		{
			"6to5",
			user(6),
			user(5),
			3,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			err := repo.AddRelations(ctx, entity.NewRelation(&cc.apprentice, &cc.master, nil))
			assert.NoError(t, err)
			r, err := repo.QueryRelations(ctx, entity.Relation{
				MasterID: cc.master.ID,
			})
			assert.NoError(t, err)
			assert.Equal(t, cc.numApprentice, len(r), fmt.Sprintf("%+v\n", r))
		})
	}
}

func TestRelationRepo_AddRelationsWithOrientation(t *testing.T) {
	setUp(t)
	defer tearDown()

	repo := RelationRepo{db}
	ctx := context.Background()

	cases := []struct {
		name       string
		apprentice entity.User
		master     entity.User
		ok         bool
	}{
		{
			"1to2",
			user(1),
			user(2),
			true,
		},
		{
			"3to2",
			user(3),
			user(2),
			true,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			repo.AddRelations(ctx, entity.NewRelation(&cc.apprentice, &cc.master, []entity.OrientationStep{
				{
					EventType: "foo",
					EventId:   1,
				},
				{
					EventType: "bar",
					EventId:   1,
				},
			}))
			var rel entity.Relation
			db.Preload("OrientationSteps").First(&rel, "master_id = ? and apprentice_id = ?", cc.master.ID, cc.apprentice.ID)

			fmt.Println(rel.OrientationSteps)
			assert.Equal(t, "foo", rel.OrientationSteps[0].EventType)
			assert.Equal(t, "bar", rel.OrientationSteps[1].EventType)
			repo.UpdateRelations(ctx, &cc.apprentice, func(relations []entity.Relation) error {
				for i := range relations {
					relations[i].CompleteStep(entity.OrientationStep{EventType: "foo", EventId: 1})
					relations[i].CompleteStep(entity.OrientationStep{EventType: "bar", EventId: 1})
				}
				return nil
			})
			db.Preload("OrientationSteps").First(&rel, "master_id = ? and apprentice_id = ?", cc.master.ID, cc.apprentice.ID)
			assert.Equal(t, true, rel.OrientationCompleted)
		})
	}
}
