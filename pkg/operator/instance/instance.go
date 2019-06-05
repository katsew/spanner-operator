package instance

import (
	cli "cloud.google.com/go/spanner/admin/instance/apiv1"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

type Operator interface {
	CreateInstance(displayName string, instanceId string, instanceConfig string, nodeCount int32) error
	GetInstance(instanceId string) (*instance.Instance, error)
	Scale(instanceId string, nodeCount int32) error
	DeleteInstance(instanceId string) error
	UpdateLabels(instanceId string, labels map[string]string) error
}

type operator struct {
	projectId string
	instanceId string
	instanceConfig string
	client *cli.InstanceAdminClient
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
		Config: fmt.Sprintf("projects/%s/instanceConfigs/%s", o.projectId, instanceConfig),
		DisplayName: displayName,
		Name: instanceName,
		NodeCount: nodeCount,
	}
	req := &instance.CreateInstanceRequest{
		Parent: fmt.Sprintf("projects/%s", o.projectId),
		InstanceId: o.instanceId,
		Instance: instanceInfo,
	}
	op, err := o.client.CreateInstance(ctx, req)
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
	i, err := o.client.GetInstance(ctx, req)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (o *operator) updateInstance(req *instance.UpdateInstanceRequest) error {
	ctx := context.Background()
	op, err := o.client.UpdateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Print("Update instance done!")
	}
	return err
}

func (o *operator) Scale(instanceId string, nodeCount int32) error {
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	instanceInfo := &instance.Instance{
		Name: instanceName,
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
	err := o.client.DeleteInstance(ctx, &instance.DeleteInstanceRequest{
		Name: instanceName,
	})
	return err
}

func (o *operator) UpdateLabels(instanceId string, labels map[string]string) error {
	instanceName := fmt.Sprintf("projects/%s/instances/%s", o.projectId, instanceId)
	instanceInfo := &instance.Instance{
		Name: instanceName,
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
