package dao

import "github.com/BlackCarDriver/StockMaster/common"

// ========================== http 接口响应结构 =======================

// http://45.push2his.eastmoney.com/api/qt/stock/kline/get?
// cb=jQuery35106032242962875369_1664948801885&
// secid=1.600036&
// ut=fa5fd1943c7b386f172d6893dbfba10b&
// fields1=f1,f2,f3,f4,f5,f6&
// fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&
// klt=15&
// fqt=1&
// beg=0&
// end=20500101&
// smplmt=460&
// lmt=1000000&
// _=1664948801985

type GetKLineResp struct {
	RC     int              `json:"rc"`
	RT     int              `json:"rt"`
	SVR    int64            `json:"svr"`
	LT     string           `json:"lt"`
	Full   int              `json:"full"`
	DlMkts string           `json:"dlmkts"`
	Data   GetKLineRespData `json:"data"`
}

type GetKLineRespData struct {
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	Decimal   int      `json:"decimal"`
	DKTotal   int      `json:"dktotal"`
	PreKPrice float64  `json:"PreKPrice"`
	KLine     []string `json:"klines"` // k线图数据
}

// https://push2.eastmoney.com/api/qt/ulist.np/get?
// cb=jQuery112305582875234802821_1665462176107&
// fltt=2&
// secids=1.000001,0.399001&
// fields=f1,f2,f3,f4,f6,f12,f13,f104,f105,f106&
// ut=b2884a393a59ad64002292a3e90d46a5&
// _=1665462176108

type GetUListParams struct {
	CallBack  string `json:"cb"`
	FLTT      int    `json:"fltt"`
	SecIDS    string `json:"secids"` // 选中的股票ID, 参考: "1.000001,0.399001"
	Fields    string `json:"fields"`
	UT        string `json:"ut"`
	Timestamp int64  `json:"_"`
}

type GetUListResp struct {
	Data GetUListRespData `json:"data"`
}

type GetUListRespData struct {
	Total int                    `json:"total"`
	Diff  []GetUListRespDataDiff `json:"diff"`
}

type GetUListRespDataDiff struct {
	F2   float64 `json:"f2"`   // 现价
	F3   float64 `json:"f3"`   // 增幅
	F4   float64 `json:"f4"`   // 增值
	F6   float64 `json:"f6"`   // 总市值
	F12  string  `json:"f12"`  // 股票代码
	F104 int     `json:"f104"` // 涨_数量
	F105 int     `json:"f105"` // 跌_数量
	F106 int     `json:"f106"` // 平_数量
}

// ========================== 图表用  =======================

type KLineData struct {
	Code       string             `json:"code"`
	Name       string             `json:"name"`
	UpdateTime int64              `json:"updateTime"` // 更新时间
	Length     int                `json:"length"`     // k线图节点数量
	From       string             `json:"from"`       // 开始时间
	To         string             `json:"to"`         // 结束时间
	KLines     []common.KLineNode `json:"kLines"`
}
