package commands

import "github.com/spf13/cobra"

func NewDatabaseCommand() *cobra.Command {
	databaseCommand := cobra.Command{
		Use:   "database",
		Short: "database management functions",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	databaseShrinkCommand := cobra.Command{
		Use:   "shrink",
		Short: "shrink/compact database",
		Run: func(cmd *cobra.Command, args []string) {
			db := OpenDatabase()
			db.Shrink()
			defer db.Close()
		},
	}

	databaseCommand.AddCommand(&databaseShrinkCommand)

	return &databaseCommand
}
