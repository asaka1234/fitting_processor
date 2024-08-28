package fitting_processor

import (
	"testing"
)

func TestProcess(t *testing.T) {
	//1. 获取有效原始数据
	gTimeDataMap, sTimeDataMap := GetLatestTickList(nil)
	//2. 处理数据
	UpdateDiffTable(nil, gTimeDataMap, sTimeDataMap)
}
