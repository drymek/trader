package request_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dryka.pl/trader/internal/application/user/request"
	"dryka.pl/trader/tests/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type AccountRequestSuite struct {
	suite.Suite
}

func TestAccountRequestSuite(t *testing.T) {
	suite.Run(t, new(AccountRequestSuite))
}

func (s *AccountRequestSuite) TestValidInput() {
	json := `{
		"id": "123",
		"owner":"Marcin Dryka", 
		"balance":"100.0",  
		"currency":"PLN",
		"account_number":123456789
	}`

	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(json))
	got, err := request.DecodeAccountRequest(mock.NewNullLogger())(context.Background(), req)

	s.Nil(err)

	want := request.AccountRequest{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		s.T().Fatalf("mismatch (-want +got):\n%v", diff)
	}
}

func (s *AccountRequestSuite) TestInvalidInput() {
	json := `
		"id": "123",
		"owner":"Marcin Dryka", 
	`

	req := httptest.NewRequest(http.MethodPost, "/stream", strings.NewReader(json))
	_, err := request.DecodeAccountRequest(mock.NewNullLogger())(context.Background(), req)

	s.Error(err)
}
