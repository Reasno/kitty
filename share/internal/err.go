package internal

import (
	"errors"
)

var (
	ErrNoRewardAvailable   = errors.New("no reward available")
	ErrFailedToDecodeToken = errors.New("cannot decode token")
	ErrFailedXtaskRequest  = errors.New("failed to request xtask")
)
