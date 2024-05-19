package commands

import "git.mkz.me/mycroft/asoai/internal/database"

// Opens a database located in current directory or given directory from
// parameters. If an error happens, it will fail and exit the process.
func OpenDatabase() *database.DB {
	return database.OpenOrFail(".")
}
