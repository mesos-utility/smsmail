package http

import (
	"fmt"
	"github.com/toolkits/file"
	"net/http"

	"github.com/soarpenguin/smsmail/g"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.Config())
	})

	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		SmsMessageDeal(w, r)
	})

	http.HandleFunc("/mail", func(w http.ResponseWriter, r *http.Request) {
		MailMessageDeal(w, r)
	})
}
