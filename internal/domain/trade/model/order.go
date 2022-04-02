package model

import "github.com/shopspring/decimal"

type Order struct {
	Quantity       decimal.Decimal
	Price          decimal.Decimal
	SourceCurrency string
	TargetCurrency string
}

func (o Order) Validate() error {
	if o.Quantity.LessThanOrEqual(decimal.NewFromInt(0)) {
		return ErrInvalidOderQuantity
	}

	if o.SourceCurrency != "BNB" {
		return ErrInvalidOderSourceCurrency
	}

	if o.TargetCurrency != "USDT" {
		return ErrInvalidOderTargetCurrency
	}

	return nil
}

func (o Order) Symbol() string {
	return o.SourceCurrency + o.TargetCurrency
}
