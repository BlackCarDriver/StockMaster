package common

type StrategyMode string // 交易模式
type ActionType string   // 账号操作

const (
	modeWait      StrategyMode = "不操作"
	modeBuy       StrategyMode = "买入"
	modeShell     StrategyMode = "卖出"
	actionBuy     ActionType   = "成功买入"
	actionShell   ActionType   = "成功卖出"
	actionEntrust ActionType   = "创建委托"
	actionGiveUp  ActionType   = "放弃交易"
)

// KLineNode K线图节点
type KLineNode struct {
	Timestamp int64   `json:"sjc"` // 时间戳
	TimeDesc  string  `json:"sj"`  // 时间 格式: 2022-09-30 15:00
	Start     float64 `json:"kpj"` // 开盘价
	End       float64 `json:"spj"` // 收盘价
	Top       float64 `json:"zgj"` // 最高
	Bottom    float64 `json:"zdj"` // 最低
	Vol       float64 `json:"cje"` // 成交额
	Vov       float64 `json:"cjl"` // 成交量
	Wave      float64 `json:"zf"`  // 振幅
	PriceWave float64 `json:"zdf"` // 涨跌幅
	PriceRise float64 `json:"zde"` // 涨跌额
	HSL       float64 `json:"hsl"` // 换手率
}
