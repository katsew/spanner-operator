package main

import "github.com/spf13/cobra"

var instanceCommand = cobra.Command{
	Use:   "instance",
}

func init()  {
	instanceCommand.AddCommand(&createInstanceCommand)
	instanceCommand.AddCommand(&deleteInstanceCommand)
}
