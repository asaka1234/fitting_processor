package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Quote struct {
	CurTime time.Time
	Bid     decimal.Decimal
	Ask     decimal.Decimal
	Source  string //来源
}

type AvgQuote struct {
	MinuteIndex     int       //5min bar在一个小时内的刻度(0-11)
	MinuteStartTime time.Time //每一段的起始时间
	AvgBid          decimal.Decimal
	AvgAsk          decimal.Decimal
}

type DiffQuote struct {
	MinuteIndex     int       //5min bar在一个小时内的刻度(0-11)
	MinuteStartTime time.Time //每一段的起始时间
	DiffBid         decimal.Decimal
	DiffAsk         decimal.Decimal
}
