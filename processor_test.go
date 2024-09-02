package fitting_processor

import (
	"github.com/asaka1234/fitting_processor/entity"
	"testing"
	"time"
)

type TestHistory struct {
}

func (aa *TestHistory) GetNextBarTime() *time.Time {
	re := time.Now()
	return &re
}

func (aa *TestHistory) GetDiffQuoteByBarTime(barTime time.Time) *entity.DiffQuote {
	return nil

}
func (aa *TestHistory) SaveDiffQuoteByBarTime(barTime time.Time, diffQuote entity.DiffQuote) {

}

func (aa *TestHistory) GetMinuteLength() int {
	return 0
}

//------------------------------------

func TestProcess(t *testing.T) {

	gauPath := "/Users/ccc/Downloads/GAUCNH.csv"
	sauPath := "/Users/ccc/Downloads/SAUCNH.csv"

	//1. 实现接口
	ins := &TestHistory{}
	//2. 获取有效原始数据
	gTimeDataMap, sTimeDataMap := GetLatestTickListFromCSV(gauPath, sauPath, ins)
	//3. 处理数据
	UpdateDiffTable(nil, gTimeDataMap, sTimeDataMap)
}
