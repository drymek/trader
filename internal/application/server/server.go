package server

import (
	"net/http"

	"dryka.pl/trader/internal/application/healthcheck"
	"dryka.pl/trader/internal/application/httpx"
	tradeEndpoint "dryka.pl/trader/internal/application/trade/endpoint"
	tradeRequest "dryka.pl/trader/internal/application/trade/request"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewServer(d Dependencies) *mux.Router {
	l := d.Logger
	streamHandler := kithttp.NewServer(
		tradeEndpoint.MakeStreamEndpoint(l, d.Service),
		tradeRequest.DecodeStreamRequest(l),
		httpx.EncodeResponse(l),
		kithttp.ServerErrorEncoder(httpx.EncodeError(l)),
	)

	healthcheckHandler := kithttp.NewServer(
		healthcheck.MakeEndpoint(l),
		healthcheck.DecodeRequest(),
		httpx.EncodeResponse(l),
		kithttp.ServerErrorEncoder(httpx.EncodeError(l)),
	)

	r := mux.NewRouter()
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(d.Config.GetNewRelicConfigAppName()),
		newrelic.ConfigLicense(d.Config.GetNewRelicConfigLicense()),
	)

	if err != nil {
		err2 := l.Log("context", "newrelic", "error", err)
		if err2 != nil {
			panic(err2)
		}
	}

	r.Handle(newrelic.WrapHandle(app, "/stream", AccessControl(streamHandler))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle(newrelic.WrapHandle(app, "/healthcheck", AccessControl(healthcheckHandler))).Methods(http.MethodGet)

	return r
}
