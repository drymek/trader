package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AccountSuite struct {
	suite.Suite
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

func (s *AccountSuite) TestValid() {
	a := Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	err := a.Validate()
	s.NoError(err)
}

func (s *AccountSuite) TestInvalidOwner() {
	a := Account{
		ID:            "123",
		Owner:         "",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	err := a.Validate()
	s.Error(err)
	s.Equal(ErrInvalidAccount, err)
}

func (s *AccountSuite) TestInvalidBalance() {
	a := Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "1 houndred",
		Currency:      "PLN",
		AccountNumber: 123456789,
	}

	err := a.Validate()
	s.Error(err)
	s.Equal(ErrInvalidAccount, err)
}

func (s *AccountSuite) TestInvalidCurrency() {
	a := Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "Dollars",
		AccountNumber: 123456789,
	}

	err := a.Validate()
	s.Error(err)
	s.Equal(ErrInvalidAccount, err)
}

func (s *AccountSuite) TestInvalidAccountNumber() {
	a := Account{
		ID:            "123",
		Owner:         "Marcin Dryka",
		Balance:       "100.0",
		Currency:      "PLN",
		AccountNumber: 0,
	}

	err := a.Validate()
	s.Error(err)
	s.Equal(ErrInvalidAccount, err)
}
