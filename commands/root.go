package commands

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/asoai/internal/database"
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

// Builds the cobra argument parsing state
func InitCommands() {
	RootCmd.AddCommand(NewChatCommand())
	RootCmd.AddCommand(NewSessionCommand())
	RootCmd.AddCommand(NewModelsCommand())
	RootCmd.AddCommand(NewDatabaseCommand())

	dbPath = RootCmd.PersistentFlags().String("db-path", "", "database file path")
}

// Opens a database located in current directory or given directory from
// parameters. If an error happens, it will fail and exit the process.
func OpenDatabase() *database.DB {
	if *dbPath != "" {
		return database.OpenOrFail(*dbPath)
	}
	return database.OpenOrFail(GetDefaultDbFilePath())
}

// Get database default directory; On error, defaults to current working directory
// and prints the error.
func GetDefaultDbFilePath() string {
	filePath, err := xdg.DataFile("asoai/data.db")

	if err != nil {
		fmt.Printf("could not find a suitable location for datbase: %v; falling back to working directory.\n", err)
		return "."
	}

	return filePath
}
