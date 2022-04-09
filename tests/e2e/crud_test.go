package e2e

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"dryka.pl/trader/internal/application/server"
	testCmpopts "dryka.pl/trader/tests/cmpopts"
	"github.com/google/go-cmp/cmp"
)

func (s *Suite) TestCreate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789,
	}`

	requestBody := []byte(accountJSON)
	res, err := http.Post(srv.URL+"/account", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"id":"123"}\n`), got)
	s.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201")
}

func (s *Suite) TestCrateDupplicate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	s.AppDependencies.AccountRepository.Create(user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       100.0,
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
	res, err := http.Post(srv.URL+"/account", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"error":"Bad request"}`), got)
	s.Equal(http.StatusBadRequest, res.StatusCode, "Expected status code 400")
}

func (s *Suite) TestFetch() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()
	account := user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       100.0,
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	s.AppDependencies.AccountRepository.Create(account)

	res, err := http.Get(srv.URL + "/account/" + account.ID)
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	transformJSON := testCmpopts.TransformJSON()
	want := account
	if diff := cmp.Diff(want, got, cmp.Options{transformJSON}); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}

	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 400")
}

func (s *Suite) TestFetchNotFound() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/account/fake-id")
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal(http.StatusNotFound, res.StatusCode, "Expected status code 400")
}

func (s *Suite) TestUpdateNewObject() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	accountJSON := `{
		"id": "123",
		"owner": "Marcin Dryka",
		"balance": "100.0",
		"currency": "PLN",
		"account_number": 123456789,
	}`

	requestBody := []byte(accountJSON)
	res, err := http.Put(srv.URL+"/account", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"id":"123"}\n`), got)
	s.Equal(http.StatusCreated, res.StatusCode, "Expected status code 201")
}

func (s *Suite) TestUpdate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	s.AppDependencies.AccountRepository.Create(user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       100.0,
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
	res, err := http.Put(srv.URL+"/account", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"id":"123"}\n`), got)
	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")
}

func (s *Suite) TestPartialUpdate() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	account := user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       100.0,
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	s.AppDependencies.AccountRepository.Create(account)

	accountJSON := `{
		"balance": "200.0",
		"currency": "USD",
	}`

	requestBody := []byte(accountJSON)
	res, err := http.Patch(srv.URL+"/account/"+account.ID, "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte(`{"id":"123"}\n`), got)

	want := user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       200.0,
		Currency:      "USD",
		AccountNumber: 123456789,
	}
	existing := s.AppDependencies.AccountRepository.Find("123")

	if diff := cmp.Diff(want, got); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}

	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")
}

func (s *Suite) TestRemove() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	s.AppDependencies.AccountRepository.Create(user.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       100.0,
		Currency:      "PLN",
		AccountNumber: 123456789,
	})

	res, err := http.Delete(srv.URL+"/account/"+account.ID, "application/json", bytes.NewBuffer(requestBody))
	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")

	s.Nil(err)

	_, err = s.AppDependencies.AccountRepository.Find("123")
	s.ErrorIs(err, user.ErrAccountNotFound)
}
