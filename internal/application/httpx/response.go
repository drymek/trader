package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"dryka.pl/trader/internal/infrastructure/logger"
	kithttp "github.com/go-kit/kit/transport/http"
)

type StatusCodeHolder interface {
	StatusCode() int
}

type HeaderHolder interface {
	Headers() map[string]string
}

type NoContent interface {
	NoContent() bool
}

func EncodeResponse(logger logger.TraderLogger) kithttp.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, response interface{}) error {
		hr, hasHeaders := response.(HeaderHolder)
		if hasHeaders {
			for key, value := range hr.Headers() {
				w.Header().Set(key, value)
			}
		}

		scr, hasStatusCode := response.(StatusCodeHolder)
		if hasStatusCode {
			w.WriteHeader(scr.StatusCode())
		}

		_, hasNoContent := response.(NoContent)
		if hasNoContent {
			return nil
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			if err2 := logger.Log("error", err); err2 != nil {
				return err2
			}
		}

		return err
	}
}
