package server

import (
	"net/http"

	"dryka.pl/trader/internal/application/healthcheck"
	"dryka.pl/trader/internal/application/httpx"
	tradeEndpoint "dryka.pl/trader/internal/application/trade/endpoint"
	tradeRequest "dryka.pl/trader/internal/application/trade/request"
	userEndpoint "dryka.pl/trader/internal/application/user/endpoint"
	userRequest "dryka.pl/trader/internal/application/user/request"
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

	accountCreateHandler := kithttp.NewServer(
		userEndpoint.MakeCreateEndpoint(l, d.CrudService),
		userRequest.DecodeAccountRequest(l),
		httpx.EncodeResponse(l),
		kithttp.ServerErrorEncoder(httpx.EncodeError(l)),
	)

	accountFetchHandler := kithttp.NewServer(
		userEndpoint.MakeFetchEndpoint(l, d.CrudService),
		userRequest.DecodeAccountIDRequest(l),
		httpx.EncodeResponse(l),
		kithttp.ServerErrorEncoder(httpx.EncodeError(l)),
	)

	accountUpdateHandler := kithttp.NewServer(
		userEndpoint.MakeUpdateEndpoint(l, d.CrudService),
		userRequest.DecodeAccountRequest(l),
		httpx.EncodeResponse(l),
		kithttp.ServerErrorEncoder(httpx.EncodeError(l)),
	)

	accountDeleteHandler := kithttp.NewServer(
		userEndpoint.MakeDeleteEndpoint(l, d.CrudService),
		userRequest.DecodeAccountIDRequest(l),
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
	r.Handle(newrelic.WrapHandle(app, "/accounts", AccessControl(accountCreateHandler))).Methods(http.MethodOptions, http.MethodPost)
	r.Handle(newrelic.WrapHandle(app, "/accounts", AccessControl(accountUpdateHandler))).Methods(http.MethodOptions, http.MethodPut)
	r.Handle(newrelic.WrapHandle(app, "/accounts/{id}", AccessControl(accountFetchHandler))).Methods(http.MethodOptions, http.MethodGet)
	r.Handle(newrelic.WrapHandle(app, "/accounts/{id}", AccessControl(accountDeleteHandler))).Methods(http.MethodOptions, http.MethodDelete)

	return r
}
