package goo_wechat

import (
	"fmt"
	"github.com/liqiongtao/goo/utils"
	"net/url"
	"strings"
	"time"
)

func JsApi(appid, secret, urlStr string) map[string]interface{} {
	ticket := CGITicket(appid, secret).Get()

	ts := time.Now().Unix()
	nonceStr := utils.NonceStr()

	urlStr, _ = url.QueryUnescape(urlStr)
	urlStr = strings.Split(urlStr, "#")[0]

	rawstr := fmt.Sprintf(jsapi_ticket_qs, ticket, nonceStr, ts, urlStr)
	rawstr = utils.SHA1([]byte(rawstr))

	params := map[string]interface{}{
		"debug":     false,
		"appId":     appid,
		"timestamp": ts,
		"nonceStr":  nonceStr,
		"signature": rawstr,
		"jsApiList": []string{"checkJsApi", "onMenuShareTimeline", "onMenuShareAppMessage", "chooseWXPay",
			"openLocation", "getLocation", "chooseImage", "previewImage", "uploadImage", "downloadImage",
			"startRecord", "stopRecord", "onVoiceRecordEnd", "playVoice", "pauseVoice", "stopVoice",
			"onVoicePlayEnd", "uploadVoice", "downloadVoice", "translateVoice", "getNetworkType", "scanQRCode",
			"addCard", "chooseCard", "openCard"},
	}

	return params
}
