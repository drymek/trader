package httpx

import (
	"net/http"
)

var ErrBadRequest = NewHttpError(http.StatusBadRequest, "Bad Request")
var ErrNotFound = NewHttpError(http.StatusNotFound, "Not Found")
var ErrInternalServerError = NewHttpError(http.StatusInternalServerError, "Internal Server Error")

type HttpError struct {
	code    int
	message string
}

func (a HttpError) StatusCode() int {
	return a.code
}

func (a HttpError) Error() string {
	return a.message
}

func NewHttpError(code int, message string) HttpError {
	return HttpError{
		code:    code,
		message: message,
	}
}
