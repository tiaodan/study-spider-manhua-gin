// db website 相关操作
package db

import (
	"reflect"
	"strings"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"
	"study-spider-manhua-gin/util/stringutil"

	// 导入 clause 包
	"gorm.io/gorm/clause"
)

// 初始化
// 为 Website 表实现完整的操作接口
type WebsiteOperations struct{}

// 增
func (w WebsiteOperations) Add(website *models.Website) error {
	// 预处理空格
	website.Name = strings.TrimSpace(website.Name)
	website.Url = strings.TrimSpace(website.Url)

	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "NameId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":       website.Name,
			"url":        website.Url,
			"need_proxy": website.NeedProxy,
			"is_https":   website.IsHttps,
		}),
	}).Create(website)
	if result.Error != nil {
		log.Error("创建失败:", result.Error)
		return result.Error
	} else {
		log.Info("创建成功:", website)
	}
	return nil
}

// 批量增
func (w WebsiteOperations) BatchAdd(websites []*models.Website) {
	for i, website := range websites {
		err := w.Add(website)
		if err == nil {
			log.Debugf("批量创建第%d条成功, website: %v", i+1, website.Name)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删-通过id
func (w WebsiteOperations) DeleteById(id uint) {
	var website models.Website
	result := DB.Delete(&website, id)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", id)
	}
}

// 删-通过 nameId
func (W WebsiteOperations) DeleteByNameId(nameId any) {
	var website models.Website
	result := DB.Where("name_id = ?", nameId).Delete(&website)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", nameId)
	}
}

// 删-通过其它
func (w WebsiteOperations) DeleteByOther(condition string, other any) {
	var website models.Website
	result := DB.Where(condition+" = ?", other).Delete(&website)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", other)
	}
}

// 批量删-通过id
func (w WebsiteOperations) BatchDeleteById(ids []uint) {
	var websites []models.Website
	result := DB.Delete(&websites, ids)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", ids)
	}
}

// 批量删-通过nameIds
func (w WebsiteOperations) BatchDeleteByNameId(nameIds []int) {
	var websites []models.Website
	result := DB.Where("name_id IN ?", nameIds).Delete(&websites)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", nameIds)
	}
}

// 批量删-通过other
func (w WebsiteOperations) BatchDeleteByOther(condition string, others []any) {
	var websites []models.Website
	// result := DB.Where(condition+" IN ?", others).Delete(&websites) // other这样写错？
	result := DB.Where(condition+" IN ?", others).Delete(&websites)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", others)
	}
}

// 改 - by Id
func (w WebsiteOperations) UpdateById(id uint, update map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(update)

	var website models.Website
	// 解决0值不更新
	result := DB.Model(&website).Where("id = ?", id).Select("name", "url", "need_proxy", "is_https").Updates(update)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", id)
	}
}

// 改 - by nameId
func (w WebsiteOperations) UpdateByNameId(nameId int, update map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(update)

	var website models.Website
	// 解决0值不更新
	result := DB.Model(&website).Where("name_id = ?", nameId).Select("name", "url", "need_proxy", "is_https").Updates(update)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", nameId)
	}
}

// 改 - by other
func (w WebsiteOperations) UpdateByOther(condition string, other any, update map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(update)

	var website models.Website
	// 解决0值不更新
	// result := DB.Model(&website).Where("name_id = ?", nameId).Select("name", "url", "need_proxy", "is_https").Updates(update)  // 之前写法
	result := DB.Model(&website).Where(condition+" = ?", other).Select("name", "url", "need_proxy", "is_https").Updates(update) // 之前写法
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", other)
	}
}

// 改 - 批量 by id
func (w WebsiteOperations) BatchUpdateById(updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where("id = ?", update["Id"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by nameId
func (w WebsiteOperations) BatchUpdateByNameId(updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where("name_id = ?", update["NameId"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by other
func (w WebsiteOperations) BatchUpdateByOther(condition string, others any, updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where(condition+" IN ?", others).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 查 - by id
func (w WebsiteOperations) QueryById(id uint) *models.Website {
	var website models.Website
	result := DB.First(&website, id)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 查 - by nameId
func (w WebsiteOperations) QueryByNameId(nameId int) *models.Website {
	var website models.Website
	result := DB.Where("name_id = ?", nameId).First(&website)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 查 - by other
func (w WebsiteOperations) QueryByOther(condition string, other any) *models.Website {
	var website models.Website
	result := DB.Where(condition+" = ?", other).First(&website)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 批量查 - by ids
func (w WebsiteOperations) BatchQueryById(ids []uint) ([]*models.Website, error) {
	var websites []*models.Website
	result := DB.Find(&websites, ids)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// 批量查 - by nameIds
func (w WebsiteOperations) BatchQueryByNameId(nameIds []int) ([]*models.Website, error) {
	var websites []*models.Website
	result := DB.Where("name_id IN ?", nameIds).Order("name_id").Find(&websites) // 默认升序
	// result := DB.Where("name_id IN ?", nameIds).Order("name_id DESC")Find(&websites) // 倒序排列
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// 批量查 - by others
// 参数：orderby 排序字符串 如: name_id   sort 排序方式，ASC 为正序，DESC 为倒序
func (w WebsiteOperations) BatchQueryByOther(condition string, others []any, orderby string, sort string) ([]*models.Website, error) {
	var websites []*models.Website
	// result := DB.Where(condition+" IN ?", others).Find(&websites)  // other这样写错？
	// result := DB.Where(condition+" IN ?", others).Order("name_id DESC").Find(&websites)
	result := DB.Where(condition+" IN ?", others).Order(orderby + " " + sort).Find(&websites)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// 批量查 - all
// 参数：orderby 排序字符串 如: name_id   sort 排序方式，ASC 为正序，DESC 为倒序
func (w WebsiteOperations) BatchQueryAll() ([]*models.Website, error) {
	var websites []*models.Website
	result := DB.Find(&websites)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// ----------- 测试用例封装 start --------------
// 根据给的对象， 生成有0值, 单个为0 的对象 (不判断Id + NameId)
// 参数: 有id/无id 的 int值全是1的对象
// 返回：有0值, 单个为0 的对象arr
// 思路：1 有id一种思路 2 无id一种思路。不用区分id
func (w WebsiteOperations) returnObjZeroOne(obj models.Website) []models.Website {
	var forAddHasZeroOneArr []models.Website
	// 复制一份
	forAddHasZeroOne := obj

	// 使用反射修改 int 类型字段为 0
	value := reflect.ValueOf(&forAddHasZeroOne).Elem()
	typ := value.Type()

	// 不判断Id(index: 0) + NameId(index: 1),i 从2开始
	typNum := typ.NumField()
	if typNum < 2 {
		return forAddHasZeroOneArr
	}
	for i := 2; i < typNum; i++ {

		fieldValue := value.Field(i)

		// 检查字段类型是否为 int
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 将字段值设置为0
			fieldValue.SetInt(0) // 一般要判断: if fieldValue.CanSet()
			// 加入数组
			forAddHasZeroOneArr = append(forAddHasZeroOneArr, forAddHasZeroOne)

			// 重置变量
			forAddHasZeroOne = obj
		}
	}
	return forAddHasZeroOneArr

	// v1 写法 -------------------- start
	// 	// // 复制一份
	// websiteForAddHasIdHasZeroOne := websiteForAddHasIdNoZero

	// // 使用反射修改 int 类型字段为 0
	// value := reflect.ValueOf(&websiteForAddHasIdHasZeroOne).Elem()
	// typ := value.Type()

	// for i := 0; i < typ.NumField(); i++ {
	// 	fieldValue := value.Field(i)

	// 	// 检查字段类型是否为 int
	// 	switch fieldValue.Kind() {
	// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 		// 将字段值设置为0
	// 		fieldValue.SetInt(0) // 一般要判断: if fieldValue.CanSet()
	// 		// 加入数组
	// 		websitesForAddHasIdHasZeroOne = append(websitesForAddHasIdHasZeroOne, websiteForAddHasIdHasZeroOne)

	// 		// 重置变量
	// 		websiteForAddHasIdHasZeroOne = websiteForAddHasIdNoZero
	// 	}
	// }
	// v1 写法 -------------------- end
}

// ----------- 测试用例封装 start --------------
// 根据给的对象， 给int类型对象，取反 (不判断Id + NameId) - 用于生成updates 相关的。生成一组测试用例
// 参数: 有id/无id 的 int值全是1的对象  map对象
// 返回：有0值, 单个为0 的对象arr
// 思路：1 有id一种思路 2 无id一种思路。不用区分id
func (w WebsiteOperations) returnObjZeroOneNegate(obj map[string]any) []map[string]any {
	var arr []map[string]any
	// 复制一份
	hasZeroOne := obj

	for key, value := range obj {
		if strings.ToLower(key) == "id" || strings.ToLower(key) == "nameid" {
			continue
		}
		// 判断key是否为int类型
		switch value.(type) { // 别的写法vType := value.(type)
		case int:
			switch value {
			case 0:
				hasZeroOne[key] = 1 // 将字段值设置为1
			case 1:
				hasZeroOne[key] = 0
			default:
				hasZeroOne[key] = 0 // 可能== 2 8 等数字
			}
			// 加入数组
			arr = append(arr, hasZeroOne)

			// 重置变量
			hasZeroOne = obj
		}
	}
	return arr
}

// 根据给的对象， 给int类型对象，取反 (不判断Id + NameId) - 用于生成updates 相关的.生成一个测试用例
// 参数: 有id/无id 的 int值全是1的对象  map对象
// 返回：有0值, 单个为0 的对象arr
// 思路：1 有id一种思路 2 无id一种思路。不用区分id
func (w WebsiteOperations) returnObjZeroAllNegate(obj map[string]any) map[string]any {
	// 复制一份
	hasZeroOne := obj

	for key, value := range obj {
		if strings.ToLower(key) == "id" || strings.ToLower(key) == "nameid" {
			continue
		}
		// 判断key是否为int类型
		switch value.(type) { // 别的写法vType := value.(type)
		case int:
			switch value {
			case 0:
				hasZeroOne[key] = 1 // 将字段值设置为1
			case 1:
				hasZeroOne[key] = 0
			default:
				hasZeroOne[key] = 0 // 可能== 2 8 等数字
			}
		}
	}
	return hasZeroOne
}

// 根据给的对象， 生成有0值, all为0 的对象 (不判断Id + NameId)
// 参数: 有id/无id 的 int值全是1的对象
// 返回：有0值, 单个为0 的对象arr
// 思路：1 有id一种思路 2 无id一种思路。不用区分id
func (w WebsiteOperations) returnObjZeroAll(obj models.Website) models.Website {
	// 复制一份
	forAddHasZeroAll := obj

	// 使用反射修改 int 类型字段为 0
	value := reflect.ValueOf(&forAddHasZeroAll).Elem()
	typ := value.Type()

	// 不判断Id(index: 0) + NameId(index: 1),i 从2开始
	typNum := typ.NumField()
	if typNum < 2 {
		return forAddHasZeroAll
	}
	for i := 2; i < typNum; i++ {
		fieldValue := value.Field(i)

		// 检查字段类型是否为 int
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 将字段值设置为0
			fieldValue.SetInt(0) // 一般要判断: if fieldValue.CanSet()
		}
	}
	return forAddHasZeroAll

	// v1.0 写法 -------------------- start
	// // 复制一份
	// websiteForAddHasIdHasZeroAll = websiteForAddHasIdNoZero

	// // 使用反射修改 int 类型字段为 0
	// value = reflect.ValueOf(&websiteForAddHasIdHasZeroAll).Elem()
	// typ = value.Type()
	// // fmt.Println("value = ", value) // {1 1 Test Website Add http://add.com 1 1}
	// // fmt.Println("typ = ", typ)     // typ =  models.Website

	// for i := 0; i < typ.NumField(); i++ {
	// 	fieldValue := value.Field(i)
	// 	// fmt.Println("fieldValue = ", fieldValue)  // 测试打印
	// 	// fmt.Println("fieldValue.Kind() = ", fieldValue.Kind()) // 测试打印

	// 	// 检查字段类型是否为 int
	// 	switch fieldValue.Kind() {
	// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 		// 将字段值设置为0
	// 		fieldValue.SetInt(0) // 一般要判断: if fieldValue.CanSet()
	// 	}
	// }
	// v1.0 写法 -------------------- end
}

// ----------- 测试用例封装 end --------------
