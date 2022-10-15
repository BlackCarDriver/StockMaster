package main

import (
	"github.com/BlackCarDriver/StockMaster/common"
	"github.com/BlackCarDriver/StockMaster/handler"
)

var log = common.GetLogger()

func main() {
	handler.TestGridStrategy()
}
