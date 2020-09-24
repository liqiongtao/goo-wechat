package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
)

func MenuCreate(appid, secret, content string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuCreateUrl := fmt.Sprintf(menu_create_url, accessToken)
	buf, err := goo.NewRequest().JsonContentType().Post(menuCreateUrl, []byte(content))
	if err != nil {
		goo.Log.Error("[wx-menu-create]", err.Error())
		return err
	}

	rst := struct {
		ErrorCode int    `json:"errorcode"`
		ErrMsg    string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo.Log.Error("[wx-menu-create]", err.Error())
		return err
	}
	if rst.ErrorCode != 0 {
		goo.Log.Error("[wx-menu-create]", rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}

	return nil
}

func MenuGet(appid, secret string) (string, error) {
	accessToken := CGIToken(appid, secret).Get()

	menuGetrl := fmt.Sprintf(menu_get_url, accessToken)
	buf, err := goo.NewRequest().JsonContentType().Get(menuGetrl)
	if err != nil {
		goo.Log.Error("[wx-menu-get]", err.Error())
		return "", err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo.Log.Error("[wx-menu-get]", err.Error())
		return "", err
	}
	if rst.ErrCode != 0 {
		goo.Log.Error("[wx-menu-get]", rst.ErrMsg)
		return "", errors.New(rst.ErrMsg)
	}

	return string(buf), nil
}

func MenuDelete(appid, secret string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuDeleteUrl := fmt.Sprintf(menu_del_url, accessToken)
	buf, err := goo.NewRequest().JsonContentType().Post(menuDeleteUrl, nil)
	if err != nil {
		goo.Log.Error("[wx-menu-del]", err.Error())
		return err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goo.Log.Error("[wx-menu-del]", err.Error())
		return err
	}
	if rst.ErrCode != 0 {
		goo.Log.Error("[wx-menu-del]", rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}
	return nil
}
