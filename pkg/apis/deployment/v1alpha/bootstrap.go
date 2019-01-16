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

package v1alpha

import (
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil"
)

const (
	// UserNameRoot root user name
	UserNameRoot = "root"
)

// PasswordSecretName contains user password secret name
type PasswordSecretName string

const (
	// PasswordSecretNameNone is magic value for no action
	PasswordSecretNameNone PasswordSecretName = "None"
	// PasswordSecretNameAuto is magic value for autogenerate name
	PasswordSecretNameAuto PasswordSecretName = "Auto"
)

// PasswordSecretNameList is a map from username to secretnames
type PasswordSecretNameList map[string]PasswordSecretName

// BootstrapSpec contains information for cluster bootstrapping
type BootstrapSpec struct {
	// PasswordSecretNames contains a map of username to password-secret-name
	PasswordSecretNames PasswordSecretNameList `json:"passwordSecretNames,omitempty"`
}

// IsNone returns ture if p is empty or None
func (p PasswordSecretName) IsNone() bool {
	return p == PasswordSecretNameNone || p == ""
}

// IsAuto returns ture if p is Auto
func (p PasswordSecretName) IsAuto() bool {
	return p == PasswordSecretNameAuto
}

// Validate validates the password secret name
func (p PasswordSecretName) Validate() error {
	if !p.IsNone() {
		if err := k8sutil.ValidateResourceName(string(p)); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

// GetSecretName returns the secret name given by the specs. Or None if not set.
func (s PasswordSecretNameList) GetSecretName(user string) PasswordSecretName {
	if s != nil {
		if secretname, ok := s[user]; ok {
			return secretname
		}
	}
	return PasswordSecretNameNone
}

// getSecretNameForUserPassword returns the default secret name for the given user
func getSecretNameForUserPassword(deploymentname, username string) PasswordSecretName {
	return PasswordSecretName(deploymentname + "-" + username + "-password")
}

// Validate the specification.
func (b *BootstrapSpec) Validate() error {
	for _, secretname := range b.PasswordSecretNames {
		if err := secretname.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// SetDefaults fills in default values when a field is not specified.
func (b *BootstrapSpec) SetDefaults(deploymentname string) {
	if b.PasswordSecretNames == nil {
		b.PasswordSecretNames = PasswordSecretNameList{
			UserNameRoot: getSecretNameForUserPassword(deploymentname, UserNameRoot),
		}
	} else {
		// Check if root is specified
		if secretname, ok := b.PasswordSecretNames[UserNameRoot]; ok {
			if secretname.IsAuto() {
				b.PasswordSecretNames[UserNameRoot] = getSecretNameForUserPassword(deploymentname, UserNameRoot)
			}
		} else {
			// implicit default
			b.PasswordSecretNames[UserNameRoot] = getSecretNameForUserPassword(deploymentname, UserNameRoot)
		}

		// Now fill in values for all users
		for user, secretname := range b.PasswordSecretNames {
			if user != UserNameRoot {
				if secretname.IsAuto() {
					b.PasswordSecretNames[user] = getSecretNameForUserPassword(deploymentname, user)
				}
			}
		}
	}
}

// NewPasswordSecretNameListOrNil returns nil if input is nil, otherwise returns a clone of the given value.
func NewPasswordSecretNameListOrNil(list PasswordSecretNameList) PasswordSecretNameList {
	if list == nil {
		return nil
	}
	var newList = make(PasswordSecretNameList)
	for k, v := range list {
		newList[k] = v
	}
	return newList
}

// SetDefaultsFrom fills unspecified fields with a value from given source spec.
func (b *BootstrapSpec) SetDefaultsFrom(source BootstrapSpec) {
	if b.PasswordSecretNames == nil {
		b.PasswordSecretNames = NewPasswordSecretNameListOrNil(source.PasswordSecretNames)
	}
}
