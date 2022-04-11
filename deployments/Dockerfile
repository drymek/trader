FROM golang:1.18-alpine as builder
RUN apk add --no-cache openssh git build-base
RUN apk add --no-cache openssl
RUN go install github.com/cespare/reflex@latest
WORKDIR /go/src/trader
COPY go.mod /go/src/trader
RUN go mod download

COPY . /go/src/trader/

RUN go build -o trader /go/src/trader/cmd/app

### Production image
FROM alpine:latest
ENV PORT 8080
EXPOSE 8080

RUN apk add --no-cache curl
RUN addgroup -S trader && adduser -S -G trader trader

COPY --chown=trader --from=builder /go/src/trader/trader /go/app/
COPY --chown=trader ./database/sqlite/database.sqlite.template /go/app/database.sqlite
ENV DATABASE_FILE /go/app/database.sqlite

USER trader

HEALTHCHECK --interval=5m --timeout=3s \
   CMD curl -f http://localhost:${PORT}/heathcheck || exit 1

CMD ["/go/app/trader"]