package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mesos-utility/smsmail/g"
	"github.com/toolkits/web/param"
)

type TokenResp struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

type MsgResp struct {
	Errcode float64 `json:"errcode"`
	Errmsg  string  `json:"errmsg"`
}

type Text struct {
	Content string `json:"content"`
}

type MessageBody struct {
	ToUser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Totag   string `json:"totag"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	Text    Text   `json:"text"`
	Safe    string `json:"safe"`
}

func WeixinGetToken(corpid, secret string) (gtoken string, err error) {
	var data TokenResp

	cfg := g.Config()
	gurl := fmt.Sprintf("%s/gettoken?corpid=%s&corpsecret=%s", cfg.Weixin.Url, corpid, secret)

	r, err := http.Get(gurl)
	if err != nil || r.StatusCode != http.StatusOK {
		return "", err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(body, &data)
	gtoken = data.AccessToken

	return gtoken, nil
}

func WeixinSendMsg(gtoken, content string) error {
	var msgresp MsgResp

	cfg := g.Config()
	m := &MessageBody{
		ToUser:  "",
		Toparty: "3",
		Totag:   "",
		Msgtype: "text",
		Agentid: cfg.Weixin.Agentid,
		Text: Text{
			Content: content,
		},
		Safe: "0",
	}

	mJson, _ := json.MarshalIndent(m, " ", "  ")
	contentReader := bytes.NewReader(mJson)
	addr := fmt.Sprintf("%s/message/send?access_token=%s", cfg.Weixin.Url, gtoken)
	//fmt.Printf("%v\n", addr)

	req, _ := http.NewRequest("POST", addr, contentReader)
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &msgresp)

	if msgresp.Errcode != 0 {
		debugInfo(cfg.Debug, string(body))

		return fmt.Errorf(msgresp.Errmsg)
	}

	return err
}

func WeixinMessageDeal(w http.ResponseWriter, r *http.Request) {
	cfg := g.Config()
	if !cfg.Weixin.Enable {
		debugInfo(cfg.Debug, "Weixin not enable...")
		return
	}

	addr := cfg.Weixin.Url
	if addr == "" {
		debugInfo(cfg.Debug, "Weixin Url is null...")
		return
	}

	if cfg.Weixin.CorpID == "" || cfg.Weixin.Secret == "" {
		debugInfo(cfg.Debug, "Weixin corpid or secret is null...")
		return
	}

	//tos := param.MustString(r, "tos")
	content := param.MustString(r, "content")

	gtoken, _ := WeixinGetToken(cfg.Weixin.CorpID, cfg.Weixin.Secret)

	err := WeixinSendMsg(gtoken, content)
	if err != nil {
		debugInfo(cfg.Debug, "Send Weixin msg failed!!!")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		debugInfo(cfg.Debug, "Send Weixin msg success!!!")
		//http.Error(w, "success", http.StatusOK)
	}
}
