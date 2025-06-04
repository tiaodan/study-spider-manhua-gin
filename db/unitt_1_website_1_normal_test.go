/*
* 功能: website 正常用例，单元测试
*
* 测试用例：- 涉及功能
* 1. 增
* 2. 删
* 3. 改
* 4. 查
* 5. 批量增
* 6. 批量删
* 7. 批量改
* 8. 批量查
*
* 测试用例：- 涉及正常操作
* 每个操作区分0值+非0值
 */
package db

import (
	"fmt"
	"strings"
	"study-spider-manhua-gin/models"
	"sync"
	"testing"

	"github.com/jinzhu/inflection"
	"gorm.io/gorm"
)

// 所有方法思路
// 1. 清空表
// 2. 添加数据
// 3. 增删改查、批量增删改查
// 4. 检测

// 测试用例
//

// ---------------------------- 变量 start ----------------------------
var onceTestCase sync.Once // 单例

// 用于 add
var websiteForAddHasIdNoZero models.Website        // 用于add, 有id, 无0值
var websitesForAddHasIdHasZeroOne []models.Website // 用于add, 有id, 无0值-单个为0
var websiteForAddHasIdHasZeroAll models.Website    // 用于add, 有id, 无0值-全为0
var websiteForAddNoIdNoZero models.Website         // 用于add, 无id, 无0值
var websitesForAddNoIdHasZeroOne []models.Website  // 用于add, 无id, 无0值-单个为0
var websiteForAddNoIdHasZeroAll models.Website     // 用于add, 有id, 无0值-全为0

// 用于 add batch
var website2ForAddHasIdNoZero models.Website        // 用于add, 有id, 无0值
var websites2ForAddHasIdHasZeroOne []models.Website // 用于add, 有id, 无0值-单个为0
var website2ForAddHasIdHasZeroAll models.Website    // 用于add, 有id, 无0值-全为0
var website2ForAddNoIdNoZero models.Website         // 用于add, 无id, 无0值
var websites2ForAddNoIdHasZeroOne []models.Website  // 用于add, 无id, 无0值-单个为0
var website2ForAddNoIdHasZeroAll models.Website     // 用于add, 有id, 无0值-全为0

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
	//---------- 用于add
	// // // 有id
	websiteForAddHasIdNoZero = models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websitesForAddHasIdHasZeroOne = WebsiteOps.returnObjZeroOne(websiteForAddHasIdNoZero) // 用于add, 有id, 无0值-单个为0 数组
	websiteForAddHasIdHasZeroAll = WebsiteOps.returnObjZeroAll(websiteForAddHasIdNoZero)  // 用于add, 有id, 无0值-全为0
	// // // 无id
	websiteForAddNoIdNoZero = models.Website{
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websitesForAddNoIdHasZeroOne = WebsiteOps.returnObjZeroOne(websiteForAddNoIdNoZero) // 用于add, 无id, 无0值-单个为0 数组
	websiteForAddNoIdHasZeroAll = WebsiteOps.returnObjZeroAll(websiteForAddNoIdNoZero)  // 用于add, 无id, 无0值-全为0

	//---------- 用于add batch
	// // // 有id
	website2ForAddHasIdNoZero = models.Website{
		Id:        2, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites2ForAddHasIdHasZeroOne = WebsiteOps.returnObjZeroOne(website2ForAddHasIdNoZero) // 用于add, 有id, 无0值-单个为0 数组
	website2ForAddHasIdHasZeroAll = WebsiteOps.returnObjZeroAll(website2ForAddHasIdNoZero)  // 用于add, 有id, 无0值-全为0

	// // // 无id
	website2ForAddNoIdNoZero = models.Website{
		NameId:    2,
		Name:      "Test Website Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}
	websites2ForAddNoIdHasZeroOne = WebsiteOps.returnObjZeroOne(website2ForAddNoIdNoZero) // 用于add, 无id, 无0值-单个为0 数组
	website2ForAddNoIdHasZeroAll = WebsiteOps.returnObjZeroAll(website2ForAddNoIdNoZero)  // 用于add, 无id, 无0值-全为0

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

// 初始化测试用例池
// 返回用例池
// func initTestCasePoll() []CaseContent { // 不用返回了，已经修改了全局变量了
func initTestCasePoll() { // 不用返回了，已经修改了全局变量了
	// 单例
	onceTestCase.Do(func() {
		// 思路：
		// 1. add
		// 2. batch add
		// 3. delete
		// 4. batch delete
		// 5. update
		// 6. batch update
		// 7. query
		// 8. batch query
		// 1) 生成用例: 有id用例, 参考xmind
		// 2) 生成用例: 无id用例, 参考xmind

		// 生成用例
		// 1. 生成用例: 有id用例, 参考xmind

		// 1. add ------------- start --------------------
		// 有id,无0值，
		// 第一行= case标题str
		// 第二行= 用例具体内容
		// 第三行= byOhter操作
		testCase := GenCaseContent("有id", "无0值", "", "", "",
			testDB, "website", "add", []models.Website{websiteForAddHasIdNoZero}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 有id,有0值，1个0
		for _, v := range websitesForAddHasIdHasZeroOne {
			testCase := GenCaseContent("有id", "有0值", "单个为0", "", "",
				testDB, "website", "add", []models.Website{v}, nil,
				false, "", nil, "", "", "")
			casePool = append(casePool, testCase)
		}

		// 有id,有0值，全0
		testCase = GenCaseContent("有id", "有0值", "全为0", "", "",
			testDB, "website", "add", []models.Website{websiteForAddHasIdHasZeroAll}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 2. 生成用例: 无id用例, 参考xmind
		// 无id,无0值
		testCase = GenCaseContent("无id", "无0值", "", "", "",
			testDB, "website", "add", []models.Website{websiteForAddNoIdNoZero}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 无id,有0值，1个0
		for _, v := range websitesForAddNoIdHasZeroOne {
			testCase := GenCaseContent("无id", "有0值", "单个为0", "", "",
				testDB, "website", "add", []models.Website{v}, nil,
				false, "", nil, "", "", "")
			casePool = append(casePool, testCase)
		}

		// 无id,有0值，全0
		testCase = GenCaseContent("无id", "有0值", "全为0", "", "",
			testDB, "website", "add", []models.Website{websiteForAddNoIdHasZeroAll}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)
		// 1. add ------------- end --------------------

		// 2. batch add ------------- start --------------------
		// 有id,无0值，1个
		testCase = GenCaseContent("有id", "无0值", "", "", "",
			testDB, "website", "batch add", []models.Website{websiteForAddHasIdNoZero, website2ForAddHasIdNoZero}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 有id,有0值，1个0
		for i, v := range websitesForAddHasIdHasZeroOne {
			testCase = GenCaseContent("有id", "有0值", "单个为0", "", "",
				testDB, "website", "batch add", []models.Website{v, websites2ForAddHasIdHasZeroOne[i]}, nil,
				false, "", nil, "", "", "")
			casePool = append(casePool, testCase)
		}

		// 有id,有0值，全0
		testCase = GenCaseContent("有id", "有0值", "全为0", "", "",
			testDB, "website", "batch add", []models.Website{websiteForAddHasIdHasZeroAll, website2ForAddHasIdHasZeroAll}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 2) 生成用例: 无id用例, 参考xmind
		// 无id,无0值
		testCase = GenCaseContent("无id", "无0值", "", "", "",
			testDB, "website", "batch add", []models.Website{websiteForAddNoIdNoZero, website2ForAddNoIdNoZero}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)

		// 无id,有0值，1个0
		for i, v := range websitesForAddNoIdHasZeroOne {
			testCase = GenCaseContent("无id", "有0值", "单个为0", "", "",
				testDB, "website", "batch add", []models.Website{v, websites2ForAddNoIdHasZeroOne[i]}, nil,
				false, "", nil, "", "", "")
			casePool = append(casePool, testCase)
		}

		// 无id,有0值，全0
		testCase = GenCaseContent("无id", "有0值", "全为0", "", "",
			testDB, "website", "batch add", []models.Website{websiteForAddNoIdHasZeroAll, website2ForAddNoIdHasZeroAll}, nil,
			false, "", nil, "", "", "")
		casePool = append(casePool, testCase)
		// 2. batch add ------------- end --------------------

		// 3. delete ------------- start --------------------
		// 3. delete ------------- end --------------------
	})
}

// ---------------------------- init end ----------------------------

// 检测函数封装, 自动判断Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheck(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	if obj.Id == 0 { // 无id
		if query.NameId != obj.NameId ||
			query.Name != obj.Name ||
			query.Url != obj.Url ||
			query.NeedProxy != obj.NeedProxy ||
			query.IsHttps != obj.IsHttps {
			// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
			t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
			t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
			ProcessFail(t, nil, "测试不通过")
		}
		return // 退出
	}
	// 判断第1个, 默认判断id
	if query.Id != obj.Id ||
		query.NameId != obj.NameId ||
		query.Name != obj.Name ||
		query.Url != obj.Url ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

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
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
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
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

// 检测函数封装,检测空白 自动对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckSpace(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	// 判断第1个
	if obj.Id == 0 { // 无id
		if query.NameId != obj.NameId ||
			query.Name != strings.TrimSpace(obj.Name) ||
			query.Url != strings.TrimSpace(obj.Url) ||
			query.NeedProxy != obj.NeedProxy ||
			query.IsHttps != obj.IsHttps {
			// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
			t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
			t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
			ProcessFail(t, nil, "测试不通过")
		}
		return
	}

	// 默认判断id
	if query.Id != obj.Id ||
		query.NameId != obj.NameId ||
		query.Name != strings.TrimSpace(obj.Name) ||
		query.Url != strings.TrimSpace(obj.Url) ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

// 检测函数封装,检测空白 对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckSpaceHasId(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.Id != obj.Id ||
		query.NameId != obj.NameId ||
		query.Name != strings.TrimSpace(obj.Name) ||
		query.Url != strings.TrimSpace(obj.Url) ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

// 检测函数封装,检测空白 不对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckSpaceNoId(query *models.Website, obj *models.Website, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.NameId != obj.NameId ||
		query.Name != strings.TrimSpace(obj.Name) ||
		query.Url != strings.TrimSpace(obj.Url) ||
		query.NeedProxy != obj.NeedProxy ||
		query.IsHttps != obj.IsHttps {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

// 检测函数封装, 删
// 参数1: 查到的指针arr 参数2: 要对比的对象指针arr
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckDelete(queries []*models.Website, objs []*models.Website, t *testing.T, errTitleStr string) {
	if len(queries) == 0 {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, queries)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, objs)
		ProcessFail(t, nil, "测试不通过")
	}
}

// 检测更新函数封装, 自动对比Id
// 参数1: 查到的指针 参数2: 更新参数 map[string]interface{}
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func WebsiteCheckUpdate(query *models.Website, obj map[string]interface{}, t *testing.T, errTitleStr string) {
	if obj["Id"] == 0 { // 无id
		if query.NameId != obj["NameId"] ||
			query.Name != obj["Name"] ||
			query.Url != obj["Url"] ||
			query.NeedProxy != obj["NeedProxy"] ||
			query.IsHttps != obj["IsHttps"] {
			// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
			t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
			t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
			ProcessFail(t, nil, "测试不通过")
		}
		return
	}

	// 判断第1个,默认检测id
	if query.Id != obj["Id"] ||
		query.NameId != obj["NameId"] ||
		query.Name != obj["Name"] ||
		query.Url != obj["Url"] ||
		query.NeedProxy != obj["NeedProxy"] ||
		query.IsHttps != obj["IsHttps"] {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
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
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
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
		t.Errorf("[ %s ] 测试不通过, got= %v", errTitleStr, query)
		t.Errorf("[ %s ] 测试不通过, obj= %v", errTitleStr, obj)
		ProcessFail(t, nil, "测试不通过")
	}
}

// ---------------------------- 阶段二：封装通用测试函数 start ----------------------------
// 参数1：t  *testing.T 。 测试对象
// 参数2：testDb。var testDB *gorm.DB。测试db对象指针
// 参数3：dbOps 。db操作对象。如WebsiteOperations
// 参数4: tableNameSingular string 。表名
// 参数5：functionName string。功能名称。增删改查、批量增删改查标签 如：add delete update batchAdd batchDelete batchUpdate
// 参数6：objs []*models.Website // 要添加 arr
// 参数7：queries []*models.Website  // 这个不用加，得进入函数后查询
// 参数8：updates []map[string]interface{} // 单个/批量修改，传的参数，单个修改用updates[0]
// 参数9：isByOther bool // 是否使用byOther，如：deleteByOther updateByOther
// 参数10：others []any  // byOther where in 传的参数。如果是单个操作，用others[0]
// 参数11：condition stirng // byOther条件，如 "NameId in (?)"中的NameId
// 参数12：others []any] // 接condition  NameId IN (?) 就是?号里要传的数组
// 参数13：queryType stirng // 查询方式， “byId” "byNameId" "byOther"
// 参数14：orderby stirng // 排序参数， Order("nameId DESC")
// 参数15：sort stirng // 排序参数， Order("nameId DESC")  这里是DESC
func commonDbTest_Website(t *testing.T, testDB *gorm.DB, tableNameSingular string, functionName string,
	objs []*models.Website, updates []map[string]interface{},
	isByOther bool, condition string, others []any, queryType string, orderby string, sort string) {

	// v2.0 写法。所有批量操作都用batchXX实现
	t.Logf("------------  %s ... start ", functionName)
	// 1. 清空表
	tableName := inflection.Plural(tableNameSingular) // 单数英文，转复数 如 website -> websites
	t.Log("清空表, tableName = ", tableName)
	TruncateTable(testDB.Table(tableName), nil) // 方式1： truncate table。 通过表名清空表

	// 2. 添加数据
	// 增删改查默认都会添加第一个数据，不用判断 方法名
	if strings.Contains(functionName, "batch") {
		WebsiteOps.BatchAdd(objs)

	} else {
		WebsiteOps.Add(objs[0])
	}

	// 判断是否插入2条数据
	checkQueries, _ := WebsiteOps.BatchQueryAll()
	if strings.Contains(functionName, "batch") {
		if len(checkQueries) != 2 {
			t.Errorf("批量操作失败, 期望返回2条数据, 实际返回%d条数据", len(checkQueries))
		}
	} else { // 判断是否插入1条数据
		if len(checkQueries) != 1 {
			t.Errorf("批量操作失败, 期望返回1条数据, 实际返回%d条数据", len(checkQueries))
		}
	}

	// 3. 增删改查、批量增删改查操作
	// 不带batch,只操作第一条数据。 带batch，操作第二条数据。
	// 所有批量操作，只判断数组中, 第一个数的Id,其它的不管
	switch {
	case strings.Contains(functionName, "delete"):
		if strings.Contains(functionName, "batch") { // 批量操作
			if isByOther { // 使用byOther删除
				WebsiteOps.BatchDeleteByOther(condition, others)
			} else if objs[0].Id == 0 { // Id字段是空的，用DeleteByNameId
				var nameIds []int
				for _, obj := range objs {
					nameIds = append(nameIds, obj.NameId)
				}
				WebsiteOps.BatchDeleteByNameId(nameIds)
			} else {
				var ids []uint
				for _, obj := range objs {
					ids = append(ids, obj.Id)
				}
				WebsiteOps.BatchDeleteById(ids) // 默认通过id删除
			}
		} else { // 只操作第1条数据
			if isByOther { // 使用byOther删除
				WebsiteOps.DeleteByOther(condition, others[0])
			} else if objs[0].Id == 0 { // Id字段是空的，用DeleteByNameId
				WebsiteOps.DeleteByNameId(objs[0].NameId)
			} else {
				WebsiteOps.DeleteById(objs[0].Id) // 默认通过id删除
			}
		}
	case strings.Contains(functionName, "update"):
		if strings.Contains(functionName, "batch") {
			if isByOther {
				WebsiteOps.BatchUpdateByOther(condition, others, updates)
			} else if objs[0].Id == 0 { // Id字段是空的，用DeleteByNameId
				WebsiteOps.BatchUpdateByNameId(updates)
			} else {
				WebsiteOps.BatchUpdateById(updates) // 默认通过id删除
			}
		} else { // 只操作第1条数据
			if isByOther {
				WebsiteOps.UpdateByOther(condition, others[0], updates[0])
			} else if objs[0].Id == 0 { // Id字段是空的，用DeleteByNameId
				WebsiteOps.UpdateByNameId(objs[0].NameId, updates[0])
			} else {
				WebsiteOps.UpdateById(objs[0].Id, updates[0]) // 默认通过id删除
			}
		}
	default:
		// t.Log("functionName 未匹配到【 删改 】操作, functionName = ", functionName)
	}

	// 4. 查询数据,赋值给query变量,必须根据obsj的nameIds 按顺序获取，不能用 queryAll -> 顺序可能是乱的

	var queries []*models.Website
	// 不区分大小写
	if strings.ToLower(queryType) == "byid" {
		var ids []uint
		for _, obj := range objs {
			ids = append(ids, obj.Id)
		}
		queries, _ = WebsiteOps.BatchQueryById(ids)
	} else if strings.ToLower(queryType) == "byother" {
		queries, _ = WebsiteOps.BatchQueryByOther(condition, others, orderby, sort)
	} else {
		// byNameId - 默认方式
		var nameIds []int
		for _, obj := range objs {
			nameIds = append(nameIds, obj.NameId)
		}
		queries, _ = WebsiteOps.BatchQueryByNameId(nameIds)
	}

	// 5. 检测对比
	switch functionName {
	case "add":
		for i, obj := range objs {
			WebsiteCheck(queries[i], obj, t, functionName)
			if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
				break
			}
		}
	case "delete":
		WebsiteCheckDelete(queries, objs, t, functionName)
	case "update":
		for i, update := range updates {
			WebsiteCheckUpdate(queries[i], update, t, functionName)
			if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
				break
			}
		}
	case "query":
		for i, obj := range objs {
			WebsiteCheck(queries[i], obj, t, functionName)
			if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
				break
			}
		}
	}

	// 检测。
	t.Logf("------------  %s ... end ", functionName)
	// v1.0写法，所有批量操作，都用for 单个实现
	/*
		t.Logf("------------  %s ... start ", functionName)
		// 1. 清空表
		tableName := inflection.Plural(tableNameSingular) // 单数英文，转复数 如 website -> websites
		t.Log("清空表, tableName = ", tableName)
		TruncateTable(testDB.Table(tableName), nil) // 方式1： truncate table。 通过表名清空表

		// 2. 添加数据
		// 增删改查默认都会添加第一个数据，不用判断 方法名
		for _, obj := range objs {
			WebsiteOps.Add(obj)
			// 判断functionName是否是批量操作，添加第二条数据
			if !strings.Contains(functionName, "batch") {
				break
			}
		}
		// 判断是否插入2条数据
		if strings.Contains(functionName, "batch") {
			t.Log("------  functionName  ", functionName)
			queries, _ := WebsiteOps.BatchQueryAll()
			t.Log("------  queries  ", queries)
			if len(queries) != 2 {
				t.Errorf("批量操作失败, 期望返回2条数据, 实际返回%d条数据", len(queries))
			}
		}

		// 3. 增删改查、批量增删改查操作
		// 不带batch,只操作第一条数据。 带batch，操作第二条数据
		switch {
		case strings.Contains(functionName, "delete"):
			for i, obj := range objs {
				if isByOther { // 使用byOther删除
					WebsiteOps.DeleteByOther(condition, others[i])
				} else if obj.Id == 0 { // Id字段是空的，用DeleteByNameId
					WebsiteOps.DeleteByNameId(obj.NameId)
				} else {
					WebsiteOps.DeleteById(obj.Id) // 默认通过id删除
				}

				// 判断是否操作 第二条数据
				if !strings.Contains(functionName, "batch") {
					break
				}
			}
		case strings.Contains(functionName, "update"):
			for i, obj := range objs {
				if isByOther { // 使用byOther删除
					WebsiteOps.UpdateByOther(condition, others[i])
				} else if obj.Id == 0 { // Id字段是空的，用DeleteByNameId
					WebsiteOps.UpdateByNameId(obj.NameId, updates[i])
				} else {
					WebsiteOps.UpdateById(obj.Id, updates[i]) // 默认通过id删除
				}

				// 判断是否操作 第二条数据
				if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
					break
				}
			}
		default:
			// t.Log("functionName 未匹配到【 删改 】操作, functionName = ", functionName)
		}

		// 4. 查询数据,赋值给query变量
		var nameIds []int
		for _, obj := range objs { // init nameids
			nameIds = append(nameIds, obj.NameId)
			if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
				break
			}
		}
		queries, _ := WebsiteOps.BatchQueryByNameId(nameIds)

		// 5. 检测对比
		switch functionName {
		case "add":
			for i, obj := range objs {
				WebsiteCheck(queries[i], obj, t, functionName)
				if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
					break
				}
			}
		case "delete":
			WebsiteCheckDelete(queries, objs, t, functionName)
		case "update":
			for i, update := range updates {
				WebsiteCheckUpdate(queries[i], update, t, functionName)
				if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
					break
				}
			}
		case "query":
			for i, obj := range objs {
				WebsiteCheck(queries[i], obj, t, functionName)
				if !strings.Contains(functionName, "batch") { // 没有批量操作，就退出
					break
				}
			}
		}

		// 检测。
		t.Logf("------------  %s ... end ", functionName)
	*/
}

// 测试通过方法
func TestCommon(t *testing.T) {
	initTestCasePoll() // 获取用例池
	for i, v := range casePool {
		objsStr := fmt.Sprintf("%v", v.objs[0])
		if len(v.objs) > 1 {
			objsStr += fmt.Sprintf(", %v", v.objs[1])
		}
		// fmt.Printf("用例池 len(pool) = %v, i=%v, pool = %v, objs=%v \n", len(casePool), i+1, v, objsStr) // 老的写法: v.objs[0]
		// fmt.Printf("用例池 len(pool) = %v, i=%v, objs=%v \n", len(casePool), i+1, objsStr) // 老的写法: v.objs[0]
		fmt.Printf("用例池 len(pool) = %v, i=%v, funcName=%s, objs=%v, case = [%s] - [%s] - [%s] - [%s] - [%s] \n",
			len(casePool), i+1, v.funcName,
			objsStr, v.caseTree1, v.caseTree2, v.caseTree3, v.caseTree4, v.caseTree5)
	}

	// 通用用例池，循环进行测试
	// for len(casePool) > 0 {  // 另一种for写法
	for i := range casePool {
		if len(casePool) == 0 {
			break
		}
		// 取出第一个用例
		current := casePool[0]
		fmt.Printf("------------------------- 当前用例=%v, funcName= [%s] case = [%s] [%s] [%s] [%s] [%s] ------------------------- \n",
			i+1, current.funcName, current.caseTree1, current.caseTree2, current.caseTree3,
			current.caseTree4, current.caseTree5)

		// 执行用例
		var objsPointer []*models.Website // 把测试用例中变量 -> 指针
		for _, obj := range current.objs {
			objsPointer = append(objsPointer, &obj)
		}
		commonDbTest_Website(t, current.db, current.tbNameSingular, current.funcName, objsPointer, current.updates,
			current.isByOther, current.condition, current.others, current.queryType, current.orderby, current.sort)

		// 从用例池中移除已执行的用例
		casePool = casePool[1:]
		fmt.Println("") // case end
	}

}

// ---------------------------- 阶段二：封装通用测试函数 end ----------------------------

// ---------------------------- 阶段一：每个用例都是一个函数 start ----------------------------
/*
// 增
func TestWebsiteAdd(t *testing.T) {
	t.Log("------------ website add ...  start ")
	// 1. 测试项1，有0值
	website := websiteForAddHasIdHasZero
	t.Log("website: ", website)
	WebsiteOps.Add(website)

	// var createdWebsite *models.Website // 手动写法
	// testDB.Where("name_id = ?", website.NameId).First(&createdWebsite) // 手动写法,不调用方法
	createdWebsite := WebsiteOps.QueryByNameId(website.NameId) // 调用方法
	WebsiteCheckNoId(createdWebsite, website, t, "【增】")        // 测试第1个

	// 2. 测试项2 无0值, 测试needProxy =1 时候
	website2 := websiteForAddHasIdNoZero
	t.Log("website2: ", website2)
	WebsiteOps.Add(website2)
	createdWebsite2 := WebsiteOps.QueryByNameId(website2.NameId) // 调用方法
	WebsiteCheckNoId(createdWebsite2, website2, t, "【增】")        // 测试第2个

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
	WebsiteOps.BatchAdd(websites)
	nameIds := []int{website.NameId, website2.NameId}
	t.Log("namedis = ", nameIds)
	createdWebsites, err := WebsiteOps.BatchQueryByNameId(nameIds) // 调用方法
	if err != nil {
		t.Errorf("【增-批量】测试不通过, 查询nil, got=  %v", createdWebsites)
		ProcessFail(t, nil, "测试不通过")
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
	WebsiteOps.Add(website)

	WebsiteOps.DeleteById(website.Id)

	var deletedWebsite models.Website
	result := testDB.First(&deletedWebsite, website.Id)
	// result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by id】测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		ProcessFail(t, nil, "测试不通过")
	}
	t.Log("------------ website delete by id... end ----------------")
}

// 删-通过 nameId
func TestWebsiteDeleteByNameId(t *testing.T) {
	t.Log("------------ website delete by nameId... start ----------------")
	website := websiteForAddNoIdNoZero
	WebsiteOps.Add(website)

	WebsiteOps.DeleteByNameId(website.NameId)

	var deletedWebsite models.Website
	result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by nameId 】 测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		ProcessFail(t, nil, "测试不通过")
	}
	t.Log("------------ website delete by nameId... end ----------------")
}

// 删-通过 其它
func TestWebsiteDeleteByOther(t *testing.T) {
	t.Log("------------ website delete by other... start ----------------")
	website := websiteForAddNoIdNoZero
	WebsiteOps.Add(website)

	WebsiteOps.DeleteByOther("name_id", website.NameId)

	var deletedWebsite models.Website
	result := testDB.Where("name_id = ?", website.NameId).First(&deletedWebsite)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by other 】测试不通过,删除后仍能查到, =  %v", deletedWebsite)
		ProcessFail(t, nil, "测试不通过")
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
	WebsiteOps.BatchAdd(websites) // 添加

	ids := []uint{website.Id, website2.Id}
	t.Log("ids = ", ids)

	// 判断是否添加了2个
	websites, err := WebsiteOps.BatchQueryById(ids)
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", websites)
		// panic("【删 批量 - by id 】测试不通过,删除后仍能查到") // 测试v原本不能用pnic
	}

	WebsiteOps.BatchDeleteById(ids) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsiteOps.BatchQueryById(ids)

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", websites)
		ProcessFail(t, nil, "测试不通过")
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
	WebsiteOps.BatchAdd(websites) // 添加

	nameIds := []int{website.NameId, website2.NameId}
	t.Log("nameIds = ", nameIds)

	// 判断是否添加了2个
	websites, err := WebsiteOps.BatchQueryByNameId(nameIds)
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", websites)
		ProcessFail(t, nil, "测试不通过")
	}

	WebsiteOps.BatchDeleteByNameId(nameIds) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsiteOps.BatchQueryByNameId(nameIds)

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", websites)
		ProcessFail(t, nil, "测试不通过")
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
	WebsiteOps.BatchAdd(websites) // 添加

	others := []any{website.NameId, website2.NameId}
	t.Log("others = ", others)

	// 判断是否添加了2个
	websites, err := WebsiteOps.BatchQueryByOther("name_id", others, "name_id", "ASC")
	if len(websites) != 2 || err != nil {
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", websites)
		ProcessFail(t, nil, "测试不通过")
	}

	WebsiteOps.BatchDeleteByOther("name_id", others) // 删除

	// 检测，如果报错，或者 结果>0
	websites, err = WebsiteOps.BatchQueryByOther("name_id", others, "name_id", "ASC")

	if len(websites) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", websites)
		ProcessFail(t, nil, "测试不通过")
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
	WebsiteOps.Add(website)

	updates := websiteForUpdateHasIdHasZero
	WebsiteOps.UpdateById(website.Id, updates)

	// 检查
	updatedWebsite := WebsiteOps.QueryById(website.Id)
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
	WebsiteOps.Add(website)

	updates := websiteForUpdateNoIdHasZero
	WebsiteOps.UpdateByNameId(website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteOps.QueryByNameId(website.NameId)
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
	WebsiteOps.Add(website)

	updates := websiteForUpdateNoIdHasZero
	WebsiteOps.UpdateByOther("name_id", website.NameId, updates)

	// 检查
	updatedWebsite := WebsiteOps.QueryByOther("name_id", website.NameId)
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
	WebsiteOps.BatchAdd(websites)

	updates := websiteForUpdateHasIdNoZero
	updates2 := website2ForUpdateHasIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsiteOps.BatchUpdateById(updatesArr)

	// 检测，如果报错
	ids := []uint{1, 2}
	websites, err := WebsiteOps.BatchQueryById(ids)
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
	WebsiteOps.BatchAdd(websites)

	updates := websiteForUpdateNoIdNoZero
	updates2 := website2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsiteOps.BatchUpdateByNameId(updatesArr)

	// 检测，如果报错
	nameIds := []int{
		updates["NameId"].(int),
		updates2["NameId"].(int),
	}
	websites, err := WebsiteOps.BatchQueryByNameId(nameIds)
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
	WebsiteOps.BatchAdd(websites)

	updates := websiteForUpdateNoIdNoZero
	updates2 := website2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	WebsiteOps.BatchUpdateByOther(updatesArr)

	// 检测，如果报错
	others := []any{
		updates["NameId"],
		updates2["NameId"],
	}
	websites, err := WebsiteOps.BatchQueryByOther("name_id", others, "name_id", "ASC")
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
	WebsiteOps.Add(website)

	queryWebsite := WebsiteOps.QueryById(website.Id)
	WebsiteCheckHasId(queryWebsite, website, t, "【 查 by id 】")
	t.Log("------------ website query by id ... start ")
}

// 查 by nameId
func TestWebsiteQueryByNameId(t *testing.T) {
	t.Log("------------ website query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	WebsiteOps.Add(website)

	queryWebsite := WebsiteOps.QueryByNameId(website.NameId)
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by nameId 】")
	t.Log("------------ website query by nameId ... start ")
}

// 查 by other
func TestWebsiteQueryByOther(t *testing.T) {
	t.Log("------------ website query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Website{}) // 方式1： truncate table

	website := websiteForAddNoIdNoZero
	WebsiteOps.Add(website)

	queryWebsite := WebsiteOps.QueryByOther("name_id", website.NameId)
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
	WebsiteOps.BatchAdd(websites)

	ids := []uint{website.Id, website2.Id}
	queryWebsites, err := WebsiteOps.BatchQueryById(ids)
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
	WebsiteOps.BatchAdd(websites)

	nameIds := []int{website.NameId, website2.NameId}
	queryWebsites, err := WebsiteOps.BatchQueryByNameId(nameIds)
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
	WebsiteOps.BatchAdd(websites)

	others := []any{website.NameId, website2.NameId}
	queryWebsites, err := WebsiteOps.BatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryWebsites)
	}

	queryWebsite := queryWebsites[0]
	queryWebsite2 := queryWebsites[1]
	WebsiteCheckNoId(queryWebsite, website, t, "【 查 by other 】")   // 判断第1个
	WebsiteCheckNoId(queryWebsite2, website2, t, "【 查 by other 】") // 判断第2个
	t.Log("------------ website batch query by other ... start ")
}
*/

// ---------------------------- 阶段一：每个用例都是一个函数 end ----------------------------
