package repository

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type FileRepo struct {
	uploader   contract.Uploader
	client     contract.HttpDoer
	once       sync.Once
	defaultImg []byte
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
		f.once.Do(func() {
			req, _ := http.NewRequest("GET", "http://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png", nil)
			req.WithContext(ctx)
			buf, _ := f.client.Do(req)
			f.defaultImg, _ = ioutil.ReadAll(buf.Body)
			buf.Body.Close()
		})
		body = ioutil.NopCloser(bytes.NewReader(f.defaultImg))
	} else {
		body = resp.Body
	}
	defer body.Close()
	return f.uploader.Upload(ctx, body)
}

func (f *FileRepo) Upload(ctx context.Context, reader io.Reader) (newUrl string, err error) {
	return f.uploader.Upload(ctx, reader)
}
