package operator

import (
	"log"
)

type spannerMockOperator struct {}

func (mc *spannerMockOperator) CreateInstance(name string) error {
	log.Print("Create instance...")
	return nil
}

func (mc *spannerMockOperator) ScaleNode(num int32) error {
	log.Printf("Scale node to %d...", num)
	return nil
}

func (mc *spannerMockOperator) DeleteInstance() error {
	log.Print("Delete instance...")
	return nil
}

func (mc *spannerMockOperator) UpdateLabels(labels map[string]string) error {
	log.Printf("Update labels to %+v...", labels)
	return nil
}
