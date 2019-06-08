package operator

import (
	"encoding/json"
	"fmt"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"io/ioutil"
	"log"
	"os"
)

type operatorMock struct {
	projectId string
	dataDir   string
}

func (om *operatorMock) IsNotFoundError(err error) bool {
	return os.IsNotExist(err)
}

func (om *operatorMock) CreateInstance(displayName string, instanceId string, instanceConfig string, nodeCount int32) error {
	log.Print("Create instance...")
	instanceName := fmt.Sprintf("projects/%s/instances/%s", om.projectId, instanceId)
	b, err := json.Marshal(instance.Instance{
		Name:        instanceName,
		Config:      instanceConfig,
		DisplayName: displayName,
		NodeCount:   nodeCount,
		State:       instance.Instance_READY,
		Labels:      map[string]string{"mock": "true"},
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/instance_%s.json", om.dataDir, om.projectId, instanceId), b, 0755)
	return err
}

func (om *operatorMock) GetInstance(instanceId string) (*instance.Instance, error) {
	log.Print("Get instance...")
	instanceName := fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId)
	_, err := os.Stat(instanceName)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(instanceName)
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

func (om *operatorMock) Scale(instanceId string, nodeCount int32) error {
	log.Printf("Scale node to %d...", nodeCount)
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId))
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
	err = ioutil.WriteFile(fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId), b, 0755)
	return nil
}

func (om *operatorMock) DeleteInstance(instanceId string) error {
	log.Print("Delete instance...")
	err := os.Remove(fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId))
	return err
}

func (om *operatorMock) UpdateLabels(instanceId string, labels map[string]string) error {
	log.Printf("Update labels to %+v...", labels)
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId))
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
	err = ioutil.WriteFile(fmt.Sprintf("%s/instance_%s.json", om.dataDir, instanceId), b, 0755)
	return nil
}

func (om *operatorMock) CreateDatabase(instanceId string, name string) error {
	log.Print("Create database...")
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", om.projectId, instanceId, name)
	b, err := json.Marshal(database.Database{
		Name: databaseName,
	})
	err = ioutil.WriteFile(fmt.Sprintf("%s/database_%s.json", om.dataDir, name), b, 0755)
	return err
}

func (om *operatorMock) GetDatabase(instanceId string, name string) (*database.Database, error) {
	log.Print("Get database...")
	databaseName := fmt.Sprintf("%s/database_%s.json", om.dataDir, name)
	_, err := os.Stat(databaseName)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(databaseName)
	if err != nil {
		return nil, err
	}
	var databaseInfo *database.Database
	err = json.Unmarshal(b, &databaseInfo)
	if err != nil {
		return nil, err
	}
	return databaseInfo, nil
}

func (om *operatorMock) DropDatabase(instanceId string, name string) error {
	log.Print("Drop database...")
	err := os.Remove(fmt.Sprintf("%s/database_%s.json", om.dataDir, name))
	return err
}
