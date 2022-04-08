package endpoint_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/application/trade/endpoint"
	request "dryka.pl/trader/internal/application/trade/request"
	"dryka.pl/trader/internal/domain/trade/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EndpointSuite struct {
	suite.Suite
}

func TestEndpointSuite(t *testing.T) {
	suite.Run(t, new(EndpointSuite))
}

func (s *EndpointSuite) TestHandleErrorFromService() {
	req := request.StreamRequest{
		UpdateId:        123,
		Symbol:          "BNBUSDT",
		BestBidPrice:    "1",
		BestBidQuantity: "1",
		BestAskPrice:    "1",
		BestAskQuantity: "1",
	}

	service := new(ServiceMock)
	service.On("Consider", mock.Anything).
		Return(fmt.Errorf("Invalid amount"))

	_, err := endpoint.MakeStreamEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusInternalServerError, err.(httpx.StatusCodeHolder).StatusCode())
}

func (s *EndpointSuite) TestHandleErrorFromTick() {
	req := request.StreamRequest{}

	service := new(ServiceMock)

	_, err := endpoint.MakeStreamEndpoint(nil, service)(context.TODO(), req)

	s.Error(err)
	s.Equal(http.StatusBadRequest, err.(httpx.StatusCodeHolder).StatusCode())
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Consider(tick model.Tick) error {
	args := s.Called(tick)
	return args.Error(0)
}
