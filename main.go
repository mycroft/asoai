package main

import (
	"fmt"
	"os"

	"git.mkz.me/mycroft/asoai/commands"
)

func init() {
	commands.InitCommands()
}

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
