package operator

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *operator) updateInstance(req *instance.UpdateInstanceRequest) error {
	ctx := context.Background()
	op, err := o.instanceAdminClient.UpdateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Print("Update instance done!")
	}
	return err
}

func (o *operator) CreateInstance(
	displayName string,
	instanceId string,
	instanceConfig string,
	nodeCount int32,
) error {

	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	instanceInfo := &instance.Instance{
		Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", o.projectId, instanceConfig),
		DisplayName: displayName,
		Name:        instanceName,
		NodeCount:   nodeCount,
	}
	req := &instance.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", o.projectId),
		InstanceId: o.instanceId,
		Instance:   instanceInfo,
	}
	op, err := o.instanceAdminClient.CreateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err == nil {
		log.Print("Create instance done!")
	}

	return err
}

func (o *operator) GetInstance(instanceId string) (*instance.Instance, error) {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	req := &instance.GetInstanceRequest{
		Name: instanceName,
	}
	i, err := o.instanceAdminClient.GetInstance(ctx, req)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (o *operator) Scale(instanceId string, nodeCount int32) error {
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	instanceInfo := &instance.Instance{
		Name:      instanceName,
		NodeCount: nodeCount,
	}
	req := &instance.UpdateInstanceRequest{
		Instance: instanceInfo,
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"node_count"},
		},
	}
	err := o.updateInstance(req)
	if err != nil {
		return err
	}
	log.Printf("Update node count to %d", nodeCount)
	return nil
}

func (o *operator) DeleteInstance(instanceId string) error {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	err := o.instanceAdminClient.DeleteInstance(ctx, &instance.DeleteInstanceRequest{
		Name: instanceName,
	})
	return err
}

func (o *operator) UpdateLabels(instanceId string, labels map[string]string) error {
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	instanceInfo := &instance.Instance{
		Name:   instanceName,
		Labels: labels,
	}
	req := &instance.UpdateInstanceRequest{
		Instance: instanceInfo,
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"labels"},
		},
	}
	err := o.updateInstance(req)
	if err != nil {
		return err
	}
	log.Printf("Update labels to %+v", labels)
	return nil
}

func (o *operator) CreateDatabase(instanceId string, name string) error {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	req := &database.CreateDatabaseRequest{
		Parent:          instanceName,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`;", name),
	}
	op, err := o.databaseAdminClient.CreateDatabase(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Print("Create database done!")
	}
	return nil
}

func (o *operator) GetDatabase(instanceId string, name string) (*database.Database, error) {
	ctx := context.Background()
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", o.projectId, instanceId, name)
	req := &database.GetDatabaseRequest{
		Name: databaseName,
	}
	return o.databaseAdminClient.GetDatabase(ctx, req)
}

func (o *operator) DropDatabase(instanceId string, name string) error {
	ctx := context.Background()
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", o.projectId, instanceId, name)
	req := &database.DropDatabaseRequest{
		Database: databaseName,
	}
	return o.databaseAdminClient.DropDatabase(ctx, req)
}

func (o *operator) IsNotFoundError(err error) bool {
	s, ok := status.FromError(err)
	return ok && s.Code() == codes.NotFound
}
