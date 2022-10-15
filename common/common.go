package common

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

type OpMode string     // 交易模式
type ActionType string // 账号操作

var log = GetLogger()

const (
	ModeWait      OpMode     = "不操作"
	ModeBuy       OpMode     = "买入"
	ModeShell     OpMode     = "卖出"
	ActionBuy     ActionType = "成功买入"
	ActionShell   ActionType = "成功卖出"
	ActionEntrust ActionType = "创建委托"
	ActionGiveUp  ActionType = "放弃交易"
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

// GetAveragePrice 获取该时段的平均价
func (m KLineNode) GetAveragePrice() (avg float64) {
	return (m.Top + m.Bottom) / 2.0
}

func init() {
	initConsoleLogger()
	//initFileLogger()
}

func initConsoleLogger() {
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLevel(logs.LevelDebug)
}

func initFileLogger() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(2)
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"./logs/%s.log"}`, time.Now().Format("20060102150405")))
}

func GetLogger() (log *logs.BeeLogger) {
	return logs.GetBeeLogger()
}

// ================ tool function ===========

// CountRiseRange 百分比涨幅计算
func CountRiseRange(before, after float64) (rise float64) {
	if before == 0 {
		return 0
	}
	rise = (after - before) / before
	return rise * 100.0
}

// RisePrizeByFlow 计算before涨跌flow个点后的结果
func RisePrizeByFlow(before float64, flow float64) (after float64) {
	after = before * ((100.0 + flow) / 100.0)
	return after
}
