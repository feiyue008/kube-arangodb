//
// DISCLAIMER
//
// Copyright 2018 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package v1alpha

import (
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/arangodb/kube-arangodb/pkg/util"
)

// ServerGroupSpec contains the specification for all servers in a specific group (e.g. all agents)
type ServerGroupSpec struct {
	// Count holds the requested number of servers
	Count *int `json:"count,omitempty"`
	// Args holds additional commandline arguments
	Args []string `json:"args,omitempty"`
	// StorageClassName specifies the classname for storage of the servers.
	StorageClassName *string `json:"storageClassName,omitempty"`
	// Resources holds resource requests & limits
	Resources v1.ResourceRequirements `json:"resource,omitempty"`
}

// GetCount returns the value of count.
func (s ServerGroupSpec) GetCount() int {
	if s.Count == nil {
		return 0
	}
	return *s.Count
}

// GetStorageClassName returns the value of count.
func (s ServerGroupSpec) GetStorageClassName() string {
	if s.StorageClassName == nil {
		return ""
	}
	return *s.StorageClassName
}

// Validate the given group spec
func (s ServerGroupSpec) Validate(group ServerGroup, used bool, mode DeploymentMode, env Environment) error {
	if used {
		minCount := 1
		if env == EnvironmentProduction {
			switch group {
			case ServerGroupSingle:
				if mode == DeploymentModeResilientSingle {
					minCount = 2
				}
			case ServerGroupAgents:
				minCount = 3
			case ServerGroupDBServers, ServerGroupCoordinators, ServerGroupSyncMasters, ServerGroupSyncWorkers:
				minCount = 2
			}
		}
		if s.GetCount() < minCount {
			return maskAny(errors.Wrapf(ValidationError, "Invalid count value %d. Expected >= %d", s.Count, minCount))
		}
		if s.GetCount() > 1 && group == ServerGroupSingle && mode == DeploymentModeSingle {
			return maskAny(errors.Wrapf(ValidationError, "Invalid count value %d. Expected 1", s.Count))
		}
	} else if s.GetCount() != 0 {
		return maskAny(errors.Wrapf(ValidationError, "Invalid count value %d for un-used group. Expected 0", s.Count))
	}
	return nil
}

// SetDefaults fills in missing defaults
func (s *ServerGroupSpec) SetDefaults(group ServerGroup, used bool, mode DeploymentMode) {
	if s.GetCount() == 0 && used {
		switch group {
		case ServerGroupSingle:
			if mode == DeploymentModeSingle {
				s.Count = util.Int(1) // Single server
			} else {
				s.Count = util.Int(2) // Resilient single
			}
		default:
			s.Count = util.Int(3)
		}
	}
	if _, found := s.Resources.Requests[v1.ResourceStorage]; !found {
		switch group {
		case ServerGroupSingle, ServerGroupAgents, ServerGroupDBServers:
			if s.Resources.Requests == nil {
				s.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
			}
			s.Resources.Requests[v1.ResourceStorage] = resource.MustParse("8Gi")
		}
	}
}

// SetDefaultsFrom fills unspecified fields with a value from given source spec.
func (s *ServerGroupSpec) SetDefaultsFrom(source ServerGroupSpec) {
	if s.Count == nil {
		s.Count = util.IntOrNil(source.Count)
	}
	if s.Args == nil {
		s.Args = source.Args
	}
	if s.StorageClassName == nil {
		s.StorageClassName = util.StringOrNil(source.StorageClassName)
	}
	// TODO Resources
}

// ResetImmutableFields replaces all immutable fields in the given target with values from the source spec.
// It returns a list of fields that have been reset.
func (s ServerGroupSpec) ResetImmutableFields(group ServerGroup, fieldPrefix string, target *ServerGroupSpec) []string {
	var resetFields []string
	if group == ServerGroupAgents {
		if s.GetCount() != target.GetCount() {
			target.Count = util.IntOrNil(s.Count)
			resetFields = append(resetFields, fieldPrefix+".count")
		}
	}
	if s.GetStorageClassName() != target.GetStorageClassName() {
		target.StorageClassName = util.StringOrNil(s.StorageClassName)
		resetFields = append(resetFields, fieldPrefix+".storageClassName")
	}
	return resetFields
}
