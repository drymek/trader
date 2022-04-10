package endpoint

import (
	"context"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	requestx "dryka.pl/trader/internal/application/user/request"
	"dryka.pl/trader/internal/application/user/response"
	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
	"dryka.pl/trader/internal/domain/user/service"
	"dryka.pl/trader/internal/infrastructure/logger"
	"github.com/go-kit/kit/endpoint"
)

func MakeFetchEndpoint(_ logger.TraderLogger, service service.CrudService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := request.(requestx.AccountIDRequest)

		account, err := service.Fetch(r.ID)

		if err != nil {
			if err == repository.ErrAccountNotFound {
				return nil, httpx.ErrNotFound
			}
			return nil, httpx.ErrInternalServerError
		}

		return response.NewFullAccountResponse(http.StatusOK, account.(*model.Account)), nil
	}
}
