package request

import (
	"context"
	"encoding/json"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/infrastructure/logger"
	httpkit "github.com/go-kit/kit/transport/http"
)

type StreamRequest struct {
	UpdateId        int64  `json:"u"`
	Symbol          string `json:"s"`
	BestBidPrice    string `json:"b"`
	BestBidQuantity string `json:"B"`
	BestAskPrice    string `json:"a"`
	BestAskQuantity string `json:"A"`
}

func DecodeStreamRequest(logger logger.TraderLogger) httpkit.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var request StreamRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			if err2 := logger.Log("function", "DecodeStreamRequest", "error", err); err2 != nil {
				return nil, err2
			}

			return nil, httpx.ErrBadRequest
		}

		return request, nil
	}
}
