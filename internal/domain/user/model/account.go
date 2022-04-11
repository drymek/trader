package model

import "github.com/shopspring/decimal"

type Account struct {
	ID            string `json:"id"`
	Owner         string `json:"owner"`
	Balance       string `json:"balance"`
	Currency      string `json:"currency"`
	AccountNumber int64  `json:"account_number"`
}

func (a *Account) GetID() string {
	return a.ID
}

func (a *Account) SetID(id string) {
	a.ID = id
}

func (a Account) Validate() error {
	if a.Owner == "" {
		return ErrInvalidAccount
	}

	_, err := decimal.NewFromString(a.Balance)
	if err != nil {
		return ErrInvalidAccount
	}

	currencySymbolLength := 3
	if len(a.Currency) > currencySymbolLength {
		return ErrInvalidAccount
	}

	if a.AccountNumber == 0 {
		return ErrInvalidAccount
	}

	return nil
}

func NewAccount(ID string, Owner string, Balance string, Currency string, AccountNumber int64) *Account {
	return &Account{
		ID:            ID,
		Owner:         Owner,
		Balance:       Balance,
		Currency:      Currency,
		AccountNumber: AccountNumber,
	}
}
