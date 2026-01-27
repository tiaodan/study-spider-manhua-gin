// 数据库操作器 - 配置驱动的数据库操作
// 支持多种插入策略和批量操作
/*
使用方式：

*/

package spider

import (
	"fmt"
	"reflect"
	"time"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

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
	log.Warn("-- 放在警告,打印 tableName = ", tableName)
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
	log.Debug("-- 放在警告,打印 tableName = ", tableName)
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

/*
// BatchExecute 批量执行多个表的操作 - v0.1 写法
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
*/

// BatchExecute 批量执行多个表的操作 - v0.2 写法，支持关联表操作
func BatchExecute(dataMap map[string]interface{}, config *config.WebsiteConfig) map[string]*DBOperationResult {
	results := make(map[string]*DBOperationResult)

	// 主表操作
	var mainData interface{}
	if data, exists := dataMap["main"]; exists {
		mainData = data
		results["main"] = ExecuteDBWithConfig(mainData, config)

		// 主表插入成功后，回填 ID（照着 V1 抄）
		if mainResult, ok := results["main"]; ok && len(mainResult.Errors) == 0 {
			if err := fillComicIdsV1(mainData); err != nil {
				log.Warnf("回填 Comic ID 失败: %v", err)
			}
		}
	}

	// 关联表操作
	if config.RelatedTables != nil {
		for tableName, tableConfig := range config.RelatedTables {
			// 从主表数据中提取关联表数据
			relatedData := extractRelatedData(mainData, tableConfig)
			if relatedData != nil {
				results[tableName] = ExecuteDBWithTableConfig(relatedData, tableConfig.Insert, tableConfig.Table)
			}
		}
	}

	return results

	/* // BatchExecute 批量执行多个表的操作 - v0.1 写法
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
	*/
}

// extractRelatedData 从主表数据中提取关联表数据
func extractRelatedData(mainData interface{}, tableConfig *config.RelatedTableConfig) interface{} {
	if mainData == nil {
		return nil
	}

	if tableConfig.Source == "field" {
		// 从主表对象的字段中提取数据
		dataValue := reflect.ValueOf(mainData)
		if dataValue.Kind() == reflect.Slice {
			// 如果是切片，遍历每个元素提取字段
			statsSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(&models.ComicSpiderStats{}).Elem()), 0, dataValue.Len())
			for i := 0; i < dataValue.Len(); i++ {
				item := dataValue.Index(i)
				if item.Kind() == reflect.Ptr {
					item = item.Elem()
				}
				if item.Kind() != reflect.Struct {
					continue
				}

				// 通过反射获取指定字段
				field := item.FieldByName(tableConfig.SourcePath)
				if field.IsValid() && !field.IsZero() {
					// 创建新的 ComicSpiderStats 对象，设置 ComicId
					stats := models.ComicSpiderStats{}
					starField := field.FieldByName("Star")
					if starField.IsValid() {
						stats.Star = starField.Float()
					}
					latestChapterNameField := field.FieldByName("LatestChapterName")
					if latestChapterNameField.IsValid() {
						stats.LatestChapterName = latestChapterNameField.String()
					}
					hitsField := field.FieldByName("Hits")
					if hitsField.IsValid() {
						stats.Hits = int(hitsField.Int())
					}
					totalChapterField := field.FieldByName("TotalChapter")
					if totalChapterField.IsValid() {
						stats.TotalChapter = int(totalChapterField.Int())
					}
					lastestChapterReleaseDateField := field.FieldByName("LastestChapterReleaseDate")
					if lastestChapterReleaseDateField.IsValid() {
						stats.LastestChapterReleaseDate = lastestChapterReleaseDateField.Interface().(time.Time)
					}
					latestChapterIdField := field.FieldByName("LatestChapterId")
					if latestChapterIdField.IsValid() && !latestChapterIdField.IsNil() {
						val := int(latestChapterIdField.Int())
						stats.LatestChapterId = &val
					}

					// 设置 ComicId（从主表的 Id 字段获取）
					idField := item.FieldByName("Id")
					if idField.IsValid() && idField.Int() > 0 {
						stats.ComicId = int(idField.Int())
					}

					// 应用数据清洗
					stats.TrimSpaces()
					stats.Trad2Simple()

					statsSlice = reflect.Append(statsSlice, reflect.ValueOf(stats))
				}
			}
			return statsSlice.Interface()
		}
	}

	return nil
}

// ExecuteDBWithTableConfig 使用表配置执行数据库操作
func ExecuteDBWithTableConfig(dataList interface{}, insertConfig *config.InsertConfig, tableName string) *DBOperationResult {
	if insertConfig == nil {
		return &DBOperationResult{
			SuccessCount: 0,
			FailedCount:  0,
			Errors:       []error{fmt.Errorf("没有数据库插入配置")},
		}
	}

	operator := NewDBOperator(insertConfig)
	return operator.Execute(dataList, tableName)
}

// fillComicIdsV1 回填 ComicSpider 的 ID（照着 V1 抄）
func fillComicIdsV1(data interface{}) error {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("数据必须是切片类型")
	}

	if dataValue.Len() == 0 {
		return nil
	}

	// 照着 V1 抄：用 Where 条件查询获取 ID
	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			continue
		}

		comic, ok := item.Addr().Interface().(*models.ComicSpider)
		if !ok {
			continue
		}

		// 照着 V1 抄：构建查询条件
		var existingComic models.ComicSpider
		condition := map[string]interface{}{
			"name":          comic.Name,
			"country_id":    comic.CountryId,
			"website_id":    comic.WebsiteId,
			"porn_type_id":  comic.PornTypeId,
			"type_id":       comic.TypeId,
			"author_concat": comic.AuthorConcat,
		}
		result := db.DBComic.Where(condition).First(&existingComic)
		if result.Error == nil {
			// 更新对象的ID为数据库中的实际ID
			comic.Id = existingComic.Id
			log.Debugf("更新comic ID: %s -> %d", comic.Name, existingComic.Id)
		} else {
			log.Errorf("查询comic失败: %s, err: %v", comic.Name, result.Error)
		}
	}

	return nil
}
