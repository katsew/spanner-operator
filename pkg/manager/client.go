package manager

import (
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"fmt"
	"context"
	adminInstance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/labstack/gommon/log"
	"google.golang.org/genproto/protobuf/field_mask"
)

type SpannerManager interface {
	CreateInstance(name string) error
	ScaleNode(num int32) error
	DeleteInstance() error
	UpdateLabels(labels map[string]string) error
}

type spannerManager struct {
	projectId string
	instanceId string
	instanceConfig string
	client *adminInstance.InstanceAdminClient
}

func (sm *spannerManager) CreateInstance(name string) error {

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

func (sm *spannerManager) ScaleNode(num int32) error {

	ctx := context.Background()
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
	op, err := sm.client.UpdateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Printf("Scale to %d done!", num)
	}
	return err

}

func (sm *spannerManager) DeleteInstance() error {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
	err := sm.client.DeleteInstance(ctx, &instance.DeleteInstanceRequest{
		Name: instanceName,
	})
	return err
}

func (sm *spannerManager) UpdateLabels(labels map[string]string) error {
	ctx := context.Background()
	instanceName := fmt.Sprintf("projects/%s/instances/%s", sm.projectId, sm.instanceId)
	instanceInfo := &instance.Instance{
		Name: instanceName,
		Labels: labels,
	}
	req := &instance.UpdateInstanceRequest{
		Instance: instanceInfo,
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"node_count"},
		},
	}
	op, err := sm.client.UpdateInstance(ctx, req)
	if err != nil {
		return err
	}
	if _, err = op.Wait(ctx); err == nil {
		log.Printf("Update labels to %+v done!", labels)
	}
	return err
}
