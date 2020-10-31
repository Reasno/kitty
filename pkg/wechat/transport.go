package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Transport struct {
	wechatAccessTokenUrl string
	wechatGetUserInfoUrl string
	appId                string
	appSecret            string
	client               contract.HttpDoer
}

type WechatConfig struct {
	WechatAccessTokenUrl string
	WeChatGetUserInfoUrl string
	AppId                string
	AppSecret            string
	Client               contract.HttpDoer
}

func NewTransport(conf *WechatConfig) *Transport {
	return &Transport{
		wechatAccessTokenUrl: conf.WechatAccessTokenUrl,
		wechatGetUserInfoUrl: conf.WeChatGetUserInfoUrl,
		appId:                conf.AppId,
		appSecret:            conf.AppSecret,
		client:               conf.Client,
	}
}

func (t *Transport) GetWechatLoginResponse(ctx context.Context, code string) (result *WxLoginResult, err error) {
	url := fmt.Sprintf(t.wechatAccessTokenUrl, t.appId, t.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot build wechat request")
	}
	res, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending wechat request")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (t *Transport) GetWechatUserInfoResult(ctx context.Context, wxLoginResult *WxLoginResult) (*WxUserInfoResult, error) {
	url := fmt.Sprintf(t.wechatGetUserInfoUrl, wxLoginResult.AccessToken, wxLoginResult.Openid)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot build wechat request")
	}
	response, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending wechat request")
	}
	body := response.Body
	defer body.Close()
	bodyByte, err := ioutil.ReadAll(body)
	var result WxUserInfoResult
	err = json.Unmarshal(bodyByte, &result)
	if err != nil {
		return &result, errors.Wrap(err, "cannot unmarshal WxUserInfoResult json")
	}
	return &result, nil
}
