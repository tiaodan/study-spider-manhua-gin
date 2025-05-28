/*
* 功能: website 异常用例，单元测试
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
* 测试用例：- 涉及异常操作
 1. 空
    - 单个值
    -- 有id/无id
    - 多个值
    -- 有id/无id
 2. 空格
    - 单个值
    -- 有id/无id
    - 多个值
    -- 有id/无id
 3. 特殊字符
    - 单个值
    - 多个值
 4. 过长
    - 单个值
    - 多个值
 5. 带单引号或双引号
    - 单个值
    - 多个值
 6. 带双引号
    - 单个值
    - 多个值
 7. 参数过多或过少
    - 过多
    - 过少
 8. 非指定类型
    - int -> string float bool
    - float -> string float bool
    - string -> int float bool
    - bool -> string float bool

* 所有方法思路
 1. 清空表
 2. 添加数据
 3. 增删改查、批量增删改查
 4. 检测

* 文件分区
- 变量 分区
- 测试函数 分区
- 被调用函数 分区
*/
package db

/*
import (
	"study-spider-manhua-gin/models"
	"testing"
)

// ---------------------------- 变量 start ----------------------------
// 1. 空
// - 单个值
// -- 有id/无id
// - 多个值
// -- 有id/无id
var var_websiteExcept_forAdd_caseNullOneValue_hasId_name *models.Website // 单个值
var var_websiteExcept_forAdd_caseNullOneValue_hasId_url *models.Website
var var_websiteExcept_forAdd_caseNullOneValue_noId_name *models.Website
var var_websiteExcept_forAdd_caseNullOneValue_noId_url *models.Website

var var_websiteExcept_forAdd_caseNullManyValue_hasId *models.Website // 多个值
var var_websiteExcept_forAdd_caseNullManyValue_noId *models.Website  // 多个值

// 1. 空格
// - 单个值
// -- 有id/无id
// - 多个值
// -- 有id/无id
var var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_name models.Website // 单个值
var var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_url models.Website
var var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_name models.Website
var var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_url models.Website

var var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_hasId models.Website // 多个值
var var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_noId models.Website  // 多个值

// ---------------------------- 变量 end ----------------------------

// ---------------------------- 初始化 start ----------------------------
func init() {
	// 1. 空
	// - 单个值
	// -- 有id/无id
	// - 多个值
	// -- 有id/无id
	var_websiteExcept_forAdd_caseNullOneValue_hasId_name = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_websiteExcept_forAdd_caseNullOneValue_hasId_url = &models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_websiteExcept_forAdd_caseNullOneValue_noId_name = &models.Website{
		NameId:    1,
		Name:      "",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_websiteExcept_forAdd_caseNullOneValue_noId_url = &models.Website{
		NameId:    1,
		Name:      "Test Website Add",
		Url:       "",
		NeedProxy: 1,
		IsHttps:   1,
	}

	var_websiteExcept_forAdd_caseNullManyValue_hasId = &models.Website{ // 多个值
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "",
		Url:       "",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_websiteExcept_forAdd_caseNullManyValue_noId = &models.Website{ // 多个值
		NameId:    1,
		Name:      "",
		Url:       "",
		NeedProxy: 1,
		IsHttps:   1,
	}

	// 1. 空格
	// - 单个值
	// -- 有id/无id
	// - 多个值
	// -- 有id/无id
	var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_name = models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      " Test Website Add ",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_url = models.Website{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Website Add",
		Url:       " http://ad d.com ",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_name = models.Website{
		NameId:    1,
		Name:      " Test Website Add ",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_url = models.Website{
		NameId:    1,
		Name:      "Test Website Add",
		Url:       " http://ad d.com ",
		NeedProxy: 1,
		IsHttps:   1,
	}

	var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_hasId = models.Website{ // 多个值
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      " Test Website Add ",
		Url:       " http://ad d.com ",
		NeedProxy: 1,
		IsHttps:   1,
	}
	var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_noId = models.Website{ // 多个值
		NameId:    1,
		Name:      " Test Website Add ",
		Url:       " http://ad d.com ",
		NeedProxy: 1,
		IsHttps:   1,
	}
}

// ---------------------------- 初始化 end ----------------------------

// ---------------------------- 测试函数 start ----------------------------
// 增
func TestWebsiteExceptAdd(t *testing.T) {
	t.Log("------------ website add ...  start ----------------")
	//  清空表
	TruncateTable(testDB, &models.Website{})
	//  调用 被调用函数
	websiteExceptAdd_null(t)  // 空
	websiteExceptAdd_space(t) // 空格

	t.Log("----------- website add ... end ----------------")
}

// ---------------------------- 测试函数 end ----------------------------

// ---------------------------- 被调用函数 start ----------------------------
// 增-空 只能string参数用
func websiteExceptAdd_null(t *testing.T) {
	// 1. 空
	// - 单个值
	// -- 有id/无id
	// - 多个值
	// -- 有id/无id

	t.Log("------------ website add except ... [test: null ] start ----------------")
	// 0 防止影响原始数据,拷贝个备份

	// 1 单个值 - 有id  -- name
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err := WebsiteOps.Add(var_websiteExcept_forAdd_caseNullOneValue_hasId_name) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query := WebsiteOps.QueryById(var_websiteExcept_forAdd_caseNullOneValue_hasId_name.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckHasId(query, var_websiteExcept_forAdd_caseNullOneValue_hasId_name, t, "【不通过】")

	// 1 单个值 - 有id  -- url
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err = WebsiteOps.Add(var_websiteExcept_forAdd_caseNullOneValue_hasId_url) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryById(var_websiteExcept_forAdd_caseNullOneValue_hasId_url.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckHasId(query, var_websiteExcept_forAdd_caseNullOneValue_hasId_url, t, "【不通过】")

	// 2 单个值 - 无id  -- name
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err = WebsiteOps.Add(var_websiteExcept_forAdd_caseNullOneValue_noId_name) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(var_websiteExcept_forAdd_caseNullOneValue_noId_name.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckNoId(query, var_websiteExcept_forAdd_caseNullOneValue_noId_name, t, "【不通过】")

	// 2 单个值 - 无id  -- url
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err = WebsiteOps.Add(var_websiteExcept_forAdd_caseNullOneValue_noId_url) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(var_websiteExcept_forAdd_caseNullOneValue_noId_url.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckNoId(query, var_websiteExcept_forAdd_caseNullOneValue_noId_url, t, "【不通过】")

	// 3 多个值 - 有id
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err = WebsiteOps.Add(var_websiteExcept_forAdd_caseNullManyValue_hasId) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryById(var_websiteExcept_forAdd_caseNullManyValue_hasId.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckHasId(query, var_websiteExcept_forAdd_caseNullManyValue_hasId, t, "【不通过】")

	// 3 多个值 - 无id
	//  清空表
	TruncateTable(testDB, &models.Website{})
	err = WebsiteOps.Add(var_websiteExcept_forAdd_caseNullManyValue_noId) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(var_websiteExcept_forAdd_caseNullManyValue_noId.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckNoId(query, var_websiteExcept_forAdd_caseNullManyValue_noId, t, "【不通过】")

	t.Log("------------ website add except ... [test: null ] end ----------------")

}

// 增-空 只能string参数用
func websiteExceptAdd_space(t *testing.T) {
	// 2. 空格
	// - 单个值
	// -- 有id/无id
	// - 多个值
	// -- 有id/无id

	t.Log("------------ website add except ... [test: space ] start ----------------")
	// 1 单个值 - 有id  -- name
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_name)
	copy := var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_name
	copyPointer := &copy
	err := WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query := WebsiteOps.QueryById(copyPointer.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	t.Log("插入后，变量有没有空格= ", query, "【】")
	// 判断插入数据对不对
	WebsiteCheckSpaceHasId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_name, t, "【单个值 有id 不通过】")

	// 1 单个值 - 有id  -- url
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_url)
	copy = var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_url
	copyPointer = &copy
	err = WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryById(copyPointer.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	t.Log("插入后，变量有没有空格= ", query, "【】")
	// 判断插入数对不对
	WebsiteCheckSpaceHasId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_hasId_url, t, "【不通过】")

	// 2 单个值 - 无id  -- name
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_name)
	copy = var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_name
	copyPointer = &copy
	err = WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(copyPointer.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	t.Log("插入后，变量有没有空格= ", query, "【】")
	// 判断插入数对不对
	WebsiteCheckSpaceNoId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_name, t, "【不通过】")

	// 2 单个值 - 无id  -- url
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_url)
	copy = var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_url
	copyPointer = &copy
	err = WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(copyPointer.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	t.Log("插入后，变量有没有空格= ", query, "【】")
	// 判断插入数对不对
	WebsiteCheckSpaceNoId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceOneValue_noId_url, t, "【不通过】")

	// 3 多个值 - 有id
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_hasId)
	copy = var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_hasId
	copyPointer = &copy
	err = WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryById(copyPointer.Id)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckSpaceHasId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_hasId, t, "【不通过】")

	// 3 多个值 - 无id
	//  清空表
	TruncateTable(testDB, &models.Website{})
	t.Log("插入前，变量有没有空格=  ", var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_noId)
	copy = var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_noId
	copyPointer = &copy
	err = WebsiteOps.Add(copyPointer) // 能加成功
	// 判断能否插入成功。如果查询空, 报错
	query = WebsiteOps.QueryByNameId(copyPointer.NameId)
	if query == nil {
		ProcessFailNoCheckErr(t, err, "没插入成功")
	}
	// 判断插入数对不对
	WebsiteCheckSpaceNoId(query, &var_noPointer_websiteExcept_forAdd_caseSpaceManyValue_noId, t, "【不通过】")

	t.Log("------------ website add except ... [test: space ] end ----------------")
}

// ---------------------------- 被调用函数 end ----------------------------
*/
