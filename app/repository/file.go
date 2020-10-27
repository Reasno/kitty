package repository

import (
	"context"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type FileRepo struct {
	uploader contract.Uploader
	client contract.HttpDoer
}

func NewFileRepo(uploader contract.Uploader, client contract.HttpDoer) *FileRepo  {
	return &FileRepo{
		uploader: uploader,
		client: client,
	}
}

func (f *FileRepo) UploadFromUrl(ctx context.Context, url string) (newUrl string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "cannot build request")
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "cannot fetch image")
	}
	body := resp.Body
	defer body.Close()
	return f.uploader.Upload(ctx, body)
}

func (f *FileRepo) Upload(ctx context.Context, reader io.Reader) (newUrl string, err error) {
	return f.uploader.Upload(ctx, reader)
}

