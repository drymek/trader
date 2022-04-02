package service

import "errors"

var ErrInvalidBestBidPrice = errors.New("invalid best bid price")
var ErrUnableToProcessOrder = errors.New("unable to process order")
var ErrInvalidBestBidQuanity = errors.New("invalid best bid quantity")
