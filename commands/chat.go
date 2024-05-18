package commands

import (
	"context"
	"fmt"
	"os"

	"git.mkz.me/mycroft/asoai/internal"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func NewChatCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "interact with chatgpt",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chat(args[0])
		},
	}
}

func chat(input string) {
	envVar := os.Getenv("OPENAI_API_KEY")
	if envVar == "" {
		fmt.Printf("could not find OPENAI_API_KEY")
		os.Exit(1)
	}

	currentSession, err := internal.DbGetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	if currentSession == "" {
		// create a new default session
		currentSession, err = SessionCreate(true)
		if err != nil {
			fmt.Printf("could not create a new default session: %v\n", err)
			os.Exit(1)
		}
	}

	session, err := internal.DbGetSession(currentSession)
	if err != nil {
		fmt.Printf("could not get session's details: %v\n", err)
		os.Exit(1)
	}

	messages := []openai.ChatCompletionMessage{}

	for _, message := range session.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    session.Model,
		Messages: messages,
	}

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	session.Messages = append(session.Messages, internal.Message{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	internal.DbSetSession(currentSession, session)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", resp.Choices[0].Message.Content)

	session.Messages = append(session.Messages, internal.Message{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	})

	internal.DbSetSession(currentSession, session)
}
