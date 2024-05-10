/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 * SPDX-License-Identifier: MIT
 */

package metrics

import (
	"time"

	"github.com/eclipse-symphony/symphony/api/constants"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/observability"
)

const (
	ValidateRuleOperation             string = "ValidateRule"
	ApplyScriptOperation              string = "ApplyScript"
	ApplyOperation                    string = "Apply"
	ApplyYamlOperation                string = "ApplyYaml"
	ApplyCustomResource               string = "ApplyCustomResource"
	ReceiveDataChannelOperation       string = "ReceiveFromDataChannel"
	ReceiveErrorChannelOperation      string = "ReceiveFromErrorChannel"
	ConvertResourceDataBytesOperation string = "ConvertResourceDataToBytes"
	ObjectOperation                   string = "Object"
	ResourceOperation                 string = "Resource"
	PullChartOperation                string = "PullChart"
	LoadChartOperation                string = "LoadChart"
	HelmChartOperation                string = "HelmChart"
	HelmActionConfigOperation         string = "HelmActionConfig"
	HelmPropertiesOperation           string = "HelmProperties"

	GetOperationType    string = "Get"
	CreateOperationType string = "Create"
	UpdateOperationType string = "Update"
	DeleteOperationType string = "Delete"
)

// Metrics is a metrics tracker for a provider operation.
type Metrics struct {
	providerOperationLatency observability.Histogram
	providerOperationErrors  observability.Counter
}

func New() (*Metrics, error) {
	observable := observability.New(constants.API)

	providerOperationLatency, err := observable.Metrics.Histogram(
		"symphony_provider_operation_latency",
		"measure of overall latency for provider operation side",
	)
	if err != nil {
		return nil, err
	}

	providerOperationErrors, err := observable.Metrics.Counter(
		"symphony_provider_operation_errors",
		"count of errors in provider operation side",
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		providerOperationLatency: providerOperationLatency,
		providerOperationErrors:  providerOperationErrors,
	}, nil
}

// Close closes all metrics.
func (m *Metrics) Close() {
	if m == nil {
		return
	}

	m.providerOperationErrors.Close()
}

// ProviderOperationLatency tracks the overall provider target latency.
func (m *Metrics) ProviderOperationLatency(
	startTime time.Time,
	providerType string,
	operation string,
	operationType string,
	functionName string,
) {
	if m == nil {
		return
	}

	m.providerOperationLatency.Add(
		latency(startTime),
		Target(
			providerType,
			functionName,
			operation,
			operationType,
			v1alpha2.OK.String(),
		),
	)
}

// ProviderOperationErrors increments the count of errors for provider target.
func (m *Metrics) ProviderOperationErrors(
	providerType string,
	functionName string,
	operation string,
	operationType string,
	errorCode string,
) {
	if m == nil {
		return
	}

	m.providerOperationErrors.Add(
		1,
		Target(
			providerType,
			functionName,
			operation,
			operationType,
			errorCode,
		),
	)
}

// Latency gets the time since the given start in milliseconds.
func latency(start time.Time) float64 {
	return float64(time.Since(start)) / float64(time.Millisecond)
}
