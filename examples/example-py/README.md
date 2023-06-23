### Python plugin example

## Setup

Enter the following commands in the `example-py` directory:

```bash
pip install virtualenv
virtualenv env
source env/bin/activate
pip install xk6-output-plugin-py
```

## Run

In the project directory:

```bash
./k6 run --out=plugin=./examples/example-py/example.py script.js
```
