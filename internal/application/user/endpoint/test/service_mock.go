package test

import "github.com/stretchr/testify/mock"

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Fetch(id string) (interface{}, error) {
	args := s.Called(id)
	return args.Get(0), args.Error(1)
}

func (s *ServiceMock) UpdateOrCreate(entity interface{}) (bool, error) {
	args := s.Called(entity)
	return args.Bool(0), args.Error(1)
}

func (s *ServiceMock) Delete(id string) error {
	args := s.Called(id)
	return args.Error(0)
}

func (s *ServiceMock) Create(entity interface{}) error {
	args := s.Called(entity)
	return args.Error(0)
}
