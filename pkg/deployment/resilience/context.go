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

package resilience

import (
	"context"

	api "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v1alpha"
	"github.com/arangodb/kube-arangodb/pkg/util/arangod"
)

// Context provides methods to the resilience package.
type Context interface {
	// GetSpec returns the current specification of the deployment
	GetSpec() api.DeploymentSpec
	// GetStatus returns the current status of the deployment
	GetStatus() api.DeploymentStatus
	// UpdateStatus replaces the status of the deployment with the given status and
	// updates the resources in k8s.
	UpdateStatus(status api.DeploymentStatus, force ...bool) error
	// GetAgencyClients returns a client connection for every agency member.
	// If the given predicate is not nil, only agents are included where the given predicate returns true.
	GetAgencyClients(ctx context.Context, predicate func(id string) bool) ([]arangod.Agency, error)
}
