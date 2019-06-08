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
	"fmt"
	"log"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	instancev1alpha1 "github.com/katsew/spanner-operator/pkg/apis/instanceadmins/v1alpha1"
	clientset "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/clientset/versioned"
	spannerscheme "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/clientset/versioned/scheme"
	informers "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/informers/externalversions/instanceadmins/v1alpha1"
	listers "github.com/katsew/spanner-operator/pkg/generated/instanceadmins/listers/instanceadmins/v1alpha1"

	"github.com/katsew/spanner-operator/pkg/operator"
)

const controllerAgentName = "spanner-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a SpannerInstance is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a SpannerInstance fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Spanner"
	// MessageResourceSynced is the message used for an Event fired when a Spanner
	// is synced successfully
	MessageResourceSynced = "SpannerInstance synced successfully"
)

// Controller is the controller implementation for SpannerInstance resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// spannerclientset is a clientset for our own API group
	spannerclientset clientset.Interface

	spannerInstanceLister  listers.SpannerInstanceLister
	spannerInstancesSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new spanner controller
func NewController(
	kubeclientset kubernetes.Interface,
	spannerclientset clientset.Interface,
	spannerInstanceInformer informers.SpannerInstanceInformer) *Controller {

	// Create event broadcaster
	// Add spanner-controller types to the default Kubernetes Scheme so Events can be
	// logged for spanner-controller types.
	utilruntime.Must(spannerscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:          kubeclientset,
		spannerclientset:       spannerclientset,
		spannerInstanceLister:  spannerInstanceInformer.Lister(),
		spannerInstancesSynced: spannerInstanceInformer.Informer().HasSynced,
		workqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Spanners"),
		recorder:               recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when SpannerInstance resources change
	spannerInstanceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueSpannerInstance,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueSpannerInstance(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting SpannerInstance controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.spannerInstancesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process SpannerInstance resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// SpannerInstance resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

var op operator.Operator

func init() {
	b := operator.NewBuilder()
	var projectId string
	projectId, err := metadata.ProjectID()
	if err != nil {
		log.Print("No projectId got from metadata server, get it from environment variables")
		projectId = os.Getenv("GCP_PROJECT_ID")
	}
	b.ProjectId(projectId)
	serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	b.ServiceAccountPath(serviceAccountPath)
	mockEnabled := os.Getenv("MOCK_ENABLED")
	if mockEnabled != "true" {
		op = b.Build()
	} else {
		dataPath := os.Getenv("MOCK_DATA_PATH")
		if dataPath == "" {
			dataPath = "/tmp/spanner-operator"
		}
		log.Printf("Mock client enabled, building mock with dataPath: %s", dataPath)
		op = b.BuildMock(dataPath)
	}
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the SpannerInstance resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {

	log.Printf("Get key: %s", key)
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the SpannerInstance resource with this namespace/name
	spannerInstance, err := c.spannerInstanceLister.SpannerInstances(namespace).Get(name)
	log.Printf("Get spanner instance %+v", spannerInstance)
	if err != nil {
		// The SpannerInstance resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			log.Printf("spannerInstance '%s' in work queue no longer exists", key)
			_, err := op.GetInstance(name)
			if err != nil && op.IsNotFoundError(err) {
				utilruntime.HandleError(fmt.Errorf("spannerInstance '%s' in work queue no longer exists", key))
				return nil
			} else if err != nil {
				return err
			}
			err = op.DeleteInstance(name)
			if err != nil {
				return err
			}
			utilruntime.HandleError(fmt.Errorf("spannerInstance '%s' in work queue no longer exists", key))
			return nil
		}
		log.Printf("Error: %s", err.Error())
		return err
	}

	inst, err := op.GetInstance(name)
	if err != nil && op.IsNotFoundError(err) {
		log.Printf("SpannerInstance does not exists on GCP, create new one with name: %s", spannerInstance.Name)
		err = op.CreateInstance(spannerInstance.Spec.DisplayName, spannerInstance.Name, spannerInstance.Spec.InstanceConfig, spannerInstance.Spec.NodeCount)
		if err != nil {
			return err
		}
		if len(spannerInstance.Labels) > 0 {
			err = op.UpdateLabels(name, spannerInstance.Labels)
			if err != nil {
				return err
			}
		}
		inst, err = op.GetInstance(name)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if spannerInstance.Spec.NodeCount != inst.NodeCount {
		log.Printf("spannerInstance nodeCount: %d is different from actual instance nodeCount: %d, fit to spannerInstance spec", spannerInstance.Spec.NodeCount, inst.NodeCount)
		err = op.Scale(spannerInstance.Name, spannerInstance.Spec.NodeCount)
		if err != nil {
			return nil
		}
	}

	labels := spannerInstance.DeepCopy().Labels
	dirty := false
	for k, v := range inst.Labels {
		if val, ok := labels[k]; ok {
			if v != val {
				dirty = true
			}
		} else {
			inst.Labels[k] = v
			dirty = true
		}
	}
	if dirty {
		log.Printf("spec labels and actual labels is different, update labels to %+v", labels)
		err = op.UpdateLabels(name, labels)
		if err != nil {
			return err
		}
	}

	// Finally, we update the status block of the SpannerInstance resource to reflect the
	// current state of the world
	err = c.updateSpannerInstanceStatus(spannerInstance)
	if err != nil {
		return err
	}

	c.recorder.Event(spannerInstance, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateSpannerInstanceStatus(spannerInstance *instancev1alpha1.SpannerInstance) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	spannerInstanceCopy := spannerInstance.DeepCopy()
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the SpannerInstance resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.spannerclientset.InstanceadminsV1alpha1().SpannerInstances(spannerInstance.Namespace).Update(spannerInstanceCopy)
	return err
}

// enqueueSpannerInstanceInstance takes a SpannerInstance resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than SpannerInstance.
func (c *Controller) enqueueSpannerInstance(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the SpannerInstance resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that SpannerInstance resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a SpannerInstance, we should not do anything more
		// with it.
		if ownerRef.Kind != "SpannerInstance" {
			return
		}

		spannerInstance, err := c.spannerInstanceLister.SpannerInstances(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of spannerInstance '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueSpannerInstance(spannerInstance)
		return
	}
}
