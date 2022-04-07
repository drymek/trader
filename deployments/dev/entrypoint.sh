#!/bin/sh
set -e

cd /go/src/trader

echo "Going to run environment ${ENVIRONMENT}"

if [ "$ENVIRONMENT" = "test" ]; then
    # shellcheck disable=SC2046
    go test $(go list ./... | grep -v /vendor/) -v
else
    cp database/sqlite/database.sqlite.template database/sqlite/database.sqlite

    openssl req  -new  -newkey rsa:2048  -nodes  -keyout configs/certificate.key  -out configs/certificate.csr  -subj "/C=PL/L=Gliwice/O=Trader/OU=Trader/CN=localhost"
    openssl  x509  -req  -days 365  -in configs/certificate.csr  -signkey configs/certificate.key  -out configs/certificate.crt

    reflex -v --start-service --regex='(\.go$|go\.mod|\.js$|\.html$)' -- sh -c 'go run /go/src/trader/cmd/app'
fi