package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var getInstanceCommand = cobra.Command{
	Use:   "get [instanceId]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		instance, err := instanceOperator.GetInstance(instanceId)
		if os.IsNotExist(err) {
			panic("Instance not exists")
		}
		if err != nil {
			panic(err)
		}
		log.Printf("Got instance: %+v", instance)
	},
}