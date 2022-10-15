package common

// BalanceInfo 账号资产信息
type BalanceInfo struct {
	BalanceRMB float64 `json:"balanceRmb"` // 可用现金
	CostRMB    float64 `json:"costRmb"`    // 持仓成本 (元)
	StockVol   int     `json:"StockVol"`   // 持有股票份额
}

type Setting struct {
	SellLock bool `json:"sellLock"` // 值为true时, 不执行卖出操作
	BuyLock  bool `json:"buyLock"`  // 值为true时, 不执行买入操作
}

// TradInfo 账户交易数据概述
type TradInfo struct {
	BuyCounter  int     `json:"buyCounter"`   // 买入次数
	SellCounter int     `json:"shellCounter"` // 卖出次数
	MaxVol      int     `json:"maxVol"`       // 历史最高持仓
	MinVol      int     `json:"minVol"`       // 历史最低持仓
	MaxCost     float64 `json:"maxCost"`      // 历史最高持仓成本
	MinCost     float64 `json:"minCost"`      // 历史最低持仓成本
	MaxValue    float64 `json:"maxValue"`     // 历史最高账户资产
	MinValue    float64 `json:"minValue"`     // 历史最低账户资产
}

// TradRecord 交易记录
type TradRecord struct {
	Timestamp int64   `json:"timestamp"` // 时间
	Mode      OpMode  `json:"mode"`      // 买或卖
	Prize     float64 `json:"prize"`     // 成交价
	Vol       int     `json:"vol"`       // 成交量
}

// Action 账号动作日志
type Action struct {
	Mode      ActionType `json:"mode"`
	Timestamp int64      `json:"timestamp"`
	Desc      string     `json:"desc"` // 具体描述
}

// Entrust 委托单
type Entrust struct {
	StarTime int64   `json:"StarTime"` // 委托时间
	DeadTime int64   `json:"deadTime"` // 失效时间 (0=一直有效)
	DealTime int64   `json:"dealTime"` // 成交时间 (0=未成交)
	Price    float64 `json:"price"`    // 委托价
	Vol      int     `json:"vol"`      // 交易份数
}
