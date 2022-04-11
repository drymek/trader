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
	"dryka.pl/trader/internal/domain/user/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DeleteSuite struct {
	suite.Suite
}

func TestCreateSuite(t *testing.T) {
	suite.Run(t, new(DeleteSuite))
}

func (s *DeleteSuite) TestHandleErrorFromService() {
	req := request.AccountIDRequest{
		ID: "123",
	}

	service := new(test.ServiceMock)
	service.On("Delete", mock.Anything).
		Return(fmt.Errorf("invalid account"))

	_, err := endpoint.MakeDeleteEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusInternalServerError, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *DeleteSuite) TestHandleNotFoundFromService() {
	req := request.AccountIDRequest{
		ID: "123",
	}

	service := new(test.ServiceMock)
	service.On("Delete", mock.Anything).
		Return(repository.ErrAccountNotFound)

	_, err := endpoint.MakeDeleteEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusNotFound, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *DeleteSuite) TestHandleDelete() {
	req := request.AccountIDRequest{
		ID: "123",
	}

	service := new(test.ServiceMock)
	service.On("Delete", mock.Anything).
		Return(nil)

	got, err := endpoint.MakeDeleteEndpoint(nil, service)(context.TODO(), req)
	s.NoError(err)

	s.Equal(true, got.(*response.AccountResponse).NoContent())
	s.Equal(http.StatusOK, got.(httpx.StatusCodeHolder).StatusCode())
}
