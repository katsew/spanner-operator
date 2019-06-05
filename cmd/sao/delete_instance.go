package main

import "github.com/spf13/cobra"

var deleteInstanceCommand = cobra.Command{
	Use:   "delete [instanceId]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := instanceOperator.DeleteInstance(args[0]); err != nil {
			panic(err)
		}
	},
}
