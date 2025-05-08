package db

import (
	"fmt"
	"os"
	"study-spider-manhua-gin/models"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

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

	website := &models.Website{
		// Id:        1, // 新增时，可以指定id,如果有会更新
		NameId:    1,
		Name:      "Test Website",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}
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
	t.Log("----------- website add ... end ----------------")
}

// 批量增
func TestWebsiteBatchAdd(t *testing.T) {
	t.Log("------------ website batch add ... start ----------------")
	website1 := &models.Website{
		NameId:    1,
		Name:      "Test Website1",
		Url:       "http://test.com1",
		NeedProxy: 1,
		IsHttps:   1,
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
	nameIds := []any{website1.NameId, website2.NameId}
	t.Log("namedis = ", nameIds)
	createdWebsites, err := WebsitesBatchQueryByNameId(nameIds) // 调用方法
	if err != nil {
		t.Errorf("【增-批量】测试不通过, got=  %v", createdWebsites)
		panic("【增-批量】测试不通过")
	}
	// 判断第1个
	createdWebsite1 := createdWebsites[0]
	createdWebsite2 := createdWebsites[1]
	t.Log("createdWebsites = ", createdWebsite1, createdWebsite2)
	if createdWebsite1.NameId != website1.NameId || createdWebsite1.Name != website1.Name ||
		createdWebsite1.Url != website1.Url || createdWebsite1.NeedProxy != website1.NeedProxy ||
		createdWebsite1.IsHttps != website1.IsHttps {
		t.Errorf("【增-批量】测试不通过, got 1= %v", createdWebsite1)
		panic("【增-批量】测试不通过")
	}
	// 判断第2个
	if createdWebsite2.NameId != website2.NameId || createdWebsite2.Name != website2.Name ||
		createdWebsite2.Url != website2.Url || createdWebsite2.NeedProxy != website2.NeedProxy ||
		createdWebsite2.IsHttps != website2.IsHttps {
		t.Errorf("【增-批量】测试不通过, got 2= %v", createdWebsite2)
		panic("【增-批量】测试不通过")
	}

	t.Log("------------ website batch add ... end ----------------")
}

func TestWebsiteDeleteById(t *testing.T) {
	t.Log("------------ website delete by id... start ----------------")
	website := &models.Website{
		Id:        1,
		NameId:    1,
		Name:      "Test Website for Delete",
		Url:       "http://delete.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	WebsiteAdd(website)

	WebsiteDeleteById(website.Id)

	var deletedWebsite models.Website
	result := testDB.First(&deletedWebsite, website.Id)
	// result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by id】测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		panic("【删 - by id】测试不通过,删除后仍能查到")
	}
	t.Log("------------ website delete by id... end ----------------")
}

// func TestWebsiteUpdate(t *testing.T) {
// 	t.Log("------------ website update ... start ")
// 	website := &models.Website{
// 		NameId:    3,
// 		Name:      "Test Website for Update",
// 		Url:       "http://update.com",
// 		NeedProxy: 1,
// 		IsHttps:   1,
// 	}
// 	WebsiteAdd(website)

// 	updates := map[string]interface{}{
// 		"Name":      "Updated Website",
// 		"Url":       "http://updated.com",
// 		"NeedProxy": 0,
// 		"IsHttps":   0,
// 	}
// 	WebsiteUpdateByOther("name_id", website.NameId, updates)

// 	var updatedWebsite models.Website
// 	testDB.Where("name_id = ?", website.NameId).First(&updatedWebsite)
// 	if updatedWebsite.Name != "Updated Website" || updatedWebsite.Url != "http://updated.com" ||
// 		updatedWebsite.NeedProxy != 0 || updatedWebsite.IsHttps != 0 {
// 		t.Errorf("【删】测试不通过, got= %v", updatedWebsite)
// 	}
// 	t.Log("------------ website update ... end ")
// }

// func TestWebsiteQueryById(t *testing.T) {
// 	t.Log("------------ website query ... start ")
// 	website := &models.Website{
// 		NameId:    4,
// 		Name:      "Test Website for Query",
// 		Url:       "http://query.com",
// 		NeedProxy: 1,
// 		IsHttps:   1,
// 	}
// 	WebsiteAdd(website)

// 	retrievedWebsite := WebsiteQueryByOther("name_id", website.NameId)
// 	if retrievedWebsite == nil || retrievedWebsite.Name != website.Name || retrievedWebsite.Url != website.Url ||
// 		retrievedWebsite.NeedProxy != website.NeedProxy || retrievedWebsite.IsHttps != website.IsHttps {
// 		t.Errorf("【删】测试不通过, got= %v", retrievedWebsite)
// 	}
// 	t.Log("------------ website query ... start ")
// }
