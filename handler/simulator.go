package handler

import (
	"fmt"
	"github.com/BlackCarDriver/GoProject-api/color"
	"github.com/BlackCarDriver/StockMaster/common"
	"github.com/BlackCarDriver/StockMaster/dao"
	"github.com/BlackCarDriver/StockMaster/strategy"
)

// Simulate 根据指定账号状态和给出的k线图数据, 按照指定交易策略遍历指数数据, 得到最终的账号状态
func Simulate(before common.Account, stockData dao.KLineData, strategy strategy.Strategy) (after *common.Account, err error) {
	account := &before
	if before.InitFundRMB <= 0.0 || before.Name == "" || len(stockData.KLines) == 0 {
		err = fmt.Errorf("unexpect params")
		return
	}
	if account.Balance.BalanceRMB == 0 {
		account.Balance.BalanceRMB = account.InitFundRMB
	}
	for i, moment := range stockData.KLines {
		account.UpdateStat(moment)

		err = strategy.Execute(account, moment)
		if err != nil {
			log.Error("execute fail: i=%d err=%v moment=%+v", i, err, moment)
			break
		}
	}
	return account, err
}

// PrintRunResult 在控制台打印模拟结果
func PrintRunResult(account *common.Account, strategy strategy.Strategy, data dao.KLineData) {
	if account == nil || account.LastPrize == nil {
		log.Warning("unexpect nil account")
		return
	}

	t := account.TradStat
	var buyList, sellList string
	var buyCount, sellCount int
	for _, item := range account.BuyEntrust {
		if item.DealTime == 0 {
			buyCount++
			buyList += fmt.Sprintf("%.2f, ", item.Price)
		}
	}
	for _, item := range account.SellEntrust {
		if item.DealTime == 0 {
			sellCount++
			sellList += fmt.Sprintf("%.2f, ", item.Price)
		}
	}

	balance := account.Balance
	canSell := float64(balance.StockVol) * account.LastPrize.End // 当前持有市值
	currentValue := balance.BalanceRMB + canSell                 // 当前总资产

	color.Blue("============ 数据描述 =============")
	color.HiBlack("名称: %s  代码: %s", data.Name, data.Code)
	color.HiBlack("数据时间范围:  %s ~ %s", data.From, data.To)
	color.HiBlack("节点长度:  %d", len(data.KLines))
	color.HiBlack("更新时间: %s", common.TimeFormat(data.UpdateTime))

	color.Blue("============ 策略描述 =============")
	color.HiBlack(strategy.GetDesc())

	color.Blue("============ 账号信息 =============")
	color.HiBlack("账号名称: %s", account.Name)
	color.HiBlack("备注信息: %s", account.Note)
	color.HiBlack("初始金额: %.2f", account.InitFundRMB)

	color.Blue("============ 操作日志 =============")
	for _, item := range account.ActionLog {
		color.HiBlack("%s, %s, %s", common.TimeFormat(item.Timestamp), item.Mode, item.Desc)
	}

	color.Blue("============ 交易记录 =============")
	for i, item := range account.TradLog {
		color.HiBlack("i=%d, %s, %s, 价格=%.2f   份额=%d ", i+1, common.TimeFormat(item.Timestamp), item.Mode, item.Prize, item.Vol)
	}

	color.Blue("============ 委托列表 =============")
	color.HiBlack("买入委托数量=%d    卖出委托数量=%d", buyCount, sellCount)
	color.HiBlack("买入价格档位： %s", buyList)
	color.HiBlack("卖出价格档位： %s", sellList)

	color.Blue("============ 过程统计 =============")
	color.HiBlack("交易次数=%d  (Buy=%d, Sell=%d)", t.SellCounter+t.BuyCounter, t.BuyCounter, t.SellCounter)
	color.HiBlack("最高持仓=%d    最低持仓=%d", t.MaxVol, t.MinVol)
	color.HiBlack("最高成本=%.2f  最低成本=%.2f", t.MaxCost, t.MinCost)
	color.HiBlack("最高资产=%.2f  最低资产=%.2f", t.MaxValue, t.MinValue)

	color.Blue("============ 最终结果 =============")
	color.HiBlack("账户总资产=%.2f", currentValue)
	color.HiBlack("可用余额=%.2f", balance.BalanceRMB)
	color.HiBlack("最新报价=%.2f", account.LastPrize.End)
	color.HiBlack("持有份额=%d", balance.StockVol)
	color.HiBlack("持仓成本=%.2f", balance.CostRMB)
	color.HiBlack("持有市值=%.2f", canSell)
	color.HiBlack("持仓盈亏=%.2f  (%.2f%%)", canSell-balance.CostRMB, common.CountRiseRange(balance.CostRMB, canSell))
	color.HiBlack("总盈亏=%.2f  (%.2f%%)", currentValue-account.InitFundRMB, common.CountRiseRange(account.InitFundRMB, currentValue))
}
