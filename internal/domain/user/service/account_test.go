package service_test

import (
	"testing"

	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestCreate() {
	r := AccountRepositoryMock{}
	r.On("Create", mock.Anything).Return(nil)
	svc := service.NewAccountService(&r)
	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "200.00",
		Currency:      "PLN",
		AccountNumber: 3456789,
	}
	err := svc.Create(account)
	s.NoError(err)
	r.AssertCalled(s.T(), "Create", account)
}

type AccountRepositoryMock struct {
	mock.Mock
}

func (a *AccountRepositoryMock) Update(entity interface{}) error {
	args := a.Called(entity.(*model.Account))
	return args.Error(0)
}

func (a *AccountRepositoryMock) Delete(ID string) error {
	args := a.Called(ID)
	return args.Error(0)
}

func (a *AccountRepositoryMock) Create(entity interface{}) error {
	args := a.Called(entity.(*model.Account))
	return args.Error(0)
}

func (a *AccountRepositoryMock) Find(ID string) (interface{}, error) {
	args := a.Called(ID)
	return args.Get(0), args.Error(1)
}
