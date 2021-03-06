/*
Copyright 2015 Home Office All rights reserved.

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
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	envRegex = regexp.MustCompile("%[[:alnum:]]+%")
)

// VaultResources is a collection of type resource
type VaultResources struct {
	// an array of resource to retrieve
	items []*VaultResource
}

// Set is the implementation for the parser
// secret:test:file=filename.test,fmt=yaml
func (r *VaultResources) Set(value string) error {
	rn := defaultVaultResource()

	// step: split on the ':'
	items := strings.Split(value, ":")
	if len(items) < 2 {
		return fmt.Errorf("invalid resource, must have at least two sections TYPE:PATH")
	}
	if len(items) > 3 {
		return fmt.Errorf("invalid resource, can only has three sections, TYPE:PATH[:OPTIONS]")
	}
	if items[0] == "" || items[1] == "" {
		return fmt.Errorf("invalid resource, neither type or path can be empty")
	}

	// step: look for any token in the resource
	tokens := envRegex.FindAllStringSubmatch(items[1], -1)
	if len(tokens) > 0 {
		for _, x := range tokens {
			// step: replace the token with the environment variable
			items[1] = strings.Replace(items[1], x[0], os.Getenv(strings.Replace(x[0], "%", "", -1)), -1)
		}
	}

	// step: extract the elements
	rn.resource = items[0]
	rn.path = items[1]
	rn.options = make(map[string]string, 0)

	// step: extract any options
	if len(items) > 2 {
		for _, x := range strings.Split(items[2], ",") {
			kp := strings.Split(x, "=")
			if len(kp) != 2 {
				return fmt.Errorf("invalid resource option: %s, must be KEY=VALUE", x)
			}
			if kp[1] == "" {
				return fmt.Errorf("invalid resource option: %s, must have a value", x)
			}
			// step: set the name and value
			name := kp[0]
			value := strings.Replace(kp[1], "|", ",", -1)

			// step: extract the control options from the path resource parameters
			switch name {
			case optionFormat:
				if matched := resourceFormatRegex.MatchString(value); !matched {
					return fmt.Errorf("unsupported output format: %s", value)
				}
				rn.format = value
			case optionUpdate:
				duration, err := time.ParseDuration(value)
				if err != nil {
					return fmt.Errorf("update option: %s is not value, should be a duration format", value)
				}
				rn.update = duration
			case optionRevoke:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the revoke option: %s is invalid, should be a boolean", value)
				}
				rn.revoked = choice
			case optionsRevokeDelay:
				duration, err := time.ParseDuration(value)
				if err != nil {
					return fmt.Errorf("the revoke delay option: %s is not value, should be a duration format", value)
				}
				rn.revokeDelay = duration
			case optionRenewal:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the renewal option: %s is invalid, should be a boolean", value)
				}
				rn.renewable = choice
			case optionCreate:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the create option: %s is invalid, should be a boolean", value)
				}
				if rn.resource != "secret" {
					return fmt.Errorf("the create option is only supported for 'cn=secret' at this time")
				}
				rn.create = choice
			case optionSize:
				size, err := strconv.ParseInt(value, 10, 16)
				if err != nil {
					return fmt.Errorf("the size option: %s is invalid, should be an integer", value)
				}
				rn.size = size
			case optionExec:
				rn.execPath = value
			case optionFilename:
				rn.filename = value
			case optionTemplatePath:
				rn.templateFile = value
			default:
				rn.options[name] = value
			}
		}
	}
	// step: append to the list of resources
	r.items = append(r.items, rn)

	return nil
}

// String returns a string representation of the struct
func (r VaultResources) String() string {
	return ""
}
