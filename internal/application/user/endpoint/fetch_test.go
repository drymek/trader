package endpoint_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/application/user/endpoint"
	"dryka.pl/trader/internal/application/user/endpoint/test"
	"dryka.pl/trader/internal/application/user/request"
	"dryka.pl/trader/internal/application/user/response"
	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FetchSuite struct {
	suite.Suite
}

func TestFetchSuite(t *testing.T) {
	suite.Run(t, new(FetchSuite))
}

func (s *FetchSuite) TestHandleErrorFromService() {
	req := request.AccountIDRequest{
		ID: "123",
	}

	service := new(test.ServiceMock)
	service.On("Fetch", mock.Anything).
		Return(nil, fmt.Errorf("invalid account"))

	_, err := endpoint.MakeFetchEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusInternalServerError, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *FetchSuite) TestHandleAccountError() {
	req := request.AccountIDRequest{}

	service := new(test.ServiceMock)

	service.On("Fetch", mock.Anything).
		Return(nil, repository.ErrAccountNotFound)

	_, err := endpoint.MakeFetchEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusNotFound, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *FetchSuite) TestHandleFetch() {
	req := request.AccountIDRequest{
		ID: "123",
	}

	service := new(test.ServiceMock)
	service.On("Fetch", mock.Anything).
		Return(&model.Account{
			ID:            "123",
			Owner:         "Marcin Dryka",
			Balance:       "100",
			Currency:      "PLN",
			AccountNumber: 123456,
		}, nil)

	got, err := endpoint.MakeFetchEndpoint(nil, service)(context.TODO(), req)
	s.NoError(err)

	s.Equal(req.ID, got.(*response.FullAccountResponse).ID)
	s.Equal("Marcin Dryka", got.(*response.FullAccountResponse).Owner)
	s.Equal(http.StatusOK, got.(httpx.StatusCodeHolder).StatusCode())
}
