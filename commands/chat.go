package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/asoai/internal/session"
)

var (
	maxTokens *int
)

func NewChatCommand() *cobra.Command {
	chatCommand := cobra.Command{
		Use:   "chat",
		Short: "interact with chatgpt",
		Long:  "query the OpenAI conversation API with current saved discussion in session",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chat(args[0])
		},
	}

	maxTokens = chatCommand.Flags().Int("max-tokens", 0, "Maximum number of tokens to return")

	return &chatCommand
}

func chat(input string) {
	envVar := os.Getenv("OPENAI_API_KEY")
	if envVar == "" {
		fmt.Printf("could not find OPENAI_API_KEY")
		os.Exit(1)
	}

	db := OpenDatabase()

	currentSessionName, err := db.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	if currentSessionName == "" {
		// create a new default session
		currentSessionName, err = SessionCreate(*model, *prompt, true)
		if err != nil {
			fmt.Printf("could not create a new default session: %v\n", err)
			os.Exit(1)
		}
	}

	currentSession, err := db.GetSession(currentSessionName)
	if err != nil {
		fmt.Printf("could not get session's details: %v\n", err)
		os.Exit(1)
	}

	messages := []openai.ChatCompletionMessage{}

	for _, message := range currentSession.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    currentSession.Model,
		Messages: messages,
	}

	if *maxTokens != 0 {
		req.MaxTokens = *maxTokens
	}

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	currentSession.Messages = append(currentSession.Messages, session.Message{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	db.SetSession(currentSessionName, currentSession)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", resp.Choices[0].Message.Content)

	currentSession.Messages = append(currentSession.Messages, session.Message{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	})

	db.SetSession(currentSessionName, currentSession)
}
