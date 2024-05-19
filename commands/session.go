package commands

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/asoai/internal/session"
)

var (
	createName   *string
	createModel  *string
	createPrompt *string

	configDescription *string
	configModel       *string
	configPrompt      *string
	configRename      *string
)

func NewSessionCommand() *cobra.Command {
	sessionCommand := cobra.Command{
		Use:   "session",
		Short: "handle sessions",
		Long:  "sessions are used to get context about chats with the AI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	newSessionCommand := cobra.Command{
		Use:   "create",
		Short: "create a new session",

		Run: func(cmd *cobra.Command, args []string) {
			sessionUuid, _, err := SessionCreate(*createName, *createModel, *createPrompt, false)
			if err != nil {
				fmt.Printf("could not create a new session: %s\n", err)
				os.Exit(1)
			}

			fmt.Println(sessionUuid)
		},
	}

	createName = newSessionCommand.Flags().String("name", "", "Session's name")
	createModel = newSessionCommand.Flags().String("model", "gpt-3.5-turbo", "Model (gpt-3.5-turbo, gpt-4-turbo, gpt-4o)")
	createPrompt = newSessionCommand.Flags().String("system-prompt", "", "Initial system prompt")
	sessionCommand.AddCommand(&newSessionCommand)

	sessionCommand.AddCommand(&cobra.Command{
		Use:   "dump",
		Short: "dump current session",
		Run: func(cmd *cobra.Command, args []string) {
			SessionDump()
		},
	})

	sessionCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list existing sessions",
		Run: func(cmd *cobra.Command, args []string) {
			SessionList()
		},
	})

	sessionCommand.AddCommand(&cobra.Command{
		Use:   "get-current",
		Short: "returns current session uuid",
		Run: func(cmd *cobra.Command, args []string) {
			SessionGetCurrent()
		},
	})

	sessionCommand.AddCommand(&cobra.Command{
		Use:   "set-current",
		Short: "set current session uuid",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			SessionSetCurrent(args[0])
		},
	})

	configCommand := cobra.Command{
		Use:   "config",
		Short: "configure the current session",
		Run: func(cmd *cobra.Command, args []string) {
			SessionConfigure()
		},
	}

	configDescription = configCommand.Flags().String("description", "", "Set a description")
	configPrompt = configCommand.Flags().String("prompt", "", "Set a prompt (gpt-3.5-turbo, gpt-4-turbo, gpt-4o)")
	configModel = configCommand.Flags().String("model", "", "Set a model")
	configRename = configCommand.Flags().String("rename", "", "Rename session")

	sessionCommand.AddCommand(&configCommand)

	return &sessionCommand
}

func SessionCreate(name, model, prompt string, setDefaultSession bool) (string, session.Session, error) {
	sessionName := uuid.New().String()
	if name != "" {
		sessionName = name
	}

	db := OpenDatabase()
	createdSession := session.NewSession(model, prompt)

	if err := db.SetSession(sessionName, createdSession); err != nil {
		return "", session.Session{}, err
	}

	if setDefaultSession {
		if err := db.SetCurrentSession(sessionName); err != nil {
			return "", session.Session{}, err
		}
	}

	return sessionName, createdSession, nil
}

func SessionList() {
	db := OpenDatabase()

	sessions, err := db.ListSessions()
	if err != nil {
		fmt.Printf("could not list sessions: %v\n", err)
		os.Exit(1)
	}

	for _, name := range sessions {
		session, err := db.GetSession(name)
		if err != nil {
			fmt.Printf("could not get session %s: %v\n", name, err)
			os.Exit(1)
		}

		output := name

		if session.Description != "" {
			output = fmt.Sprintf("%s - %s", name, session.Description)
		}

		fmt.Println(output)
	}
}

func SessionGetCurrent() {
	db := OpenDatabase()

	currentSessionName, err := db.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(currentSessionName)
}

func SessionSetCurrent(session string) error {
	return OpenDatabase().SetCurrentSession(session)
}

func SessionDump() {
	db := OpenDatabase()

	currentSessionName, err := db.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	session, err := db.GetSession(currentSessionName)
	if err != nil {
		fmt.Printf("could not retrieve session details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current session: %s\n", currentSessionName)
	fmt.Printf("Model: %s\n", session.Model)

	if session.Description != "" {
		fmt.Printf("Description: %s\n", session.Description)
	}

	fmt.Println()

	for _, message := range session.Messages {
		fmt.Printf("%s> %s\n", message.Role, message.Content)
	}
}

func SessionConfigure() error {
	db := OpenDatabase()

	currentSessionName, err := db.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	session, err := db.GetSession(currentSessionName)
	if err != nil {
		fmt.Printf("could not retrieve session details: %v\n", err)
		os.Exit(1)
	}

	if *configDescription != "" {
		session.Description = *configDescription
	}

	if *configModel != "" {
		session.Model = *configModel
	}

	if *configPrompt != "" {
		session.Messages[0].Content = *configPrompt
	}

	if *configRename != "" {
		if err = db.DeleteSession(currentSessionName); err != nil {
			fmt.Printf("could not rename session: %v\n", err)
			os.Exit(1)
		}

		currentSessionName = *configRename
	}

	err = db.SetSession(currentSessionName, session)
	if err != nil {
		fmt.Printf("could not save session: %v\n", err)
		os.Exit(1)
	}

	db.SetCurrentSession(currentSessionName)

	return nil
}
