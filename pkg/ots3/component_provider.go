package ots3

import (
	"fmt"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"net/url"
)

func provideUploadManager(conf contract.ConfigReader) *Manager {
	return NewManager(
		conf.String("s3.accessKey"),
		conf.String("s3.accessSecret"),
		conf.String("s3.endpoint"),
		conf.String("s3.region"),
		conf.String("s3.bucket"),
		WithLocationFunc(func(location string) (uri string) {
			u, err := url.Parse(location)
			if err != nil {
				return location
			}
			return fmt.Sprintf(conf.String("s3.cdnUrl"), u.Path[1:])
		}),
	)
}
