package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	dbPath *string
)

var RootCmd = &cobra.Command{
	Use:   "asoai",
	Short: "asoai is another stupid OpenAI client",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

func InitCommands() {
	RootCmd.AddCommand(NewChatCommand())
	RootCmd.AddCommand(NewSessionCommand())

	dbPath = RootCmd.PersistentFlags().String("db-path", "", "database file path")
}
