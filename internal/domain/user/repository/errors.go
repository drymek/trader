package repository

import "errors"

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountNotFound = errors.New("account not found")
var ErrInvalidModel = errors.New("invalid model")
