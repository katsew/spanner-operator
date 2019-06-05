package main

import (
	"github.com/katsew/spanner-operator/pkg/operator/instance"
	"github.com/spf13/cobra"
	"log"
)

var instanceOperator instance.Operator
var instanceCommand = cobra.Command{
	Use:   "instance",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		builder := instance.NewBuilder()

		if projectId != "" {
			builder.ProjectId(projectId)
		} else {
			panic("No projectId provided")
		}

		if serviceAccountPath != "" {
			builder.ServiceAccountPath(serviceAccountPath)
		}
		if useMock {
			log.Print("Using mock client to execute")
			instanceOperator = builder.BuildMock("/tmp/spanner-instance-operator")
		} else {
			instanceOperator = builder.Build()
		}
	},
}

func init()  {
	instanceCommand.AddCommand(&createInstanceCommand)
	instanceCommand.AddCommand(&deleteInstanceCommand)
	instanceCommand.AddCommand(&scaleCommand)
	instanceCommand.AddCommand(&getInstanceCommand)
}
