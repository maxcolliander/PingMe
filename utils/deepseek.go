package utils

import (
	"context"
	"log"
	"fmt"
	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
)

type DeepSeekClient struct {
	client deepseek.Client
}

func NewDeepSeekClient() *DeepSeekClient {
	apikey := GetEnv("DEEPSEEK_API_KEY")

	if apikey == "" {
		log.Fatal("DEEPSEEK KEY NOT FOUND")
	}

	client, err := deepseek.NewClient(apikey)
	if err != nil {
		log.Fatal("Error creating DeepSeek client: ", err)
	}

	fmt.Println("Connected to DeepSeek")
	return &DeepSeekClient{
		client: client,
	}
}

func (d *DeepSeekClient) Chat(message string) (string, error) {
	chatReq := &request.ChatCompletionsRequest{
		Model: deepseek.DEEPSEEK_CHAT_MODEL,
		Stream: false,
		Messages: []*request.Message{
			{
				Role: "user",
				Content: message,
			},
		},
	}
	chatResp, err := d.client.CallChatCompletionsChat(context.Background(), chatReq)
	if err != nil {
		return "", err
	}
	return chatResp.Choices[0].Message.Content, nil
}