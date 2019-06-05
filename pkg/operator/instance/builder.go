package instance

import (
	spanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"context"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
)

type Builder interface {
	ProjectId(projectId string) Builder
	ServiceAccountPath(path string) Builder
	Build() Operator
	BuildMock(dataPath string) *instanceMock
}

type builder struct {
	projectId string
	serviceAccountPath string
}

func NewBuilder() *builder {
	return &builder{}
}

func (b *builder) ProjectId(projectId string) Builder {
	b.projectId = projectId
	return b
}

func (b *builder) ServiceAccountPath(path string) Builder {
	b.serviceAccountPath = path
	return b
}

func (b *builder) Build() Operator {

	ctx := context.Background()
	var client *spanner.InstanceAdminClient
	var err error
	if b.serviceAccountPath != "" {
		data, err := ioutil.ReadFile(b.serviceAccountPath)
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

	return &operator{
		projectId: b.projectId,
		client: client,
	}
}

func (b *builder) BuildMock(dataPath string) *instanceMock {
	err := os.MkdirAll(dataPath, 744)
	if err != nil {
		log.Warn(err.Error())
		log.Warnf("If you use mock client, you should create directory to dataPath: %s", dataPath)
	}
	return &instanceMock{
		projectId: b.projectId,
		dataPath: dataPath,
	}
}