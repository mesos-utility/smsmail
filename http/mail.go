package http

import (
	"net/http"
	"strings"

	"github.com/soarpenguin/smsmail/g"
	"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
)

func MailMessageDeal(w http.ResponseWriter, r *http.Request) {
	cfg := g.Config()
	if !cfg.Mail.Enable {
		debugInfo(cfg.Debug, "Mail not enable...")
		return
	}

	addr := cfg.Mail.Addr
	if addr == "" {
		debugInfo(cfg.Debug, "Mail Addr is null...")
		return
	}

	tos := param.MustString(r, "tos")
	subject := param.MustString(r, "subject")
	content := param.MustString(r, "content")
	tos = strings.Replace(tos, ",", ";", -1)

	s := smtp.New(cfg.Mail.Addr, cfg.Mail.Username, cfg.Mail.Password)
	err := s.SendMail(cfg.Mail.From, tos, subject, content)
	if err != nil {
		debugInfo(cfg.Debug, "Send mail failed!!!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		debugInfo(cfg.Debug, "Send mail success!!!")
		http.Error(w, "success", http.StatusOK)
	}
}
