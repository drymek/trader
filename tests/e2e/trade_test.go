package e2e

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"dryka.pl/trader/internal/application/server"
)

func (s *Suite) TestStream() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	streamJSON := `{ 
		"u":400900223,
		"s":"BNBUSDT",
		"b":"42.0", 
		"B":"5.0", 
		"a":"43.0",
		"A":"10"
	}`

	requestBody := []byte(streamJSON)
	res, err := http.Post(srv.URL+"/stream", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal([]byte("{}\n"), got)
	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")
}

func (s *Suite) TestInvalidJsonStream() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	invalidStreamJSON := `
		"u":400900223,
		"s":"BNBUSDT",
	`

	requestBody := []byte(invalidStreamJSON)
	res, err := http.Post(srv.URL+"/stream", "application/json", bytes.NewBuffer(requestBody))
	s.Nil(err)

	_, err = ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	s.Equal(http.StatusInternalServerError, res.StatusCode, "Expected status code 500")
}
