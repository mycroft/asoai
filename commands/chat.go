package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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
		currentSessionName, err = SessionCreate(*createName, *createModel, *createPrompt, true)
		if err != nil {
			fmt.Printf("could not create a new session: %v\n", err)
			os.Exit(1)
		}

		// set session as default
		err = SessionSetCurrent(currentSessionName)
		if err != nil {
			fmt.Printf("could not set new session as default: %v\n", err)
			os.Exit(1)
		}
	}

	currentSession, err := db.GetSession(currentSessionName)
	if err != nil {
		fmt.Printf("could not get session's details: %v\n", err)
		os.Exit(1)
	}

	// Prior dealing with the API, finding out if there is some stdin
	stdinData := []string{}
	stdinMessage := ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdinData = append(stdinData, scanner.Text())
		}
	}

	if len(stdinData) > 0 {
		stdinData = append([]string{"```"}, stdinData...)
		stdinData = append(stdinData, "```")

		stdinMessage = strings.Join(stdinData, "\n")
	}

	if len(input) > 0 {
		if len(stdinMessage) > 0 {
			input = strings.Join([]string{input, stdinMessage}, "\n")
		}
	} else {
		input = stdinMessage
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
