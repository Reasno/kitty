package contract

import (
	"context"
	"io"
)

type Uploader interface {
	Upload(ctx context.Context, reader io.Reader) (newUrl string, err error)
}
