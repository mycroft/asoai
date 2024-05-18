package internal

import "github.com/sashabaranov/go-openai"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	Model    string    `json:"model"`
	Messages []Message `json:"message"`
}

func NewSession() Session {
	return Session{
		Model: openai.GPT3Dot5Turbo,
		Messages: []Message{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are chatgpt, a large language model trained by OpenAI, based on the GPT-4 architecture.",
			},
		},
	}
}
