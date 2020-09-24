package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
	"time"
)

type cgiToken struct {
	Appid  string
	Secret string
}

func CGIToken(appid, secret string) *cgiToken {
	return &cgiToken{Appid: appid, Secret: secret}
}

func (this *cgiToken) Get() string {
	key := fmt.Sprintf(cgi_token_key, this.Appid)
	return __cache.Get(key).Val()
}

func (this *cgiToken) TTL() time.Duration {
	key := fmt.Sprintf(cgi_token_key, this.Appid)
	return __cache.TTL(key).Val()
}

func (this *cgiToken) Set() error {
	buf, _ := goo.NewRequest().Get(fmt.Sprintf(cgi_token_url, this.Appid, this.Secret))

	rsp := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}{}

	if err := json.Unmarshal(buf, &rsp); err != nil {
		goo.Log.Error("[wx-cgi-access-token]", err.Error())
		return err
	}
	if errCode := rsp.ErrCode; errCode != 0 {
		goo.Log.Error("[wx-cgi-access-token]", rsp.ErrMsg)
		return errors.New(rsp.ErrMsg)
	}

	goo.Log.Debug("[wx-cgi-access-token]", rsp)

	key := fmt.Sprintf(cgi_token_key, this.Appid)
	return __cache.Set(key, rsp.AccessToken, time.Duration(rsp.ExpiresIn)*time.Second).Err()
}
