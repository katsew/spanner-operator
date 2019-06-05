package main

import (
	"github.com/katsew/spanner-operator/pkg/helper/gcloud"
	"github.com/spf13/cobra"
)

var useMock bool
var projectId string
var serviceAccountPath string

func main() {
	var io = &cobra.Command{
		Use: "panop",
	}
	pid, _, err := gcloud.GetDefaults()
	if err != nil {
		panic(err)
	}
	io.PersistentFlags().BoolVar(&useMock, "use-mock", false, "Use mock client")
	io.PersistentFlags().StringVarP(&projectId, "project-id", "p", pid, "GCP project ID")
	io.PersistentFlags().StringVarP(&serviceAccountPath, "service-account-path", "s", "", "Path to GCP ServiceAccount")
	io.AddCommand(&instanceCommand)

	if err := io.Execute(); err != nil {
		panic(err)
	}

}