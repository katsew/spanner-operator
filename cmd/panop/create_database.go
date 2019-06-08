package main

import "github.com/spf13/cobra"

var createDatabaseCommand = cobra.Command{
	Use:  "create [instanceId] [databaseName]",
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
		if err := op.CreateDatabase(instanceId, databaseName); err != nil {
			panic(err)
		}
	},
}
