package response

import "dryka.pl/trader/internal/domain/user/model"

type FullAccountResponse struct {
	ID            string `json:"id"`
	code          int
	Owner         string `json:"owner"`
	Balance       string `json:"balance"`
	Currency      string `json:"currency"`
	AccountNumber int64  `json:"account_number"`
}

func (r *FullAccountResponse) StatusCode() int {
	return r.code
}

func NewFullAccountResponse(code int, account *model.Account) *FullAccountResponse {
	return &FullAccountResponse{
		ID:            account.ID,
		Owner:         account.Owner,
		Balance:       account.Balance,
		Currency:      account.Currency,
		AccountNumber: account.AccountNumber,
		code:          code,
	}
}
