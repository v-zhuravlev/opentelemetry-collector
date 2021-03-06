// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internaldata

import (
	"testing"

	occommon "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	ocmetrics "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	ocresource "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/internal/data"
	"github.com/open-telemetry/opentelemetry-collector/internal/data/testdata"
)

func TestOCToMetricData(t *testing.T) {
	tests := []struct {
		name     string
		oc       consumerdata.MetricsData
		internal data.MetricData
	}{
		{
			name:     "empty",
			oc:       consumerdata.MetricsData{},
			internal: data.NewMetricData(),
		},

		{
			name: "empty-metrics",
			oc: consumerdata.MetricsData{
				Node:     &occommon.Node{},
				Resource: &ocresource.Resource{},
			},
			internal: testdata.GenerateMetricDataOneEmptyResourceMetrics(),
		},

		{
			name:     "no-libraries",
			oc:       generateOCTestDataNoMetrics(),
			internal: testdata.GenerateMetricDataNoLibraries(),
		},

		{
			name:     "no-points",
			oc:       generateOCTestDataNoPoints(),
			internal: testdata.GenerateMetricDataAllTypesNoDataPoints(),
		},

		{
			name:     "no-labels-metric",
			oc:       generateOCTestDataNoLabels(),
			internal: testdata.GenerateMetricDataOneMetricNoLabels(),
		},

		{
			name: "int64-metric",
			oc: consumerdata.MetricsData{
				Resource: generateOCTestResource(),
				Metrics:  []*ocmetrics.Metric{generateOCTestMetricInt()},
			},
			internal: testdata.GenerateMetricDataOneMetric(),
		},

		{
			name: "int64-and-double-metrics",
			oc: consumerdata.MetricsData{
				Resource: generateOCTestResource(),
				Metrics:  []*ocmetrics.Metric{generateOCTestMetricInt(), generateOCTestMetricDouble()},
			},
			internal: testdata.GenerateMetricDataTwoMetrics(),
		},

		{
			name: "sample-metric",
			oc: consumerdata.MetricsData{
				Resource: generateOCTestResource(),
				Metrics: []*ocmetrics.Metric{
					generateOCTestMetricInt(),
					generateOCTestMetricDouble(),
					generateOCTestMetricHistogram(),
					generateOCTestMetricSummary(),
				},
			},
			internal: testdata.GenerateMetricDataWithCountersHistogramAndSummary(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := OCToMetricData(test.oc)
			assert.EqualValues(t, test.internal, got)
		})
	}
}

func BenchmarkMetricIntOCToInternal(b *testing.B) {
	ocMetric := consumerdata.MetricsData{
		Resource: generateOCTestResource(),
		Metrics: []*ocmetrics.Metric{
			generateOCTestMetricInt(),
			generateOCTestMetricInt(),
			generateOCTestMetricInt(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		OCToMetricData(ocMetric)
	}
}

func BenchmarkMetricDoubleOCToInternal(b *testing.B) {
	ocMetric := consumerdata.MetricsData{
		Resource: generateOCTestResource(),
		Metrics: []*ocmetrics.Metric{
			generateOCTestMetricDouble(),
			generateOCTestMetricDouble(),
			generateOCTestMetricDouble(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		OCToMetricData(ocMetric)
	}
}

func BenchmarkMetricHistogramOCToInternal(b *testing.B) {
	ocMetric := consumerdata.MetricsData{
		Resource: generateOCTestResource(),
		Metrics: []*ocmetrics.Metric{
			generateOCTestMetricHistogram(),
			generateOCTestMetricHistogram(),
			generateOCTestMetricHistogram(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		OCToMetricData(ocMetric)
	}
}

func BenchmarkMetricSummaryOCToInternal(b *testing.B) {
	ocMetric := consumerdata.MetricsData{
		Resource: generateOCTestResource(),
		Metrics: []*ocmetrics.Metric{
			generateOCTestMetricSummary(),
			generateOCTestMetricSummary(),
			generateOCTestMetricSummary(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		OCToMetricData(ocMetric)
	}
}

func generateOCTestResource() *ocresource.Resource {
	return &ocresource.Resource{
		Labels: map[string]string{
			"resource-attr": "resource-attr-val-1",
		},
	}
}
