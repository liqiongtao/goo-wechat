package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
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
	buf, _ := goo_http_request.Get(fmt.Sprintf(cgi_token_url, this.Appid, this.Secret))

	rsp := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}{}

	if err := json.Unmarshal(buf, &rsp); err != nil {
		goo_log.Error(err.Error())
		return err
	}
	if errCode := rsp.ErrCode; errCode != 0 {
		goo_log.Error(rsp.ErrMsg)
		return errors.New(rsp.ErrMsg)
	}

	goo_log.WithField("access_token", rsp.AccessToken).WithField("expire_in", rsp.ExpiresIn).Debug()

	key := fmt.Sprintf(cgi_token_key, this.Appid)
	return __cache.Set(key, rsp.AccessToken, time.Duration(rsp.ExpiresIn)*time.Second).Err()
}
