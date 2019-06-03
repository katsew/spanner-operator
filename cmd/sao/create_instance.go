package main

import "github.com/spf13/cobra"

var createInstanceCommand = cobra.Command{
	Use:   "create [instance name]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := SaoClient.CreateInstance(args[0]); err != nil {
			panic(err)
		}
	},
}
