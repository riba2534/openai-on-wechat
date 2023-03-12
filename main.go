package main

import (
	"io"
	"log"
	"os"

	"github.com/riba2534/openai-on-wechat/ai"
	"github.com/riba2534/openai-on-wechat/bot"
	"github.com/riba2534/openai-on-wechat/config"
	"github.com/riba2534/openai-on-wechat/utils"
)

func Init() {
	// 1. `config.json` init
	config.Init()
	// 2. log init
	f, _ := os.OpenFile("run.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetPrefix("[openai-on-wechat] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// 3. OpenAI API init
	ai.Init()
	// 4. Wechat bot init
	if err := bot.Init(); err != nil {
		log.Fatalf("微信登录失败, 错误信息为: %v", err)
	}
	log.Println("登录成功")
}

func main() {
	Init()
	// 获取登陆的用户
	self, err := bot.Bot.GetCurrentUser()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	log.Printf("self=%s", utils.MarshalAnyToString(self))
	bot.Bot.MessageHandler = MessageHandler // 微信消息回调注册
	bot.Bot.Block()
}
