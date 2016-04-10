package http

import (
	"bytes"
	"fmt"
	"net"
	"net/http"

	"github.com/mesos-utility/smsmail/g"
	"github.com/toolkits/web/param"
)

func SmsMessageDeal(w http.ResponseWriter, r *http.Request) {
	cfg := g.Config()
	if !cfg.Sms.Enable {
		debugInfo(cfg.Debug, "Sms not enable...")
		return
	}

	addr := cfg.Sms.Url
	if addr == "" {
		debugInfo(cfg.Debug, "Sms Url is null...")
		return
	}

	tos := param.MustString(r, "tos")
	content := param.MustString(r, "content")

	var buf, msg bytes.Buffer
	if conn, err := net.Dial("udp", addr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(&msg, "m:%s c:%s", tos, content)
		fmt.Fprintf(&buf, "s:%04d %s", msg.Len(), msg.String())
		debugInfo(cfg.Debug, "message:"+buf.String())

		if _, err = buf.WriteTo(conn); err != nil {
			debugInfo(cfg.Debug, "Send sms msg failed!!!")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			debugInfo(cfg.Debug, "Send sms msg success!!!")
			http.Error(w, "success", http.StatusOK)
		}
		buf.Reset()
	}
}
