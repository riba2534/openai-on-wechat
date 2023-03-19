package config

import (
	"io/ioutil"
	"log"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/riba2534/openai-on-wechat/utils"
)

var (
	C      *Config
	Prompt string
	once   sync.Once
)

type Config struct {
	WechatConfig  *WechatConfig  `json:"wechat_config"`
	ContextConfig *ContextConfig `json:"context_config"`
}

type AuthConfig struct {
	OpenApiUrl    string `json:"openapi_url"`
	AuthToken     string `json:"auth_token"`
	TriggerPrefix string `json:"trigger_prefix"`
}

type WechatConfig struct {
	TextConfig  *AuthConfig `json:"text_config"`
	ImageConfig *AuthConfig `json:"image_config"`
}

type ContextConfig struct {
	SwitchOn    bool `json:"switch_on"`
	CacheMinute int  `json:"cache_minute"`
}

func (c *Config) IsValid() bool {
	if c.WechatConfig == nil || c.ContextConfig == nil {
		return false
	}

	authConfigs := []*AuthConfig{
		c.WechatConfig.TextConfig,
		c.WechatConfig.ImageConfig,
	}

	for _, authConfig := range authConfigs {
		if authConfig == nil || authConfig.OpenApiUrl == "" || authConfig.AuthToken == "" || authConfig.TriggerPrefix == "" {
			return false
		}
	}
	if c.ContextConfig.CacheMinute <= 0 {
		return false
	}
	return true
}

func init() {
	once.Do(func() {
		// 1. 读取 `config.json`
		data, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Fatalf("读取配置文件失败，请检查配置文件 `config.json` 的配置, 错误信息: %+v\n", err)
		}
		config := Config{}
		if err = jsoniter.Unmarshal(data, &config); err != nil {
			log.Fatalf("读取配置文件失败，请检查配置文件 `config.json` 的格式, 错误信息: %+v\n", err)
		}
		if !config.IsValid() {
			log.Fatal("配置文件校验失败，请检查 `config.json`")
		}
		C = &config
		// 2. 读取 prompt.txt
		prompt, err := ioutil.ReadFile("prompt.txt")
		if err != nil {
			log.Fatalf("读取配置文件失败，请检查配置文件 `prompt.txt` 的配置, 错误信息: %+v\n", err)
		}
		Prompt = string(prompt)
		log.Printf("配置加载成功, `config.json` is \n%s\n`prompt.txt` is \n%s\n", utils.MarshalAnyToString(C), Prompt)
	})
}
