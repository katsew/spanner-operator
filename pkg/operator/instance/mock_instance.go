package instance

import (
	"encoding/json"
	"fmt"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"io/ioutil"
	"log"
	"os"
)

type instanceMock struct {
	projectId string
	dataPath string
}

func (im *instanceMock) CreateInstance(displayName string, instanceId string, instanceConfig string, nodeCount int32) error {
	log.Print("Create instance...")
	instanceName := fmt.Sprintf("projects/<project>/instances/%s", instanceId)
	b, err := json.Marshal(instance.Instance{
		Name: instanceName,
		Config: instanceConfig,
		DisplayName: displayName,
		NodeCount: nodeCount,
		State: instance.Instance_READY,
		Labels: map[string]string{"mock": "true"},
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId), b, 0755)
	return err
}

func (im *instanceMock) GetInstance(instanceId string) (*instance.Instance, error) {
	log.Print("Get instance...")
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId))
	if err != nil {
		return nil, err
	}
	var instanceInfo *instance.Instance
	err = json.Unmarshal(b, &instanceInfo)
	if err != nil {
		return nil, err
	}
	return instanceInfo, nil
}

func (im *instanceMock) Scale(instanceId string, nodeCount int32) error {
	log.Printf("Scale node to %d...", nodeCount)
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId))
	if err != nil {
		return err
	}
	var instanceInfo *instance.Instance
	err = json.Unmarshal(b, &instanceInfo)
	if err != nil {
		return err
	}
	instanceInfo.NodeCount = nodeCount
	b, err = json.Marshal(instanceInfo)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId), b, 0755)
	return nil
}

func (im *instanceMock) DeleteInstance(instanceId string) error {
	log.Print("Delete instance...")
	err := os.Remove(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId))
	return err
}

func (im *instanceMock) UpdateLabels(instanceId string, labels map[string]string) error {
	log.Printf("Update labels to %+v...", labels)
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId))
	if err != nil {
		return err
	}
	var instanceInfo *instance.Instance
	err = json.Unmarshal(b, &instanceInfo)
	if err != nil {
		return err
	}
	instanceInfo.Labels = labels
	b, err = json.Marshal(instanceInfo)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", im.dataPath, instanceId), b, 0755)
	return nil
}
