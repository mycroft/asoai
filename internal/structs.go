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

func NewSession(model, prompt string) Session {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}
	if prompt == "" {
		prompt = "You are chatgpt, a large language model trained by OpenAI, based on the GPT-4 architecture."
	}

	return Session{
		Model: model,
		Messages: []Message{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
		},
	}
}
