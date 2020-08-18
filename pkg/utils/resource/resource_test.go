/*
Copyright 2020 Google LLC

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

package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestResourceRequirementsBuilder(t *testing.T) {
	tests := []struct {
		name          string
		cpuRequest    string
		cpuLimit      string
		memoryRequest string
		memoryLimit   string
		expected      *corev1.ResourceRequirements
	}{
		{
			name:          "All specified requirements are correctly represented",
			cpuRequest:    "1500m",
			cpuLimit:      "1500m",
			memoryRequest: "500Mi",
			memoryLimit:   "3000Mi",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("3000Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name:          "CPU request requirement is dropped when not specified",
			cpuRequest:    "",
			cpuLimit:      "1500m",
			memoryRequest: "500Mi",
			memoryLimit:   "3000Mi",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("3000Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name:          "CPU limit requirement is dropped when not specified",
			cpuRequest:    "1500m",
			cpuLimit:      "",
			memoryRequest: "500Mi",
			memoryLimit:   "3000Mi",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("3000Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name:          "Memory request requirement is dropped when not specified",
			cpuRequest:    "1500m",
			cpuLimit:      "1500m",
			memoryRequest: "",
			memoryLimit:   "3000Mi",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("3000Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU: resource.MustParse("1500m"),
				},
			},
		},
		{
			name:          "Memory limit requirement is dropped when not specified",
			cpuRequest:    "1500m",
			cpuLimit:      "1500m",
			memoryRequest: "500Mi",
			memoryLimit:   "",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU: resource.MustParse("1500m"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name:          "CPU limit requirement is dropped when the value is not valid",
			cpuRequest:    "1500m",
			cpuLimit:      "invalid",
			memoryRequest: "500Mi",
			memoryLimit:   "3000Mi",
			expected: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("3000Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
			},
		},
		{
			name:          "Clean requirements collection is passed when all values are empty or invalid",
			cpuRequest:    "",
			cpuLimit:      "",
			memoryRequest: "invalid",
			memoryLimit:   "",
			expected: &corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{},
				Requests: corev1.ResourceList{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BuildResourceRequirements(test.cpuRequest, test.cpuLimit, test.memoryRequest, test.memoryLimit)

			if diff := cmp.Diff(test.expected, &result); diff != "" {
				t.Errorf("failed to get expected (-want, +got) = %v", diff)
			}
		})
	}
}

func TestMultiplyQuantity(t *testing.T) {
	tests := []struct {
		name     string
		initial  resource.Quantity
		multiple float64
		expected resource.Quantity
	}{
		{
			name:     "A quantity in the format 'Mi' is multiplied correctly",
			initial:  resource.MustParse("1000Mi"),
			multiple: float64(2),
			expected: resource.MustParse("2000Mi"),
		},
		{
			name:     "A quantity in the format 'Gi' is multiplied correctly",
			initial:  resource.MustParse("100Gi"),
			multiple: float64(2),
			expected: resource.MustParse("200Gi"),
		},
		{
			name:     "Raw quantity values are supported",
			initial:  resource.MustParse("256000"),
			multiple: float64(2),
			expected: resource.MustParse("512000"),
		},
		{
			name:     "Fractional multiplications are supported",
			initial:  resource.MustParse("2000Mi"),
			multiple: float64(0.5),
			expected: resource.MustParse("1000Mi"),
		},
		{
			name:     "Large quantities are multiplied correctly",
			initial:  resource.MustParse("1000000000Gi"),
			multiple: float64(2),
			expected: resource.MustParse("2000000000Gi"),
		},
		{
			name:     "The result preserves the original format",
			initial:  resource.MustParse("2Mi"),
			multiple: float64(1000000000),
			expected: resource.MustParse("2000000000Mi"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MultiplyQuantity(test.initial, test.multiple)

			if diff := cmp.Diff(&test.expected, result); diff != "" {
				t.Errorf("failed to get expected (-want, +got) = %v", diff)
			}
		})
	}
}
