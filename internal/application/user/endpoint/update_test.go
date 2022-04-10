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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UpdateSuite struct {
	suite.Suite
}

func TestUpdateSuite(t *testing.T) {
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) TestHandleErrorFromService() {
	req := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	service := new(test.ServiceMock)
	service.On("UpdateOrCreate", mock.Anything).
		Return(false, fmt.Errorf("invalid account"))

	_, err := endpoint.MakeUpdateEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusInternalServerError, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *UpdateSuite) TestHandleAccountError() {
	req := request.AccountRequest{}

	service := new(test.ServiceMock)

	_, err := endpoint.MakeUpdateEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusBadRequest, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *UpdateSuite) TestHandleUpdate() {
	req := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	service := new(test.ServiceMock)
	service.On("UpdateOrCreate", mock.Anything).
		Return(true, nil)

	got, err := endpoint.MakeUpdateEndpoint(nil, service)(context.TODO(), req)
	s.NoError(err)

	s.Equal(req.ID, got.(*response.AccountResponse).ID)
	s.Equal(http.StatusCreated, got.(httpx.StatusCodeHolder).StatusCode())
}

func (s *UpdateSuite) TestHandleCreate() {
	req := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	service := new(test.ServiceMock)
	service.On("UpdateOrCreate", mock.Anything).
		Return(true, nil)

	got, err := endpoint.MakeUpdateEndpoint(nil, service)(context.TODO(), req)
	s.NoError(err)

	s.Equal(req.ID, got.(*response.AccountResponse).ID)
	s.Equal(http.StatusCreated, got.(httpx.StatusCodeHolder).StatusCode())
}
