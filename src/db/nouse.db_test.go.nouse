// 测试主函数封装
package db

import (
	"fmt"
	"os"
	"study-spider-manhua-gin/src/models"
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
	tbNameSingular string           // 表名 tableNameSingular
	funcName       string           // 方法名
	objs           []models.Website // 不能用指针，测试用例，变量容易被多次调用，用指针，很容易把最原始变量改了
	updates        []map[string]interface{}
	isByOther      bool   // 是否用byOther
	condition      string // other 条件
	// other          any    // other 条件 参数 nameId == ? -》 ?中内容,不用这个字段，用others[0]可以替代
	others    []any  // other 条件 参数 nameId IN ? -》 ?中内容
	queryType string // 查询类型 "byId" "byNameId" "byOther"
	orderby   string // 根据什么排序
	sort      string // 排序方式 ASC DESC
	caseTree1 string // 用例树顶层名字
	caseTree2 string // 用例树2层名字
	caseTree3 string // 用例树3层名字
	caseTree4 string
	caseTree5 string
}

// 测试主函数
func TestMain(m *testing.M) {
	// 使用 MySQL 数据库进行测试

	// 设置全局 db 变量，防止调用方法DB.xx报错
	InitDB("mysql", "comic_test", "root", "password")
	testDB = DBComic

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

// 生成测试用例
/*
testCase := CaseContent{
	caseTree1:      "有id",
	caseTree2:      "无0值",
	db:             testDB,    // db Pointer
	tbNameSingular: "website", // tbName
	funcName:       "add",
	objs:           []*models.Website{&websiteForAddHasIdNoZero},
	updates:        []map[string]interface{}{},
}
*/
func GenCaseContent(caseTree1, caseTree2, caseTree3, caseTree4, caseTree5 string,
	db *gorm.DB, tbNameSingular, funcName string,
	objs []models.Website, updates []map[string]interface{},
	isByOther bool, condition string, others []any, queryType, orderby, sort string) CaseContent {

	testCase := CaseContent{
		caseTree1:      caseTree1,
		caseTree2:      caseTree2,
		caseTree3:      caseTree3,
		caseTree4:      caseTree4,
		caseTree5:      caseTree5,
		db:             db,             // db Pointer
		tbNameSingular: tbNameSingular, // tbName
		funcName:       funcName,
		objs:           objs,
		updates:        updates,
		isByOther:      isByOther,
		condition:      condition,
		others:         others,
		queryType:      queryType,
		orderby:        orderby,
		sort:           sort,
	}
	return testCase
}
