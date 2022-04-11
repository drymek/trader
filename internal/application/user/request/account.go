package request

import (
	"context"
	"encoding/json"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/infrastructure/logger"
	httpkit "github.com/go-kit/kit/transport/http"
)

type AccountRequest struct {
	ID            string `json:"id"`
	Owner         string `json:"owner"`
	Balance       string `json:"balance"`
	Currency      string `json:"currency"`
	AccountNumber int64  `json:"account_number"`
}

func DecodeAccountRequest(logger logger.TraderLogger) httpkit.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var request AccountRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			if err2 := logger.Log("function", "DecodeAccountRequest", "error", err); err2 != nil {
				return nil, err2
			}

			return nil, httpx.ErrBadRequest
		}

		return request, nil
	}
}
