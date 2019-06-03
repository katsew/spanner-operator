package main

import (
	ic "github.com/katsew/spanner-operator/pkg/config/instance_config"
	"github.com/katsew/spanner-operator/pkg/helper/gcloud"
	"github.com/katsew/spanner-operator/pkg/operator"
	"github.com/spf13/cobra"
)

var SaoClient operator.SpannerOperator
var projectId string
var instanceId string
var instanceConfig string
var serviceAccountPath string

func main() {
	var sao = &cobra.Command{
		Use: "sao",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			builder := operator.New()

			if projectId != "" {
				builder.ProjectId(projectId)
			} else {
				panic("No projectId provided")
			}

			if instanceId != "" {
				builder.InstanceId(instanceId)
			}  else {
				panic("No instanceId provided")
			}

			if instanceConfig != "" {
				builder.InstanceConfig(ic.FindByName(instanceConfig))
			}

			if serviceAccountPath != "" {
				builder.ServiceAccountPath(serviceAccountPath)
			}

			SaoClient = builder.Build()
		},
	}
	pid, icfg, err := gcloud.GetDefaults()
	if err != nil {
		panic(err)
	}
	sao.PersistentFlags().StringVarP(&projectId, "projectId", "p", pid, "GCP project ID")
	sao.PersistentFlags().StringVarP(&instanceId, "instanceId", "i", "", "Cloud Spanner instance ID")
	sao.PersistentFlags().StringVarP(&instanceConfig, "instanceConfig", "c", icfg, "Cloud Spanner instance config")
	sao.PersistentFlags().StringVarP(&serviceAccountPath, "serviceAccountPath", "s", "", "Path to GCP ServiceAccount")
	sao.AddCommand(&instanceCommand)

	if err := sao.Execute(); err != nil {
		panic(err)
	}

}