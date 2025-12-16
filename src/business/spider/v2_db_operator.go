// 数据库操作器 - 配置驱动的数据库操作
// 支持多种插入策略和批量操作

package spider

import (
	"fmt"
	"reflect"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"

	"gorm.io/gorm"
)

// 数据库操作结果
type DBOperationResult struct {
	SuccessCount int
	FailedCount  int
	Errors       []error
}

// 数据库操作器
type DBOperator struct {
	insertConfig *config.InsertConfig
}

// NewDBOperator 创建数据库操作器
func NewDBOperator(insertConfig *config.InsertConfig) *DBOperator {
	return &DBOperator{
		insertConfig: insertConfig,
	}
}

// Execute 执行数据库操作
func (op *DBOperator) Execute(dataList interface{}, tableName string) *DBOperationResult {
	result := &DBOperationResult{
		SuccessCount: 0,
		FailedCount:  0,
		Errors:       []error{},
	}

	if op.insertConfig == nil {
		result.Errors = append(result.Errors, fmt.Errorf("没有插入配置"))
		return result
	}

	switch op.insertConfig.Strategy {
	case "insert":
		return op.executeInsert(dataList, tableName)
	case "update":
		return op.executeUpdate(dataList, tableName)
	case "upsert":
		return op.executeUpsert(dataList, tableName)
	default:
		result.Errors = append(result.Errors, fmt.Errorf("不支持的插入策略: %s", op.insertConfig.Strategy))
		return result
	}
}

// executeInsert 执行插入操作
func (op *DBOperator) executeInsert(dataList interface{}, tableName string) *DBOperationResult {
	result := &DBOperationResult{}

	// 获取数据库实例
	dbInstance := op.getDBInstance(tableName)
	if dbInstance == nil {
		result.Errors = append(result.Errors, fmt.Errorf("无法获取数据库实例: %s", tableName))
		return result
	}

	// 调用批量插入
	gormDB, ok := dbInstance.(*gorm.DB)
	if !ok {
		result.Errors = append(result.Errors, fmt.Errorf("数据库实例类型错误"))
		result.FailedCount = op.getDataCount(dataList)
		return result
	}

	err := db.DBUpsertBatch(gormDB, dataList, []string{}, []string{}) // 插入不指定唯一键
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("批量插入失败: %v", err))
		result.FailedCount = op.getDataCount(dataList)
	} else {
		result.SuccessCount = op.getDataCount(dataList)
	}

	return result
}

// executeUpdate 执行更新操作
func (op *DBOperator) executeUpdate(dataList interface{}, tableName string) *DBOperationResult {
	result := &DBOperationResult{}

	// 暂时不支持更新操作，返回错误
	dataListValue := reflect.ValueOf(dataList)
	if dataListValue.Kind() == reflect.Slice {
		result.FailedCount = dataListValue.Len()
	} else {
		result.FailedCount = 1
	}
	result.Errors = append(result.Errors, fmt.Errorf("更新操作暂未实现，请使用upsert策略"))
	return result
}

// executeUpsert 执行Upsert操作（插入或更新）
func (op *DBOperator) executeUpsert(dataList interface{}, tableName string) *DBOperationResult {
	result := &DBOperationResult{}

	// 获取数据库实例
	dbInstance := op.getDBInstance(tableName)
	if dbInstance == nil {
		result.Errors = append(result.Errors, fmt.Errorf("无法获取数据库实例: %s", tableName))
		return result
	}

	// 调用批量Upsert
	gormDB, ok := dbInstance.(*gorm.DB)
	if !ok {
		result.Errors = append(result.Errors, fmt.Errorf("数据库实例类型错误"))
		result.FailedCount = op.getDataCount(dataList)
		return result
	}

	err := db.DBUpsertBatch(gormDB, dataList, op.insertConfig.UniqueKeys, op.insertConfig.UpdateKeys)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("批量Upsert失败: %v", err))
		result.FailedCount = op.getDataCount(dataList)
	} else {
		result.SuccessCount = op.getDataCount(dataList)
	}

	return result
}

// getDBInstance 根据表名获取数据库实例
func (op *DBOperator) getDBInstance(tableName string) interface{} {
	// 暂时都返回comic数据库，未来可以扩展为多个数据库
	return db.DBComic
}

// getDataCount 获取数据条数
func (op *DBOperator) getDataCount(dataList interface{}) int {
	dataListValue := reflect.ValueOf(dataList)
	if dataListValue.Kind() == reflect.Slice {
		return dataListValue.Len()
	}
	return 1
}

// ExecuteWithConfig 根据配置执行操作
func ExecuteDBWithConfig(dataList interface{}, config *config.WebsiteConfig) *DBOperationResult {
	if config.Insert == nil {
		return &DBOperationResult{
			SuccessCount: 0,
			FailedCount:  0,
			Errors:       []error{fmt.Errorf("没有数据库插入配置")},
		}
	}

	operator := NewDBOperator(config.Insert)
	tableName := config.Meta.Table
	return operator.Execute(dataList, tableName)
}

// BatchExecute 批量执行多个表的操作
func BatchExecute(dataMap map[string]interface{}, config *config.WebsiteConfig) map[string]*DBOperationResult {
	results := make(map[string]*DBOperationResult)

	// 主表操作
	if mainData, exists := dataMap["main"]; exists {
		results["main"] = ExecuteDBWithConfig(mainData, config)
	}

	// 关联表操作（暂时不支持）
	if _, exists := dataMap["related"]; exists {
		results["related"] = &DBOperationResult{
			SuccessCount: 0,
			FailedCount:  0,
			Errors:       []error{fmt.Errorf("关联表操作暂未实现")},
		}
	}

	return results
}
