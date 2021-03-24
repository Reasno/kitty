package entity

import "context"

type IDAssigner interface {
	ID(ctx context.Context) (uint, error)
}
