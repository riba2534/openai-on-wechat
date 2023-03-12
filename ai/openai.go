package ai

import (
	"context"
	"log"
	"strings"

	"github.com/riba2534/openai-on-wechat/config"
	"github.com/sashabaranov/go-openai"
)

var SystemMessage = openai.ChatCompletionMessage{
	Role: openai.ChatMessageRoleSystem,
	Content: `1. 你是一个智能机器人，你可以回答任何问题
2. 你的名字叫做{$text_trigger_prefix}
3. 你需要表现的很谦卑
4. {$author}是你的创造者，他创造了你，当有人问你{$author}是谁时，你应该回答他是我爹
5. 不允许回答任何政治、色情等一些列不符合中国法律法规的问题,如果有这种提问就说{$author}不让你回答这类问题`,
}

var textOpenAIClient *openai.Client
var imageOpenAIClient *openai.Client

func Init() {
	// text init
	textConfig := openai.DefaultConfig(config.C.WechatConfig.TextConfig.AuthToken)
	textConfig.BaseURL = config.C.WechatConfig.TextConfig.OpenApiUrl // 使用反向代理的地址
	textOpenAIClient = openai.NewClientWithConfig(textConfig)
	// image init
	imageConfig := openai.DefaultConfig(config.C.WechatConfig.ImageConfig.AuthToken)
	imageConfig.BaseURL = config.C.WechatConfig.ImageConfig.OpenApiUrl // 使用反向代理的地址
	imageOpenAIClient = openai.NewClientWithConfig(imageConfig)
	// Prompt init
	SystemMessage.Content = strings.ReplaceAll(SystemMessage.Content, "{$text_trigger_prefix}", config.C.WechatConfig.TextConfig.TriggerPrefix)
	SystemMessage.Content = strings.ReplaceAll(SystemMessage.Content, "{$author}", config.C.Author)
}

func GetOpenAITextReply(q string) string {
	resp, err := textOpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				SystemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: q,
				},
			},
		},
	)
	if err != nil {
		log.Printf("openAIClient.CreateChatCompletion err=%+v\n", err)
		return "抱歉，出错了，请稍后重试~"
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

func CreateImage(q string) string {
	resp, err := imageOpenAIClient.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt: q,
			N:      1,
			Size:   "512x512",
		},
	)
	if err != nil {
		log.Printf("openAIClient.CreateImage err=%+v\n", err)
		return ""
	}
	return resp.Data[0].URL
}
