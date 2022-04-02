package model

import (
	"github.com/shopspring/decimal"
)

type Tick struct {
	UpdateId        int64
	Symbol          string
	BestBidPrice    decimal.Decimal
	BestBidQuantity decimal.Decimal
	BestAskPrice    decimal.Decimal
	BestAskQuantity decimal.Decimal
}

func NewTick(
	updateId int64,
	symbol string,
	bestBidPrice string,
	bestBidQuantity string,
	bestAskPrice string,
	bestAskQuantity string,
) (Tick, error) {
	bbp, err := decimal.NewFromString(bestBidPrice)
	if err != nil {
		return Tick{}, err
	}
	bbq, err := decimal.NewFromString(bestBidQuantity)
	if err != nil {
		return Tick{}, err
	}
	bap, err := decimal.NewFromString(bestAskPrice)
	if err != nil {
		return Tick{}, err
	}
	baq, err := decimal.NewFromString(bestAskQuantity)
	if err != nil {
		return Tick{}, err
	}
	return Tick{
		UpdateId:        updateId,
		Symbol:          symbol,
		BestBidPrice:    bbp,
		BestBidQuantity: bbq,
		BestAskPrice:    bap,
		BestAskQuantity: baq,
	}, nil
}

func (t Tick) Validate() error {
	if t.UpdateId <= 0 {
		return ErrInvalidTickUpdateId
	}

	if t.Symbol != "BNBUSDT" {
		return ErrInvalidTickSymbol
	}

	if t.BestBidQuantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidTickBestBidQuantity
	}

	return nil
}
