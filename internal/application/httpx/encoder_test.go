package httpx_test

import (
	"context"
	"net/http"

	"dryka.pl/trader/internal/application/httpx"
	mockx "dryka.pl/trader/tests/mock"
)

func (s *Suite) TestSimpleError() {
	w := new(Writer)
	w.On("Write", []byte(`{"error":"simple error"}`)).Return(3, nil)
	w.On("WriteHeader", http.StatusInternalServerError).Return(http.Header{})

	httpx.EncodeError(mockx.NewNullLogger())(context.TODO(), SimpleError{}, w)
	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 0)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 1)

	w.AssertCalled(s.T(), "Write", []byte(`{"error":"simple error"}`))
	w.AssertCalled(s.T(), "WriteHeader", http.StatusInternalServerError)
}

func (s *Suite) TestStatusCodeError() {
	w := new(Writer)
	w.On("Write", []byte(`{"error":"status code error"}`)).Return(3, nil)
	w.On("WriteHeader", http.StatusOK).Return(http.Header{})

	httpx.EncodeError(mockx.NewNullLogger())(context.TODO(), NewStatusCodeError(http.StatusOK), w)

	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 0)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 1)

	w.AssertCalled(s.T(), "Write", []byte(`{"error":"status code error"}`))
	w.AssertCalled(s.T(), "WriteHeader", http.StatusOK)
}

func (s *Suite) TestHeaderError() {
	w := new(Writer)
	w.On("Write", []byte(`{"error":"header error"}`)).Return(3, nil)
	w.On("WriteHeader", http.StatusInternalServerError).Return(http.Header{})

	h := make(http.Header)
	w.On("Header").Return(h)
	httpx.EncodeError(mockx.NewNullLogger())(context.TODO(), NewHeaderError("Content-Type", "application/json"), w)
	w.AssertNumberOfCalls(s.T(), "Write", 1)
	w.AssertNumberOfCalls(s.T(), "Header", 1)
	w.AssertNumberOfCalls(s.T(), "WriteHeader", 1)
}

func NewHeaderError(name string, value string) HeaderError {
	headers := make(map[string]string)
	headers[name] = value
	return HeaderError{headers: headers}
}

type HeaderError struct {
	headers map[string]string
}

func (r HeaderError) Error() string {
	return "header error"
}

func (r HeaderError) Headers() map[string]string {
	return r.headers
}

type SimpleError struct {
}

func (s SimpleError) Error() string {
	return "simple error"
}

type StatusCodeError struct {
	statusCode int
}

func (r StatusCodeError) Error() string {
	return "status code error"
}

func (r StatusCodeError) StatusCode() int {
	return r.statusCode
}

func NewStatusCodeError(statusCode int) StatusCodeError {
	r := StatusCodeError{statusCode: statusCode}
	return r
}
