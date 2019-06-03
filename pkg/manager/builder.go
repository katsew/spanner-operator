package manager

import (
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"context"
	spanner "cloud.google.com/go/spanner/admin/instance/apiv1"

	config "github.com/katsew/spanner-operator/pkg/config/instance_config"
)

type SpannerManagerBuilder interface {
	ProjectId(projectId string) SpannerManagerBuilder
	InstanceId(instanceId string) SpannerManagerBuilder
	InstanceConfig(config config.InstanceConfig) SpannerManagerBuilder
	ServiceAccountPath(path string) SpannerManagerBuilder
	Build() SpannerManager
	BuildMock() SpannerManager
}

type spannerManagerBuilder struct {
	projectId string
	instanceId string
	instanceConfig string
	serviceAccountPath string
}

func New() *spannerManagerBuilder {
	return &spannerManagerBuilder{}
}

func (sb *spannerManagerBuilder) ProjectId(projectId string) SpannerManagerBuilder {
	sb.projectId = projectId
	return sb
}

func (sb *spannerManagerBuilder) InstanceId(instanceId string) SpannerManagerBuilder {
	sb.instanceId = instanceId
	return sb
}

func (sb *spannerManagerBuilder) InstanceConfig(config config.InstanceConfig) SpannerManagerBuilder {
	sb.instanceConfig = config.String()
	return sb
}

func (sb *spannerManagerBuilder) ServiceAccountPath(path string) SpannerManagerBuilder {
	sb.serviceAccountPath = path
	return sb
}

func (sb *spannerManagerBuilder) Build() SpannerManager {

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

	return &spannerManager{
		projectId: sb.projectId,
		instanceId: sb.instanceId,
		instanceConfig: sb.instanceConfig,
		client: client,
	}
}

func (sb *spannerManagerBuilder) BuildMock() SpannerManager {
	return &mockClient{}
}