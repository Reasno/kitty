package taobao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Transport struct {
	taobaoAccessTokenUrl string
	taobaoGetUserInfoUrl string
	appKey               string
	partnerId            string
	appSecret            string
	client               contract.HttpDoer
}

type TaobaoConfig struct {
	TaobaoAccessTokenUrl string
	TaobaoGetUserInfoUrl string
	PartnerId            string
	AppKey               string
	AppSecret            string
	Client               contract.HttpDoer
}

func NewTransport(conf *TaobaoConfig) *Transport {
	return &Transport{
		taobaoAccessTokenUrl: conf.TaobaoAccessTokenUrl,
		taobaoGetUserInfoUrl: conf.TaobaoGetUserInfoUrl,
		partnerId:            conf.PartnerId,
		appKey:               conf.AppKey,
		appSecret:            conf.AppSecret,
		client:               conf.Client,
	}
}

// 获取签名
func (t *Transport) getSign(args url.Values) string {
	// 获取Key
	keys := []string{}
	for k := range args {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	query := t.appSecret
	for _, k := range keys {
		query += k + args.Get(k)
	}
	query += t.appSecret
	// 使用MD5加密
	signBytes := md5.Sum([]byte(query))
	// 把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(signBytes[:]))
}

func (t *Transport) GetLoginResponse(ctx context.Context, code string) (result interface{}, err error) {
	var data url.Values
	data.Add("app_key", t.appKey)
	data.Add("format", "json")
	data.Add("method", "taobao.top.auth.token.create")
	data.Add("partner_id", t.partnerId)
	data.Add("sign_method", "md5")
	data.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	data.Add("v", "2")
	data.Add("code", code)
	data.Add("sign", t.getSign(data))

	req, err := http.NewRequestWithContext(ctx, "POST", t.taobaoAccessTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot build taobao request")
	}
	res, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending taobao request")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//func (t *Transport) GetUserInfoResult(ctx context.Context, wxLoginResult *WxLoginResult) (*WxUserInfoResult, error) {
//	url := fmt.Sprintf(t.taobaoGetUserInfoUrl, wxLoginResult.AccessToken, wxLoginResult.Openid)
//	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot build wechat request")
//	}
//	response, err := t.client.Do(req)
//	if err != nil {
//		return nil, errors.Wrap(err, "error sending wechat request")
//	}
//	body := response.Body
//	defer body.Close()
//	bodyByte, err := ioutil.ReadAll(body)
//	var result WxUserInfoResult
//	err = json.Unmarshal(bodyByte, &result)
//	if err != nil {
//		return &result, errors.Wrap(err, "cannot unmarshal WxUserInfoResult json")
//	}
//	return &result, nil
//}
