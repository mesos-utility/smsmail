package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type MailConfig struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type SmsConfig struct {
	Enable bool   `json:"enable"`
	Url    string `json:"url"`
}

type WeixinConfig struct {
	Enable  bool           `json:"enable"`
	Url     string         `json:"url"`
	Agentid string         `json:"agentid"`
	CorpID  string         `json:"corpid"`
	Secret  string         `json:"secret"`
	Partys  map[string]int `json:"partys"`
	Tags    map[string]int `json:"tags"`
}

type GlobalConfig struct {
	Debug  bool          `json:"debug"`
	Http   *HttpConfig   `json:"http"`
	Mail   *MailConfig   `json:"mail"`
	Sms    *SmsConfig    `json:"sms"`
	Weixin *WeixinConfig `json:"weixin"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g:ParseConfig, ok, ", cfg)
}
