package httpx_test

import (
	"context"
	"net/http"
	"testing"

	"dryka.pl/trader/internal/application/httpx"
	mockx "dryka.pl/trader/tests/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestSimpleResponse() {
	w := new(Writer)
	w.On("Write", []byte("{}\n")).Return(3, nil)
	err := httpx.EncodeResponse(mockx.NewNullLogger())(context.TODO(), w, SimpleResponse{})
	s.NoError(err)
	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 0)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 0)

	w.AssertCalled(s.T(), "Write", []byte("{}\n"))
}

func (s *Suite) TestStatusCodeResponse() {
	w := new(Writer)
	w.On("Write", []byte("{}\n")).Return(3, nil)
	w.On("WriteHeader", http.StatusOK).Return(http.Header{})

	err := httpx.EncodeResponse(mockx.NewNullLogger())(context.TODO(), w, NewStatusCodeResponse(http.StatusOK))

	s.NoError(err)
	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 0)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 1)

	w.AssertCalled(s.T(), "Write", []byte("{}\n"))
}

func (s *Suite) TestHeaderResponse() {
	w := new(Writer)
	w.On("Write", []byte("{}\n")).Return(3, nil)
	h := make(http.Header)
	w.On("Header").Return(h)
	err := httpx.EncodeResponse(mockx.NewNullLogger())(context.TODO(), w, NewHeaderResponse("Content-Type", "application/json"))
	s.NoError(err)
	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 1)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 0)
}

func (s *Suite) TestEmptyResponse() {
	w := new(Writer)
	//w.On("Write", []byte("{}\n")).Return(3, nil)
	err := httpx.EncodeResponse(mockx.NewNullLogger())(context.TODO(), w, EmptyResponse{})
	s.NoError(err)
	w.AssertNumberOfCalls(s.T(), "Write", 0)
	w.AssertNumberOfCalls(s.T(), "Header", 0)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 0)

	w.AssertNotCalled(s.T(), "Write", mock.Anything)
}

func NewHeaderResponse(name string, value string) HeaderResponse {
	headers := make(map[string]string)
	headers[name] = value
	return HeaderResponse{headers: headers}
}

type EmptyResponse struct {
}

func (r EmptyResponse) NoContent() bool {
	return true
}

type HeaderResponse struct {
	headers map[string]string
}

func (r HeaderResponse) Headers() map[string]string {
	return r.headers
}

type SimpleResponse struct {
}

type StatusCodeResponse struct {
	statusCode int
}

func (r StatusCodeResponse) StatusCode() int {
	return r.statusCode
}

func NewStatusCodeResponse(statusCode int) StatusCodeResponse {
	r := StatusCodeResponse{statusCode: statusCode}
	return r
}

type Writer struct {
	mock.Mock
}

func (w *Writer) Header() http.Header {
	args := w.Called()
	return args.Get(0).(http.Header)
}

func (w *Writer) Write(v []byte) (int, error) {
	args := w.Called(v)
	return args.Int(0), args.Error(1)
}

func (w *Writer) WriteHeader(statusCode int) {
	w.Called(statusCode)
}
