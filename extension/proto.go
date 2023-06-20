// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package extension

import (
	"github.com/szkiba/xk6-output-plugin-go/output"
	"go.k6.io/k6/metrics"
)

func mapMetric(metric *metrics.Metric, thresholds *metrics.Thresholds) *output.Metric {
	met := new(output.Metric)
	met.Name = metric.Name
	met.Type = mapMetricType(metric.Type)
	met.Contains = mapValueType(metric.Contains)
	met.Tainted = metric.Tainted.Bool

	if thresholds != nil {
		met.Thresholds = make([]string, 0, len(thresholds.Thresholds))

		for _, t := range thresholds.Thresholds {
			met.Thresholds = append(met.Thresholds, t.Source)
		}
	}

	for _, sub := range metric.Submetrics {
		smet := new(output.Submetric)
		smet.Name = sub.Name
		smet.Suffix = sub.Suffix
		smet.Tags = sub.Tags.Map()
		smet.Metric = mapMetric(sub.Metric, thresholds)

		met.Submetrics = append(met.Submetrics, smet)
	}

	return met
}

func mapSample(sample *metrics.Sample) *output.Sample {
	sam := new(output.Sample)

	sam.Metric = sample.Metric.Name
	sam.Time = sample.Time.UnixMilli()
	sam.Value = sample.Value
	sam.Metadata = sample.Metadata

	if sample.Tags != nil {
		sam.Tags = sample.Tags.Map()
	}

	return sam
}

// the order is depends on k6's source...
var metricTypes = []output.MetricType{
	output.MetricType_COUNTER,
	output.MetricType_GAUGE,
	output.MetricType_TREND,
	output.MetricType_RATE,
}

func mapMetricType(metricType metrics.MetricType) output.MetricType {
	idx := int(metricType)
	if idx >= len(metricTypes) || idx < 0 {
		return output.MetricType_METRIC_TYPE_UNSPECIFIED
	}

	return metricTypes[idx]
}

// the order is depends on k6's source...
var valueTypes = []output.ValueType{
	output.ValueType_DEFAULT,
	output.ValueType_TIME,
	output.ValueType_DATA,
}

func mapValueType(valueType metrics.ValueType) output.ValueType {
	idx := int(valueType)
	if idx >= len(valueTypes) || idx < 0 {
		return output.ValueType_VALUE_TYPE_UNSPECIFIED
	}

	return valueTypes[idx]
}
