package endpoint

import (
	"context"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	requestx "dryka.pl/trader/internal/application/user/request"
	"dryka.pl/trader/internal/application/user/response"
	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/service"
	"dryka.pl/trader/internal/infrastructure/logger"
	"github.com/go-kit/kit/endpoint"
)

func MakeUpdateEndpoint(_ logger.TraderLogger, service service.CrudService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := request.(requestx.AccountRequest)

		account := model.NewAccount(r.ID, r.Owner, r.Balance, r.Currency, r.AccountNumber)
		err := account.Validate()
		if err != nil {
			return nil, httpx.ErrBadRequest
		}

		isNew, err := service.UpdateOrCreate(account)

		if err != nil {
			return nil, httpx.ErrInternalServerError
		}

		if isNew {
			return response.NewAccountResponse(http.StatusCreated, account), nil
		}

		return response.NewAccountResponse(http.StatusOK, account), nil
	}
}
