// 网格交易策略

package strategy

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
