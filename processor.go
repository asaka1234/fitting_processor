package fitting_processor

import (
	"encoding/csv"
	"fmt"
	"github.com/asaka1234/fitting_processor/entity"
	"github.com/asaka1234/fitting_processor/utils"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"sort"
	"time"
)

type HistoryDataInterface interface {
	GetMinuteLength() int
	GetNextBarTime() *time.Time
	GetDiffQuoteByBarTime(barTime time.Time) *entity.DiffQuote
	SaveDiffQuoteByBarTime(barTime time.Time, diffQuote entity.DiffQuote)
}

// ------------------------------------------

// 获取最新的数据(要做增量才可以)
func GetLatestTickListFromCSV(gauPath, sauPath string, history HistoryDataInterface) (map[time.Time][]entity.Quote, map[time.Time][]entity.Quote) {

	//获取最新的时间刻度
	nextBarTime := history.GetNextBarTime()

	barLength := history.GetMinuteLength()

	// 解析数据源
	gTimeDataMap, _ := parseCSV(gauPath, barLength, nextBarTime) //最后一个参数是：查找的最小时间的数据
	sTimeDataMap, _ := parseCSV(sauPath, barLength, nextBarTime) //最后一个参数是：查找的最小时间的数据

	return gTimeDataMap, sTimeDataMap
}

// 更新5分钟刻度的码表,插入新的数据
func UpdateDiffTable(history HistoryDataInterface, gTimeDataMap, sTimeDataMap map[time.Time][]entity.Quote) {

	barLength := history.GetMinuteLength()

	//需要做数据增补. 所以要把所有时间段都涵盖进来
	minBarTime, maxBarTime := getScope(lo.Keys(gTimeDataMap), lo.Keys(sTimeDataMap))

	for barTime := minBarTime; barTime.Before(maxBarTime) || barTime.Equal(maxBarTime); barTime = barTime.Add(time.Minute * time.Duration(barLength)) {

		if lo.HasKey(gTimeDataMap, barTime) && lo.HasKey(sTimeDataMap, barTime) {
			//2个都有, 那就计算差值

			//2024-08-22 22:25:00
			//if barTime.Format("2006-01-02 15:04:05") != "2024-08-22 22:25:00" {
			//	continue
			//}

			//2.1 拿到这一刻度的quote列表
			gQuoteList := gTimeDataMap[barTime]
			sQuoteList := sTimeDataMap[barTime]

			//fmt.Printf("--q1-----%d\n", len(gQuoteList))
			//fmt.Printf("--q2-----%d, %+v, %s\n", len(sQuoteList), sQuoteList, barTime)

			//fmt.Printf("---> raw: %+v\n", gQuoteList)

			//2.2 计算这一刻度的均值
			gAvgQuote := calculateAveragePrice(barTime, gQuoteList, barLength)
			sAvgQuote := calculateAveragePrice(barTime, sQuoteList, barLength)
			//sAvgQuote := gAvgQuote

			//fmt.Printf("---> cnt: %d,  barTime:%s\n", len(gQuoteList), barTime)

			//fmt.Printf("---> barTime------%+v\n", barTime)

			//fmt.Printf("---> gAvgQuote------%+v\n", gAvgQuote)
			//fmt.Printf("---> sAvgQuote------%+v\n", sAvgQuote)

			//2.3 计算这一刻度两者的差值(gau-sau)
			diffQuote := calculateDiffPrice(gAvgQuote, sAvgQuote, barLength)

			//fmt.Printf("diffQuote------%+v\n", diffQuote)

			//fmt.Printf("---wsx12w----%+v\n", diffQuote)
			//2.4 赋值(直接入库合适.便于第一次批量处理)
			history.SaveDiffQuoteByBarTime(barTime, *diffQuote)
			//fmt.Printf("---%s---%s\n", barTime, *diffQuote)

			//os.Exit(1)
		} else {
			//查找替换的diff后替换下,找不到就算了.
			re := searchReplaceBar(history, barTime)
			if re != nil && !re.MinuteStartTime.IsZero() {
				history.SaveDiffQuoteByBarTime(barTime, *re)
				//fmt.Printf("---%s---%s\n", barTime, re)
			}
		}
	}
}

// 寻找替换的bar
func searchReplaceBar(history HistoryDataInterface, missTime time.Time) *entity.DiffQuote {
	//往前一天找一下(只往前)->增量的
	replaceTime := missTime.Add(-time.Hour * time.Duration(24))
	preFindResult := history.GetDiffQuoteByBarTime(replaceTime)
	return preFindResult
}

// 计算两边指定hour内的：各个slice的价差
func calculateDiffPrice(gauAvg, sauAvg entity.AvgQuote, barLength int) *entity.DiffQuote {
	//1. param check
	if !sauAvg.MinuteStartTime.Equal(gauAvg.MinuteStartTime) {
		return nil
	}
	//2. calculate
	return &entity.DiffQuote{
		sauAvg.MinuteStartTime.Minute() / barLength,
		sauAvg.MinuteStartTime,
		gauAvg.AvgBid.Sub(sauAvg.AvgBid),
		gauAvg.AvgAsk.Sub(sauAvg.AvgAsk),
	}
}

// 计算每一段的平均值
func calculateAveragePrice(startTime time.Time, tickList []entity.Quote, barLength int) entity.AvgQuote {

	sumBid := decimal.Zero
	sumAsk := decimal.Zero
	for _, item := range tickList {
		sumBid = sumBid.Add(item.Bid)
		sumAsk = sumAsk.Add(item.Ask)
	}

	avgBid := sumBid.Div(decimal.NewFromInt(int64(len(tickList))))
	avgAsk := sumAsk.Div(decimal.NewFromInt(int64(len(tickList))))

	//fmt.Printf("-->calculateAveragePrice--%s----sumBid:%+v, len:%d, avgBid:%+v\n", startTime, sumBid, int64(len(tickList)), avgBid)
	//fmt.Printf("-->calculateAveragePrice--%s----sumAsk:%+v, len:%d, avgAsk:%+v\n", startTime, sumAsk, int64(len(tickList)), avgAsk)

	//fmt.Printf("---sumBid:%s,  cnt:%d, avgBid:%s\n", sumBid, int64(len(tickList)), avgBid)
	//fmt.Printf("---sumAsk:%s,  cnt:%d, avgAsk:%s\n", sumAsk, int64(len(tickList)), avgAsk)

	//获取这5分钟这一段的平均价格-------
	return entity.AvgQuote{
		startTime.Minute() / barLength,
		startTime,
		avgBid,
		avgAsk,
	}
}

// 返回开始/结束时间
func getScope(a []time.Time, b []time.Time) (time.Time, time.Time) {
	ccc := lo.Union(a, b)
	sort.Slice(ccc, func(i, j int) bool { return ccc[i].Before(ccc[j]) })
	return ccc[0], ccc[len(ccc)-1]
}

// 解析excel
// key是5min的起点
// sliceLength : 5分钟一段
// nextBarTime 所有数据的时间都要晚于或者等于 这个时间
func parseCSV(filePath string, barLength int, nextBarTime *time.Time) (map[time.Time][]entity.Quote, error) {

	result := make(map[time.Time][]entity.Quote, 0)

	// Load a csv file.
	f, err1 := os.Open(filePath)
	if err1 != nil {
		panic(err1)
	}

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if len(record) != 4 {
			fmt.Printf("err[%s] data[%s]\n", err.Error(), record[0])
			continue
		}

		ctime, err := time.Parse("2006.01.02 15:04:05", record[0])
		if err != nil {
			fmt.Printf("err[%s] data[%s]\n", err.Error(), record[0])
			continue
		}
		bid, err := decimal.NewFromString(record[1])
		if err != nil {
			fmt.Printf("err[%s] data[%s]\n", err.Error(), record[1])
			continue
		}
		ask, err := decimal.NewFromString(record[2])
		if err != nil {
			fmt.Printf("err[%s] data[%s]\n", err.Error(), record[2])
			continue
		}

		//构造对应的tick
		tick := entity.Quote{
			CurTime: ctime,
			Bid:     bid,
			Ask:     ask,
			Source:  record[3],
		}

		if nextBarTime != nil && ctime.Before(*nextBarTime) {
			//早期的数据都计算过了, 所以不用再次计算了
			continue
		}

		//5分钟的起点时间
		minuteTime := utils.GetSliceStartTime(ctime, barLength)

		//只保留小时
		//hourTime := time.Date(ctime.Year(), ctime.Month(), ctime.Day(), ctime.Hour(), 0, 0, 0, ctime.Location())
		//看看是否存在了
		exist := lo.HasKey(result, minuteTime)
		if exist {
			//存在
			result[minuteTime] = append(result[minuteTime], tick)
		} else {
			//不存在
			result[minuteTime] = []entity.Quote{tick}
		}
	}
	return result, nil
}
