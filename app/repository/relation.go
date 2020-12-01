package repository

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"gorm.io/gorm"
)

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
		newRelations      []entity.Relation
		ancestor          entity.Relation
		secondaryAncestor entity.Relation
		err               error
	)

	if err := candidate.Validate(); err != nil {
		return errors.WithStack(err)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		tx.WithContext(ctx).Where(&entity.Relation{
			MasterID:     candidate.MasterID,
			ApprenticeID: candidate.ApprenticeID,
		}).First(&ancestor)

		tx.WithContext(ctx).Preload("Master").Where(&entity.Relation{
			ApprenticeID: candidate.MasterID,
			Depth:        1,
		}).First(&secondaryAncestor)

		tx.WithContext(ctx).Select("apprentice_id", "depth").Where(&entity.Relation{
			MasterID: candidate.ApprenticeID,
		}).Find(&descendants)

		if ancestor.ID != 0 {
			return entity.ErrRelationExist
		}

		newRelations, err = candidate.Connect(&secondaryAncestor.Master, descendants)
		if err != nil {
			return err
		}

		// save new relations
		err = tx.WithContext(ctx).Omit("Master").Omit("Apprentice").Create(&newRelations).Error
		if err != nil {
			return errors.Wrap(err, "unable to create relations")
		}
		return nil
	})
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
		tx.FullSaveAssociations = true
		err = tx.WithContext(ctx).Save(&updated).Error
		if err != nil {
			return errors.Wrap(err, "unable to save relations")
		}
		return nil
	})

}
