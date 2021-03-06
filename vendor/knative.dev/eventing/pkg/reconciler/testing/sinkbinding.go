/*
Copyright 2020 The Knative Authors

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

package testing

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sourcesv1alpha2 "knative.dev/eventing/pkg/apis/sources/v1alpha2"

	legacysourcesv1alpha1 "knative.dev/eventing/pkg/apis/legacysources/v1alpha1"
	sourcesv1alpha1 "knative.dev/eventing/pkg/apis/sources/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"
)

// SinkBindingV1Alpha1Option enables further configuration of a SinkBinding.
type SinkBindingV1Alpha1Option func(*sourcesv1alpha1.SinkBinding)

// SinkBindingV1Alpha2Option enables further configuration of a SinkBinding.
type SinkBindingV1Alpha2Option func(*sourcesv1alpha2.SinkBinding)

// LegacySinkBindingOption enables further configuration of a SinkBinding.
type LegacySinkBindingOption func(*legacysourcesv1alpha1.SinkBinding)

// NewLegacySinkBinding creates a SinkBinding with SinkBindingOptions
func NewLegacySinkBinding(name, namespace string, o ...LegacySinkBindingOption) *legacysourcesv1alpha1.SinkBinding {
	c := &legacysourcesv1alpha1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range o {
		opt(c)
	}
	//c.SetDefaults(context.Background()) // TODO: We should add defaults and validation.
	return c
}

// NewSinkBindingV1Alpha1 creates a SinkBinding with SinkBindingOptions
func NewSinkBindingV1Alpha1(name, namespace string, o ...SinkBindingV1Alpha1Option) *sourcesv1alpha1.SinkBinding {
	c := &sourcesv1alpha1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range o {
		opt(c)
	}
	//c.SetDefaults(context.Background()) // TODO: We should add defaults and validation.
	return c
}

// NewSinkBindingV1Alpha2 creates a SinkBinding with SinkBindingOptions
func NewSinkBindingV1Alpha2(name, namespace string, o ...SinkBindingV1Alpha2Option) *sourcesv1alpha2.SinkBinding {
	c := &sourcesv1alpha2.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range o {
		opt(c)
	}
	//c.SetDefaults(context.Background()) // TODO: We should add defaults and validation.
	return c
}

// WithLegacySubject assigns the subject of the SinkBinding.
func WithLegacySubject(subject tracker.Reference) LegacySinkBindingOption {
	return func(sb *legacysourcesv1alpha1.SinkBinding) {
		sb.Spec.Subject = subject
	}
}

// WithLegacySink assigns the sink of the SinkBinding.
func WithLegacySink(sink duckv1.Destination) LegacySinkBindingOption {
	return func(sb *legacysourcesv1alpha1.SinkBinding) {
		sb.Spec.Sink = sink
	}
}

// WithSubjectV1A1 assigns the subject of the SinkBinding.
func WithSubjectV1A1(subject tracker.Reference) SinkBindingV1Alpha1Option {
	return func(sb *sourcesv1alpha1.SinkBinding) {
		sb.Spec.Subject = subject
	}
}

// WithSinkV1A1 assigns the sink of the SinkBinding.
func WithSinkV1A1(sink duckv1.Destination) SinkBindingV1Alpha1Option {
	return func(sb *sourcesv1alpha1.SinkBinding) {
		sb.Spec.Sink = sink
	}
}

// WithCloudEventOverridesV1A1 assigns the CloudEventsOverrides of the SinkBinding.
func WithCloudEventOverridesV1A1(overrides duckv1.CloudEventOverrides) SinkBindingV1Alpha1Option {
	return func(sb *sourcesv1alpha1.SinkBinding) {
		sb.Spec.CloudEventOverrides = &overrides
	}
}

// WithSubjectV1A2 assigns the subject of the SinkBinding.
func WithSubjectV1A2(subject tracker.Reference) SinkBindingV1Alpha2Option {
	return func(sb *sourcesv1alpha2.SinkBinding) {
		sb.Spec.Subject = subject
	}
}

// WithSinkV1A2 assigns the sink of the SinkBinding.
func WithSinkV1A2(sink duckv1.Destination) SinkBindingV1Alpha2Option {
	return func(sb *sourcesv1alpha2.SinkBinding) {
		sb.Spec.Sink = sink
	}
}

// WithCloudEventOverridesV1A1 assigns the CloudEventsOverrides of the SinkBinding.
func WithCloudEventOverridesV1A2(overrides duckv1.CloudEventOverrides) SinkBindingV1Alpha2Option {
	return func(sb *sourcesv1alpha2.SinkBinding) {
		sb.Spec.CloudEventOverrides = &overrides
	}
}
