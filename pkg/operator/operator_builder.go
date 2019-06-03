package operator

import (
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"context"
	spanner "cloud.google.com/go/spanner/admin/instance/apiv1"

	config "github.com/katsew/spanner-operator/pkg/config/instance_config"
)

type SpannerOperatorBuilder interface {
	ProjectId(projectId string) SpannerOperatorBuilder
	InstanceId(instanceId string) SpannerOperatorBuilder
	InstanceConfig(config config.InstanceConfig) SpannerOperatorBuilder
	ServiceAccountPath(path string) SpannerOperatorBuilder
	Build() SpannerOperator
	BuildMock() *spannerMockOperator
}

type spannerOperatorBuilder struct {
	projectId string
	instanceId string
	instanceConfig string
	serviceAccountPath string
}

func New() *spannerOperatorBuilder {
	return &spannerOperatorBuilder{}
}

func (sb *spannerOperatorBuilder) ProjectId(projectId string) SpannerOperatorBuilder {
	sb.projectId = projectId
	return sb
}

func (sb *spannerOperatorBuilder) InstanceId(instanceId string) SpannerOperatorBuilder {
	sb.instanceId = instanceId
	return sb
}

func (sb *spannerOperatorBuilder) InstanceConfig(config config.InstanceConfig) SpannerOperatorBuilder {
	sb.instanceConfig = config.String()
	return sb
}

func (sb *spannerOperatorBuilder) ServiceAccountPath(path string) SpannerOperatorBuilder {
	sb.serviceAccountPath = path
	return sb
}

func (sb *spannerOperatorBuilder) Build() SpannerOperator {

	ctx := context.Background()
	var client *spanner.InstanceAdminClient
	var err error
	if sb.serviceAccountPath != "" {
		data, err := ioutil.ReadFile(sb.serviceAccountPath)
		if err != nil {
			panic(err)
		}
		conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spanner.admin", "https://www.googleapis.com/auth/spanner.data")
		if err != nil {
			panic(err)
		}
		opt := option.WithTokenSource(conf.TokenSource(ctx))
		client, err = spanner.NewInstanceAdminClient(ctx, opt)
	} else {
		client, err = spanner.NewInstanceAdminClient(ctx)
	}

	if err != nil {
		panic(err)
	}

	return &spannerOperator{
		projectId: sb.projectId,
		instanceId: sb.instanceId,
		instanceConfig: sb.instanceConfig,
		client: client,
	}
}

func (sb *spannerOperatorBuilder) BuildMock() *spannerMockOperator {
	return &spannerMockOperator{}
}