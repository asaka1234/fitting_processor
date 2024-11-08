package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

// 数据源
type Quote struct {
	CurTime time.Time
	Bid     decimal.Decimal
	Ask     decimal.Decimal
	Source  string //来源
	//-------------------------------------------
	Symbol string
}

type AvgQuote struct {
	MinuteIndex     int       //5min bar在一个小时内的刻度(0-11)
	MinuteStartTime time.Time //每一段的起始时间
	//-------各自的均值计算中间结果-----------------
	BidSum   decimal.Decimal
	BidCount int
	BidAvg   decimal.Decimal
	//-------------------------------------------
	AskSum   decimal.Decimal
	AskCount int
	AskAvg   decimal.Decimal
	//-------------------------------------------
	Symbol string
}

type DiffQuote struct {
	MinuteIndex     int       //5min bar在一个小时内的刻度(0-11)
	MinuteStartTime time.Time //每一段的起始时间

	//-------各自的均值计算中间结果-----------------
	BidSumGau   decimal.Decimal
	BidCountGau int
	BidAvgGau   decimal.Decimal
	//-------各自的均值计算中间结果-----------------
	BidSumSau   decimal.Decimal
	BidCountSau int
	BidAvgSau   decimal.Decimal
	//-------------------------------------------
	BidDiff decimal.Decimal //差值(avg之差)
	//=====================================================
	AskSumGau   decimal.Decimal
	AskCountGau int
	AskAvgGau   decimal.Decimal
	//-------------------------------------------
	AskSumSau   decimal.Decimal
	AskCountSau int
	AskAvgSau   decimal.Decimal
	//-------------------------------------------
	AskDiff decimal.Decimal
	//-------------------------------------------
	Symbol string
}
