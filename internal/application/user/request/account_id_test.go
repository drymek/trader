package request_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dryka.pl/trader/internal/application/httpx"
	"dryka.pl/trader/internal/application/user/request"
	"dryka.pl/trader/tests/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type AccountIDRequestSuite struct {
	suite.Suite
}

func TestAccountIDRequestSuite(t *testing.T) {
	suite.Run(t, new(AccountIDRequestSuite))
}

func (s *AccountIDRequestSuite) TestNoIdInput() {
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{}`))
	ctx := context.Background()
	values := map[string]string{}
	mux.SetURLVars(req, values)
	_, err := request.DecodeAccountIDRequest(mock.NewNullLogger())(ctx, req)

	s.Error(err)
	s.Equal(httpx.ErrNotFound, err)
}

func (s *AccountIDRequestSuite) TestValidInput() {
	req := httptest.NewRequest(http.MethodPost, "/accounts/{id}", strings.NewReader(`{}`))
	ctx := context.Background()
	values := map[string]string{
		"id": "123",
	}
	req = mux.SetURLVars(req, values)
	got, err := request.DecodeAccountIDRequest(mock.NewNullLogger())(ctx, req)

	s.Nil(err)

	want := request.AccountIDRequest{
		ID: "123",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		s.T().Fatalf("mismatch (-want +got):\n%v", diff)
	}
}

func (s *AccountIDRequestSuite) TestInvalidInput() {
	json := `
		"id": "123",
	`

	req := httptest.NewRequest(http.MethodPost, "/stream", strings.NewReader(json))
	_, err := request.DecodeAccountIDRequest(mock.NewNullLogger())(context.Background(), req)

	s.Error(err)
}
