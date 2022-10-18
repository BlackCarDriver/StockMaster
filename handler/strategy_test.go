package handler

import (
	"fmt"
	"github.com/BlackCarDriver/StockMaster/common"
	"github.com/BlackCarDriver/StockMaster/dao"
	"github.com/BlackCarDriver/StockMaster/strategy"
	"testing"
)

var account1 = common.Account{
	Name:        "GridTest1",
	Note:        "网格策略测试账号",
	InitFundRMB: 100000.0,
}

var gridStrategy1 = strategy.GridStrategy{
	FlowStepUp:   11,
	FlowStepDown: -1,
	StartTime:    0,
	EndTime:      0,
	FirstPrize:   0.0,
	FirstVol:     3000,
	MaxCost:      100000.0,
	MinRetain:    100,
	Vol:          200,
	ExpireDay:    120,
}

func TestGridStrategy(t *testing.T) {
	// step1: 读取K线图
	path := fmt.Sprintf("../dao/mockdata/%s.json", "510500_1day")
	mkData, err := dao.ReadKLineMockData(path)
	if err != nil {
		log.Error("read fail: err=%v", err)
		return
	}
	log.Info("read mockData success: n_data=%d", len(mkData.KLines))

	// step2: 执行策略
	after, err := Simulate(account1, mkData, &gridStrategy1)
	if err != nil {
		log.Error("simulate fail: err=%v", err)
		return
	}
	log.Info("simulate success")
	PrintRunResult(after, &gridStrategy1, mkData)
}
