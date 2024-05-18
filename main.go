package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"git.mkz.me/mycroft/asoai/commands"
)

var rootCmd = &cobra.Command{
	Use:   "asoai",
	Short: "asoai is another stupid OpenAI client",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(commands.ChatCmd)
	rootCmd.AddCommand(commands.SessionCmd)

	commands.SessionCmd.AddCommand(commands.SessionCreateCmd)
	commands.SessionCmd.AddCommand(commands.SessionDumpCmd)
	commands.SessionCmd.AddCommand(commands.SessionListCmd)
	commands.SessionCmd.AddCommand(commands.SessionGetCurrentCmd)
	commands.SessionCmd.AddCommand(commands.SessionSetCurrentCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
