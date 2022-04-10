package e2e

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"dryka.pl/trader/internal/application/config"
	"dryka.pl/trader/internal/application/server"
	"dryka.pl/trader/internal/domain/user/model"
	repository2 "dryka.pl/trader/internal/domain/user/repository"
	"dryka.pl/trader/internal/domain/user/service"
	"dryka.pl/trader/internal/infrastructure/persistence/inmemory/repository"
	testCmpopts "dryka.pl/trader/tests/cmpopts"
	"dryka.pl/trader/tests/mock"
	"github.com/caarlos0/env/v6"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type CrudSuite struct {
	suite.Suite
	AppDependencies server.Dependencies
}

func TestCrudSuite(t *testing.T) {
	s := new(CrudSuite)
	opts := env.Options{Environment: map[string]string{
		"DATABASE_FILE": "database/sqlite/database_test.sqlite",
	}}
	c, err := config.NewConfig(opts)
	if err != nil {
		t.Fatal("invalid config")
	}

	accountRepository := repository.NewAccountRepository()
	s.AppDependencies = server.Dependencies{
		Logger:            mock.NewNullLogger(),
		Config:            c,
		CrudService:       service.NewAccountService(accountRepository),
		AccountRepository: accountRepository,
	}

	suite.Run(t, s)
}

func (s *CrudSuite) SetupTest() {
	account, err := s.AppDependencies.AccountRepository.Find("123")
	if err != nil {
		return
	}
	s.AppDependencies.AccountRepository.Delete(account.(*model.Account).ID)
}

func (s *CrudSuite) TestCreate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789
	}`

	requestBody := []byte(accountJSON)
	res, err := http.Post(srv.URL+"/accounts", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte("{\"id\":\"123\"}\n"), got)
	s.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201")

	want := model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	existing, err := s.AppDependencies.AccountRepository.Find("123")
	s.Nil(err)

	if diff := cmp.Diff(want, *existing.(*model.Account)); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}
}

func (s *CrudSuite) TestCrateDupplicate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	s.AppDependencies.AccountRepository.Create(model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	})

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789,
	}`

	requestBody := []byte(accountJSON)
	res, err := http.Post(srv.URL+"/accounts", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"error":"Bad Request"}`), got)
	s.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400")
}

func (s *CrudSuite) TestFetch() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()
	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	err := s.AppDependencies.AccountRepository.Create(account)
	s.Nil(err)
	res, err := http.Get(srv.URL + "/accounts/" + account.ID)
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	transformJSON := testCmpopts.TransformJSON()
	want := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789
	}`
	if diff := cmp.Diff([]byte(want), got, cmp.Options{transformJSON}); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}

	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 400")

}

func (s *CrudSuite) TestFetchNotFound() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/accounts/fake-id")
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"error":"Not Found"}`), got)
	s.Equal(http.StatusNotFound, res.StatusCode, "Expected status code 400")
}

func (s *CrudSuite) TestUpdateNewObject() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789
	}`

	requestBody := []byte(accountJSON)
	req, err := http.NewRequest(http.MethodPut, srv.URL+"/accounts", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)
	res, err := http.DefaultClient.Do(req)
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte("{\"id\":\"123\"}\n"), got)
	s.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201")
}

func (s *CrudSuite) TestUpdate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	err := s.AppDependencies.AccountRepository.Create(account)
	s.NoError(err)

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789
	}`

	requestBody := []byte(accountJSON)
	req, err := http.NewRequest(http.MethodPut, srv.URL+"/accounts", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)
	res, err := http.DefaultClient.Do(req)
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte("{\"id\":\"123\"}\n"), got)
	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")
}

func (s *CrudSuite) TestRemove() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	s.AppDependencies.AccountRepository.Create(account)

	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/accounts/"+account.ID, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)
	res, err := http.DefaultClient.Do(req)
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(""), got)

	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")

	_, err = s.AppDependencies.AccountRepository.Find("123")
	s.ErrorIs(err, repository2.ErrAccountNotFound)
}

func (s *CrudSuite) TestRemoveFakeId() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/accounts/fake-id", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)
	res, err := http.DefaultClient.Do(req)
	s.Nil(err)

	s.Equal(http.StatusNotFound, res.StatusCode, "Expected status code 404")
}
