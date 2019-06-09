/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"cloud.google.com/go/compute/metadata"
	"context"
	"flag"
	"github.com/katsew/spanner-operator/pkg/controllers/databaseadmins"
	"github.com/katsew/spanner-operator/pkg/operator"
	"log"
	"os"
	"sync"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	databaseadminsClientset "github.com/katsew/spanner-operator/pkg/generated/databaseadmins/clientset/versioned"
	databaseadminsInformers "github.com/katsew/spanner-operator/pkg/generated/databaseadmins/informers/externalversions"
	instanceadminsClientset "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/clientset/versioned"
	instanceadminsInformers "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/informers/externalversions"
	"github.com/katsew/spanner-operator/pkg/signals"

	_ "github.com/katsew/spanner-operator/pkg/controllers/databaseadmins"
	"github.com/katsew/spanner-operator/pkg/controllers/instanceadmins"
)

var (
	masterURL          string
	kubeconfig         string
	debuggable         bool
	mockEnabled        bool
	projectId          string
	serviceAccountPath string
	op                 operator.Operator
)

func main() {

	flag.Parse()
	log.Printf(`
kubeconfig: %s
masterURL: %s
projectId: %s
serviceAccountPath: %s
mockEnabled: %v
debuggable: %v
`, kubeconfig, masterURL, projectId, serviceAccountPath, mockEnabled, debuggable)
	if debuggable {
		log.Print("Enable debugging!")
		klog.InitFlags(nil)
	}
	b := operator.NewBuilder()
	projectId, err := metadata.ProjectID()
	if err != nil {
		log.Print("No projectId got from metadata server, get it from environment variables")
		projectId = os.Getenv("GCP_PROJECT_ID")
	}
	b.ProjectId(projectId)
	serviceAccountPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	b.ServiceAccountPath(serviceAccountPath)
	if !mockEnabled {
		op = b.Build()
	} else {
		dataPath := os.Getenv("MOCK_DATA_PATH")
		if dataPath == "" {
			dataPath = "/tmp/spanner-operator"
		}
		log.Printf("Mock client enabled, building mock with dataPath: %s", dataPath)
		op = b.BuildMock(dataPath)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// set up signals so we handle the first shutdown signal gracefully
	signals.SetupSignalHandler(cancelFunc)

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	instanceadminsCtrl, err := instanceadminsClientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}
	databaseadminsCtrl, err := databaseadminsClientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	var wg sync.WaitGroup

	instanceadminsInformerFactory := instanceadminsInformers.NewSharedInformerFactory(instanceadminsCtrl, time.Second*30)
	instanceadminsController := instanceadmins.NewController(kubeClient, instanceadminsCtrl,
		instanceadminsInformerFactory.Instanceadmins().V1alpha1().SpannerInstances(), op)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = instanceadminsController.Run(2, ctx.Done()); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}()
	databaseadminsInformerFactory := databaseadminsInformers.NewSharedInformerFactory(databaseadminsCtrl, time.Second*30)
	databaseadminsController := databaseadmins.NewController(kubeClient, databaseadminsCtrl,
		databaseadminsInformerFactory.Databaseadmins().V1alpha1().SpannerDatabases(), op)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = databaseadminsController.Run(2, ctx.Done()); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}()

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	go kubeInformerFactory.Start(ctx.Done())
	go instanceadminsInformerFactory.Start(ctx.Done())
	go databaseadminsInformerFactory.Start(ctx.Done())

	<-ctx.Done()

	log.Print("Waiting for all controllers to shut down gracefully")
	wg.Wait()
}

func init() {

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.BoolVar(&debuggable, "debuggable", false, "Enable debug flag.")
	flag.BoolVar(&mockEnabled, "use-mock", false, "Enable mock client.")

}
