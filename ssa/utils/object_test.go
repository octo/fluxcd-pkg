/*
Copyright 2023 The Flux authors

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

package utils

import (
	"strings"
	"testing"
)

func TestReadObjects_DropsInvalid(t *testing.T) {
	testCases := []struct {
		name      string
		resources string
		expected  int
	}{
		{
			name: "valid resources",
			resources: `
---
apiVersion: v1
kind: Secret
metadata:
  name: test
  namespace: default
immutable: true
stringData:
  key: "private-key"
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: crossplane-provider-aws2
spec:
  package: crossplane/provider-aws:v0.23.0
  controllerConfigRef:
    name: provider-aws
`,
			expected: 2,
		},
		{
			name: "some invalid resources",
			resources: `
---
piVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: crossplane-provider-aws1
spec:
  package: crossplane/provider-aws:v0.23.0
  controllerConfigRef:
    name: provider-aws
---
apiVersion: v1
kind: Secret
metadata:
  name: test
  namespace: default
immutable: true
stringData:
  key: "private-key"
`,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			objects, err := ReadObjects(strings.NewReader(tc.resources))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(objects) != tc.expected {
				t.Errorf("unexpected number of objects in %v", objects)
			}
		})
	}
}
