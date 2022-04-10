package request

import (
	"context"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/infrastructure/logger"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type AccountIDRequest struct {
	ID string `json:"id"`
}

func DecodeAccountIDRequest(logger logger.TraderLogger) httpkit.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {

		id, ok := mux.Vars(r)["id"]
		if !ok {
			return nil, httpx.ErrNotFound
		}

		request := AccountIDRequest{
			ID: id,
		}

		return request, nil
	}
}
