/*
Copyright 2019 Google LLC

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
	"context"
	"testing"

	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	apisv1alpha1 "knative.dev/pkg/apis/v1alpha1"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
)

var (
	// Bare minimum is Bucket and Sink
	minimalStorageSpec = StorageSpec{
		Bucket: "my-test-bucket",
		SourceSpec: duckv1beta1.SourceSpec{
			Sink: apisv1alpha1.Destination{
				ObjectReference: &corev1.ObjectReference{
					APIVersion: "foo",
					Kind:       "bar",
					Namespace:  "baz",
					Name:       "qux",
				},
			},
		},
	}

	// Bucket, Sink and GCSSecret
	withGCSSecret = StorageSpec{
		Bucket: "my-test-bucket",
		SourceSpec: duckv1beta1.SourceSpec{
			Sink: apisv1alpha1.Destination{
				ObjectReference: &corev1.ObjectReference{
					APIVersion: "foo",
					Kind:       "bar",
					Namespace:  "baz",
					Name:       "qux",
				},
			},
		},
		GCSSecret: corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "secret-name",
			},
			Key: "secret-key",
		},
	}

	// Bucket, Sink, GCSSecret, and PullSubscriptionSecret
	withPullSubscriptionSecret = StorageSpec{
		Bucket: "my-test-bucket",
		SourceSpec: duckv1beta1.SourceSpec{
			Sink: apisv1alpha1.Destination{
				ObjectReference: &corev1.ObjectReference{
					APIVersion: "foo",
					Kind:       "bar",
					Namespace:  "baz",
					Name:       "qux",
				},
			},
		},
		GCSSecret: corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "gcs-secret-name",
			},
			Key: "gcs-secret-key",
		},
		PullSubscriptionSecret: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "pullsubscription-secret-name",
			},
			Key: "pullsubscription-secret-key",
		},
	}
)

func TestValidationFields(t *testing.T) {
	testCases := []struct {
		name string
		s    *Storage
		want *apis.FieldError
	}{{
		name: "empty",
		s:    &Storage{Spec: StorageSpec{}},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("spec.bucket", "spec.sink")
			return fe
		}(),
	}, {
		name: "missing sink",
		s:    &Storage{Spec: StorageSpec{Bucket: "foo"}},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("spec.sink")
			return fe
		}(),
	}}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.Validate(context.TODO())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Errorf("%s: Validate StorageSpec (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestSpecValidationFields(t *testing.T) {
	testCases := []struct {
		name string
		spec *StorageSpec
		want *apis.FieldError
	}{{
		name: "empty",
		spec: &StorageSpec{},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("bucket", "sink")
			return fe
		}(),
	}, {
		name: "missing sink",
		spec: &StorageSpec{Bucket: "foo"},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("sink")
			return fe
		}(),
	}, {
		name: "missing bucket",
		spec: &StorageSpec{
			SourceSpec: duckv1beta1.SourceSpec{
				Sink: apisv1alpha1.Destination{
					ObjectReference: &corev1.ObjectReference{
						APIVersion: "foo",
						Kind:       "bar",
						Namespace:  "baz",
						Name:       "qux",
					},
				},
			},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("bucket")
			return fe
		}(),
	}, {
		name: "invalid gcs secret, missing name",
		spec: &StorageSpec{
			Bucket: "my-test-bucket",
			SourceSpec: duckv1beta1.SourceSpec{
				Sink: apisv1alpha1.Destination{
					ObjectReference: &corev1.ObjectReference{
						APIVersion: "foo",
						Kind:       "bar",
						Namespace:  "baz",
						Name:       "qux",
					},
				},
			},
			GCSSecret: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{},
				Key:                  "secret-test-key",
			},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("gcsSecret.name")
			return fe
		}(),
	}, {
		name: "invalid gcs secret, missing key",
		spec: &StorageSpec{
			Bucket: "my-test-bucket",
			SourceSpec: duckv1beta1.SourceSpec{
				Sink: apisv1alpha1.Destination{
					ObjectReference: &corev1.ObjectReference{
						APIVersion: "foo",
						Kind:       "bar",
						Namespace:  "baz",
						Name:       "qux",
					},
				},
			},
			GCSSecret: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: "gcs-test-secret"},
			},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("gcsSecret.key")
			return fe
		}(),
	}, {
		name: "invalid pullsubscription secret, missing name",
		spec: &StorageSpec{
			Bucket: "my-test-bucket",
			SourceSpec: duckv1beta1.SourceSpec{
				Sink: apisv1alpha1.Destination{
					ObjectReference: &corev1.ObjectReference{
						APIVersion: "foo",
						Kind:       "bar",
						Namespace:  "baz",
						Name:       "qux",
					},
				},
			},
			PullSubscriptionSecret: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{},
				Key:                  "secret-test-key",
			},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("pullSubscriptionSecret.name")
			return fe
		}(),
	}, {
		name: "invalid gcs secret, missing key",
		spec: &StorageSpec{
			Bucket: "my-test-bucket",
			SourceSpec: duckv1beta1.SourceSpec{
				Sink: apisv1alpha1.Destination{
					ObjectReference: &corev1.ObjectReference{
						APIVersion: "foo",
						Kind:       "bar",
						Namespace:  "baz",
						Name:       "qux",
					},
				},
			},
			PullSubscriptionSecret: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: "gcs-test-secret"},
			},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("pullSubscriptionSecret.key")
			return fe
		}(),
	}}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := test.spec.Validate(context.TODO())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Errorf("%s: Validate StorageSpec (-want, +got) = %v", test.name, diff)
			}
		})
	}

}