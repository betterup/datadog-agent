// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver
// +build kubeapiserver

package kubernetesapiserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

func TestUnbundledEventsTransform(t *testing.T) {
	ts := metav1.Time{Time: time.Now()}
	pod := v1.ObjectReference{
		UID:       "foobar",
		Kind:      "Pod",
		Namespace: "default",
		Name:      "redis",
	}

	tests := []struct {
		name     string
		event    *v1.Event
		expected []metrics.Event
	}{
		{
			name: "event is filtered out",
			event: &v1.Event{
				InvolvedObject: pod,
				Type:           "Info",
				Reason:         "SandboxChanged",
				Message:        "Pod sandbox changed, it will be killed and re-created.",
				Source: v1.EventSource{
					Component: "kubelet",
					Host:      "test-host",
				},
				FirstTimestamp: ts,
				LastTimestamp:  ts,
				Count:          1,
			},
			expected: nil,
		},
		{
			name: "event is collected",
			event: &v1.Event{
				InvolvedObject: pod,
				Type:           "Warning",
				Reason:         "Failed",
				Message:        "All containers terminated",
				Source: v1.EventSource{
					Component: "kubelet",
					Host:      "test-host",
				},
				FirstTimestamp: ts,
				LastTimestamp:  ts,
				Count:          1,
			},
			expected: []metrics.Event{
				{
					Title:    "Pod default/redis: Failed",
					Text:     "All containers terminated",
					Ts:       ts.Time.Unix(),
					Priority: metrics.EventPriorityNormal,
					Host:     "test-host-test-cluster",
					Tags: []string{
						"kube_kind:Pod",
						"kube_name:redis",
						"kubernetes_kind:Pod",
						"name:redis",
						"kube_namespace:default",
						"namespace:default",
						"pod_name:redis",
						"source_component:kubelet",
						"event_reason:Failed",
					},
					AlertType:      metrics.EventAlertTypeWarning,
					AggregationKey: "kubernetes_apiserver:foobar",
					SourceTypeName: "kubernetes",
					EventType:      "kubernetes_apiserver",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectedTypes := map[string][]string{
				"pod": {"Failed"},
			}
			transformer := newUnbundledTransformer("test-cluster", collectedTypes)

			events, errors := transformer.Transform([]*v1.Event{tt.event})

			assert.Empty(t, errors)
			assert.Equal(t, tt.expected, events)
		})
	}
}