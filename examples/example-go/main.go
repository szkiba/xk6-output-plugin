package main

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/szkiba/xk6-output-plugin-go/output"
)

type example struct{}

func (e *example) Init(ctx context.Context, params *output.Params) (*output.Info, error) {
	hclog.L().Info("init")

	return &output.Info{Description: "example-go plugin"}, nil // nolint:exhaustruct
}

func (e *example) Start(ctx context.Context) error {
	hclog.L().Info("start")

	return nil
}

func (e *example) Stop(ctx context.Context) error {
	hclog.L().Info("stop")

	return nil
}

func (e *example) AddMetrics(ctx context.Context, metrics []*output.Metric) error {
	hclog.L().Info("metrics")

	for _, metric := range metrics {
		hclog.L().Info(metric.Name,
			"metric.type", metric.Type.String(),
			"metric.contains", metric.Contains.String(),
		)
	}

	return nil
}

func (e *example) AddSamples(ctx context.Context, samples []*output.Sample) error {
	hclog.L().Info("samples")

	for _, sample := range samples {
		hclog.L().Info(sample.Metric,
			"sample.time", time.UnixMilli(sample.Time).Format(time.RFC3339),
			"sample.value", sample.Value,
		)
	}

	return nil
}

func main() {
	output.Serve(new(example))
}
