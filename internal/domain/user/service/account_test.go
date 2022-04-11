package service_test

import (
	"testing"

	"dryka.pl/trader/internal/domain/user/model"
	"dryka.pl/trader/internal/domain/user/repository"
	"dryka.pl/trader/internal/domain/user/service"
	"github.com/google/go-cmp/cmp"
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

func (s *ServiceSuite) TestCreateNoId() {
	r := AccountRepositoryMock{}
	r.On("Create", mock.Anything).Return(nil)
	svc := service.NewAccountService(&r)
	account := &model.Account{
		ID:            "",
		Owner:         "Marcin Dryka",
		Balance:       "200.00",
		Currency:      "PLN",
		AccountNumber: 3456789,
	}
	err := svc.Create(account)
	s.NoError(err)
	s.NotEmpty(r.Calls[0].Arguments[0].(*model.Account).ID)
	r.AssertCalled(s.T(), "Create", account)
}

func (s *ServiceSuite) TestDelete() {
	r := AccountRepositoryMock{}
	r.On("Delete", "123").Return(nil)
	svc := service.NewAccountService(&r)

	err := svc.Delete("123")
	s.NoError(err)
	r.AssertCalled(s.T(), "Delete", "123")
}

func (s *ServiceSuite) TestFetch() {
	r := AccountRepositoryMock{}
	account := model.Account{}
	r.On("Find", "123").Return(account, nil)
	svc := service.NewAccountService(&r)

	got, err := svc.Fetch("123")
	want := account
	s.NoError(err)

	if diff := cmp.Diff(want, got); diff != "" {
		s.Failf("value mismatch", "(-want +got):\n%v", diff)
	}
	r.AssertCalled(s.T(), "Find", "123")
}

func (s *ServiceSuite) TestUpdateNew() {
	r := AccountRepositoryMock{}
	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "123.00",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	r.On("Find", "123").Return(nil, repository.ErrAccountNotFound)
	r.On("Create", account).Return(nil)

	svc := service.NewAccountService(&r)

	got, err := svc.UpdateOrCreate(account)
	s.True(got)
	s.NoError(err)
	r.AssertCalled(s.T(), "Create", account)
}

func (s *ServiceSuite) TestUpdateExisting() {
	r := AccountRepositoryMock{}
	account := &model.Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "123.00",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}
	accountOld := &model.Account{
		ID:            "123",
		Owner:         "Marcin",
		Balance:       "1.00",
		Currency:      "USD",
		AccountNumber: 123456789,
	}

	r.On("Find", "123").Return(accountOld, nil)
	r.On("Update", account).Return(nil)

	svc := service.NewAccountService(&r)

	got, err := svc.UpdateOrCreate(account)
	s.False(got)
	s.NoError(err)
	r.AssertCalled(s.T(), "Update", account)
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
