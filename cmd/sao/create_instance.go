package main

import "github.com/spf13/cobra"

var createInstanceCommand = cobra.Command{
	Use:   "create [instanceId] [instanceConfig]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		displayName := cmd.Flags().StringP("displayName", "p", args[0], "Display name for UI")
		nodeCount := cmd.Flags().Int32P("nodeCount", "n", 1, "Number of nodes to allocate")
		if err := instanceOperator.CreateInstance(*displayName, args[0], args[1], *nodeCount); err != nil {
			panic(err)
		}
	},
}