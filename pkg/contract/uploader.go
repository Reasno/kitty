package contract

import (
	"context"
	"io"
)

type Uploader interface {
	UploadFromIOReader(ctx context.Context, reader io.Reader) (newUrl string, err error)
	UploadFromUrl(ctx context.Context, url string) (newUrl string, err error)
}
