// 网格交易策略

package strategy

import (
	"fmt"
	"github.com/BlackCarDriver/StockMaster/common"
	"time"
)

// GridStrategy 网格交易策略
type GridStrategy struct {
	FlowStepUp   float64 `json:"flowStepUp"`   // 委托卖出价的距今涨幅
	FlowStepDown float64 `json:"flowStepDown"` // 委托买入价的距今涨幅 (正常应该是负数)
	StartTime    int64   `json:"startTime"`    // 最早执行时间 (超过该事件开始交易, 0-不限制)
	EndTime      int64   `json:"endTime"`      // 最晚执行时间 (超过该事件不再交易, 0-不限制)
	FirstPrize   float64 `json:"firstPrize"`   // 建仓限价, (大于0时跌破该价格时买入第一笔, 等于0时按照开盘价买入)
	FirstVol     int     `json:"firstVol"`     // 建仓交易量
	MaxCost      float64 `json:"maxCost"`      // 最大持有成本 (0-不限制)
	MinRetain    int     `json:"minRetain"`    // 最低保留份额 (卖出时保证最少剩余多少份额)
	Vol          int     `json:"vol"`          // 每次委托买入或卖出的数量
}

func (g *GridStrategy) Execute(account *common.Account, moment common.KLineNode) (err error) {
	// 先跳过一些不执行操作的场景
	if g.StartTime > 0 && g.StartTime > moment.Timestamp {
		log.Debug("%s, 未到执行时间", moment.TimeDesc)
		return
	}
	if g.EndTime > 0 && g.EndTime < moment.Timestamp {
		log.Debug("%s, 已过执行时间", moment.TimeDesc)
		return
	}

	// 建仓
	if account.LastDeal == nil {
		if g.FirstPrize > 0 && g.FirstPrize > moment.Top { // 跌破指定价格按现价买入
			log.Debug("%s 未达到触发价", moment.TimeDesc)
			return
		}
		isFirstDealOk, reason := false, ""
		var dealPrize float64  // 成交价
		if g.FirstPrize == 0 { // 到了指定时间, 按现价买入
			dealPrize = moment.Start
			isFirstDealOk, reason = account.Trad(common.ModeBuy, dealPrize, g.FirstVol, moment)
		}
		if g.FirstPrize > 0 && g.FirstPrize < moment.Top { // 跌破了触发价,以指定价格买入
			dealPrize = g.FirstPrize
			isFirstDealOk, reason = account.Trad(common.ModeBuy, dealPrize, g.FirstVol, moment)
		}
		if !isFirstDealOk {
			log.Warning("create first deal fail: reason=%v", reason)
			return
		}
		nextSellPrize := common.RisePrizeByFlow(dealPrize, g.FlowStepUp)
		nextBuyPrize := common.RisePrizeByFlow(dealPrize, g.FlowStepDown)
		account.CreateEntrust(common.ModeShell, nextSellPrize, g.Vol, moment.Timestamp, 0)
		account.CreateEntrust(common.ModeBuy, nextBuyPrize, g.Vol, moment.Timestamp, 0)
		return
	}

	account.Setting.BuyLock = false
	account.Setting.SellLock = false
	if account.Balance.CostRMB+moment.Start*float64(g.Vol) > g.MaxCost {
		account.Setting.BuyLock = true
	}
	if account.Balance.StockVol-g.Vol < g.MinRetain {
		account.Setting.SellLock = true
	}

	// 等待条件单触发
	mode, record := account.ExecuteEntrust(moment)
	if mode == common.ModeWait {
		return
	}
	if mode == common.ModeBuy || mode == common.ModeShell {
		nextSellPrize := common.RisePrizeByFlow(record.Price, g.FlowStepUp)
		nextBuyPrize := common.RisePrizeByFlow(record.Price, g.FlowStepDown)
		account.CreateEntrust(common.ModeShell, nextSellPrize, g.Vol, moment.Timestamp, 0)
		account.CreateEntrust(common.ModeBuy, nextBuyPrize, g.Vol, moment.Timestamp, 0)
	}

	return
}

func (g *GridStrategy) GetDesc() (desc string) {
	startTime, endTime, firstPrize, maxCost := "不限制", "不限制", "不限制", "不限制"
	if g.StartTime > 0 {
		startTime = time.Unix(g.StartTime, 0).Format("2006-01-02 15:04")
	}
	if g.EndTime > 0 {
		endTime = time.Unix(g.EndTime, 0).Format("2006-01-02 15:04")
	}
	if g.FirstPrize > 0 {
		firstPrize = fmt.Sprintf("%.2f", g.FirstPrize)
	}
	if g.MaxCost > 0 {
		maxCost = fmt.Sprintf("%.2f", g.MaxCost)
	}
	return fmt.Sprintf("生效时间=[%s~%s] \n建仓价格=%s 建仓交易额=%d \n买入跌幅=%.2f%% 卖出涨幅=%.2f%% \n保留额度=%d 限制市值=%s \n",
		startTime, endTime, firstPrize, g.FirstVol, -g.FlowStepDown, g.FlowStepUp, g.MinRetain, maxCost)
}
