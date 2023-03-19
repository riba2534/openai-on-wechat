package ai

import (
	"context"
	"log"
	"strings"

	"github.com/riba2534/openai-on-wechat/config"
	"github.com/riba2534/openai-on-wechat/consts"
	"github.com/riba2534/openai-on-wechat/utils"
	"github.com/sashabaranov/go-openai"
)

func getOpenAIClient(model string) *openai.Client {
	var c openai.ClientConfig
	switch model {
	case openai.GPT3Dot5Turbo:
		c = openai.DefaultConfig(config.C.WechatConfig.TextConfig.AuthToken)
		c.BaseURL = config.C.WechatConfig.TextConfig.OpenApiUrl
	case "image":
		c = openai.DefaultConfig(config.C.WechatConfig.ImageConfig.AuthToken)
		c.BaseURL = config.C.WechatConfig.ImageConfig.OpenApiUrl
	default:
		return nil
	}
	return openai.NewClientWithConfig(c)
}

func CreateChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage) string {
	client := getOpenAIClient(model)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		log.Printf("openAIClient.CreateChatCompletion err=%+v\n", err)
		return consts.ErrTips
	}
	if len(resp.Choices) == 0 {
		log.Printf("resp is err, resp=%s", utils.MarshalAnyToString(resp))
		return consts.ErrTips
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

func CreateImageReply(ctx context.Context, q string) string {
	resp, err := getOpenAIClient("image").CreateImage(ctx, openai.ImageRequest{
		Prompt: q,
		N:      1,
		Size:   "512x512",
	})
	if err != nil {
		log.Printf("openAIClient.CreateImage err=%+v\n", err)
		return ""
	}
	return resp.Data[0].URL
}
