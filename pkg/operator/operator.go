package operator

import (
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"fmt"
	"context"
	cli "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/labstack/gommon/log"
	"google.golang.org/genproto/protobuf/field_mask"
)

type SpannerOperator interface {
	CreateInstance(name string) error
	ScaleNode(num int32) error
	DeleteInstance() error
	UpdateLabels(labels map[string]string) error
}

type spannerOperator struct {
	projectId string
	instanceId string
	instanceConfig string
	client *cli.InstanceAdminClient
}

func (sm *spannerOperator) CreateInstance(name string) error {

	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
	instanceInfo := &instance.Instance{
		Config: fmt.Sprintf("projects/%s/instanceConfigs/%s", sm.projectId, sm.instanceConfig),
		DisplayName: name,
		Name: instanceName,
		NodeCount: 1,
	}
	req := &instance.CreateInstanceRequest{
		Parent: fmt.Sprintf("projects/%s", sm.projectId),
		InstanceId: sm.instanceId,
		Instance: instanceInfo,
	}
	op, err := sm.client.CreateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err == nil {
		log.Print("Create instance done!")
	}

	return err
}

func (sm *spannerOperator) updateInstance(req *instance.UpdateInstanceRequest) error {
	ctx := context.Background()
	op, err := sm.client.UpdateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Print("Update instance done!")
	}
	return err
}

func (sm *spannerOperator) ScaleNode(num int32) error {

	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
	instanceInfo := &instance.Instance{
		Name: instanceName,
		NodeCount: num,
	}
	req := &instance.UpdateInstanceRequest{
		Instance: instanceInfo,
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"node_count"},
		},
	}
	err := sm.updateInstance(req)
	if err != nil {
		return err
	}
	log.Printf("Update node count to %d", num)
	return nil
}

func (sm *spannerOperator) DeleteInstance() error {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
	err := sm.client.DeleteInstance(ctx, &instance.DeleteInstanceRequest{
		Name: instanceName,
	})
	return err
}

func (sm *spannerOperator) UpdateLabels(labels map[string]string) error {
	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
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
	err := sm.updateInstance(req)
	if err != nil {
		return err
	}
	log.Printf("Update labels to %+v", labels)
	return nil
}
