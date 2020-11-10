package sms

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const DATETIMESTRING = "20060102150405"

type Transport struct {
	tag        string
	sendUrl    string
	balanceUrl string
	userName   string
	password   string
	keyPrefix  string
	client     contract.HttpDoer
}

type TransportConfig struct {
	Tag        string
	SendUrl    string
	BalanceUrl string
	UserName   string
	Password   string
	Client     contract.HttpDoer
}

func NewTransport(config *TransportConfig) *Transport {
	if config.Client == nil {
		config.Client = http.DefaultClient
	}
	return &Transport{
		tag:        config.Tag,
		sendUrl:    config.SendUrl,
		balanceUrl: config.BalanceUrl,
		userName:   config.UserName,
		password:   config.Password,
		client:     config.Client,
	}
}

func (s *Transport) Send(ctx context.Context, mobile string, content string) error {
	now := time.Now().Format(DATETIMESTRING)
	args := url.Values{}
	args.Add("content", fmt.Sprintf(s.tag,content))
	args.Add("mobile", mobile)
	args.Add("tkey", now)
	args.Add("username", s.userName)
	args.Add("password", md5(md5(s.password)+now))

	argsStr := args.Encode()
	req, err := http.NewRequestWithContext(ctx, "POST", s.sendUrl, bytes.NewBuffer([]byte(argsStr)))
	if err != nil {
		return errors.Wrap(err, "cannot create post request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot send post request in sender")
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (s *Transport) GetBalance(ctx context.Context) (int64, error) {
	now := time.Now().Format(DATETIMESTRING)
	args := url.Values{}

	args.Add("tkey", now)
	args.Add("username", s.userName)
	args.Add("password", md5(md5(s.password)+now))
	argsStr := args.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", s.balanceUrl, bytes.NewBuffer([]byte(argsStr)))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(string(r), 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}
