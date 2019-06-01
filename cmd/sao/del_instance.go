package main

import "github.com/spf13/cobra"

var deleteInstanceCommand = cobra.Command{
	Use:   "delete",
	Run: func(cmd *cobra.Command, args []string) {
		SaoClient.DeleteInstance()
	},
}
