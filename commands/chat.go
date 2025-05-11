package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	asoai_chat "git.mkz.me/mycroft/asoai/internal/chat"
	"git.mkz.me/mycroft/asoai/internal/database"
	"git.mkz.me/mycroft/asoai/internal/session"
)

var (
	maxTokens  *int
	useStream  *bool
	newSession *bool
	replMode   *bool

	chatModel       *string
	chatName        *string
	chatDescription *string
	chatPrompt      *string
	chatOutput      *string
)

func NewChatCommand() *cobra.Command {
	chatCommand := cobra.Command{
		Use:   "chat",
		Short: "interact with chatgpt",
		Long:  "query the OpenAI conversation API with current saved discussion in session",
		Run: func(cmd *cobra.Command, args []string) {
			chat(args)
		},
	}

	maxTokens = chatCommand.Flags().Int("max-tokens", 0, "Maximum number of tokens to return")
	useStream = chatCommand.Flags().Bool("stream", false, "Stream response from API")
	newSession = chatCommand.Flags().Bool("new-session", false, "Force creating a new session")
	replMode = chatCommand.Flags().Bool("repl", false, "Enable Repeat Evaluate Print Loop mode")

	chatName = chatCommand.Flags().String("name", "", "Session's name (if created, else ignored)")
	chatDescription = chatCommand.Flags().String("description", "", "Session's description (if created, else ignored)")
	chatModel = chatCommand.Flags().String("model", "gpt-3.5-turbo", "Model (gpt-3.5-turbo, gpt-4-turbo, gpt-4o)")
	chatPrompt = chatCommand.Flags().String("system-prompt", "", "Set system prompt")
	chatOutput = chatCommand.Flags().String("output", "", "Output file path (if not set, output to stdout)")

	return &chatCommand
}

func chat(args []string) {
	var currentSession session.Session

	envVar := os.Getenv("OPENAI_API_KEY")
	if envVar == "" {
		fmt.Printf("could not find OPENAI_API_KEY")
		os.Exit(1)
	}

	input := strings.Join(args, " ")

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	db := database.OpenDatabase(*dbPath)
	defer db.Close()

	currentSessionName, err := db.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	if currentSessionName == "" || *newSession {
		// create a new default session
		currentSessionName, currentSession, err = SessionCreate(db, *chatName, *chatModel, *chatPrompt, true)
		if err != nil {
			fmt.Printf("could not create a new session: %v\n", err)
			os.Exit(1)
		}

		if *chatDescription != "" {
			currentSession.Description = *chatDescription
		}
	} else {
		currentSession, err = db.GetSession(currentSessionName)
		if err != nil {
			fmt.Printf("could not get %s session's details: %v\n", currentSessionName, err)
			os.Exit(1)
		}
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

	if len(input) == 0 && !*replMode {
		fmt.Println("no input; exiting")
		os.Exit(1)
	}

	for {
		messages := []openai.ChatCompletionMessage{}

		for _, message := range currentSession.Messages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    message.Role,
				Content: message.Content,
			})
		}

		// overwrite system prompt, if needed
		if *chatPrompt != "" {
			messages[0].Content = *chatPrompt
		}

		model := currentSession.Model

		if *chatModel != "" {
			// model was changed but session is not updated.
			model = *chatModel
		}

		req := openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		}

		if *maxTokens != 0 {
			req.MaxTokens = *maxTokens
		}

		req.Stream = *useStream

		if *replMode {
			// Read input
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("user> ")
			input, err = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if err != nil || len(input) == 0 {
				break
			}
		}

		// Patch input to handle inserting files
		input, err = asoai_chat.PatchInput(input)
		if err != nil {
			fmt.Printf("error while patching input: %v\n", err)
			os.Exit(1)
		}

		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		currentSession.Messages = append(currentSession.Messages, session.Message{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		// Save session, as we added an input
		db.SetSession(currentSessionName, currentSession)

		returnedRole := ""
		returnedContent := ""

		if !*useStream {
			resp, err := client.CreateChatCompletion(context.Background(), req)
			if err != nil {
				fmt.Printf("ChatCompletion error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("assistant> %s\n", resp.Choices[0].Message.Content)

			currentSession.Messages = append(currentSession.Messages, session.Message{
				Role:    resp.Choices[0].Message.Role,
				Content: resp.Choices[0].Message.Content,
			})

			if *chatOutput != "" {
				f, err := os.OpenFile(*chatOutput, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("could not create output file: %v\n", err)
					os.Exit(1)
				}
				_, err = f.WriteString(resp.Choices[0].Message.Content)
				if err != nil {
					fmt.Printf("could not write to output file: %v\n", err)
					os.Exit(1)
				}

				f.Close()
			}
		} else {
			resp, err := client.CreateChatCompletionStream(context.Background(), req)
			if err != nil {
				fmt.Printf("ChatCompletionStream error: %v\n", err)
			}
			defer resp.Close()

			fmt.Printf("assistant> ")

			for {
				content, err := resp.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Printf("error while streaming response...")
					os.Exit(1)
				}

				if content.Choices[0].Delta.Role != "" {
					returnedRole = content.Choices[0].Delta.Role
					continue
				}

				returnedContent += content.Choices[0].Delta.Content

				fmt.Print(content.Choices[0].Delta.Content)
			}

			fmt.Println()

			currentSession.Messages = append(currentSession.Messages, session.Message{
				Role:    returnedRole,
				Content: returnedContent,
			})

			if *chatOutput != "" {
				f, err := os.OpenFile(*chatOutput, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("could not create output file: %v\n", err)
					os.Exit(1)
				}
				_, err = f.WriteString(returnedContent)
				if err != nil {
					fmt.Printf("could not write to output file: %v\n", err)
					os.Exit(1)
				}

				f.Close()
			}
		}

		if !*replMode {
			break
		}
	}

	db.SetSession(currentSessionName, currentSession)
}
