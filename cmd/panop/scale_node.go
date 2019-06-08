package main

import (
	"github.com/spf13/cobra"
	"strconv"
)

var scaleCommand = cobra.Command{
	Use:  "scale [instanceId] [nodeCount]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		instanceId := args[0]
		if instanceId == "" {
			panic("No instanceId provided")
		}
		nodeCount, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			panic(err)
		}
		if err := op.Scale(args[0], int32(nodeCount)); err != nil {
			panic(err)
		}
	},
}
