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
// 全局变量

// 用于 add
var comicTypeForAddHasIdNoZero *models.Type  // 用于add, 有id, 无0值
var comicTypeForAddHasIdHasZero *models.Type // 用于add, 有id, 无0值
var comicTypeForAddNoIdNoZero *models.Type   // 用于add, 无id, 无0值
var comicTypeForAddNoIdHasZero *models.Type  // 用于add, 无id, 无0值

// 用于 add batch
var comicType2ForAddHasIdNoZero *models.Type  // 用于add, 有id, 无0值
var comicType2ForAddHasIdHasZero *models.Type // 用于add, 有id, 无0值
var comicType2ForAddNoIdNoZero *models.Type   // 用于add, 无id, 无0值
var comicType2ForAddNoIdHasZero *models.Type  // 用于add, 无id, 无0值

// 用于 update
var comicTypeForUpdateHasIdNoZero map[string]interface{}  // 用于update, 有id, 无0值
var comicTypeForUpdateHasIdHasZero map[string]interface{} // 用于update, 有id, 有0值

var comicTypeForUpdateNoIdNoZero map[string]interface{}  // 用于update, 无id, 无0值
var comicTypeForUpdateNoIdHasZero map[string]interface{} // 用于update, 无id, 有0值

// 用于 update batch
var comicType2ForUpdateHasIdNoZero map[string]interface{}  // 用于update, 有id, 无0值
var comicType2ForUpdateHasIdHasZero map[string]interface{} // 用于update, 有id, 有0值

var comicType2ForUpdateNoIdNoZero map[string]interface{}  // 用于update, 无id, 无0值
var comicType2ForUpdateNoIdHasZero map[string]interface{} // 用于update, 无id, 有0值

// ---------------------------- 变量 end ----------------------------

// ---------------------------- init start ----------------------------
func init() {
	// 用于add, 有id
	comicTypeForAddHasIdNoZero = &models.Type{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Type Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	comicTypeForAddHasIdHasZero = &models.Type{
		Id:        1, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    1,
		Name:      "Test Type Add",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add, 无id
	comicTypeForAddNoIdNoZero = &models.Type{
		NameId:    1,
		Name:      "Test Type Add",
		Url:       "http://add.com",
		NeedProxy: 1,
		IsHttps:   1,
	}

	comicTypeForAddNoIdHasZero = &models.Type{
		NameId:    1,
		Name:      "Test Type Add",
		Url:       "http://add.com",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add batch, 有id
	comicType2ForAddHasIdNoZero = &models.Type{
		Id:        2, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    2,
		Name:      "Test Type Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}

	comicType2ForAddHasIdHasZero = &models.Type{
		Id:        2, // 新增时,可以指定id,gorm会插入指定id,而不是自增
		NameId:    2,
		Name:      "Test Type Add2",
		Url:       "http://add.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于add batch, 无id
	comicType2ForAddNoIdNoZero = &models.Type{
		NameId:    2,
		Name:      "Test Type Add2",
		Url:       "http://add.com2",
		NeedProxy: 1,
		IsHttps:   1,
	}

	comicType2ForAddNoIdHasZero = &models.Type{
		NameId:    2,
		Name:      "Test Type Add2",
		Url:       "http://add.com2",
		NeedProxy: 0,
		IsHttps:   0,
	}

	// 用于update
	comicTypeForUpdateHasIdNoZero = map[string]interface{}{
		"Id":        uint(1),
		"NameId":    1,
		"Name":      "Updated Type",
		"Url":       "http://updated.com",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	comicTypeForUpdateHasIdHasZero = map[string]interface{}{
		"Id":        uint(1),
		"NameId":    1,
		"Name":      "Updated Type",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	// 无id
	comicTypeForUpdateNoIdNoZero = map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Type",
		"Url":       "http://updated.com",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	comicTypeForUpdateNoIdHasZero = map[string]interface{}{
		"NameId":    1,
		"Name":      "Updated Type",
		"Url":       "http://updated.com",
		"NeedProxy": 0,
		"IsHttps":   0,
	}

	// 用于update batch
	comicType2ForUpdateHasIdNoZero = map[string]interface{}{
		"Id":        uint(2),
		"NameId":    2,
		"Name":      "Updated Type2",
		"Url":       "http://updated.com2",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	comicType2ForUpdateHasIdHasZero = map[string]interface{}{
		"Id":        uint(2),
		"NameId":    2,
		"Name":      "Updated Type2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
	// 无id
	comicType2ForUpdateNoIdNoZero = map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Type2",
		"Url":       "http://updated.com2",
		"NeedProxy": 1,
		"IsHttps":   1,
	}

	comicType2ForUpdateNoIdHasZero = map[string]interface{}{
		"NameId":    2,
		"Name":      "Updated Type2",
		"Url":       "http://updated.com2",
		"NeedProxy": 0,
		"IsHttps":   0,
	}
}

// ---------------------------- init end ----------------------------

// 检测函数封装, 对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func TypeCheckHasId(query *models.Type, obj *models.Type, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.Id != obj.Id ||
		query.NameId != obj.NameId ||
		query.Name != obj.Name ||
		query.Level != obj.Level ||
		query.Parent != obj.Parent {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 检测函数封装, 不对比Id
// 参数1: 查到的指针 参数2: 要对比的对象指针
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func TypeCheckNoId(query *models.Type, obj *models.Type, t *testing.T, errTitleStr string) {
	// 判断第1个
	if query.NameId != obj.NameId ||
		query.Name != obj.Name ||
		query.Url != obj.Url ||
		query.Level != obj.Level ||
		query.Parent != obj.Parent {
		// t.Errorf("【查 by nameId 】测试不通过, got= %v", query)
		t.Errorf(" %s 测试不通过, got= %v", errTitleStr, query)
	}
}

// 检测更新函数封装, 对比Id
// 参数1: 查到的指针 参数2: 更新参数 map[string]interface{}
// 参数3: 测试对象指针 t *testing.T  参数4:错误标题字符串，如: 【查 by nameId】中括号里内容
func TypeCheckUpdateHasId(query *models.Type, obj map[string]interface{}, t *testing.T, errTitleStr string) {
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
func TypeCheckUpdateNoId(query *models.Type, obj map[string]interface{}, t *testing.T, errTitleStr string) {
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
func TestTypeAdd(t *testing.T) {
	t.Log("------------ comicType add ...  start ")

	// 1. 测试项1，有0值
	comicType := comicTypeForAddHasIdHasZero
	t.Log("comicType: ", comicType)
	TypeAdd(comicType)

	// var createdType *models.Type // 手动写法
	// testDB.Where("name_id = ?", comicType.NameId).First(&createdType) // 手动写法,不调用方法
	createdType := TypeQueryByNameId(comicType.NameId) // 调用方法
	TypeCheckNoId(createdType, comicType, t, "【增】")    // 测试第1个

	// 2. 测试项2 无0值, 测试needProxy =1 时候
	comicType2 := comicTypeForAddHasIdNoZero
	t.Log("comicType2: ", comicType2)
	TypeAdd(comicType2)
	createdType2 := TypeQueryByNameId(comicType2.NameId) // 调用方法
	TypeCheckNoId(createdType2, comicType2, t, "【增】")    // 测试第2个

	t.Log("----------- comicType add ... end ----------------")
}

// 批量增
func TestTypeBatchAdd(t *testing.T) {
	t.Log("------------ comicType batch add ... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{})

	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero

	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)
	nameIds := []int{comicType.NameId, comicType2.NameId}
	t.Log("namedis = ", nameIds)
	createdTypes, err := TypesBatchQueryByNameId(nameIds) // 调用方法
	if err != nil {
		t.Errorf("【增-批量】测试不通过, 查询nil, got=  %v", createdTypes)
		panic("【增-批量】测试不通过")
	}

	// 判断第1个
	createdType := createdTypes[0]
	createdType2 := createdTypes[1]
	t.Log("查询结果 createdTypes = ", createdType, createdType2)
	TypeCheckNoId(createdType, comicType, t, "【增-批量 】")   // 判断第1个
	TypeCheckNoId(createdType2, comicType2, t, "【增-批量 】") // 判断第2个
	t.Log("------------ comicType batch add ... end ----------------")
}

// 删-通过id
func TestTypeDeleteById(t *testing.T) {
	t.Log("------------ comicType delete by id... start ----------------")
	comicType := comicTypeForAddHasIdNoZero
	TypeAdd(comicType)

	TypeDeleteById(comicType.Id)

	var deletedType models.Type
	result := testDB.First(&deletedType, comicType.Id)
	// result := testDB.Where("name_id = ?", comicType.NameId).First(&deletedType)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by id】测试不通过,删除后仍能查到, =  %v", deletedType)
		panic("【删 - by id 】测试不通过,删除后仍能查到")
	}
	t.Log("------------ comicType delete by id... end ----------------")
}

// 删-通过 nameId
func TestTypeDeleteByNameId(t *testing.T) {
	t.Log("------------ comicType delete by nameId... start ----------------")
	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	TypeDeleteByNameId(comicType.NameId)

	var deletedType models.Type
	result := testDB.Where("name_id = ?", comicType.NameId).First(&deletedType)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by nameId 】 测试不通过,删除后仍能查到, =  %v", deletedType)
		panic("【删 - by nameId 】 测试不通过,删除后仍能查到")
	}
	t.Log("------------ comicType delete by nameId... end ----------------")
}

// 删-通过 其它
func TestTypeDeleteByOther(t *testing.T) {
	t.Log("------------ comicType delete by other... start ----------------")
	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	TypeDeleteByOther("name_id", comicType.NameId)

	var deletedType models.Type
	result := testDB.Where("name_id = ?", comicType.NameId).First(&deletedType)
	if result.Error == nil { // err是空, 说明记录存在
		t.Errorf("【删 - by other 】测试不通过,删除后仍能查到, =  %v", deletedType)
		panic("【删 - by other 】测试不通过,删除后仍能查到")
	}
	t.Log("------------ comicType delete by other... end ----------------")
}

// 删-批量 通过id
func TestTypesBatchDeleteById(t *testing.T) {
	t.Log("------------ comicType batch delete by id... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	comicType := comicTypeForAddHasIdNoZero

	comicType2 := &models.Type{
		Id:        2,
		NameId:    2,
		Name:      "Test Type for Delete By Id 2",
		Url:       "http://delete.com id 2",
		NeedProxy: 0,
		IsHttps:   0,
	}
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes) // 添加

	ids := []uint{comicType.Id, comicType2.Id}
	t.Log("ids = ", ids)

	// 判断是否添加了2个
	comicTypes, err := TypesBatchQueryById(ids)
	if len(comicTypes) != 2 || err != nil {
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", comicTypes)
		// panic("【删 批量 - by id 】测试不通过,删除后仍能查到") // 测试v原本不能用pnic
	}

	TypesBatchDeleteById(ids) // 删除

	// 检测，如果报错，或者 结果>0
	comicTypes, err = TypesBatchQueryById(ids)

	if len(comicTypes) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by id】测试不通过,删除后仍能查到, got %v", comicTypes)
	}

	t.Log("------------ comicType batch delete by id... end ----------------")
}

// 删-批量 通过 nameId
func TestTypesBatchDeleteByNameId(t *testing.T) {
	t.Log("------------ comicType batch delete by nameId... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero

	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes) // 添加

	nameIds := []int{comicType.NameId, comicType2.NameId}
	t.Log("nameIds = ", nameIds)

	// 判断是否添加了2个
	comicTypes, err := TypesBatchQueryByNameId(nameIds)
	if len(comicTypes) != 2 || err != nil {
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", comicTypes)
	}

	TypesBatchDeleteByNameId(nameIds) // 删除

	// 检测，如果报错，或者 结果>0
	comicTypes, err = TypesBatchQueryByNameId(nameIds)

	if len(comicTypes) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by nameId 】测试不通过,删除后仍能查到, got %v", comicTypes)
	}

	t.Log("------------ comicType batch delete by nameId... end ----------------")
}

// 删-批量 通过 other
func TestTypesBatchDeleteByOther(t *testing.T) {
	t.Log("------------ comicType batch delete by other... start ----------------")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 删除数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero

	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes) // 添加

	others := []any{comicType.NameId, comicType2.NameId}
	t.Log("others = ", others)

	// 判断是否添加了2个
	comicTypes, err := TypesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if len(comicTypes) != 2 || err != nil {
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", comicTypes)
	}

	TypesBatchDeleteByOther("name_id", others) // 删除

	// 检测，如果报错，或者 结果>0
	comicTypes, err = TypesBatchQueryByOther("name_id", others, "name_id", "ASC")

	if len(comicTypes) > 0 || err != nil { // 判断错放后面，因为是 ||, 第一个不通过，就不判断第2个
		t.Errorf("【删 批量- by other 】测试不通过,删除后仍能查到, got %v", comicTypes)
	}

	t.Log("------------ comicType batch delete by other... end ----------------")
}

// 改 by id
func TestTypeUpdateById(t *testing.T) {
	t.Log("------------ comicType update by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddHasIdNoZero
	TypeAdd(comicType)

	updates := comicTypeForUpdateHasIdHasZero
	TypeUpdateById(comicType.Id, updates)

	// 检查
	updatedType := TypeQueryById(comicType.Id)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedType =", updatedType)
	t.Log("更新后 查的 updatedType.== =", updatedType.Id == updates["Id"]) // 得转成uint
	t.Log("更新后 查的 updatedType.== =", updatedType.IsHttps == updates["IsHttps"])
	t.Log("更新后 查的 updatedType.== =", updatedType.NameId == updates["NameId"])
	t.Log("更新后 查的 updatedType.== =", updatedType.Name == updates["Name"])
	t.Log("更新后 查的 updatedType.== =", updatedType.Url == updates["Url"])
	t.Log("更新后 查的 updatedType.== =", updatedType.NeedProxy == updates["NeedProxy"])
	t.Log("更新后 查的 updatedType.== =", updatedType.IsHttps == updates["IsHttps"])
	TypeCheckUpdateHasId(updatedType, updates, t, "【改 by id 】")
	t.Log("------------ comicType update by id ... end ")
}

// 改 by nameId
func TestTypeUpdateByNameId(t *testing.T) {
	t.Log("------------ comicType update by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	updates := comicTypeForUpdateNoIdHasZero
	TypeUpdateByNameId(comicType.NameId, updates)

	// 检查
	updatedType := TypeQueryByNameId(comicType.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedType =", updatedType)
	TypeCheckUpdateNoId(updatedType, updates, t, "【改 by nameId 】")
	t.Log("------------ comicType update by nameId ... end ")
}

// 改 by other
func TestTypeUpdateByOther(t *testing.T) {
	t.Log("------------ comicType update by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	updates := comicTypeForUpdateNoIdHasZero
	TypeUpdateByOther("name_id", comicType.NameId, updates)

	// 检查
	updatedType := TypeQueryByOther("name_id", comicType.NameId)
	t.Log("更新后 原始数据 updates =", updates)
	t.Log("更新后 查的 updatedType =", updatedType)
	TypeCheckUpdateNoId(updatedType, updates, t, "【改 by other 】")
	t.Log("------------ comicType update by other ... end ")
}

// 改-批量 by id
func TestTypeBatchUpdateById(t *testing.T) {
	t.Log("------------ comicType batch update by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddHasIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	updates := comicTypeForUpdateHasIdNoZero
	updates2 := comicType2ForUpdateHasIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	TypesBatchUpdateById(updatesArr)

	// 检测，如果报错
	ids := []uint{1, 2}
	comicTypes, err := TypesBatchQueryById(ids)
	if err != nil {
		t.Errorf("【改 by id 】测试不通过, got= %v", comicTypes)
	}

	updatedType := comicTypes[0]
	updatedType2 := comicTypes[1]
	TypeCheckUpdateHasId(updatedType, updates, t, "【改 by id 】")   // 检测第1个
	TypeCheckUpdateHasId(updatedType2, updates2, t, "【改 by id 】") // 检测第1个
	t.Log("------------ comicType batch update by id ... end ")
}

// 改-批量 by nameId
func TestTypeBatchUpdateByNameId(t *testing.T) {
	t.Log("------------ comicType batch update by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	updates := comicTypeForUpdateNoIdNoZero
	updates2 := comicType2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	TypesBatchUpdateByNameId(updatesArr)

	// 检测，如果报错
	nameIds := []int{
		updates["NameId"].(int),
		updates2["NameId"].(int),
	}
	comicTypes, err := TypesBatchQueryByNameId(nameIds)
	if err != nil {
		t.Errorf("【改 by nameId 】测试不通过, got= %v", comicTypes)
	}

	updatedType := comicTypes[0]
	updatedType2 := comicTypes[1]
	TypeCheckUpdateNoId(updatedType, updates, t, "【改 by nameId 】")   // 检测第1个
	TypeCheckUpdateNoId(updatedType2, updates2, t, "【改 by nameId 】") // 检测第1个
	t.Log("------------ comicType batch update by nameId ... end ")
}

// 改-批量 by other
func TestTypeBatchUpdateByOther(t *testing.T) {
	t.Log("------------ comicType batch update by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table
	// DeleteTableAllData(testDB, &models.Type{}) // 方式2: delete 数据
	// 2. 增加数据
	// 3. 修改数据
	// 4. 判断
	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	updates := comicTypeForUpdateNoIdNoZero
	updates2 := comicType2ForUpdateNoIdHasZero
	updatesArr := []map[string]interface{}{updates, updates2}
	TypesBatchUpdateByOther(updatesArr)

	// 检测，如果报错
	others := []any{
		updates["NameId"],
		updates2["NameId"],
	}
	comicTypes, err := TypesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【改 by other 】测试不通过, got= %v", comicTypes)
	}

	updatedType := comicTypes[0]
	updatedType2 := comicTypes[1]
	TypeCheckUpdateNoId(updatedType, updates, t, "【改 by other 】")   // 检测第1个
	TypeCheckUpdateNoId(updatedType2, updates2, t, "【改 by other 】") // 检测第1个
	t.Log("------------ comicType batch update by other ... end ")
}

// 查 by id
func TestTypeQueryById(t *testing.T) {
	t.Log("------------ comicType query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddHasIdNoZero
	TypeAdd(comicType)

	queryType := TypeQueryById(comicType.Id)
	TypeCheckHasId(queryType, comicType, t, "【 查 by id 】")
	t.Log("------------ comicType query by id ... start ")
}

// 查 by nameId
func TestTypeQueryByNameId(t *testing.T) {
	t.Log("------------ comicType query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	queryType := TypeQueryByNameId(comicType.NameId)
	TypeCheckNoId(queryType, comicType, t, "【 查 by nameId 】")
	t.Log("------------ comicType query by nameId ... start ")
}

// 查 by other
func TestTypeQueryByOther(t *testing.T) {
	t.Log("------------ comicType query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddNoIdNoZero
	TypeAdd(comicType)

	queryType := TypeQueryByOther("name_id", comicType.NameId)
	TypeCheckNoId(queryType, comicType, t, "【 查 by other 】")
	t.Log("------------ comicType query by other ... start ")
}

// 查 - 批量 by id
func TestTypeBatchQueryById(t *testing.T) {
	t.Log("------------ comicType batch query by id ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddHasIdNoZero
	comicType2 := comicType2ForAddHasIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	ids := []uint{comicType.Id, comicType2.Id}
	queryTypes, err := TypesBatchQueryById(ids)
	if err != nil {
		t.Errorf("【查 by id 】测试不通过, got= %v", queryTypes)
	}

	queryType := queryTypes[0]
	queryType2 := queryTypes[1]
	TypeCheckHasId(queryType, comicType, t, "【 查 by id 】")   // 判断第1个
	TypeCheckHasId(queryType2, comicType2, t, "【 查 by id 】") // 判断第2个
	t.Log("------------ comicType batch query by id ... start ")
}

// 查 - 批量 by nameId
func TestTypeBatchQueryByNameId(t *testing.T) {
	t.Log("------------ comicType batch query by nameId ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	nameIds := []int{comicType.NameId, comicType2.NameId}
	queryTypes, err := TypesBatchQueryByNameId(nameIds)
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryTypes)
	}

	queryType := queryTypes[0]
	queryType2 := queryTypes[1]
	TypeCheckNoId(queryType, comicType, t, "【 查 by nameId 】")   // 判断第1个
	TypeCheckNoId(queryType2, comicType2, t, "【 查 by nameId 】") // 判断第2个
	t.Log("------------ comicType batch query by nameId ... start ")
}

// 查 - 批量 by other
func TestTypeBatchQueryByOther(t *testing.T) {
	t.Log("------------ comicType batch query by other ... start ")
	// 1. 清空数据
	TruncateTable(testDB, &models.Type{}) // 方式1： truncate table

	comicType := comicTypeForAddNoIdNoZero
	comicType2 := comicType2ForAddNoIdHasZero
	comicTypes := []*models.Type{comicType, comicType2}
	TypeBatchAdd(comicTypes)

	others := []any{comicType.NameId, comicType2.NameId}
	queryTypes, err := TypesBatchQueryByOther("name_id", others, "name_id", "ASC")
	if err != nil {
		t.Errorf("【查 by nameId 】测试不通过, got= %v", queryTypes)
	}

	queryType := queryTypes[0]
	queryType2 := queryTypes[1]
	TypeCheckNoId(queryType, comicType, t, "【 查 by other 】")   // 判断第1个
	TypeCheckNoId(queryType2, comicType2, t, "【 查 by other 】") // 判断第2个
	t.Log("------------ comicType batch query by other ... start ")
}
