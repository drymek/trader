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

type CreateSuite struct {
	suite.Suite
}

func TestEndpointSuite(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

func (s *CreateSuite) TestHandleErrorFromService() {
	req := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	service := new(test.ServiceMock)
	service.On("Create", mock.Anything).
		Return(fmt.Errorf("invalid account"))

	_, err := endpoint.MakeCreateEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusInternalServerError, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *CreateSuite) TestHandleAccountError() {
	req := request.AccountRequest{}

	service := new(test.ServiceMock)

	_, err := endpoint.MakeCreateEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusBadRequest, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *CreateSuite) TestHandleCreate() {
	req := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	service := new(test.ServiceMock)
	service.On("Create", mock.Anything).
		Return(nil)

	got, err := endpoint.MakeCreateEndpoint(nil, service)(context.TODO(), req)
	s.NoError(err)

	s.Equal(req.ID, got.(*response.AccountResponse).ID)
	s.Equal(http.StatusCreated, got.(httpx.StatusCodeHolder).StatusCode())
}
