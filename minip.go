package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

// ---------------------------------
// -- 小程序登录
// ---------------------------------

type JsCode2SessionResponse struct {
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

func JsCode2Session(appid, secret, code string) (*JsCode2SessionResponse, error) {
	jscode2sess_url := fmt.Sprintf(sns_jsscode2sess_url, appid, secret, code)
	buf, err := goo_http_request.Get(jscode2sess_url)
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	rsp := &JsCode2SessionResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}
	if rsp.Errcode != 0 {
		goo_log.Error(rsp.Errmsg)
		return nil, errors.New(rsp.Errmsg)
	}

	goo_log.
		WithField("openid", rsp.Openid).
		WithField("unionid", rsp.Unionid).
		WithField("session_key", rsp.SessionKey).
		Debug()

	return rsp, nil
}

// ---------------------------------
// -- 解析用户数据
// ---------------------------------

type MinipUserInfoResponse struct {
	Openid   string `json:"openid"`
	Unionid  string `json:"unionid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatarUrl"`
	Gender   int    `json:"gender"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func MinipUserInfo(sessionKey, encryptedData, iv string) (*MinipUserInfoResponse, error) {
	data := goo_utils.Base64Decode(encryptedData)
	key := goo_utils.Base64Decode(sessionKey)

	buf, err := goo_utils.AESCBCDecrypt(data, key, goo_utils.Base64Decode(iv))
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	rsp := &MinipUserInfoResponse{}
	if err = json.Unmarshal(buf, rsp); err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	goo_log.
		WithField("openid", rsp.Openid).
		WithField("unionid", rsp.Unionid).
		WithField("nickname", rsp.Nickname).
		WithField("gender", rsp.Gender).
		WithField("avatar", rsp.Avatar).
		WithField("province", rsp.Province).
		WithField("city", rsp.City).
		Debug()

	return rsp, nil
}

// ---------------------------------
// -- 解析手机号
// ---------------------------------

type WXMobileData struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Appid     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

func MinipMobile(sessionKey, encryptedData, iv string) (*WXMobileData, error) {
	// 解析数据
	encryptedDataRaw := goo_utils.Base64Decode(encryptedData)
	ivRaw := goo_utils.Base64Decode(iv)
	key := goo_utils.Base64Decode(sessionKey)
	buf, err := goo_utils.AESCBCDecrypt(encryptedDataRaw, key, ivRaw)
	if err != nil {
		return nil, err
	}
	// 获取手机号
	dt := &WXMobileData{}
	if err = json.Unmarshal(buf, dt); err != nil {
		return nil, err
	}
	return dt, nil
}

// ---------------------------------
// -- 发送模板消息
// ---------------------------------

func SendTemplateMessage(appid, secret, openid, templateId, page, formId string, data interface{}) error {
	accessToken := CGIToken(appid, secret).Get()
	params := map[string]interface{}{
		"access_token": accessToken,
		"touser":       openid,
		"template_id":  templateId,
		"page":         page,
		"form_id":      formId,
		"data":         data,
	}

	buf, _ := json.Marshal(params)

	messageTplSendUrl := fmt.Sprintf(message_tpl_send_url, accessToken)
	buf, err := goo_http_request.PostJson(messageTplSendUrl, buf)
	if err != nil {
		goo_log.Error(err.Error())
		return err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, rst); err != nil {
		goo_log.Error(err.Error())
		return err
	}
	if rst.ErrCode != 0 {
		goo_log.Error(err.Error())
		return errors.New(rst.ErrMsg)
	}

	return nil
}
