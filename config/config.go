package config

import (
	"log"

	"github.com/riba2534/openai-on-wechat/utils"
	"github.com/spf13/viper"
)

var C Config

type Config struct {
	Author        string        `json:"author"`
	WechatConfig  WechatConfig  `json:"wechat_config"`
	ContextConfig ContextConfig `json:"context_config"`
}

type AuthConfig struct {
	OpenApiUrl    string `json:"openapi_url"`
	AuthToken     string `json:"auth_token"`
	TriggerPrefix string `json:"trigger_prefix"`
}

type WechatConfig struct {
	TextConfig  AuthConfig `json:"text_config"`
	ImageConfig AuthConfig `json:"image_config"`
}

type ContextConfig struct {
	SwitchOn    bool `json:"switch_on"`
	CacheMinute int  `json:"cache_minute"`
}

func Init() {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// 设置配置文件类型为 JSON
	viper.SetConfigType("json")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置失败，请检查配置文件 `config.json` 的配置, 错误信息: %+v\n", err)
		return
	}
	// 将配置绑定到指定结构体上
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("读取配置失败，请检查配置文件 `config.json` 的配置, 错误信息: %+v\n", err)
		return
	}
	log.Printf("配置加载成功, `config.json` is %s", utils.MarshalAnyToString(C))
}
