package request_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dryka.pl/trader/internal/application/trade/request"
	"dryka.pl/trader/tests/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type RequestSuite struct {
	suite.Suite
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

func (s *RequestSuite) TestValidInput() {
	json := `{
		"u":400900235,
		"s":"BNBUSDT", 
		"b":"42.5",  
		"B":"30.0", 
		"a":"43.0",  
		"A":"10" 
	}`

	req := httptest.NewRequest(http.MethodPost, "/stream", strings.NewReader(json))
	got, err := request.DecodeStreamRequest(mock.NewNullLogger())(context.Background(), req)

	s.Nil(err)

	want := request.StreamRequest{
		UpdateId:        400900235,
		Symbol:          "BNBUSDT",
		BestBidPrice:    "42.5",
		BestBidQuantity: "30.0",
		BestAskPrice:    "43.0",
		BestAskQuantity: "10",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		s.T().Fatalf("mismatch (-want +got):\n%v", diff)
	}
}

func (s *RequestSuite) TestInvalidInput() {
	json := `
		"u":400900235
		"A":"10" 
	`

	req := httptest.NewRequest(http.MethodPost, "/stream", strings.NewReader(json))
	_, err := request.DecodeStreamRequest(mock.NewNullLogger())(context.Background(), req)

	s.Error(err)
}
