package main

import (
	"github.com/spf13/cobra"
	"log"
)

var getInstanceCommand = cobra.Command{
	Use:  "get [instanceId]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		instance, err := op.GetInstance(instanceId)
		if err != nil && op.IsNotFoundError(err) {
			log.Print("Instance does not exists, should create first")
			return
		} else if err != nil {
			panic(err)
		}
		log.Printf("Got instance: %+v", instance)
	},
}
