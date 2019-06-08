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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:method=GetScale,verb=get,subresource=scale,result=SpannerInstance
// +genclient:method=UpdateScale,verb=update,subresource=scale,input=k8s.io/kubernetes/pkg/apis/autoscaling.Scale,result=SpannerInstance

// SpannerInstance is a specification for a SpannerInstance resource
type SpannerInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpannerInstanceSpec   `json:"spec"`
	Status SpannerInstanceStatus `json:"status"`
}

// SpannerInstanceSpec is the spec for a SpannerInstance resource
type SpannerInstanceSpec struct {
	DisplayName    string `json:"displayName"`
	InstanceConfig string `json:"instanceConfig"`
	NodeCount      int32  `json:"nodeCount"`
}

// SpannerInstanceStatus is the status for a SpannerInstance resource
type SpannerInstanceStatus struct {
	AvailableNodes int32             `json:"availableNodes"`
	InstanceLabels map[string]string `json:"instanceLabels"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpannerInstanceList is a list of SpannerInstance resources
type SpannerInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SpannerInstance `json:"items"`
}
