package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
	"github.com/liqiongtao/goo/utils"
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
	buf, err := goo.NewRequest().Get(jscode2sess_url)
	if err != nil {
		return nil, err
	}

	rsp := &JsCode2SessionResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		return nil, err
	}
	if rsp.Errcode != 0 {
		return nil, errors.New(rsp.Errmsg)
	}

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
	data := utils.Base64Decode(encryptedData)
	key := utils.Base64Decode(sessionKey)

	buf, err := utils.AESCBCEncrypt(data, key, utils.Base64Decode(iv))
	if err != nil {
		return nil, err
	}

	userInfo := &MinipUserInfoResponse{}
	if err = json.Unmarshal(buf, userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
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
	buf, err := goo.NewRequest().JsonContentType().Post(messageTplSendUrl, buf)
	if err != nil {
		return err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, rst); err != nil {
		return err
	}
	if rst.ErrCode != 0 {
		return errors.New(rst.ErrMsg)
	}

	return nil
}
