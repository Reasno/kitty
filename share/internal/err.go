package internal

import (
	"errors"
)

var (
	ErrNoRewardAvailable  = errors.New("no reward available")
	ErrFailedXtaskRequest = errors.New("failed to request xtask")
)
