package fitting_processor

import (
	"github.com/asaka1234/fitting_processor/entity"
	"testing"
	"time"
)

type TestHistory struct {
}

func (aa *TestHistory) GetNextBarTime() time.Time {
	return time.Now()
}

func (aa *TestHistory) GetDiffQuoteByBarTime(barTime time.Time) *entity.DiffQuote {
	return nil

}
func (aa *TestHistory) SaveDiffQuoteByBarTime(barTime time.Time, diffQuote entity.DiffQuote) {

}

//------------------------------------

func TestProcess(t *testing.T) {

	gauPath := "/Users/ccc/Downloads/GAUCNH.csv"
	sauPath := "/Users/ccc/Downloads/SAUCNH.csv"
	barLength := 5

	//1. 实现接口
	ins := &TestHistory{}
	//2. 获取有效原始数据
	gTimeDataMap, sTimeDataMap := GetLatestTickListFromCSV(gauPath, sauPath, ins, barLength)
	//3. 处理数据
	UpdateDiffTable(nil, barLength, gTimeDataMap, sTimeDataMap)
}
