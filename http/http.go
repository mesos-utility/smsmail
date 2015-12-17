package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/soarpenguin/smsmail/g"
	"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
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
	cfg := g.Config()
	if !cfg.Sms.Enable {
		if cfg.Debug {
			log.Println("Sms not enable...")
		}
		return
	}

	addr := cfg.Sms.Url
	if addr == "" {
		if cfg.Debug {
			log.Println("Sms Url is null...")
		}
		return
	}

	var msg, tos string
	var errMsg string
	r.ParseForm()
	if len(r.Form["content"]) > 0 {
		msg = r.Form["content"][0]
	}
	if len(r.Form["tos"]) > 0 {
		tos = r.Form["tos"][0]
	}

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
		log.Printf("message: %s\n", buf.String())

		if _, err = buf.WriteTo(conn); err != nil {
			log.Printf("Error: %s\n", err)
			RenderMsgJson(w, err.Error())
		} else {
			log.Println("Send message success.")
			RenderDataJson(w, content.String())
		}
		buf.Reset()
	}
}

func MailMessageDeal(w http.ResponseWriter, r *http.Request) {
	cfg := g.Config()
	if !cfg.Mail.Enable {
		if cfg.Debug {
			log.Println("Mail not enable...")
		}
		return
	}

	addr := cfg.Mail.Addr
	if addr == "" {
		if cfg.Debug {
			log.Println("Mail Addr is null...")
		}
		return
	}

	tos := param.MustString(r, "tos")
	subject := param.MustString(r, "subject")
	content := param.MustString(r, "content")
	tos = strings.Replace(tos, ",", ";", -1)

	s := smtp.New(cfg.Mail.Addr, cfg.Mail.Username, cfg.Mail.Password)
	err := s.SendMail(cfg.Mail.From, tos, subject, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, "success", http.StatusOK)
	}
}
