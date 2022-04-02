package healthcheck

import (
	"context"
	"net/http"

	httpkit "github.com/go-kit/kit/transport/http"
)

type Request struct {
}

func DecodeRequest() httpkit.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		return Request{}, nil
	}
}
