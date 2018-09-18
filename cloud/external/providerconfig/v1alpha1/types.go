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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExtMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// A list of roles for this Machine to use.
	Roles []MachineRole `json:"roles,omitempty"`

	Zone        string `json:"zone"`
	MachineType string `json:"machineType"`

	// The name of the OS to be installed on the machine.
	OS             string      `json:"os"`
	Disks          []Disk      `json:"disks"`
	CRUDPrimitives *CRUDConfig `json:"crudPrimitives"`
}

// The MachineRole indicates the purpose of the Machine, and will determine
// what software and configuration will be used when provisioning and managing
// the Machine. A single Machine may have more than one role, and the list and
// definitions of supported roles is expected to evolve over time.
//
// Currently, only two roles are supported: Master and Node. In the future, we
// expect user needs to drive the evolution and granularity of these roles,
// with new additions accommodating common cluster patterns, like dedicated
// etcd Machines.
//
//                 +-----------------------+------------------------+
//                 | Master present        | Master absent          |
// +---------------+-----------------------+------------------------|
// | Node present: | Install control plane | Join the cluster as    |
// |               | and be schedulable    | just a node            |
// |---------------+-----------------------+------------------------|
// | Node absent:  | Install control plane | Invalid configuration  |
// |               | and be unschedulable  |                        |
// +---------------+-----------------------+------------------------+
type MachineRole string

const (
	MasterRole MachineRole = "Master"
	NodeRole   MachineRole = "Node"
)

type Disk struct {
	InitializeParams DiskInitializeParams `json:"initializeParams"`
}

type DiskInitializeParams struct {
	DiskSizeGb int64  `json:"diskSizeGb"`
	DiskType   string `json:"diskType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExtClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	Project        string       `json:"project"`
	CRUDPrimitives []CRUDConfig `json:"crudPrimitives"`
}

type CRUDConfig struct {
	metav1.ObjectMeta `json:",inline"`

	// Query that specifies which node(s) this config applies to
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Container that handles machine operations
	Container *v1.Container `json:"container"`

	// Optional command to be used instead of the default when
	// handling machine Create operations (power-on/provisioning)
	CreateCmd []string `json:"createCmd,omitempty"`

	// Optional command to be used instead of the default when
	// handling machine Delete operations (power-off/deprovisioning)
	DeleteCmd []string `json:"deleteCmd,omitempty"`

	// Optional command to be used instead of the default when
	// handling machine Reboot operations
	RebootCmd []string `json:"rebootCmd,omitempty"`

	// How Secrets and DynamicConfig should be passed to the
	// container: ([env], cli)
	ArgumentFormat string `json:"argumentFormat"`

	// Parameters to use for automatic variables
	PassActionAs string `json:"passActionAs,omitempty"`
	PassTargetAs string `json:"passTargetAs,omitempty"`

	// Additional parameters that may be passed as either name/value pairs or
	// "--name value" depending on the value of ArgumentFormat
	Config map[string]string `json:"config"`

	// Parameters who’s value changes depending on the affected node
	DynamicConfig []DynamicConfigElement `json:"dynamicConfig"`

	// A list of Kubernetes secrets to securely pass to the container
	Secrets map[string]string `json:"secrets"`

	// How long to wait for the Job to complete
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`

	// How long to wait before retrying failed Jobs
	RetrySeconds *int32 `json:"retrySeconds,omitempty"`

	// How long to wait before retrying failed Jobs
	Retries *int32 `json:"retries,omitempty"`
}

type DynamicConfigElement struct {
	Field   string            `json:"field"`
	Default string            `json:"default"`
	Values  map[string]string `json:"values"`
}

func (dc *DynamicConfigElement) Lookup(key string) (string, bool) {
	if val, ok := dc.Values[key]; ok {
		return val, true
	} else if len(dc.Default) > 0 {
		return dc.Default, true
	}
	return "", false
}
