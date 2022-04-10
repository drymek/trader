package model

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

func (a Account) Validate() interface{} {
	// TODO implement validation & tests
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
