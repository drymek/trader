#!/bin/sh
set -e

cd /go/src/trader

echo "Going to run environment ${ENVIRONMENT}"

if [ "$ENVIRONMENT" = "test" ]; then
    # shellcheck disable=SC2046
    go test $(go list ./... | grep -v /vendor/) -v
else
    cp database/sqlite/database.sqlite.template database/sqlite/database.sqlite
    reflex -v --start-service --regex='(\.go$|go\.mod|\.js$|\.html$)' -- sh -c 'go run /go/src/trader/cmd/app'
fi