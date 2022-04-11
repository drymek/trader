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

# Run production image

Generate certificate:

```bash
    mkdir configs
    openssl req  -new  -newkey rsa:2048  -nodes  -keyout configs/certificate.key  -out configs/certificate.csr  -subj "/C=PL/L=Gliwice/O=Trader/OU=Trader/CN=localhost"
    openssl  x509  -req  -days 365  -in configs/certificate.csr  -signkey configs/certificate.key  -out configs/certificate.crt
```

Run image

```bash
docker run -v $(pwd)/configs:/go/app/configs \
  -e KEY_FILE=/go/app/configs/certificate.key \
  -e CRT_FILE=/go/app/configs/certificate.crt \
  -p 8080:8080  \
  drymek/trader:latest 
```

# Usage examples

```bash
$ echo '{ "u":400900235, "s":"BNBUSDT", "b":"42.5", "B":"1.0", "a":"43.0", "A":"10" }' | http --verify=no https://localhost:8080/stream
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 3
Content-Type: text/plain; charset=utf-8
Date: Sat, 02 Apr 2022 19:19:33 GMT

{}
```

## CRUD operations

### Create

```bash
$ echo '{"id": "1", "owner": "Marcin", "balance": "100.02", "currency": "PLN", "account_number": 1234}' | http --verify=no https://localhost:8080/accounts
HTTP/1.1 201 Created
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Sun, 10 Apr 2022 21:14:49 GMT

{
    "id": "1"
}
```

### Read

```bash
$ http --verify=no https://localhost:8080/accounts/1
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 86
Content-Type: text/plain; charset=utf-8
Date: Sun, 10 Apr 2022 21:15:33 GMT

{
    "account_number": 1234,
    "balance": "100.02",
    "currency": "PLN",
    "id": "1",
    "owner": "Marcin"
}
```

### Update

```bash
$ echo '{"id": "1", "owner": "John Doe", "balance": "100.02", "currency": "PLN", "account_number": 1234}' | http --verify=no PUT https://localhost:8080/accounts
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Sun, 10 Apr 2022 21:16:11 GMT

{
    "id": "1"
}
```

### Delete

```bash
$ http --verify=no DELETE https://localhost:8080/accounts/1
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Origin, Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Content-Length: 0
Date: Sun, 10 Apr 2022 21:16:33 GMT
```