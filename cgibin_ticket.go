package goo_wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type cgiTicket struct {
	Appid  string
	Secret string
}

func CGITicket(appid, secret string) *cgiTicket {
	return &cgiTicket{Appid: appid, Secret: secret}
}

func (this *cgiTicket) Get() string {
	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.Get(key).Val()
}

func (this *cgiTicket) TTL() time.Duration {
	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.TTL(key).Val()
}

func (this *cgiTicket) Set() error {
	accessToken := CGIToken(this.Appid, this.Secret).Get()
	buf, _ := goo_http_request.Get(fmt.Sprintf(cgi_ticket_url, accessToken))

	rsp := struct {
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
	}{}

	if err := json.Unmarshal(buf, &rsp); err != nil {
		return err
	}
	if errCode := rsp.ErrCode; errCode != 0 {
		return errors.New(rsp.ErrMsg)
	}

	goo_log.WithField("ticket", rsp.Ticket).WithField("expire_in", rsp.ExpiresIn).Debug()

	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.Set(key, rsp.Ticket, time.Duration(rsp.ExpiresIn)*time.Second).Err()
}
