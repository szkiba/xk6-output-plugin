// SPDX-FileCopyrightText: 2021 - 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package extension

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/szkiba/xk6-output-plugin-go/output"
	"go.k6.io/k6/metrics"
	k6output "go.k6.io/k6/output"
)

const (
	defaultBuffering = 1000
	minBuffering     = 200
)

type Extension struct {
	buffer *k6output.SampleBuffer

	flusher *k6output.PeriodicFlusher
	logger  logrus.FieldLogger

	output output.Output
	info   *output.Info

	seenMetrics map[string]struct{}
	thresholds  map[string]metrics.Thresholds
}

var _ k6output.Output = (*Extension)(nil)

func New(params k6output.Params) (*Extension, error) { //nolint:ireturn
	out, err := newPlugin(params.ConfigArgument, params.Logger)
	if err != nil {
		return nil, err
	}

	conf, err := out.Init(context.Background(), &output.Params{Environment: params.Environment})
	if err != nil {
		return nil, err
	}

	ext := &Extension{
		logger:      params.Logger,
		buffer:      nil,
		flusher:     nil,
		thresholds:  nil,
		info:        conf,
		output:      out,
		seenMetrics: make(map[string]struct{}),
	}

	return ext, nil
}

func (ext *Extension) Description() string {
	return ext.info.Description
}

func (ext *Extension) Start() error {
	var err error

	ext.buffer = new(k6output.SampleBuffer)

	period := ext.info.Buffering
	if period == 0 {
		period = defaultBuffering
	}

	if period < minBuffering {
		ext.logger.Warnf("The requested buffering period (%dms) is too small, the minimum (%dms) will be used instead", period, minBuffering)
	}

	ext.flusher, err = k6output.NewPeriodicFlusher(time.Millisecond*time.Duration(period), ext.flush)
	if err != nil {
		return err
	}

	return ext.output.Start(context.Background())
}

func (ext *Extension) Stop() error {
	ext.flusher.Stop()

	return ext.output.Stop(context.Background())
}

func (ext *Extension) SetThresholds(thresholds map[string]metrics.Thresholds) {
	if len(thresholds) == 0 {
		return
	}

	ext.thresholds = make(map[string]metrics.Thresholds, len(thresholds))
	for name, t := range thresholds {
		ext.thresholds[name] = t
	}
}

func (ext *Extension) AddMetricSamples(samples []metrics.SampleContainer) {
	ext.buffer.AddMetricSamples(samples)
}

func (ext *Extension) flush() {
	buffered := ext.buffer.GetBufferedSamples()
	allSample := make([]*output.Sample, 0, len(buffered))
	allMetric := make([]*output.Metric, 0)

	for _, sc := range buffered {
		sc := sc
		for _, sample := range sc.GetSamples() {
			sample := sample

			if _, ok := ext.seenMetrics[sample.Metric.Name]; !ok {
				var trs *metrics.Thresholds

				if t, ok := ext.thresholds[sample.Metric.Name]; ok {
					trs = &t
				}

				allMetric = append(allMetric, mapMetric(sample.Metric, trs))
				ext.seenMetrics[sample.Metric.Name] = struct{}{}
			}

			allSample = append(allSample, mapSample(&sample))
		}
	}

	err := ext.output.AddMetrics(context.Background(), allMetric)
	if err != nil {
		ext.logger.WithError(err).Warn()
	}

	err = ext.output.AddSamples(context.Background(), allSample)
	if err != nil {
		ext.logger.WithError(err).Warn()
	}
}
