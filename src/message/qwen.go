package message

import (
	"context"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Qwen(data string) string {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("DASHSCOPE_API_KEY")),
		option.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1/"),
	)
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("你是一个专业的AI助手,擅长用中文回答问题,请务必只用纯文本回答,不要使用Markdown、代码块或格式化符号且字数不超过300字。"),
					openai.UserMessage(data),
				},
			),
			Model: openai.F("qwen-plus"),
		},
	)

	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}