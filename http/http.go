package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"

	"github.com/soarpenguin/smsmail/g"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// start http server
func Start() {
	go startHttpServer()
}

func configRoutes() {
	configCommonRoutes()
}

func startHttpServer() {
	if !g.Config().Http.Enable {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	// init url mapping
	configRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.startHttpServer ok, listening ", addr)
	log.Fatalln(s.ListenAndServe())
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, err string) {
	RenderJson(w, map[string]string{"msg": "failed", "data": err})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

func SmsMessageDeal(w http.ResponseWriter, r *http.Request) {
	if !g.Config().Sms.Enable {
		if g.Config().Debug {
			log.Println("Sms not enable...")
		}
		return
	}

	addr := g.Config().Sms.Url
	if addr == "" {
		if g.Config().Debug {
			log.Println("Sms Url is null...")
		}
		return
	}

	msg := r.URL.Query().Get("content")
	tos := r.URL.Query().Get("tos")
	var errMsg string

	switch {
	case len(msg) == 0:
		errMsg = "content must not be null."
	case len(tos) == 0:
		errMsg = "tos must not be null."
	}

	if len(errMsg) != 0 {
		RenderMsgJson(w, errMsg)
		return
	}

	var buf, content bytes.Buffer
	if conn, err := net.Dial("udp", addr); err != nil {
		RenderMsgJson(w, err.Error())
		return
	} else {
		fmt.Fprintf(&content, "m:%s c:%s", tos, msg)
		fmt.Fprintf(&buf, "s:%04d %s", content.Len(), content.String())
		log.Println("message: %s", buf.String())

		if _, err = buf.WriteTo(conn); err != nil {
			log.Println("Error: %s", err)
			RenderMsgJson(w, err.Error())
		} else {
			log.Println("Send message success.")
			RenderDataJson(w, content.String())
		}
		buf.Reset()
	}
}

func MailMessageDeal(w http.ResponseWriter, r *http.Request) {
	if !g.Config().Mail.Enable {
		if g.Config().Debug {
			log.Println("Mail not enable...")
		}
		return
	}

	addr := g.Config().Mail.Url
	if addr == "" {
		if g.Config().Debug {
			log.Println("Mail Url is null...")
		}
		return
	}

	msg := r.URL.Query().Get("content")
	subject := r.URL.Query().Get("subject")
	tos := r.URL.Query().Get("tos")
	var errMsg string

	switch {
	case len(msg) == 0:
		errMsg = "content must not be null."
	case len(tos) == 0:
		errMsg = "tos must not be null."
	case len(subject) == 0:
		errMsg = "subject must not be null."
	}

	if len(errMsg) != 0 {
		RenderMsgJson(w, errMsg)
		return
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(addr)
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	defer c.Close()

	// Set the sender and recipient.
	c.Mail("zhuyefeng@youku.com")
	c.Rcpt(tos)
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	defer wc.Close()

	buf := bytes.NewBufferString(msg)
	if _, err = buf.WriteTo(wc); err != nil {
		RenderMsgJson(w, err.Error())
		return
	} else {
		RenderDataJson(w, buf.String())
	}
}
