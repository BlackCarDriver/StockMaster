package strategy

import "github.com/BlackCarDriver/StockMaster/common"

var log = common.GetLogger()

// Strategy 交易策略_接口描述
type Strategy interface {
	Execute(account *common.Account, stock common.KLineNode) (err error) // 根据账号状况和最新指数执行策略
	GetDesc() string                                                     // 获取策略的具体行为描述
}
