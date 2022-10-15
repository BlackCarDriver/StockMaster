package common

import (
	"fmt"
	"sort"
)

// Account 交易账号
type Account struct {
	Name        string       `json:"name"`
	Note        string       `json:"note"`                  // 备注
	InitFundRMB float64      `json:"InitFundRMB"`           // 初始总资产
	TargetStock string       `json:"targetSock"`            // 目标股票代码
	Balance     BalanceInfo  `json:"BalanceInfo"`           // 账户余额信息
	TradStat    TradInfo     `json:"TradInfo"`              // 交易过程统计数据
	TradLog     []TradRecord `json:"tradLogList,omitempty"` // 交易记录
	ActionLog   []Action     `json:"actionLog"`             // 操作日志
	Setting     Setting      `json:"-"`                     // 过程变量
	BuyEntrust  []Entrust    `json:"-"`                     // 买入委托单
	SellEntrust []Entrust    `json:"-"`                     // 卖出委托单
	LastPrize   *KLineNode   `json:"-"`                     // 最新股票状况
	LastDeal    *KLineNode   `json:"-"`                     // 上次交易时的股票状况
}

// Trad 交易
func (a *Account) Trad(mode OpMode, prize float64, vol int, moment KLineNode) (isOk bool, reason string) {
	if vol <= 0 || prize == 0 || mode == ModeWait {
		log.Warning("不合法的输入: prize=%f vol=%d moment=%+v mode=%v ", prize, vol, moment, mode)
		reason = "参数错误"
		return
	}

	// 放弃交易的情况
	value := prize * float64(vol) // 交易金额
	if mode == ModeBuy && value > a.Balance.BalanceRMB {
		reason = fmt.Sprintf("余额不足,无法买入: 均价=%.2f 委托价=%.2f 请求扣费=%.2f 余额=%.2f ",
			moment.GetAveragePrice(), prize, value, a.Balance.BalanceRMB)
		a.recordAction(moment.Timestamp, ActionGiveUp, reason)
		return
	}
	if mode == ModeShell && a.Balance.StockVol-vol < 0 {
		reason = fmt.Sprintf("份额不足,无法卖出: 均价=%.2f 委托价=%.2f 请求卖出=%d 持有份额=%d",
			moment.GetAveragePrice(), prize, vol, a.Balance.StockVol)
		a.recordAction(moment.Timestamp, ActionGiveUp, reason)
		return
	}
	if mode == ModeShell && a.Setting.SellLock {
		reason = fmt.Sprintf("卖出操作被禁止,无法卖出: 均价=%.2f 委托价=%.2f", moment.GetAveragePrice(), prize)
	}
	if mode == ModeBuy && a.Setting.BuyLock {
		reason = fmt.Sprintf("买入操作被禁止,无法买入: 均价=%.2f 委托价=%.2f", moment.GetAveragePrice(), prize)
	}

	// 买入或卖出指定数量,更新账号信息
	if mode == ModeBuy {
		a.Balance.BalanceRMB -= value
		a.Balance.CostRMB += value
		a.Balance.StockVol += vol
		a.recordAction(moment.Timestamp, ActionBuy, fmt.Sprintf("成功买入, 价格区间=[%.2f~%.2f] 成交份额=%d 成交价=%.2f 成交金额=%.2f",
			moment.Bottom, moment.Top, vol, prize, value))
		a.recordTradLog(moment.Timestamp, ModeBuy, prize, vol)
		a.LastDeal = &moment
		a.maintainTradStat(ModeBuy)
	}
	if mode == ModeShell {
		a.Balance.StockVol -= vol
		a.Balance.BalanceRMB += value
		a.Balance.CostRMB -= value
		a.recordAction(moment.Timestamp, ActionShell, fmt.Sprintf("成功卖出, 价格区间=[%.2f~%.2f] 成交份额=%d 成交价=%.2f 成交金额=%.2f",
			moment.Bottom, moment.Top, vol, prize, value))
		a.recordTradLog(moment.Timestamp, ModeShell, prize, vol)
		a.LastDeal = &moment
		a.maintainTradStat(ModeShell)
	}

	// 维护其他交易信息
	return true, ""
}

// CreateEntrust 创建条件单
func (a *Account) CreateEntrust(mode OpMode, prize float64, vol int, startTime int64, deadTime int64) {
	item := Entrust{
		StarTime: startTime,
		DeadTime: deadTime,
		Price:    prize,
		Vol:      vol,
		DealTime: 0,
	}
	if mode == ModeBuy {
		a.BuyEntrust = append(a.BuyEntrust, item)
		a.recordAction(startTime, ActionEntrust, fmt.Sprintf("创建条件单, 价格下破%.2f时买入%d份", prize, vol))
		sort.Slice(a.BuyEntrust, func(i, j int) bool {
			return a.BuyEntrust[i].Price > a.BuyEntrust[j].Price
		})
	}
	if mode == ModeShell {
		a.SellEntrust = append(a.SellEntrust, item)
		a.recordAction(startTime, ActionEntrust, fmt.Sprintf("创建条件单, 价格上穿%.2f时卖出%d份", prize, vol))
		sort.Slice(a.SellEntrust, func(i, j int) bool {
			return a.SellEntrust[i].Price < a.SellEntrust[j].Price
		})
	}
}

// ExecuteEntrust 检查是否有可触发的条件单,自动执行并返回对应的委托信息(卖出优先)
func (a *Account) ExecuteEntrust(moment KLineNode) (mode OpMode, record *Entrust) {
	// 尝试执行达到条件的卖出委托
	for i, entrust := range a.SellEntrust {
		if a.Setting.SellLock {
			return
		}
		if entrust.DeadTime > moment.Timestamp || entrust.DealTime > 0 {
			continue
		}
		if moment.Top < entrust.Price { // 价格未上穿卖出价
			continue
		}
		isDeal, _ := a.Trad(ModeShell, entrust.Price, entrust.Vol, moment)
		if !isDeal {
			break
		}
		a.SellEntrust[i].DealTime = moment.Timestamp // 执行成功
		mode = ModeShell
		record = &a.SellEntrust[i]
		return
	}
	// 尝试执行达到条件的买入委托
	for i, entrust := range a.BuyEntrust {
		if a.Setting.BuyLock {
			return
		}
		if entrust.DeadTime > moment.Timestamp || entrust.DealTime > 0 {
			continue
		}
		if moment.Bottom > entrust.Price { // 未跌破到指定买入价
			continue
		}
		isDeal, _ := a.Trad(ModeBuy, entrust.Price, entrust.Vol, moment)
		if !isDeal {
			break
		}
		a.BuyEntrust[i].DealTime = moment.Timestamp
		mode = ModeBuy
		record = &a.BuyEntrust[i]
		return
	}
	mode = ModeWait
	return
}

// UpdateStat 更新统计信息维护状态变量
func (a *Account) UpdateStat(moment KLineNode) {
	a.LastPrize = &moment
	a.maintainTradStat(ModeWait)
}

// 保存交易记录
func (a *Account) recordTradLog(timestamp int64, mode OpMode, prize float64, vol int) {
	record := TradRecord{
		Timestamp: timestamp,
		Mode:      mode,
		Prize:     prize,
		Vol:       vol,
	}
	a.TradLog = append(a.TradLog, record)
}

// 记录账号的交易或委托等操作记录
func (a *Account) recordAction(timestamp int64, actionType ActionType, desc string) {
	action := Action{
		Mode:      actionType,
		Timestamp: timestamp,
		Desc:      desc,
	}
	a.ActionLog = append(a.ActionLog, action)
}

// 更新统计信息
func (a *Account) maintainTradStat(mode OpMode) {
	if a.LastPrize == nil {
		return
	}
	if mode == ModeBuy {
		a.TradStat.BuyCounter++
	}
	if mode == ModeShell {
		a.TradStat.SellCounter++
	}
	if a.TradStat.BuyCounter == 1 {
		a.TradStat.MinVol = a.Balance.StockVol
		a.TradStat.MinCost = a.Balance.CostRMB
	}

	if a.Balance.StockVol > a.TradStat.MaxVol {
		a.TradStat.MaxVol = a.Balance.StockVol
	}
	if a.Balance.StockVol < a.TradStat.MinVol {
		a.TradStat.MinVol = a.Balance.StockVol
	}
	if a.Balance.CostRMB > a.TradStat.MaxCost {
		a.TradStat.MaxCost = a.Balance.CostRMB
	}
	if a.Balance.CostRMB < a.TradStat.MinCost {
		a.TradStat.MinCost = a.Balance.CostRMB
	}

	total := a.LastPrize.End*float64(a.Balance.StockVol) + a.Balance.BalanceRMB // 当前账号总资产

	if total > a.TradStat.MaxValue {
		a.TradStat.MaxValue = total
	}
	if total < a.TradStat.MinValue || (a.TradStat.BuyCounter+a.TradStat.SellCounter) == 1 {
		a.TradStat.MinValue = total
	}
}
