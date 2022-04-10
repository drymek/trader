package response

import "dryka.pl/trader/internal/domain/user/model"

type AccountResponse struct {
	ID        string `json:"id,omitempty"`
	code      int
	noContent bool
}

func (r *AccountResponse) StatusCode() int {
	return r.code
}

func (r *AccountResponse) NoContent() bool {
	return r.noContent
}

func NewAccountResponse(code int, account *model.Account) *AccountResponse {
	if account == nil {
		return &AccountResponse{
			code:      code,
			noContent: true,
		}
	}

	return &AccountResponse{
		ID:   account.ID,
		code: code,
	}
}
