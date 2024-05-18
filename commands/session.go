package commands

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/asoai/internal"
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

	sessionCommand.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "create a new session",

		Run: func(cmd *cobra.Command, args []string) {
			sessionUuid, err := SessionCreate(false)
			if err != nil {
				fmt.Printf("could not create a new session: %s\n", err)
				os.Exit(1)
			}

			fmt.Println(sessionUuid)
		},
	})

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

	return &sessionCommand
}

func SessionCreate(setDefaultSession bool) (string, error) {
	uuid := uuid.New()
	if err := internal.DbCreateSession(uuid.String(), internal.NewSession()); err != nil {
		return "", err
	}

	if setDefaultSession {
		if err := internal.SetCurrentSession(uuid.String()); err != nil {
			return "", err
		}
	}

	return uuid.String(), nil
}

func SessionList() {
	sessions, err := internal.DbListSessions()
	if err != nil {
		fmt.Printf("could not list sessions: %v\n", err)
		os.Exit(1)
	}

	for _, sessionUuid := range sessions {
		fmt.Println(sessionUuid)
	}
}

func SessionGetCurrent() {
	currentSessionUuid, err := internal.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(currentSessionUuid)
}

func SessionSetCurrent(session string) error {
	return internal.SetCurrentSession(session)
}

func SessionDump() {
	currentSessionUuid, err := internal.GetCurrentSession()
	if err != nil {
		fmt.Printf("could not get current session: %v\n", err)
		os.Exit(1)
	}

	session, err := internal.DbGetSession(currentSessionUuid)
	if err != nil {
		fmt.Printf("could not retrieve session details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current session: %s\n", currentSessionUuid)
	fmt.Printf("Model: %s\n\n", session.Model)

	for _, message := range session.Messages {
		fmt.Printf("%s> %s\n", message.Role, message.Content)
	}
}
