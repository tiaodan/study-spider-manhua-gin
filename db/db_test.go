// 测试主函数封装
package db

import (
	"fmt"
	"os"
	"study-spider-manhua-gin/models"
	"testing"

	"gorm.io/gorm"
)

// 全局变量

// db包 全局变量
var testDB *gorm.DB
var casePool []CaseContent // 测试用例池

// 测试case 类型定义
type CaseContent struct {
	db             *gorm.DB
	tbNameSingular string // 表名 tableNameSingular
	funcName       string // 方法名
	objs           []*models.Website
	updates        []map[string]interface{}
	caseTree1      string // 用例树顶层名字
	caseTree2      string // 用例树2层名字
	caseTree3      string // 用例树3层名字
	caseTree4      string
	caseTree5      string
}

// 测试主函数
func TestMain(m *testing.M) {
	// 使用 MySQL 数据库进行测试

	// 设置全局 db 变量，防止调用方法DB.xx报错
	InitDB("mysql", "comic_test", "root", "password")
	testDB = DB

	// 自动迁移表结构
	// testDB.AutoMigrate(&models.Website{}, &models.Country{}, &models.Category{}, &models.Type{}, &models.Comic{})
	testDB.AutoMigrate(&models.Website{}, &models.Type{})

	// 运行测试
	os.Exit(m.Run())
}

// 测试日志打印
func TestLog(t *testing.T) {
	t.Log("----------- 测试能不能打印日志 --------------")
	fmt.Println("----------- 测试能不能打印日志 fmt.Println --------------")
}

// 封装中断进程函数,如果有错误就t.FailNow()
func ProcessFail(t *testing.T, err error, errTitleStr string) {
	if err != nil {
		t.Error(errTitleStr, err)
		t.FailNow()
	}
}

// 封装中断进程函数,不判断err
func ProcessFailNoCheckErr(t *testing.T, err error, errTitleStr string) {
	t.Error(errTitleStr, err)
	// t.FailNow()  // 不用了, 好像不管用
	// panic("-----------------") // 测试的时候用panic,实际使用时注释掉
}
