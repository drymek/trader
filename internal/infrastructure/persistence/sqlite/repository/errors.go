package repository

import "errors"

var ErrPersistenceCannotAddAccount = errors.New("persistence cannot add account")
var ErrPersistencePrepareError = errors.New("persistence prepare error")
var ErrPersistenceCannotAddLog = errors.New("persistence cannot add log")
var ErrPersistenceNotFound = errors.New("persistence not found")
var ErrPersistenceError = errors.New("persistence error")
