package dao

import "github.com/BlackCarDriver/GoProject-api/common/util"

// ReadKLineMockData 从指定文件中读取模拟数据
func ReadKLineMockData(path string) (mockData KLineData, err error) {
	err = util.UnmarshalJsonFromFile(path, &mockData)
	return
}
