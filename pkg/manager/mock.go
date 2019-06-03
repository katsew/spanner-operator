package manager

import (
	"log"
)

type mockClient struct {}

func (mc *mockClient) CreateInstance(name string) error {
	log.Print("Create instance...")
	return nil
}

func (mc *mockClient) ScaleNode(num int32) error {
	log.Printf("Scale node to %d...", num)
	return nil
}

func (mc *mockClient) DeleteInstance() error {
	log.Print("Delete instance...")
	return nil
}

func (mc *mockClient) UpdateLabels(labels map[string]string) error {
	log.Printf("Update labels to %+v...", labels)
	return nil
}
