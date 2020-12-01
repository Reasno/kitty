package service

import "github.com/go-kit/kit/log"

func ProvideService(logger log.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}
