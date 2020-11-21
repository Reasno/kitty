package repository

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"gorm.io/gorm"
)

var ErrRelationArgument = errors.New("错误的关系参数")
var ErrRelationExist = errors.New("关系已经存在")
var ErrRelationCircled = errors.New("关系中不能有环")

type RelationRepo struct {
	db *gorm.DB
}

func NewRelationRepo(db *gorm.DB) *RelationRepo {
	return &RelationRepo{db: db}
}

func (r *RelationRepo) QueryRelations(ctx context.Context, condition entity.Relation) ([]entity.Relation, error) {
	var relations []entity.Relation
	err := r.db.
		WithContext(ctx).
		Preload("Apprentice").
		Preload("Master").
		Preload("OrientationSteps").
		Where(&condition).
		Order("reward_claimed desc, orientation_completed, created_at desc").
		Find(&relations).Error

	return relations, err
}

func (r *RelationRepo) AddRelations(
	ctx context.Context,
	candidate *entity.Relation,
) error {
	var (
		descendants       []entity.Relation
		ancestor          entity.Relation
		secondaryAncestor entity.Relation
		grandMaster       *entity.User
		err               error
	)

	if candidate.MasterID == 0 {
		return ErrRelationArgument
	}
	if candidate.ApprenticeID == 0 {
		return ErrRelationArgument
	}
	if candidate.ApprenticeID == candidate.MasterID {
		return ErrRelationArgument
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		tx.WithContext(ctx).Where(&entity.Relation{
			MasterID:     candidate.MasterID,
			ApprenticeID: candidate.ApprenticeID,
		}).First(&ancestor)

		tx.WithContext(ctx).Preload("Apprentice").Preload("Master").Where(&entity.Relation{
			ApprenticeID: candidate.MasterID,
			Depth:        1,
		}).First(&secondaryAncestor)

		tx.WithContext(ctx).Select("apprentice_id", "depth").Where(&entity.Relation{
			MasterID: candidate.ApprenticeID,
		}).Find(&descendants)

		if ancestor.ID != 0 {
			return ErrRelationExist
		}

		newRelations := []entity.Relation{*candidate}

		if secondaryAncestor.ID != 0 {
			grandMaster = &secondaryAncestor.Master
			grandMaster.ID = secondaryAncestor.MasterID
			newRelations = append(newRelations, *entity.NewIndirectRelation(&candidate.Apprentice, grandMaster, candidate.OrientationSteps))
		}

		if circleDetected(&candidate.Master, grandMaster, descendants) {
			return ErrRelationCircled
		}

		for _, descendant := range descendants {
			if descendant.Depth == 2 {
				continue
			}
			apprentice := entity.User{Model: gorm.Model{ID: descendant.ApprenticeID}}
			newRelations = append(newRelations, *entity.NewIndirectRelation(&apprentice, &candidate.Master, candidate.OrientationSteps))
		}

		// save new relations
		err = tx.WithContext(ctx).Create(&newRelations).Error
		if err != nil {
			return errors.Wrap(err, "unable to create relations")
		}
		return nil
	})
}

func circleDetected(master, grandMaster *entity.User, descendant []entity.Relation) bool {
	if grandMaster != nil {
		return in(grandMaster, descendant) || in(master, descendant)
	}
	return in(master, descendant)
}

func in(user *entity.User, descendant []entity.Relation) bool {
	for _, v := range descendant {
		if user.ID == v.ApprenticeID {
			return true
		}
	}
	return false
}

func (r *RelationRepo) UpdateRelations(
	ctx context.Context,
	apprentice *entity.User,
	existingRelationCallback func(relations []entity.Relation) error,
) error {
	var (
		err               error
		ancestor          entity.Relation
		secondaryAncestor entity.Relation
		updated           []entity.Relation
	)
	return r.db.Transaction(func(tx *gorm.DB) error {
		ptx := tx.WithContext(ctx).Preload("Apprentice").Preload("Master").Preload("OrientationSteps")
		ptx.Where(&entity.Relation{
			ApprenticeID: apprentice.ID,
		}).Find(&ancestor)

		if ancestor.ID != 0 {
			updated = []entity.Relation{ancestor}
		}

		ptx.Where(&entity.Relation{
			ApprenticeID: ancestor.MasterID,
			Depth:        1,
		}).Find(&secondaryAncestor)

		if secondaryAncestor.ID != 0 {
			updated = append(updated, secondaryAncestor)
		}

		err = existingRelationCallback(updated)
		if err != nil {
			return errors.Wrap(err, "existingRelationCallback")
		}

		// save new relations
		err = tx.WithContext(ctx).Save(&updated).Error
		if err != nil {
			return errors.Wrap(err, "unable to save relations")
		}
		return nil
	})

}
