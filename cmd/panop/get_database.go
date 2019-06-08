package main

import (
	"github.com/spf13/cobra"
	"log"
)

var getDatabaseCommand = cobra.Command{
	Use:  "get [instanceId] [databaseName]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		databaseName := args[1]
		if databaseName == "" {
			panic("No databaseName provided")
		}
		database, err := op.GetDatabase(instanceId, databaseName)
		if err != nil && op.IsNotFoundError(err) {
			log.Printf("Database does not exists, should create first")
			return
		} else if err != nil {
			panic(err)
		}
		log.Printf("Got database: %+v", database)
	},
}
