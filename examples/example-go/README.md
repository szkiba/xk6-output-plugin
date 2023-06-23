### Go plugin example

## Setup

Enter the following commands in the `example-go` directory:

```bash
go build  -ldflags="-s -w" -o xk6-output-plugin-example .
```

## Run

In the project directory:

```bash
./k6 run --out=plugin=./examples/example-go/example script.js
```
