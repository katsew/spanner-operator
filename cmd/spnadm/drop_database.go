package main

import "github.com/spf13/cobra"

var dropDatabaseCommand = cobra.Command{
	Use:  "drop [instanceId] [databaseName]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		databaseName := args[1]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		if databaseName == "" {
			panic("No databaseName provided")
		}
		if err := op.DropDatabase(instanceId, databaseName); err != nil {
			panic(err)
		}
	},
}
