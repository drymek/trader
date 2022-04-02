package endpoint

import (
	"context"
	"net/http"

	requestx "dryka.pl/trader/internal/application/trade/request"
	"dryka.pl/trader/internal/application/trade/response"
	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/service"
	"dryka.pl/trader/internal/infrastructure/logger"
	"github.com/go-kit/kit/endpoint"
)

func MakeStreamEndpoint(_ logger.TraderLogger, service service.OrderService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := request.(requestx.StreamRequest)

		tick, err := model.NewTick(r.UpdateId, r.Symbol, r.BestBidPrice, r.BestBidQuantity, r.BestAskPrice, r.BestAskQuantity)
		err2 := tick.Validate()
		if err != nil || err2 != nil {
			return response.NewStreamResponse(http.StatusBadRequest), nil
		}

		err = service.Consider(tick)

		if err != nil {
			return response.NewStreamResponse(http.StatusBadRequest), nil
		}

		return response.NewStreamResponse(http.StatusOK), nil
	}
}
