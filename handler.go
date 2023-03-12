package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/riba2534/openai-on-wechat/ai"
	"github.com/riba2534/openai-on-wechat/config"
	"github.com/riba2534/openwechat"
)

func MessageHandler(msg *openwechat.Message) {
	if !msg.IsText() {
		return
	}
	switch {
	case strings.HasPrefix(msg.Content, config.C.WechatConfig.TextConfig.TriggerPrefix):
		// 文字回复
		if config.C.ContextConfig.SwitchOn {
			go textSessionReplyHandler(msg)
		} else {
			go textReplyHandler(msg)
		}
	case strings.HasPrefix(msg.Content, config.C.WechatConfig.ImageConfig.TriggerPrefix):
		// 图片回复
		go imageReplyHandler(msg)
	}
}

// 文字回复
func textReplyHandler(msg *openwechat.Message) {
	log.Printf("[text] Request: %s", msg.Content) // 输出请求消息到日志
	reply := ai.GetOpenAITextReply(strings.TrimSpace(msg.Content))
	log.Printf("[text] Response: %s", reply) // 输出回复消息到日志
	_, err := msg.ReplyText(reply)
	if err != nil {
		log.Printf("msg.ReplyText Error: %+v", err)
	}
}

// 带有上下文的文字回复
func textSessionReplyHandler(msg *openwechat.Message) {
	log.Printf("[text] Request: %s", msg.Content) // 输出请求消息到日志
	user := func() string {
		s := msg.FromUserName
		if msg.IsSendBySelf() {
			s = msg.ToUserName
		}
		return s
	}()
	reply := ai.GetSessionOpenAITextReply(strings.TrimSpace(msg.Content), user)
	log.Printf("[text] Response: %s", reply) // 输出回复消息到日志
	_, err := msg.ReplyText(reply)
	if err != nil {
		log.Printf("msg.ReplyText Error: %+v", err)
	}
}

// 回复图片
func imageReplyHandler(msg *openwechat.Message) {
	log.Printf("[image] Request: %s", msg.Content)
	url := ai.CreateImage(strings.TrimSpace(strings.TrimPrefix(msg.Content, config.C.WechatConfig.ImageConfig.TriggerPrefix)))
	if url == "" {
		log.Printf("[image] Response: url 为空")
		msg.ReplyText("抱歉，出错了，请稍后重试~")
		return
	}
	log.Printf("[image] Response: url = %s", url)
	image, err := downloadImage(url)
	if err != nil {
		log.Printf("[image] downloadImage err, err=%+v", err)
		msg.ReplyText("抱歉，出错了，请稍后重试~")
		return
	}
	_, err = msg.ReplyImage(image)
	if err != nil {
		log.Printf("msg.ReplyImage Error: %+v", err)
	}
}

func downloadImage(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("downloadImage failed, err=%+v", err)
		return nil, err
	}
	return response.Body, nil
}
