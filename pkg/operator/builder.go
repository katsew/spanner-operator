package operator

import (
	databaseAdmin "cloud.google.com/go/spanner/admin/database/apiv1"
	instanceAdmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/apiv1"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"io/ioutil"
	"os"
)

type Builder interface {
	ProjectId(projectId string) Builder
	ServiceAccountPath(path string) Builder
	Build() Operator
	BuildMock(dataDir string) *operatorMock
}

type builder struct {
	projectId          string
	serviceAccountPath string
}

type Operator interface {
	// InstanceAdmin method
	CreateInstance(displayName string, instanceId string, instanceConfig string, nodeCount int32) error
	GetInstance(instanceId string) (*instance.Instance, error)
	Scale(instanceId string, nodeCount int32) error
	DeleteInstance(instanceId string) error
	UpdateLabels(instanceId string, labels map[string]string) error

	// DatabaseAdmin method
	CreateDatabase(instanceId string, name string) error
	GetDatabase(instanceId string, name string) (*database.Database, error)
	DropDatabase(instanceId string, name string) error

	// Error handle method
	IsNotFoundError(err error) bool
}

type operator struct {
	projectId           string
	instanceId          string
	instanceConfig      string
	instanceAdminClient *instanceAdmin.InstanceAdminClient
	databaseAdminClient *databaseAdmin.DatabaseAdminClient
	client              *spanner.Client
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

	instanceAdminCtx := context.Background()
	databaseAdminCtx := context.Background()
	clientCtx := context.Background()
	var instanceAdminClient *instanceAdmin.InstanceAdminClient
	var databaseAdminClient *databaseAdmin.DatabaseAdminClient
	var client *spanner.Client
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
		ctx := context.Background()
		opt := option.WithTokenSource(conf.TokenSource(ctx))
		instanceAdminClient, err = instanceAdmin.NewInstanceAdminClient(instanceAdminCtx, opt)
		databaseAdminClient, err = databaseAdmin.NewDatabaseAdminClient(databaseAdminCtx, opt)
		client, err = spanner.NewClient(clientCtx, opt)
	} else {
		instanceAdminClient, err = instanceAdmin.NewInstanceAdminClient(instanceAdminCtx)
		databaseAdminClient, err = databaseAdmin.NewDatabaseAdminClient(databaseAdminCtx)
		client, err = spanner.NewClient(clientCtx)
	}

	if err != nil {
		panic(err)
	}

	return &operator{
		projectId:           b.projectId,
		instanceAdminClient: instanceAdminClient,
		databaseAdminClient: databaseAdminClient,
		client:              client,
	}
}

func (b *builder) BuildMock(dataPath string) *operatorMock {
	dataDir := fmt.Sprintf("%s/%s", dataPath, b.projectId)
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Warn(err.Error())
		log.Warnf("If you use mock client, you should create directory to dataDir: %s", dataDir)
	}
	return &operatorMock{
		projectId: b.projectId,
		dataDir:   dataPath,
	}
}
