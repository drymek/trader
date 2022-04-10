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

COPY --from=builder /go/src/trader/trader /go/app/

USER trader

HEALTHCHECK --interval=5m --timeout=3s \
   CMD curl -f http://localhost:${PORT}/heathcheck || exit 1

CMD ["/go/app/trader"]