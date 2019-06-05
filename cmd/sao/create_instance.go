package main

import "github.com/spf13/cobra"

var createInstanceCommand = cobra.Command{
	Use:   "create [instanceId] [instanceConfig]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		instanceConfig := args[1]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		if instanceConfig == "" {
			panic("No instanceConfig provided")
		}
		displayName := cmd.Flags().String("display-name", instanceId, "Display name for UI")
		nodeCount := cmd.Flags().Int32P("node-count", "n", 1, "Number of nodes to allocate")
		if err := instanceOperator.CreateInstance(*displayName, instanceId, instanceConfig, *nodeCount); err != nil {
			panic(err)
		}
	},
}