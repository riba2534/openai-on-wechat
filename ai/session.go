package ai

import (
	"context"
	"time"

	"github.com/riba2534/openai-on-wechat/config"
	"github.com/sashabaranov/go-openai"
)

var chat = NewChat()

type UserMessage struct {
	User string                       // 用户
	Time time.Time                    // 当前消息的时间
	Msg  openai.ChatCompletionMessage // 当前消息体
}

func NewUserMessage(user string, msg openai.ChatCompletionMessage) *UserMessage {
	return &UserMessage{
		User: user,
		Time: time.Now(),
		Msg:  msg,
	}
}

type Chat struct {
	UserMessagesMap map[string][]*UserMessage // 记录用户消息上下文
}

func NewChat() *Chat {
	return &Chat{
		UserMessagesMap: map[string][]*UserMessage{},
	}
}

func (c *Chat) Add(userMessage *UserMessage) {
	c.UserMessagesMap[userMessage.User] = append(c.UserMessagesMap[userMessage.User], userMessage)
}

func (c *Chat) Clear(user string) {
	now := time.Now()
	result := []*UserMessage{}
	for _, userMessage := range c.UserMessagesMap[user] {
		if now.Sub(userMessage.Time) < time.Duration(config.C.ContextConfig.CacheMinute)*time.Minute {
			result = append(result, userMessage)
		}
	}
	c.UserMessagesMap[user] = result
}

func (c *Chat) BuildMessage(userKey, systemPrompt string) []openai.ChatCompletionMessage {
	result := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}
	for _, userMessage := range c.UserMessagesMap[userKey] {
		result = append(result, userMessage.Msg)
	}
	return result
}

// q 代表本次问题; user 代表用户key
func GetSessionOpenAITextReply(ctx context.Context, q, userKey, model, systemPrompt string) string {
	// 1. 清理过期消息
	chat.Clear(userKey)
	// 2. 添加本次对话上下文
	chat.Add(NewUserMessage(userKey, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: q,
	}))
	// 3. 获取 OpenAI 回复
	reply := CreateChatCompletion(ctx, model, chat.BuildMessage(userKey, systemPrompt))
	// 4. 把回复添加进上下文
	chat.Add(NewUserMessage(userKey, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply,
	}))
	// 5. 返回结果
	return reply
}
