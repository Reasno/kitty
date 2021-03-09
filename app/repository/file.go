package repository

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type FileRepo struct {
	uploader contract.Uploader
	client   contract.HttpDoer
}

func NewFileRepo(uploader contract.Uploader, client contract.HttpDoer) *FileRepo {
	return &FileRepo{
		uploader: uploader,
		client:   client,
	}
}

func (f *FileRepo) UploadFromUrl(ctx context.Context, url string) (newUrl string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "cannot build request")
	}
	TimeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	req = req.WithContext(TimeoutCtx)
	resp, err := f.client.Do(req)
	var body io.ReadCloser
	if err != nil {
		body = ioutil.NopCloser(bytes.NewReader(nil))
	} else {
		body = resp.Body
	}
	defer body.Close()
	return f.uploader.Upload(ctx, body)
}

func (f *FileRepo) Upload(ctx context.Context, reader io.Reader) (newUrl string, err error) {
	return f.uploader.Upload(ctx, reader)
}
