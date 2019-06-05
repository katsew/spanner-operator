package main

import (
	"github.com/katsew/spanner-operator/pkg/helper/gcloud"
	"github.com/katsew/spanner-operator/pkg/operator/instance"
	"github.com/spf13/cobra"
	"log"
)

var instanceOperator instance.Operator
var useMock bool
var projectId string
var instanceId string
var instanceConfig string
var serviceAccountPath string

func main() {
	var io = &cobra.Command{
		Use: "io",
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
	pid, _, err := gcloud.GetDefaults()
	if err != nil {
		panic(err)
	}
	io.PersistentFlags().BoolVar(&useMock, "use-mock", false, "Use mock client")
	io.PersistentFlags().StringVarP(&projectId, "project-id", "p", pid, "GCP project ID")
	io.PersistentFlags().StringVarP(&serviceAccountPath, "sa-path", "s", "", "Path to GCP ServiceAccount")
	io.AddCommand(&instanceCommand)

	if err := io.Execute(); err != nil {
		panic(err)
	}

}