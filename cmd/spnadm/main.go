package main

import (
	"github.com/katsew/spanner-operator/cmd/helper"
	"github.com/katsew/spanner-operator/pkg/operator"
	"github.com/spf13/cobra"
	"log"
)

var useMock bool
var projectId string
var serviceAccountPath string
var op operator.Operator

func main() {
	var cli = &cobra.Command{
		Use: "spnadm",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			builder := operator.NewBuilder()
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
				op = builder.BuildMock("/tmp/spnadm")
			} else {
				op = builder.Build()
			}
		},
	}
	pid, _, err := helper.GetGCPDefaults()
	if err != nil {
		panic(err)
	}
	cli.PersistentFlags().BoolVar(&useMock, "use-mock", false, "Use mock client")
	cli.PersistentFlags().StringVarP(&projectId, "project-id", "p", pid, "GCP project ID")
	cli.PersistentFlags().StringVarP(&serviceAccountPath, "service-account-path", "s", "", "Path to GCP ServiceAccount")
	instanceCommand := cobra.Command{
		Use: "instance",
	}
	instanceCommand.AddCommand(
		&createInstanceCommand,
		&deleteInstanceCommand,
		&scaleCommand,
		&getInstanceCommand,
	)
	databaseCommand := cobra.Command{
		Use: "database",
	}
	databaseCommand.AddCommand(
		&createDatabaseCommand,
		&getDatabaseCommand,
		&dropDatabaseCommand,
	)
	cli.AddCommand(
		&instanceCommand,
		&databaseCommand,
	)

	if err := cli.Execute(); err != nil {
		panic(err)
	}

}
