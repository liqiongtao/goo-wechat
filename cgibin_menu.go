package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func MenuCreate(appid, secret, content string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuCreateUrl := fmt.Sprintf(menu_create_url, accessToken)
	buf, err := goo_http_request.PostJson(menuCreateUrl, []byte(content))
	if err != nil {
		goo_log.Error(err.Error())
		return err
	}

	rst := struct {
		ErrorCode int    `json:"errorcode"`
		ErrMsg    string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo_log.Error(err.Error())
		return err
	}
	if rst.ErrorCode != 0 {
		goo_log.Error(rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}

	return nil
}

func MenuGet(appid, secret string) (string, error) {
	accessToken := CGIToken(appid, secret).Get()

	menuGetrl := fmt.Sprintf(menu_get_url, accessToken)
	buf, err := goo_http_request.Get(menuGetrl)
	if err != nil {
		goo_log.Error(err.Error())
		return "", err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo_log.Error(err.Error())
		return "", err
	}
	if rst.ErrCode != 0 {
		goo_log.Error(rst.ErrMsg)
		return "", errors.New(rst.ErrMsg)
	}

	return string(buf), nil
}

func MenuDelete(appid, secret string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuDeleteUrl := fmt.Sprintf(menu_del_url, accessToken)
	buf, err := goo_http_request.PostJson(menuDeleteUrl, nil)
	if err != nil {
		goo_log.Error(err.Error())
		return err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo_log.Error(err.Error())
		return err
	}
	if rst.ErrCode != 0 {
		goo_log.Error(rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}
	return nil
}
