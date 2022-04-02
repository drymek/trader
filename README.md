# Usage

Project contains Taskfile with available actions. To use Taskfile, you need to install task:

```bash
brew install go-task/tap/go-task
```

More information: https://taskfile.dev/#/installation

To view all available tasks just type:

```bash
task
```

# Development environment

To run development environment run:

```bash
task dev
```

# Tests

To execute tests, run:

```bash
task test
```

# Usage examples

```bash
$ echo '{ "u":400900235, "s":"BNBUSDT", "b":"42.5", "B":"1.0", "a":"43.0", "A":"10" }' | http :8080/stream
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 3
Content-Type: text/plain; charset=utf-8
Date: Sat, 02 Apr 2022 19:19:33 GMT

{}
```