package e2e

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"dryka.pl/trader/internal/application/server"
	testCmpopts "dryka.pl/trader/tests/cmpopts"
	"github.com/google/go-cmp/cmp"
)

func (s *Suite) TestHealthcheck() {
	srv := httptest.NewServer(server.NewServer(s.AppDependencies))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/healthcheck")
	s.Nil(err)

	got, err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	s.Nil(err)
	s.Nil(err2)

	want := []byte(`{}`)

	transformJSON := testCmpopts.TransformJSON()

	if diff := cmp.Diff(want, got, cmp.Options{transformJSON}); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}
	s.Equal(http.StatusOK, res.StatusCode, "Expected status code 200")
}
