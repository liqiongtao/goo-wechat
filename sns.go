package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
	"net/url"
	"strings"
	"time"
)

// 获取H5授权链接
// authorizeUrl: 授权地址
// originUrl: 授权成功后的回跳地址
func Oauth2AuthorizeUrl(appid, authorizeUrl, originUrl string) string {
	var redirectUrl string
	if strings.Index(authorizeUrl, "?") != -1 {
		redirectUrl = url.QueryEscape(authorizeUrl + "&redirect_url=" + url.QueryEscape(originUrl))
	} else {
		redirectUrl = url.QueryEscape(authorizeUrl + "?redirect_url=" + url.QueryEscape(originUrl))
	}
	oauth2Url := fmt.Sprintf(oauth2_authorize_url, appid, redirectUrl, state)
	goo.Log.Debug("[wx-h5-oauth2]", map[string]interface{}{
		"authorizeUrl": authorizeUrl,
		"originUrl":    originUrl,
		"redirectUrl":  redirectUrl,
		"oauth2Url":    oauth2Url,
	})
	return oauth2Url
}

type Oauth2AccessTokenResponse struct {
	AccessToken  string        `json:"access_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	Openid       string        `json:"openid"`
	Unionid      string        `json:"unionid"`
	Scope        string        `json:"scope"`
	Errcode      int           `json:"errcode"`
	Errmsg       string        `json:"errmsg"`
}

func Oauth2AccessToken(appid, secret, code string) (*Oauth2AccessTokenResponse, error) {
	accessTokenUrl := fmt.Sprintf(sns_oauth2_accessToken_url, appid, secret, code)
	buf, err := goo.NewRequest().Get(accessTokenUrl)
	if err != nil {
		goo.Log.Error("[wx-oauth2-access-token]", err.Error())
		return nil, err
	}

	rsp := &Oauth2AccessTokenResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		goo.Log.Error("[wx-oauth2-access-token]", err.Error())
		return nil, err
	}
	if rsp.Errcode != 0 {
		goo.Log.Error("[wx-oauth2-access-token]", rsp.Errmsg)
		return nil, errors.New(rsp.Errmsg)
	}

	goo.Log.Error("[wx-oauth2-access-token]", rsp)

	return rsp, nil
}

type SnsUserInfoResponse struct {
	Openid   string `json:"openid"`
	Unionid  string `json:"unionid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
	Gender   int    `json:"sex"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
}

func SnsUserInfo(accessToken, openid string) (*SnsUserInfoResponse, error) {
	userInfoUrl := fmt.Sprintf(sns_userinfo_url, accessToken, openid)
	buf, err := goo.NewRequest().Get(userInfoUrl)
	if err != nil {
		goo.Log.Error("[wx-sns-userinfo]", err.Error())

		return nil, err
	}

	rsp := &SnsUserInfoResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		goo.Log.Error("[wx-sns-userinfo]", err.Error())
		return nil, err
	}
	if rsp.Errcode != 0 {
		goo.Log.Error("[wx-sns-userinfo]", rsp.Errmsg)
		return nil, errors.New(rsp.Errmsg)
	}

	goo.Log.Error("[wx-sns-userinfo]", rsp)

	return rsp, nil
}
