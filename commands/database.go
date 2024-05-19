package commands

import (
	"fmt"

	"git.mkz.me/mycroft/asoai/internal/database"
	"github.com/adrg/xdg"
)

// Opens a database located in current directory or given directory from
// parameters. If an error happens, it will fail and exit the process.
func OpenDatabase() *database.DB {
	return database.OpenOrFail(GetDefaultDirectory())
}

// Get database default directory; On error, defaults to current working directory
// and prints the error.
func GetDefaultDirectory() string {
	filePath, err := xdg.DataFile("asoai/")

	if err != nil {
		fmt.Printf("could not find a suitable location for datbase: %v; falling back to working directory.\n", err)
		return "."
	}

	return filePath
}
