package httpx

import (
	"context"
	"net/http"

	"dryka.pl/trader/internal/infrastructure/logger"
	kithttp "github.com/go-kit/kit/transport/http"
)

func EncodeError(logger logger.TraderLogger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		hr, hasHeaders := err.(HeaderHolder)
		if hasHeaders {
			for key, value := range hr.Headers() {
				w.Header().Set(key, value)
			}
		}

		scr, hasStatusCode := err.(StatusCodeHolder)
		if hasStatusCode {
			w.WriteHeader(scr.StatusCode())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}
