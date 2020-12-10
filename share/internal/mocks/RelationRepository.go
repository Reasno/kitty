// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "glab.tagtic.cn/ad_gains/kitty/app/entity"

	mock "github.com/stretchr/testify/mock"
)

// RelationRepository is an autogenerated mock type for the RelationRepository type
type RelationRepository struct {
	mock.Mock
}

// AddRelations provides a mock function with given fields: ctx, candidate
func (_m *RelationRepository) AddRelations(ctx context.Context, candidate *entity.Relation) error {
	ret := _m.Called(ctx, candidate)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Relation) error); ok {
		r0 = rf(ctx, candidate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryRelations provides a mock function with given fields: ctx, condition
func (_m *RelationRepository) QueryRelations(ctx context.Context, condition entity.Relation) ([]entity.Relation, error) {
	ret := _m.Called(ctx, condition)

	var r0 []entity.Relation
	if rf, ok := ret.Get(0).(func(context.Context, entity.Relation) []entity.Relation); ok {
		r0 = rf(ctx, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Relation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, entity.Relation) error); ok {
		r1 = rf(ctx, condition)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRelations provides a mock function with given fields: ctx, apprentice, existingRelationCallback
func (_m *RelationRepository) UpdateRelations(ctx context.Context, apprentice *entity.User, existingRelationCallback func([]entity.Relation) error) error {
	ret := _m.Called(ctx, apprentice, existingRelationCallback)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User, func([]entity.Relation) error) error); ok {
		r0 = rf(ctx, apprentice, existingRelationCallback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}