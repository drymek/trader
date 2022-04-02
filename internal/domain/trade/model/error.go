package model

import "errors"

var ErrInvalidOderQuantity = errors.New("invalid order quantity")
var ErrInvalidOderSourceCurrency = errors.New("invalid order base currency")
var ErrInvalidOderTargetCurrency = errors.New("invalid order target currency")
var ErrInvalidTickUpdateId = errors.New("invalid update id")
var ErrInvalidTickSymbol = errors.New("invalid symbol")
var ErrInvalidTickBestBidQuantity = errors.New("invalid best bid quantity")
