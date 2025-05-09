package db

import (
	"fmt"
	"os"
	"study-spider-manhua-gin/models"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 所有方法思路
// 1. 清空表
// 2. 添加数据
// 3. 增删改查、批量增删改查
// 4. 检测

// ---------------------------- 变量 start ----------------------------
// 全局变量

// 局部变量
var testDB *gorm.DB

var websiteForAddWithIdNoZero *models.Website  // 用于add, 有id, 无0值
var websiteForAddWithIdHasZero *models.Website // 用于add, 有id, 无0值

// ---------------------------- 变量 end ----------------------------

// ---------------------------- init start ----------------------------
func init() {
	websiteForAddWithIdNoZero = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	websiteForAddWithIdHasZero = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}
}

// ---------------------------- init end ----------------------------
// 测试主函数
func TestMain(m *testing.M) {
	// 使用 MySQL 数据库进行测试
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/comic_test?charset=utf8mb4&parseTime=True&loc=Local"
	testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("----------- 连接测试数据失败 test database: ", err)
		panic(err)
	}

	// 设置全局 db 变量，防止调用方法DB.xx报错
	DB = testDB

	// 自动迁移表结构
	testDB.AutoMigrate(&models.Website{})

	// 运行测试
	os.Exit(m.Run())
}
func TestLog(t *testing.T) {
	// t.Log("----------- 测试能不能打印日志 --------------")
	fmt.Println("----------- 测试能不能打印日志 fmt.Println --------------")
}

// 增
func TestWebsiteAdd(t *testing.T) {
	t.Log("------------ website add ...  start ")

	// 1. 测试项1
	website := websiteForAddWithIdNoZero
	t.Log("website: ", website)
	WebsiteAdd(website)

	// var createdWebsite *models.Website // 手动写法
	// testDB.Where("name_id = ?", website.NameId).First(&createdWebsite) // 手动写法,不调用方法
	createdWebsite := WebsiteQueryByNameId(website.NameId) // 调用方法
	if createdWebsite.NameId != website.NameId || createdWebsite.Name != website.Name ||
		createdWebsite.Url != website.Url || createdWebsite.NeedProxy != website.NeedProxy ||
		createdWebsite.IsHttps != website.IsHttps {
		t.Errorf("【增】测试不通过, got= %v", createdWebsite)
		panic("【增】测试不通过")
	}

	// 2. 测试项2 测试needProxy =1 时候
	website2 := &models.Website{
		// Id:        2, // 新增时，可以指定id,如果有会更新
		NameId:    2,
		Name:      "Test Website Add 2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	t.Log("website2: ", website2)
	WebsiteAdd(website2)
	createdWebsite2 := WebsiteQueryByNameId(website2.NameId) // 调用方法
	if createdWebsite2.NameId != website2.NameId || createdWebsite2.Name != website2.Name ||
		createdWebsite2.Url != website2.Url || createdWebsite2.NeedProxy != website2.NeedProxy ||
		createdWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【增】测试不通过, got= %v", createdWebsite2)
		panic("【增】测试不通过")
	}

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
	website1 := &models.Website{
		NameId:    1,
		Name:      "Test Website1",
		Url:       "http://test.com1",
		NeedProxy: 0,
		IsHttps:   0,
	}

	website2 := &models.Website{
		NameId:    2,
		Name:      "Test Website2",
		Url:       "http://test.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites := []*models.Website{website1, website2}
	WebsiteBatchAdd(websites)
	nameIds := []int{website1.NameId, website2.NameId}
	t.Log("namedis = ", nameIds)
	createdWebsites, err := WebsitesBatchQueryByNameId(nameIds) // 调用方法
	if err != nil {
		t.Errorf("【增-批量】测试不通过, 查询nil, got=  %v", createdWebsites)
		panic("【增-批量】测试不通过")
	}

	// 判断第1个
	createdWebsite1 := createdWebsites[0]
	createdWebsite2 := createdWebsites[1]
	t.Log("查询结果 createdWebsites = ", createdWebsite1, createdWebsite2)
	t.Log("website1 = ", website1)
	if createdWebsite1.NameId != website1.NameId || createdWebsite1.Name != website1.Name ||
		createdWebsite1.Url != website1.Url || createdWebsite1.NeedProxy != website1.NeedProxy ||
		createdWebsite1.IsHttps != website1.IsHttps {
		t.Errorf("【增-批量 】测试不通过, got 1= %v", createdWebsite1)
		panic("【增-批量 】测试不通过")
	}
	// 判断第2个
	if createdWebsite2.NameId != website2.NameId || createdWebsite2.Name != website2.Name ||
		createdWebsite2.Url != website2.Url || createdWebsite2.NeedProxy != website2.NeedProxy ||
		createdWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【增-批量 】测试不通过, got 2= %v", createdWebsite2)
		panic("【增-批量 】测试不通过")
	}

	t.Log("------------ website batch add ... end ----------------")
}

// 删-通过id
func TestWebsiteDeleteById(t *testing.T) {
	t.Log("------------ website delete by id... start ----------------")
	website := websiteForAddWithIdNoZero
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
	website := websiteForAddWithIdNoZero
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
	website := websiteForAddWithIdNoZero
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
	website := websiteForAddWithIdNoZero

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
	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Delete By nameId",
		Url:       "http://delete.com id",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2 := &models.Website{
		NameId:    2,
		Name:      "Test Website for Delete By nameId 2",
		Url:       "http://delete.com id 2",
		NeedProxy: 0,
		IsHttps:   0,
	}
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
	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Delete By other",
		Url:       "http://delete.com id",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2 := &models.Website{
		NameId:    2,
		Name:      "Test Website for Delete By other 2",
		Url:       "http://delete.com id 2",
		NeedProxy: 0,
		IsHttps:   0,
	}
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
	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Update",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	updates := map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	WebsiteUpdateById(website.Id, updates)

	// 检查
	updatedWebsite := WebsiteQueryById(website.Id)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by id 】测试不通过, got= %v", updatedWebsite)
	}
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
	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Update nameId ",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	updates := map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	WebsiteUpdateByNameId(website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteQueryByNameId(website.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by nameId 】测试不通过, got= %v", updatedWebsite)
	}
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
	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Update other ",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	updates := map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	WebsiteUpdateByOther("name_id", website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteQueryByOther("name_id", website.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedWebsite =", updatedWebsite)
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by other 】测试不通过, got= %v", updatedWebsite)
	}
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
	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Update",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2 := &models.Website{
		Id:        2,
		NameId:    2,
		Name:      "Test Website for Update 2",
		Url:       "http://update.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := map[string]interface{}{
		"Id":        1,
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}

	updates2 := map[string]interface{}{
		"Id":        2,
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
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
	// 检测第1个
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by id 】测试不通过, got= %v", updatedWebsite)
	}
	// 检测第2个
	if updatedWebsite2.Name != updates2["Name"] ||
		updatedWebsite2.Url != updates2["Url"] ||
		updatedWebsite2.NeedProxy != updates2["NeedProxy"] ||
		updatedWebsite2.IsHttps != updates2["IsHttps"] {
		t.Errorf("【改 by id 】测试不通过, got= %v", updatedWebsite)
	}
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
	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Update",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2 := &models.Website{
		Id:        2,
		NameId:    2,
		Name:      "Test Website for Update 2",
		Url:       "http://update.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}

	updates2 := map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
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
	// 检测第1个
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by nameId 】测试不通过, got= %v", updatedWebsite)
	}
	// 检测第2个
	if updatedWebsite2.Name != updates2["Name"] ||
		updatedWebsite2.Url != updates2["Url"] ||
		updatedWebsite2.NeedProxy != updates2["NeedProxy"] ||
		updatedWebsite2.IsHttps != updates2["IsHttps"] {
		t.Errorf("【改 by nameId 】测试不通过, got= %v", updatedWebsite)
	}
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
	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Update",
		Url:       "http://update.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	website2 := &models.Website{
		Id:        2,
		NameId:    2,
		Name:      "Test Website for Update 2",
		Url:       "http://update.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	updates := map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Website",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}

	updates2 := map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Website2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
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
	// 检测第1个
	if updatedWebsite.Name != updates["Name"] ||
		updatedWebsite.Url != updates["Url"] ||
		updatedWebsite.NeedProxy != updates["NeedProxy"] ||
		updatedWebsite.IsHttps != updates["IsHttps"] {
		t.Errorf("【改 by other 】测试不通过, got= %v", updatedWebsite)
	}
	// 检测第2个
	if updatedWebsite2.Name != updates2["Name"] ||
		updatedWebsite2.Url != updates2["Url"] ||
		updatedWebsite2.NeedProxy != updates2["NeedProxy"] ||
		updatedWebsite2.IsHttps != updates2["IsHttps"] {
		t.Errorf("【改 by other 】测试不通过, got= %v", updatedWebsite)
	}
	t.Log("------------ website batch update by other ... end ")
}

// 查 by id
func TestWebsiteQueryById(t *testing.T) {
	t.Log("------------ website query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryById(website.Id)
	if queryWebsite.Id != website.Id ||
		queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website query by id ... start ")
}

// 查 by nameId
func TestWebsiteQueryByNameId(t *testing.T) {
	t.Log("------------ website query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryByNameId(website.NameId)
	if queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website query by nameId ... start ")
}

// 查 by other
func TestWebsiteQueryByOther(t *testing.T) {
	t.Log("------------ website query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	queryWebsite := WebsiteQueryByOther("name_id", website.NameId)
	if queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by other 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website query by other ... start ")
}

// 查 - 批量 by id
func TestWebsiteBatchQueryById(t *testing.T) {
	t.Log("------------ website batch query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	website2 := &models.Website{
		Id:        2,
		NameId:    2,
		Name:      "Test Website for Query2",
		Url:       "http://query.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	ids := []uint{website.Id, website2.Id}
	queryWebsites, err := WebsitesBatchQueryById(ids)
	if err != nil {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	// 判断第1个
	if queryWebsite.Id != website.Id ||
		queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryWebsite)
	}
	// 判断第2个
	if queryWebsite2.Id != website2.Id ||
		queryWebsite2.NameId != website2.NameId ||
		queryWebsite2.Name != website2.Name ||
		queryWebsite2.Url != website2.Url ||
		queryWebsite2.NeedProxy != website2.NeedProxy ||
		queryWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website batch query by id ... start ")
}

// 查 - 批量 by nameId
func TestWebsiteBatchQueryByNameId(t *testing.T) {
	t.Log("------------ website batch query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	website2 := &models.Website{
		NameId:    2,
		Name:      "Test Website for Query2",
		Url:       "http://query.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	nameIds := []int{website.NameId, website2.NameId}
	queryWebsites, err := WebsitesBatchQueryByNameId(nameIds)
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	// 判断第1个
	if queryWebsite.Id != website.Id ||
		queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsite)
	}
	// 判断第2个
	if queryWebsite2.Id != website2.Id ||
		queryWebsite2.NameId != website2.NameId ||
		queryWebsite2.Name != website2.Name ||
		queryWebsite2.Url != website2.Url ||
		queryWebsite2.NeedProxy != website2.NeedProxy ||
		queryWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website batch query by nameId ... start ")
}

// 查 - 批量 by other
func TestWebsiteBatchQueryByOther(t *testing.T) {
	t.Log("------------ website batch query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := &models.Website{
		NameId:    1,
		Name:      "Test Website for Query",
		Url:       "http://query.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	website2 := &models.Website{
		NameId:    2,
		Name:      "Test Website for Query2",
		Url:       "http://query.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}
	websites := []*models.Website{website, website2}
	WebsiteBatchAdd(websites)

	others := []any{website.NameId, website2.NameId}
	queryWebsites, err := WebsitesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	// 判断第1个
	if queryWebsite.Id != website.Id ||
		queryWebsite.NameId != website.NameId ||
		queryWebsite.Name != website.Name ||
		queryWebsite.Url != website.Url ||
		queryWebsite.NeedProxy != website.NeedProxy ||
		queryWebsite.IsHttps != website.IsHttps {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsite)
	}
	// 判断第2个
	if queryWebsite2.Id != website2.Id ||
		queryWebsite2.NameId != website2.NameId ||
		queryWebsite2.Name != website2.Name ||
		queryWebsite2.Url != website2.Url ||
		queryWebsite2.NeedProxy != website2.NeedProxy ||
		queryWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsite)
	}
	t.Log("------------ website batch query by other ... start ")
}
