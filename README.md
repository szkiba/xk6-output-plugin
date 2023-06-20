# xk6-output-plugin

Write k6 output extension as dynamically loadable plugin

The xk6-output-plugin is a [k6 output extension](https://k6.io/docs/extensions/get-started/create/output-extensions/) that allows the use of dynamically loaded output plugins.

k6 provides many options for [output management](https://k6.io/docs/results-output/overview/). In addition, [output extensions](https://k6.io/docs/extensions/get-started/create/output-extensions/) can be made to meet individual needs. However, a custom k6 build is required to use the output extensions, and the extensions can only be created in the go programming language. The xk6-output-plugin makes custom output handling simpler and more convenient.

**Features**

- output plugins can be created using your favorite programming language (e.g. go, JavaScript, Python)
- output plugins can be created and used without rebuilding the k6 executable
- output plugins do not increase the size of the k6 executable
- the output plugins can be distributed independently of the k6 binary

## Why not JSON?

A similar approach to the output plugin can also be achieved by processing the output generated by the k6 [json output extension](https://k6.io/docs/results-output/real-time/json/) by an external program.

Why use the output plugin instead of processing JSON output?

- type safe API
- simpler callback-based processing
- logging to k6 log output
- real time result processing

## How It Works?

The plugin is a local [gRPC](https://grpc.io/) service that the xk6-output-plugin starts automatically during k6 startup and stops before k6 stops.
 
Output plugins are managed using the [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin).

## Performance

Metric samples are buffered by the xk6-output-plugin for one second (by default) and then the plugin is called. The call per second is frequent enough for real-time processing, yet infrequent enough not to cause performance issues while running the test. The grpc call itself takes place locally through a continuously open connection, so its overhead is minimal.

The duration of buffering can be set by the plugin with the `buffering` variable in the `Init` response. Its value specifies the buffering duration in milliseconds. The default is `1000ms`, which is one second. Its minimum value is `200ms`.

## Poliglot

One of the big advantages of output plugins is that practically any programming language can be used, which is supported by grpc for server writing.

Although the go programming language is popular, its popularity does not reach the popularity of, for example, Python or JavaScript languages. With the support of these languages, the range of developers who can create an output plugin in their favorite programming language has been significantly expanded.

## Helpers

Based on the [protobuf descriptor](https://github.com/szkiba/xk6-output-plugin-proto/blob/master/output.proto) and the [HashiCorp go-plugin documentation](https://github.com/hashicorp/go-plugin/blob/main/docs/guide-plugin-write-non-go.md), the output plugin can be created in any programming [language supported by gRPC](https://grpc.io/docs/languages/). Plugin development is facilitated by a helper package in the following programming languages:

- [x] go: https://github.com/szkiba/xk6-output-plugin-go
- [ ] Python: coming soon
- [ ] JavaScript: coming soon

A helper package for additional programming languages will be created in the future based on community votes. What programming language would you like to use for output plugin development? [Please vote here](https://github.com/szkiba/xk6-output-plugin/discussions/1)

## Usage

The output plugin is an executable program whose name must (for security reasons) begin with the string `xk6-output-plugin-`. This is followed by the actual name of the plugin, which can be used to refer to it. The reference can of course also contain the full name with the prefix `xk6-output-plugin-`. The two reference below specify the same plugin:

```bash
./k6 run --out=example script.js
./k6 run --out=xk6-output-plugin-example script.js
```

The plugin is run by taking the `PATH` environment variable into account, unless the reference also contains a path. The reference below runs the file named `xk6-output-plugin-example` from the `plugins` directory in the current directory:

```bash
./k6 run --out=./plugins/example script.js
```

## Configuration

Output plugins can be configured using command line arguments. Arguments can be passed to the plugin in the usual way after the name of the plugin. In this case, the entire k6 output extension parameter must be protected with apostrophes.

```bash
./k6 run --out='plugin=example -flag value' script.js
```

The output plugin can also be configured with environment variables. The environment variables specified during the execution of the k6 command are received by the plugin via the `Init()` call. These variables are available in the map named `Environment` of the `Params` parameter.

```bash
./k6 run -e var=value --out=plugin=example script.js
```

## Example

**go example**
```go
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
```

There are more examples in the [examples](examples/) directory (coming soon).