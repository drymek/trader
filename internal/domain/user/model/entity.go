package model

import "github.com/google/uuid"

type Entity interface {
	GetID() string
	SetID(id string)
}

func GenerateID() string {
	return uuid.New().String()
}
