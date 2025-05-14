package db

import (
	"study-spider-manhua-gin/models"
	"testing"
)

// 所有方法思路
// 1. 清空表
// 2. 添加数据
// 3. 增删改查、批量增删改查
// 4. 检测

// ---------------------------- 变量 start ----------------------------

// 用于 add
var websiteForAddHasIdNoZero *models.Website  // 用于add, 有id, 无0值
var websiteForAddHasIdHasZero *models.Website // 用于add, 有id, 无0值
var websiteForAddNoIdNoZero *models.Website   // 用于add, 无id, 无0值
var websiteForAddNoIdHasZero *models.Website  // 用于add, 无id, 无0值

// 用于 add batch
var website2ForAddHasIdNoZero *models.Website  // 用于add, 有id, 无0值
var website2ForAddHasIdHasZero *models.Website // 用于add, 有id, 无0值
var website2ForAddNoIdNoZero *models.Website   // 用于add, 无id, 无0值
var website2ForAddNoIdHasZero *models.Website  // 用于add, 无id, 无0值

// 用于 update
var websiteForUpdateHasIdNoZero map[string]interface{}  // 用于update, 有id, 无0值
var websiteForUpdateHasIdHasZero map[string]interface{} // 用于update, 有id, 有0值

var websiteForUpdateNoIdNoZero map[string]interface{}  // 用于update, 无id, 无0值
var websiteForUpdateNoIdHasZero map[string]interface{} // 用于update, 无id, 有0值

// 用于 update batch
var website2ForUpdateHasIdNoZero map[string]interface{}  // 用于update, 有id, 无0值
var website2ForUpdateHasIdHasZero map[string]interface{} // 用于update, 有id, 有0值

var website2ForUpdateNoIdNoZero map[string]interface{}  // 用于update, 无id, 无0值
var website2ForUpdateNoIdHasZero map[string]interface{} // 用于update, 无id, 有0值

// ---------------------------- 变量 end ----------------------------

// ---------------------------- init start ----------------------------
func init() {
	// 用于add, 有id
	websiteForAddHasIdNoZero = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	websiteForAddHasIdHasZero = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add, 无id
	websiteForAddNoIdNoZero = &models.Website{
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	websiteForAddNoIdHasZero = &models.Website{
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add batch, 有id
	website2ForAddHasIdNoZero = &models.Website{
		Id:        2, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2ForAddHasIdHasZero = &models.Website{
		Id:        2, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add batch, 无id
	website2ForAddNoIdNoZero = &models.Website{
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2ForAddNoIdHasZero = &models.Website{
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于update
	websiteForUpdateHasIdNoZero = map[string]interface{}{
		"Id":        uint(1),
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	websiteForUpdateHasIdHasZero = map[string]interface{}{
		"Id":        uint(1),
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	// 无id
	websiteForUpdateNoIdNoZero = map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	websiteForUpdateNoIdHasZero = map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}

	// 用于update batch
	website2ForUpdateHasIdNoZero = map[string]interface{}{
		"Id":        uint(2),
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	website2ForUpdateHasIdHasZero = map[string]interface{}{
		"Id":        uint(2),
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	// 无id
	website2ForUpdateNoIdNoZero = map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	website2ForUpdateNoIdHasZero = map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
}

// ---------------------------- init end ----------------------------

// 检测函数封装, 对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckHasId(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.Id != obj.Id ||
		query.NameId != obj.NameId ||
		query.Name != obj.Name ||
		query.Url != obj.Url ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 检测函数封装, 不对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckNoId(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.NameId != obj.NameId ||
		query.Name != obj.Name ||
		query.Url != obj.Url ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 检测更新函数封装, 对比Id
// 参数1: 查到的指针 参数2: 更新参数 map[string]interface{}
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckUpdateHasId(query *models.Website, obj map[string]interface{}, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.Id != obj["Id"] ||
		query.NameId != obj["NameId"] ||
		query.Name != obj["Name"] ||
		query.Url != obj["Url"] ||
		query.NeedProxy != obj["NeedProxy"] ||
		query.IsHttps != obj["IsHttps"] {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 检测更新函数封装, 不对比Id
// 参数1: 查到的指针 参数2: 更新参数 map[string]interface{}
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckUpdateNoId(query *models.Website, obj map[string]interface{}, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.NameId != obj["NameId"] ||
		query.Name != obj["Name"] ||
		query.Url != obj["Url"] ||
		query.NeedProxy != obj["NeedProxy"] ||
		query.IsHttps != obj["IsHttps"] {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 增
func TestWebsiteAdd(t *testing.T) {
	t.Log("------------ website add ...  start ")

	// 1. 测试项1，有0值
	website := websiteForAddHasIdHasZero
	t.Log("website: ", website)
	WebsiteAdd(website)

	// var createdWebsite *models.Website // 手动写法
	// testDB.Where("name_id = ?", website.NameId).First(&createdWebsite) // 手动写法,不调用方法
	createdWebsite := WebsiteQueryByNameId(website.NameId) // 调用方法
	WebsiteCheckNoId(createdWebsite, website, t, "【增】")    // 测试第1个

	// 2. 测试项2 无0值, 测试needProxy =1 时候
	website2 := websiteForAddHasIdNoZero
	t.Log("website2: ", website2)
	WebsiteAdd(website2)
	createdWebsite2 := WebsiteQueryByNameId(website2.NameId) // 调用方法
	WebsiteCheckNoId(createdWebsite2, website2, t, "【增】")    // 测试第2个

	t.Log("----------- website add ... end ----------------")
}

// 批量增
func TestWebsiteBatchAdd(t *testing.T) {
	t.Log("------------ website batch add ... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{})

	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero

	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)
	nameIds := []int{website.NameId, website2.NameId}
	t.Log("namedis = ", nameIds)
	createdWebsites, err := WebsitesBatchQueryByNameId(nameIds) // 调用方法
	if err != nil {
		t.Errorf("【增-批量】测试不通过, 查询nil, got=  %v", createdWebsites)
		panic("【增-批量】测试不通过")
	}

	// 判断第1个
	createdWebsite := createdWebsites[0]
	createdWebsite2 := createdWebsites[1]
	t.Log("查询结果 createdWebsites = ", createdWebsite, createdWebsite2)
	WebsiteCheckNoId(createdWebsite, website, t, "【增-批量 】")   // 判断第1个
	WebsiteCheckNoId(createdWebsite2, website2, t, "【增-批量 】") // 判断第2个
	t.Log("------------ website batch add ... end ----------------")
}

// 删-通过id
func TestWebsiteDeleteById(t *testing.T) {
	t.Log("------------ website delete by id... start ----------------")
	website := websiteForAddHasIdNoZero
	WebsiteAdd(website)

	WebsiteDeleteById(website.Id)

	var deletedWebsite models.Website
	result := testDB.First(&deletedWebsite, website.Id)
	// result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by id】测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		panic("【删 - by id 】测试不通过,删除后仍能查到")
	}
	t.Log("------------ website delete by id... end ----------------")
}

// 删-通过 nameId
func TestWebsiteDeleteByNameId(t *testing.T) {
	t.Log("------------ website delete by nameId... start ----------------")
	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	WebsiteDeleteByNameId(website.NameId)

	var deletedWebsite models.Website
	result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by nameId 】 测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		panic("【删 - by nameId 】 测试不通过,删除后仍能查到")
	}
	t.Log("------------ website delete by nameId... end ----------------")
}

// 删-通过 其它
func TestWebsiteDeleteByOther(t *testing.T) {
	t.Log("------------ website delete by other... start ----------------")
	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	WebsiteDeleteByOther("name_id", website.NameId)

	var deletedWebsite models.Website
	result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by other 】测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		panic("【删 - by other 】测试不通过,删除后仍能查到")
	}
	t.Log("------------ website delete by other... end ----------------")
}

// 删-批量 通过id
func TestWebsitesBatchDeleteById(t *testing.T) {
	t.Log("------------ website batch delete by id... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	website := websiteForAddHasIdNoZero

	website2 := &models.Website{
		Id:        2,
		NameId:    2,
		Name:      "Test Website for Delete By Id 2",
		Url:       "http://delete.com id 2",
		NeedProxy: 0,
		IsHttps:   0,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites) // 添加

	ids := []uint{website.Id, website2.Id}
	t.Log("ids = ", ids)

	// 判断是否添加了2个
	websites, err := WebsitesBatchQueryById(ids)
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", websites)
		// panic("【删 批量 - by id 】测试不通过,删除后仍能查到") // 测试v原本不能用pnic
	}

	WebsitesBatchDeleteById(ids) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsitesBatchQueryById(ids)

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", websites)
	}

	t.Log("------------ website batch delete by id... end ----------------")
}

// 删-批量 通过 nameId
func TestWebsitesBatchDeleteByNameId(t *testing.T) {
	t.Log("------------ website batch delete by nameId... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero

	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites) // 添加

	nameIds := []int{website.NameId, website2.NameId}
	t.Log("nameIds = ", nameIds)

	// 判断是否添加了2个
	websites, err := WebsitesBatchQueryByNameId(nameIds)
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", websites)
	}

	WebsitesBatchDeleteByNameId(nameIds) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsitesBatchQueryByNameId(nameIds)

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", websites)
	}

	t.Log("------------ website batch delete by nameId... end ----------------")
}

// 删-批量 通过 other
func TestWebsitesBatchDeleteByOther(t *testing.T) {
	t.Log("------------ website batch delete by other... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero

	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites) // 添加

	others := []any{website.NameId, website2.NameId}
	t.Log("others = ", others)

	// 判断是否添加了2个
	websites, err := WebsitesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", websites)
	}

	WebsitesBatchDeleteByOther("name_id", others) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsitesBatchQueryByOther("name_id", others, "name_id", "ASC")

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", websites)
	}

	t.Log("------------ website batch delete by other... end ----------------")
}

// 改 by id
func TestWebsiteUpdateById(t *testing.T) {
	t.Log("------------ website update by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddHasIdNoZero
	WebsiteAdd(website)

	updates := websiteForUpdateHasIdHasZero
	WebsiteUpdateById(website.Id, updates)

	// 检查
	updatedWebsite := WebsiteQueryById(website.Id)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	t.Log("更新后 查的 updatedWebsite.== =", updatedWebsite.Id == updates["Id"]) // 得转成uint
	WebsiteCheckUpdateHasId(updatedWebsite, updates, t, "【改 by id 】")
	t.Log("------------ website update by id ... end ")
}

// 改 by nameId
func TestWebsiteUpdateByNameId(t *testing.T) {
	t.Log("------------ website update by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	updates := websiteForUpdateNoIdHasZero
	WebsiteUpdateByNameId(website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteQueryByNameId(website.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	WebsiteCheckUpdateNoId(updatedWebsite, updates, t, "【改 by nameId 】")
	t.Log("------------ website update by nameId ... end ")
}

// 改 by other
func TestWebsiteUpdateByOther(t *testing.T) {
	t.Log("------------ website update by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	updates := websiteForUpdateNoIdHasZero
	WebsiteUpdateByOther("name_id", website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteQueryByOther("name_id", website.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	WebsiteCheckUpdateNoId(updatedWebsite, updates, t, "【改 by other 】")
	t.Log("------------ website update by other ... end ")
}

// 改-批量 by id
func TestWebsiteBatchUpdateById(t *testing.T) {
	t.Log("------------ website batch update by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddHasIdNoZero
	website2 := website2ForAddNoIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := websiteForUpdateHasIdNoZero
	updates2 := website2ForUpdateHasIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsitesBatchUpdateById(updatesArr)

	// 检测，如果报错
	ids := []uint{1, 2}
	websites, err := WebsitesBatchQueryById(ids)
	if err != nil {
		t.Errorf("【改 by id 】测试不通过, got= %v", websites)
	}

	updatedWebsite := websites[0]
	updatedWebsite2 := websites[1]
	WebsiteCheckUpdateHasId(updatedWebsite, updates, t, "【改 by id 】")   // 检测第1个
	WebsiteCheckUpdateHasId(updatedWebsite2, updates2, t, "【改 by id 】") // 检测第1个
	t.Log("------------ website batch update by id ... end ")
}

// 改-批量 by nameId
func TestWebsiteBatchUpdateByNameId(t *testing.T) {
	t.Log("------------ website batch update by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := websiteForUpdateNoIdNoZero
	updates2 := website2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsitesBatchUpdateByNameId(updatesArr)

	// 检测，如果报错
	nameIds := []int{
		updates["NameId"].(int),
		updates2["NameId"].(int),
	}
	websites, err := WebsitesBatchQueryByNameId(nameIds)
	if err != nil {
		t.Errorf("【改 by nameId 】测试不通过, got= %v", websites)
	}

	updatedWebsite := websites[0]
	updatedWebsite2 := websites[1]
	WebsiteCheckUpdateNoId(updatedWebsite, updates, t, "【改 by nameId 】")   // 检测第1个
	WebsiteCheckUpdateNoId(updatedWebsite2, updates2, t, "【改 by nameId 】") // 检测第1个
	t.Log("------------ website batch update by nameId ... end ")
}

// 改-批量 by other
func TestWebsiteBatchUpdateByOther(t *testing.T) {
	t.Log("------------ website batch update by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Website{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := websiteForUpdateNoIdNoZero
	updates2 := website2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsitesBatchUpdateByOther(updatesArr)

	// 检测，如果报错
	others := []any{
		updates["NameId"],
		updates2["NameId"],
	}
	websites, err := WebsitesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【改 by other 】测试不通过, got= %v", websites)
	}

	updatedWebsite := websites[0]
	updatedWebsite2 := websites[1]
	WebsiteCheckUpdateNoId(updatedWebsite, updates, t, "【改 by other 】")   // 检测第1个
	WebsiteCheckUpdateNoId(updatedWebsite2, updates2, t, "【改 by other 】") // 检测第1个
	t.Log("------------ website batch update by other ... end ")
}

// 查 by id
func TestWebsiteQueryById(t *testing.T) {
	t.Log("------------ website query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddHasIdNoZero
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryById(website.Id)
	WebsiteCheckHasId(queryWebsite, website, t, "【 查 by id 】")
	t.Log("------------ website query by id ... start ")
}

// 查 by nameId
func TestWebsiteQueryByNameId(t *testing.T) {
	t.Log("------------ website query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryByNameId(website.NameId)
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by nameId 】")
	t.Log("------------ website query by nameId ... start ")
}

// 查 by other
func TestWebsiteQueryByOther(t *testing.T) {
	t.Log("------------ website query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryByOther("name_id", website.NameId)
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by other 】")
	t.Log("------------ website query by other ... start ")
}

// 查 - 批量 by id
func TestWebsiteBatchQueryById(t *testing.T) {
	t.Log("------------ website batch query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddHasIdNoZero
	website2 := website2ForAddHasIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	ids := []uint{website.Id, website2.Id}
	queryWebsites, err := WebsitesBatchQueryById(ids)
	if err != nil {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	WebsiteCheckHasId(queryWebsite, website, t, "【 查 by id 】")   // 判断第1个
	WebsiteCheckHasId(queryWebsite2, website2, t, "【 查 by id 】") // 判断第2个
	t.Log("------------ website batch query by id ... start ")
}

// 查 - 批量 by nameId
func TestWebsiteBatchQueryByNameId(t *testing.T) {
	t.Log("------------ website batch query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	nameIds := []int{website.NameId, website2.NameId}
	queryWebsites, err := WebsitesBatchQueryByNameId(nameIds)
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by nameId 】")   // 判断第1个
	WebsiteCheckNoId(queryWebsite2, website2, t, "【 查 by nameId 】") // 判断第2个
	t.Log("------------ website batch query by nameId ... start ")
}

// 查 - 批量 by other
func TestWebsiteBatchQueryByOther(t *testing.T) {
	t.Log("------------ website batch query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	website2 := website2ForAddNoIdHasZero
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	others := []any{website.NameId, website2.NameId}
	queryWebsites, err := WebsitesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by other 】")   // 判断第1个
	WebsiteCheckNoId(queryWebsite2, website2, t, "【 查 by other 】") // 判断第2个
	t.Log("------------ website batch query by other ... start ")
}
