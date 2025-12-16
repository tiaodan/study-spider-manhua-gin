// 爬虫执行器 - 配置驱动的完整爬虫流程执行器
// 整合6个步骤：找到网站 -> 爬取 -> 提取 -> 清洗 -> 验证 -> 插入DB

package spider

import (
	"fmt"
	"reflect"
	"strings"

	"study-spider-manhua-gin/src/config"
	"study-spider-manhua-gin/src/log"
)

// 执行上下文
type ExecutionContext struct {
	Website     string
	Config      *config.WebsiteConfig
	RequestData []byte
	Urls        []string
	Params      map[string]interface{} // 前端传入的参数

	// 执行结果
	RawData           interface{} // 爬取的原始数据
	ProcessedData     interface{} // 处理后的数据
	ValidationResults []*ValidationResult
	DBResults         map[string]*DBOperationResult
}

// 执行结果
type ExecutionResult struct {
	Success        bool
	Message        string
	ProcessedCount int
	ErrorDetails   []string
	DBResults      map[string]*DBOperationResult
}

// 爬虫执行器
type SpiderExecutor struct {
	configLoader    *config.SpiderConfigLoader
	strategyFactory *SpiderStrategyFactory
	fieldMapper     *FieldMapper
	dataValidator   *DataValidator
}

// NewSpiderExecutor 创建爬虫执行器
func NewSpiderExecutor() *SpiderExecutor {
	return &SpiderExecutor{
		configLoader:    config.GetSpiderConfigLoader(),
		strategyFactory: &SpiderStrategyFactory{},
		fieldMapper:     NewFieldMapper(),
		dataValidator:   NewDataValidator(),
	}
}

// Execute 执行完整的爬虫流程
func (e *SpiderExecutor) Execute(website string, requestData []byte, urls []string, params map[string]interface{}) *ExecutionResult {
	context := &ExecutionContext{
		Website:     website,
		RequestData: requestData,
		Urls:        urls,
		Params:      params,
		DBResults:   make(map[string]*DBOperationResult),
	}

	result := &ExecutionResult{
		Success:      true,
		ErrorDetails: []string{},
		DBResults:    make(map[string]*DBOperationResult),
	}

	// 步骤1: 找到目标网站配置
	if err := e.step1_LoadConfig(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤1失败: %v", err))
		return result
	}

	// 步骤2: 爬取数据
	if err := e.step2_CrawlData(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤2失败: %v", err))
		return result
	}

	// 步骤3: 提取和转换数据
	if err := e.step3_ExtractData(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤3失败: %v", err))
		return result
	}

	// 步骤4: 数据清洗和赋值
	if err := e.step4_CleanData(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤4失败: %v", err))
		return result
	}

	// 步骤5: 数据验证
	if err := e.step5_ValidateData(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤5失败: %v", err))
		return result
	}

	// 步骤6: 插入数据库
	if err := e.step6_InsertDatabase(context); err != nil {
		result.Success = false
		result.ErrorDetails = append(result.ErrorDetails, fmt.Sprintf("步骤6失败: %v", err))
		return result
	}

	// 汇总结果
	result.ProcessedCount = e.getProcessedCount(context.ProcessedData)
	result.DBResults = context.DBResults
	result.Message = fmt.Sprintf("成功处理%d条数据", result.ProcessedCount)

	return result
}

// 步骤1: 加载网站配置
func (e *SpiderExecutor) step1_LoadConfig(context *ExecutionContext) error {
	config, err := e.configLoader.GetWebsiteConfig(context.Website)
	if err != nil {
		return fmt.Errorf("加载网站配置失败: %v", err)
	}
	context.Config = config
	return nil
}

// 步骤2: 爬取数据
func (e *SpiderExecutor) step2_CrawlData(context *ExecutionContext) error {
	strategy, err := e.strategyFactory.GetStrategy(context.Config.Crawl.Type)
	if err != nil {
		return fmt.Errorf("获取爬虫策略失败: %v", err)
	}

	rawData, err := strategy.Crawl(context.RequestData, context.Config, context.Urls)
	if err != nil {
		return fmt.Errorf("数据爬取失败: %v", err)
	}

	context.RawData = rawData
	return nil
}

// 步骤3: 提取和转换数据
func (e *SpiderExecutor) step3_ExtractData(context *ExecutionContext) error {
	// 将原始数据转换为结构体列表
	processedData, err := e.convertRawDataToStructs(context.RawData, context.Config)
	if err != nil {
		return fmt.Errorf("数据转换失败: %v", err)
	}

	log.Infof("v2 step3 扁平化后数据条数: %d", sliceLen(processedData))
	context.ProcessedData = processedData
	return nil
}

// 步骤4: 数据清洗和赋值
func (e *SpiderExecutor) step4_CleanData(context *ExecutionContext) error {
	if context.Config.Clean == nil {
		return nil // 没有清洗配置，跳过
	}

	// 设置外键字段
	if err := e.setForeignKeys(context.ProcessedData, context.Config.Clean.ForeignKeys, context.Params); err != nil {
		return fmt.Errorf("设置外键失败: %v", err)
	}

	// 设置默认值
	if err := e.setDefaultValues(context.ProcessedData, context.Config.Clean.Defaults); err != nil {
		return fmt.Errorf("设置默认值失败: %v", err)
	}

	return nil
}

// 步骤5: 数据验证
func (e *SpiderExecutor) step5_ValidateData(context *ExecutionContext) error {
	before := sliceLen(context.ProcessedData)
	validationResults := e.dataValidator.ValidateBatchWithConfig(context.ProcessedData, context.Config)
	context.ValidationResults = validationResults

	// 过滤出有效数据
	validData := e.dataValidator.FilterValid(context.ProcessedData, validationResults)
	context.ProcessedData = validData
	after := sliceLen(context.ProcessedData)
	log.Infof("v2 step5 校验前条数=%d，校验后有效条数=%d", before, after)

	return nil
}

// 步骤6: 插入数据库
func (e *SpiderExecutor) step6_InsertDatabase(context *ExecutionContext) error {
	// 如果没有可入库的数据，直接返回错误，阻止后续 upsert
	if isEmptySlice(context.ProcessedData) {
		log.Warnf("v2 step6 ProcessedData 为空，类型=%T", context.ProcessedData)
		return fmt.Errorf("无有效数据，跳过数据库操作")
	}

	log.Infof("v2 step6 入库数据条数: %d", sliceLen(context.ProcessedData))
	dbResult := ExecuteDBWithConfig(context.ProcessedData, context.Config)
	context.DBResults["main"] = dbResult

	if len(dbResult.Errors) > 0 {
		return fmt.Errorf("数据库操作失败: %v", dbResult.Errors[0])
	}

	return nil
}

// convertRawDataToStructs 将原始数据转换为结构体列表
func (e *SpiderExecutor) convertRawDataToStructs(rawData interface{}, config *config.WebsiteConfig) (interface{}, error) {
	// 将 [][]T 扁平化为 []T，便于后续校验、入库
	v := reflect.ValueOf(rawData)
	if v.Kind() != reflect.Slice {
		return rawData, nil
	}

	// 判断是否是嵌套切片
	if v.Len() == 0 {
		return rawData, nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Slice {
		// 扁平化
		elemType := first.Type().Elem()
		flat := reflect.MakeSlice(reflect.SliceOf(elemType), 0, v.Len()*2)
		for i := 0; i < v.Len(); i++ {
			inner := v.Index(i)
			if inner.Kind() != reflect.Slice {
				continue
			}
			for j := 0; j < inner.Len(); j++ {
				flat = reflect.Append(flat, inner.Index(j))
			}
		}
		return flat.Interface(), nil
	}

	return rawData, nil
}

// setForeignKeys 设置外键字段
func (e *SpiderExecutor) setForeignKeys(data interface{}, foreignKeys map[string]string, params map[string]interface{}) error {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("数据必须是切片类型")
	}

	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		if item.Kind() != reflect.Struct {
			continue
		}

		// 为每个结构体设置外键
		for fieldName, paramKey := range foreignKeys {
			if paramValue, exists := params[paramKey]; exists {
				field, ok := getFieldByConfigName(item, fieldName)
				if !ok || !field.CanSet() {
					continue
				}
				convertedValue, err := e.convertValueForField(paramValue, field.Type())
				if err != nil {
					return fmt.Errorf("转换外键字段 %s 失败: %v", fieldName, err)
				}
				field.Set(reflect.ValueOf(convertedValue))
			}
		}
	}

	return nil
}

// setDefaultValues 设置默认值
func (e *SpiderExecutor) setDefaultValues(data interface{}, defaults map[string]interface{}) error {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("数据必须是切片类型")
	}

	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		if item.Kind() != reflect.Struct {
			continue
		}

		// 为每个结构体设置默认值
		for fieldName, defaultValue := range defaults {
			field, ok := getFieldByConfigName(item, fieldName)
			if !ok || !field.CanSet() || !field.IsZero() {
				continue
			}
			convertedValue, err := e.convertValueForField(defaultValue, field.Type())
			if err != nil {
				return fmt.Errorf("转换默认值字段 %s 失败: %v", fieldName, err)
			}
			field.Set(reflect.ValueOf(convertedValue))
		}
	}

	return nil
}

// convertValueForField 类型转换辅助函数
func (e *SpiderExecutor) convertValueForField(value interface{}, targetType reflect.Type) (interface{}, error) {
	valueType := reflect.TypeOf(value)

	if valueType.AssignableTo(targetType) {
		return value, nil
	}

	// 处理常见转换
	switch targetType.Kind() {
	case reflect.Int:
		switch v := value.(type) {
		case float64:
			return int(v), nil
		case int64:
			return int(v), nil
		case int32:
			return int(v), nil
		case int16:
			return int(v), nil
		case int8:
			return int(v), nil
		case uint64:
			return int(v), nil
		case uint32:
			return int(v), nil
		case uint16:
			return int(v), nil
		case uint8:
			return int(v), nil
		}
	case reflect.String:
		return fmt.Sprintf("%v", value), nil
	}

	return value, nil // 如果无法转换，返回原值
}

// getProcessedCount 获取处理的数据条数
func (e *SpiderExecutor) getProcessedCount(data interface{}) int {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Slice {
		return dataValue.Len()
	}
	return 0
}

// sliceLen 安全获取切片长度
func sliceLen(data interface{}) int {
	if data == nil {
		return 0
	}
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return 0
	}
	return v.Len()
}

// getFieldByConfigName 兼容配置字段名（小写/下划线）到导出字段名
func getFieldByConfigName(item reflect.Value, name string) (reflect.Value, bool) {
	if !item.IsValid() || item.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	if field := item.FieldByName(name); field.IsValid() {
		return field, true
	}
	camel := toExportedCamelExec(name)
	if field := item.FieldByName(camel); field.IsValid() {
		return field, true
	}
	return reflect.Value{}, false
}

// toExportedCamelExec 将 name 或 spider_end_status 转换为 Name、SpiderEndStatus
func toExportedCamelExec(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	joined := strings.Join(parts, "")
	if joined == "" {
		return strings.ToUpper(name[:1]) + name[1:]
	}
	return strings.ToUpper(joined[:1]) + joined[1:]
}

// isEmptySlice 判断数据是否为 nil 或空切片
func isEmptySlice(data interface{}) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return false
	}
	return v.Len() == 0
}
