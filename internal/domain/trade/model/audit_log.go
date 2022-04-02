package model

import "time"

type AuditLog struct {
	SourceQuantity string
	SourceCurrency string
	TargetCurrency string
	TargetQuantity string
	Timestamp      time.Time
	UpdateId       int64
}
