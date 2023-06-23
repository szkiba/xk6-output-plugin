#!/bin/env python

import datetime
import logging

from xk6_output_plugin_py.output import serve, Output, Info, MetricType, ValueType


class Example(Output):
    def Init(self, params):
        logging.info("init")

        return Info(description="example-py plugin")

    def Start(self):
        logging.info("start")

    def Stop(self):
        logging.info("stop")

    def AddMetrics(self, metrics):
        logging.info("metrics")
        for metric in metrics:
            logging.info(
                metric.name,
                extra={
                    "metric.type": MetricType.Name(metric.type),
                    "metric.contains": ValueType.Name(metric.contains),
                },
            )

    def AddSamples(self, samples):
        logging.info("samples")
        for sample in samples:
            t = datetime.datetime.fromtimestamp(
                sample.time / 1000.0, tz=datetime.timezone.utc
            )

            logging.info(
                sample.metric,
                extra={"sample.time": t, "sample.value": sample.value},
            )


if __name__ == "__main__":
    serve(Example())
