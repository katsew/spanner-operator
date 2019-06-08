package main

import "github.com/spf13/cobra"

var deleteInstanceCommand = cobra.Command{
	Use:  "delete [instanceId]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		if err := op.DeleteInstance(instanceId); err != nil {
			panic(err)
		}
	},
}
