/*
作用：db操作 通用模板
目标：1套代码，解决所有项目 增删改查操作

数据库操作详细日志：不在此go文件打。原因：
  - 此文件如果出错，已经返给上级错误原因了
*/
package db

import (
	"errors"
	"reflect"
	"study-spider-manhua-gin/src/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 增 upsert : 插入或更新
/*
作用简单说：
  - 插入或更新 1条数据

作用详细说:
  - 插入或更新 1条数据
    - 不存在 唯一索引，插入新数据
    - 存在 唯一索引，更新数据

核心思路:
	1. 插入或更新1条数据
		- 不存在 唯一索引，插入新数据
		- 存在 唯一索引，更新数据
	2. 转成唯一索引，的写法有2种
		- 写法1： []clause.Column{{Name: "Name"}, {Name: "Id"}}, // 判断唯一索引: Name+Id
		- 写法2： 传参[]string{"Name", "Id"}, 然后调用方法toGormColumns() 转成gorm写法。推荐!!

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1 modelObj any类型 //需要插入或更新的,数据模型对象 如 comic *models.Comic
		- model 是一条数据对象，而不是 表名
		- model 可以是对象指针，也可以是对象。一般是直接传指针
	2 uniqueIndexArr []string 类型 // 用Model里定义的字段，不用数据库真实列名。唯一索引字段,可以是多个 如 []string{"Name", "Id"}
	3 updateColumnsMap map[string]any 类型 // 更新的字段，可以是多个 如 map[string]any{"Name": "comic.Name", "Id": comic.Id}
	updateColumnsMap []string  类型 // 数据库真实列名，可以是多个 如 [""Name", "Id"]
		- 一种方式是传map,弃用了，这种还得手动往里面塞值
		- 一种方式是：直接传 数据库真实列名，不是model里定义的列名 !!!!!!!!!!! -》 推荐 !!!!!!
*/
// func DBUpsert(modelObj any, uniqueIndexArr []string, updateColumnsMap map[string]any) error {  // 写法1 : 更新内容用map - 弃用
func DBUpsert(modelObj any, uniqueIndexArr []string, updateDBColumnRealNameArr []string) error { // 写法2 : 更新内容用 数据库真实字段名
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 要更新的列 数据。弃用，用方法里的参数就行
	// updateDataMap := map[string]any{
	// 	"country_id":       comic.CountryId,
	// 	"website_id":       comic.WebsiteId,
	// 	"category_id":      comic.CategoryId,
	// 	"type_id":          comic.TypeId,
	// 	"update":           comic.Update,
	// 	"hits":             comic.Hits,
	// 	"comic_url_api_path":        comic.ComicUrlApiPath,
	// 	"cover_url_api_path":        comic.CoverUrlApiPath,
	// 	"brief_short":      comic.BriefShort,
	// 	"brief_long":       comic.BriefLong,
	// 	"end":              comic.End,
	// 	"star":             comic.Star,
	// 	"need_tcp":         comic.NeedTcp,
	// 	"cover_need_tcp":   comic.CoverNeedTcp,
	// 	"spider_end":       comic.SpiderEnd,
	// 	"download_end":     comic.DownloadEnd,
	// 	"upload_aws_end":   comic.UploadAwsEnd,
	// 	"upload_baidu_end": comic.UploadBaiduEnd,
	// }

	// 4. 数据库执行
	// 准备唯一索引数据，写法1 - 不推荐: []clause.Column{{Name: "Name"}, {Name: "Id"}}, // 判断唯一索引: Name+Id
	// result := DB.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{ {Name: "Name"}, {Name: "Id"} }, // 判断唯一索引: Name
	// 	DoUpdates: clause.Assignments(updateDataMap),
	// }).Create(model)

	// 准备唯一索引数据，写法2-推荐: 传参[]string{"Name", "Id"}, 然后调用方法toGormColumns() 转成gorm写法。
	result := DBComic.Clauses(clause.OnConflict{
		Columns: toGormColumns(uniqueIndexArr), // 判断唯一索引: 如：Name + Id。 解释：如果唯一索引冲突
		// DoUpdates: clause.Assignments(updateColumnsMap),  // 写法1 弃用
		DoUpdates: clause.AssignmentColumns(updateDBColumnRealNameArr), // 写法2 推荐，只传数据库 真实列名。解释：就更新这些列。AssignmentColumns 直接从对象里取数据
	}).Create(modelObj)

	if result.Error != nil {
		// log.Error("创建失败: ", result.Error)  // 此文件不打日志，错误已经返回给上级
		return result.Error
	}
	// 5. 返回结果
	return nil // 创建成功
}

// 增-批量 upsertBatch : 批量插入或更新。批量操作，涉及到数据回滚问题
/*
updateDBColumnRealNameArr 必须传数据库真实字段，全小写带_ 的那种

作用简单说：
  - 批量插入或更新 N条数据

作用详细说:
  - 批量插入或更新 N条数据
    - 不存在 唯一索引，插入新数据
    - 存在 唯一索引，更新数据

核心思路:
	1. 批量插入或更新N条数据
		- 不存在 唯一索引，插入新数据
		- 存在 唯一索引，更新数据
	2. 转成唯一索引，的写法有2种
		- 写法1： []clause.Column{{Name: "Name"}, {Name: "Id"}}, // 判断唯一索引: Name+Id
		- 写法2： 传参[]string{"Name", "Id"}, 然后调用方法toGormColumns() 转成gorm写法。推荐!!

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1. DBConnObj *gorm.DB 类型 // 数据库连接对象
		- 比如连了 comic 数据库，就传 comidcDB 对象
		- 比如连了 audiobook 数据库，就传 		- 比如连了 audiobook 数据库，就传 comidcDB 对象
	2 modelObjs any //需要插入或更新的,数据模型对象 如 comic *models.Comic
		- model 是一条数据对象，而不是 表名
		- model 可以是对象指针，也可以是对象。一般是直接传指针
	3 uniqueIndexArr []string 类型 // 用Model里定义的字段，不用数据库真实列名。 唯一索引字段,可以是多个 如 []string{"Name", "Id"}
	4 updateColumnsMap map[string]any 类型 // 更新的字段，可以是多个 如 map[string]any{"Name": "comic.Name", "Id": comic.Id}
	updateColumnsMap []string  类型 // 数据库真实列名，可以是多个 如 [""Name", "Id"]
		- 一种方式是传map,弃用了，这种还得手动往里面塞值
		- 一种方式是：直接传 数据库真实列名，不是model里定义的列名 !!!!!!!!!!! -》 推荐 !!!!!!
*/
func DBUpsertBatch(DBConnObj *gorm.DB, modelObjs any, uniqueIndexArr []string, updateDBColumnRealNameArr []string) error {

	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数

	// 4. 数据库执行

	// 准备唯一索引数据，写法2-推荐: 传参[]string{"Name", "Id"}, 然后调用方法toGormColumns() 转成gorm写法。
	// 事务里 包Upsert操作 。 db.Transcation就是创建事务
	err := DBConnObj.Transaction(func(tx *gorm.DB) error { // 原本写法应该是这样，但是DB用的全局参数，所以tx那里就不用传参了. tx是事务对象
		// 批量插入 + 冲突更新
		result := tx.Omit(clause.Associations).Clauses(clause.OnConflict{
			Columns:   toGormColumns(uniqueIndexArr), // 判断唯一索引: 如：Name + Id。多个条件是 并且的关系。Omit(clause.Associations) -》 为了不更新关联表，只更新主表
			DoUpdates: clause.AssignmentColumns(updateDBColumnRealNameArr),
		}).Select("*").Create(modelObjs) // 等价于 CreateInBatches(users, 1000). Select("*"， “Stats”) -》 为了更新关联表 Stats.Select(*)必须保留，因为Select("*") + Omit(clause.Associations)才能保证不更新关联表

		if result.Error != nil {
			// log.Error("批量插入失败, err = ", result.Error)  // 此文件不打日志，错误已经返回给上级
			return result.Error // 返回错误，事务回滚
		}

		// 执行到这里，说明没问题，事务提交
		return nil // 事务函数返回结果
	})

	if err != nil {
		// log.Error("批量插入失败，事务已回滚:, err = ", err)  // 此文件不打日志，错误已经返回给上级
		return err
	}

	// 5. 返回结果
	return nil // 批量插入成功

	// v0.1 写法 批量插入，循环调用单个插入方法
	// 不推荐。因为：
	// - 不能确保完全插入进去。1000个数据，可能成功1部分，失败1部分
	// - 失败，不能回滚。不能确保 1000个数据完全插入，或者完全未插入
	// - 性能低。1000个数据，循环调用1000次插入方法
	/*
		func ComicBatchAdd(comics []*models.Comic) {
			for i, comic := range comics {
				err := ComicUpsert(comic)
				if err == nil {
					log.Debugf("批量创建第%d条成功, comic: %v", i+1, &comic)
				} else {
					log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
				}
			}
		}
	*/

	// v0.2 写法 事务 里包裹 Upsert，并且每次 按gorm默认实现：插入500条
}

// 配合 DBUpsert 方法使用，把wtring数组，转成gorm 类型的column
/*
执行下面这一步，用的到
Columns:   []clause.Column{{Name: "Name"}}, // 判断唯一索引: Name

作用简单说：
  - 配合 DBUpsert 方法使用，把wtring数组，转成gorm 类型的column

作用详细说:

核心思路:

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 执行核心逻辑
	4. 返回结果

参数：
	columns []string 类型 // 需要插入或更新的,唯一索引字段,可以是多个 如 []string{"Name", "Id"}

返回：
	[]clause.Column 类型 // gorm 类型的column 数组
*/
func toGormColumns(columns []string) []clause.Column {
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 执行核心逻辑
	gormCols := make([]clause.Column, len(columns))
	for i, col := range columns {
		gormCols[i] = clause.Column{Name: col}
	}

	// 4. 返回结果
	return gormCols
}

// 删
/*
作用简单说：
	- 删除1条数据

作用详细说:

核心思路:
	1. 准备删除参数
	- 从方法参数拿id
	2. 调用删除方法
	- 根据id删除 对应表数据 (这里有2种写法，因为DB.Delete要传一个数据表对象)
		写法1: 直接传表对象的指针，这样写的代码更少，更简洁。推荐！！
		写法2：先var一个表 空对象，再用空对象作为参数
	3. 返回错误信息

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1. model any类型 // 数据库表，的模型对象/模型指针
		- 方式1：模型指针。如 comic *models.Comic{} -》 模型指针 -》 推荐，因为写法最简单。不用var 一个对象再传
		- 方式2：模型对象。如 var comic models.Comic{}, 再通过 DB.Delete(模型对象, id)
	2. id int 类型 // 需要删除的,数据id
		- 其实uint 类型，更符合实际情况，因为id都是 正整数
		- 但是int 写法更简单。写代码不用为了防止报错，强制转换类型 -》 推荐
		- 并且 uint，如果传负数，强转 uint( int ) 数值不对
	为什么用int类型，不用 uint。原因如下
		- 如果是负数，强转成uint，结果完全不一样。如 -1 → 18446744073709551615，-100 → 18446744073709551516，-88 → 18446744073709551528
		- 写法麻烦。有的时候，为了不报错，写法要加上 uint(id)
		- 因此，只要人工判断 id >0 就行

返回：
	error 类型 // 错误信息

注意：
	DB.Delete(&models.Comic{}, id) 这种传id写法，gorm会不同处理。建议校验id >0。因为为负数，会执行sql，以防万一
	参数      是否执行SQL      SQL内容
	- 0      否             无操作
	- 1      是             DELETE FROM `comics` WHERE `id` = 1
	- -1     是             DELETE FROM `comics` WHERE `id` = -1
*/
func DBDeleteById(model any, id int) error {
	// 1. 校验传参
	if id <= 0 {
		// log.Error("DB删除失败: id不合法, id <= 0")  // 此文件不打日志，错误已经返回给上级
		return errors.New("DB删除失败: id不合法, id <= 0")
	}

	// 2. 数据清洗
	// 3. 准备删除参数
	// log.Debug("删除漫画, 参数id= ", id)  // 此文件不打日志，错误已经返回给上级

	// 4. 数据库执行
	// -- 写法2：先var一个表 空对象，再用空对象作为参数 --》 不推荐
	// var comic models.Comic
	// result := DB.Delete(&comic, id)

	// -- 写法1: 直接传表对象的指针，这样写的代码更少，更简洁。推荐！！
	result := DBComic.Delete(&models.ComicSpider{}, id)
	if result.Error != nil {
		// log.Error("删除失败: ", result.Error)  // 此文件不打日志，错误已经返回给上级
		return result.Error
	}

	// 5. 返回结果
	return nil // 删除成功
}

// 改 - 根据id, 排除唯一索引 参数用结构体
/*
疑问: 为什么要排除唯一索引字段?
答: 唯一索引很关键,作用比id还重要。防止误更新 唯一索引字段

作用简单说：
  - 更新
  	- 只更新 指定字段，如DB.Model().Select(指定字段).Update()，中Select()中的字段
	- 不更新 唯一索引字段。如唯一索引叫 name, 写代码的时候要排除它
	- 参数中，有0值，也会更新

作用详细说:

核心思路:
	1. 准备要用的参数
	2. 调用DB方法
	3. 返回错误信息

更新操作，并排除唯一索引，一般有4种写法：
	方式1：只调Updates()方法，不调用Select()方法。-》 不推荐，原因见下面
		举例：DB.Model().Updates() -》 问题：如果字段是0值，不会更新该字段 (因为：gorm默认就这么实现的)
	方式2：Select(要更新字段).Updates() -》 也不推荐。原因：写法乱，见方法内代码
	方式3：只调Updates()方法，传入 map[string]interface{} -》 推荐！！
	方式4：DB.Model().Omit(要排除字段).Updates() -》 不推荐。因为不安全。具体原因见下面
		原因：有风险！。因为如果有的列忘记传数据了，会更新成默认值

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1. model any类型 // 数据库表，的模型对象   不能用 模型指针，因为更新参数，需要对象才能调出来。具体看 updateDataMap代码
		- 方式1：模型指针。如 comic *models.Comic{} -》 模型指针 -》 推荐，因为写法最简单。不用var 一个对象再传
		- 方式2：模型对象。如 var comic models.Comic{}, 再通过 DB.Delete(模型对象, id)
	2. id int 类型 // 需要修改的,数据id
		- 其实uint 类型，更符合实际情况，因为id都是 正整数
		- 但是int 写法更简单。写代码不用为了防止报错，强制转换类型 -》 推荐
		- 并且 uint，如果传负数，强转 uint( int ) 数值不对
	为什么用int类型，不用 uint。原因如下
		- 如果是负数，强转成uint，结果完全不一样。如 -1 → 18446744073709551615，-100 → 18446744073709551516，-88 → 18446744073709551528
		- 写法麻烦。有的时候，为了不报错，写法要加上 uint(id)
		- 因此，只要人工判断 id >0 就行
	3. updateDataMap map[string]any 类型 // 需要更新的数据，key是字段名，value是字段值

返回：
	error 类型 // 错误信息

// id 可以int,可以string。go默认定义的 any = interface{},忘了写这个注释啥意思
*/
func DBUpdateByIdOmitIndex(model any, id int, updateDataMap map[string]any) error {
	// 1. 校验传参
	if id <= 0 {
		// log.Error("DB修改失败: id不合法, id <= 0")  // 此文件不打日志，错误已经返回给上级
		return errors.New("DB修改失败: id不合法, id <= 0")
	}

	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 更新参数 - 参数里有了，下面弃用了
	/*
		updateDataMap := map[string]any{
			"country_id":       comic.CountryId,
			"website_id":       comic.WebsiteId,
			"category_id":      comic.CategoryId,
			"type_id":          comic.TypeId,
			"update":           comic.Update,
			"hits":             comic.Hits,
			"comic_url_api_path":        comic.ComicUrlApiPath,
			"cover_url_api_path":        comic.CoverUrlApiPath,
			"brief_short":      comic.BriefShort,
			"brief_long":       comic.BriefLong,
			"end":              comic.End,
			"need_tcp":         comic.NeedTcp,
			"cover_need_tcp":   comic.CoverNeedTcp,
			"spider_end":       comic.SpiderEnd,
			"download_end":     comic.DownloadEnd,
			"upload_aws_end":   comic.UploadAwsEnd,
			"upload_baidu_end": comic.UploadBaiduEnd,
		}
	*/

	// 4. 数据库执行
	// -- 调用DB方法 方式：
	// 方式4，不推荐
	// result := DB.Model(&comic).Where("id = ?", comicId).Omit("name").Updates(comic)
	// 方式2，也不推荐。安全，写法没问题，就是乱
	// result := DB.Model(&comic).Where("id = ?", comicId).Select("country_id", "website_id",
	// 	"category_id", "type_id", "update", "hits", "comic_url_api_path",
	// 	"cover_url_api_path", "brief_short", "brief_long", "end", "need_tcp", "cover_need_tcp",
	// 	"spider_end", "download_end", "upload_aws_end", "upload_baidu_end").Updates(comic)

	// 方式3：只调Updates()方法，传入 map[string]interface{} -》 推荐！！
	result := DBComic.Model(model).Where("id = ?", id).Updates(updateDataMap)
	if result.Error != nil {
		// log.Error("修改失败: ", result.Error)  // 此文件不打日志，错误已经返回给上级
		return result.Error
	}
	// 5. 返回结果
	return nil // 修改成功
}

// 查 - 分页查询。方式：返回 any类型数据. 第2推荐这种方式，因为: 不是最安全。非常不推荐！！写法麻烦！！
/*
作用简单说：
	- 分页查询

作用详细说:

核心思路:
	1.
	2.
	3.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1. model any类型 // 数据库表，的模型对象/模型指针
		- 方式1：模型指针-推荐！！。如 comic *models.Comic{} -》 模型指针 -》 推荐，因为写法最简单。不用var 一个对象再传
		- 方式2：模型对象-不推荐。如 var comic models.Comic
	2. pageNum int 类型 // 页码，实际sql是从0开始，但是要求用户传1开始
		- 用户传值 从1开始
	3. pageSize int 类型 // 每页数量

返回：
	[]*models.Comic 类型 // 数据列表 , any类型
	error 类型 // 错误信息

注意：
	要返回 结果数据，有2种实现方式：
	- 方式1：使用any传参    -》 不推荐。因为 类型不安全 -》就是类型有错会触发panic，可能影响程序
	- 方式2：使用泛型传参 	-》 推荐！！因为 类型安全 -》就是类型有错，可能会触发panic，但更安全
*/
func DBPageQueryReturnTypeAny(model any, pageNum, pageSize int) ([]any, error) {
	//

	// 1. 校验传参
	// 校验 <=0
	if pageNum <= 0 || pageSize <= 0 {
		// log.Error("DB分页查询失败: pageNum或pageSize参数值错误, 应>0")  // 此文件不打日志，错误已经返回给上级
		return nil, errors.New("DB分页查询失败: pageNum或pageSize参数值错误, 应>0")
	}

	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行

	// -- 根据传的表类型参数，获取到是什么表， 比如传 models.Comic{}，就知道是 comic表
	// var modelObjs []model  // v0.1 这种写法根本不行，必须用反射获取到 表类型

	// v0.2 反射获取到 表类型
	// 反射创建对应类型的 slice
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	sliceType := reflect.SliceOf(modelType)
	sliceValue := reflect.New(sliceType)

	// -- 执行db操作
	result := DBComic.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(sliceValue.Interface())

	// v0.1 写法，会报错，弃用
	/*
		if result.Error != nil {
			log.Error("分页查询失败: ", result.Error)
			return modelObjs, result.Error
		}
		log.Infof("分页查询成功, 查询到 %d 条记录", len(modelObjs))
	*/
	// v0.2 写法，不会报错,反射实现
	if result.Error != nil {
		// log.Error("分页查询失败: ", result.Error)  // 此文件不打日志，错误已经返回给上级
		return nil, result.Error
	}

	// 将结果转为 []any
	slice := sliceValue.Elem()
	results := make([]any, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		results[i] = slice.Index(i).Interface()
	}

	// log.Infof("分页查询成功, 查询到 %d 条记录", len(results))  // 此文件不打日志，错误已经返回给上级

	// 5. 返回结果
	return results, nil

	// ------------------ 如何调用本方法
	// 缺点：返回的是 any，调用方需要做类型断言
	/*
		comics, _ := DBPageQueryReturnTypeAny(1, 10)
		for _, c := range comics {
			comic := c.(*models.Comic) // ⚠️ 如果断言失败会 panic。这一行就是类型断言
		}
	*/
}

// 查 - 分页查询。方式：返回 T 泛型数据。推荐
/*
作用简单说：
	- 分页查询

作用详细说:

核心思路:
	1.
	2.
	3.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1. model any类型 // 数据库表，的模型对象/模型指针
		- 方式1：模型指针-推荐！！。如 comic *models.Comic{} -》 模型指针 -》 推荐，因为写法最简单。不用var 一个对象再传
		- 方式2：模型对象-不推荐。如 var comic models.Comic
	2. pageNum int 类型 // 页码，实际sql是从0开始，但是要求用户传1开始
		- 用户传值 从1开始
	3. pageSize int 类型 // 每页数量

返回：
	[]*models.Comic 类型 // 数据列表 , any类型
	error 类型 // 错误信息

注意：
	要返回 结果数据，有2种实现方式：
	- 方式1：使用any传参    -》 不推荐。因为 类型不安全 -》就是类型有错会触发panic，可能影响程序
	- 方式2：使用泛型传参 	-》 推荐！！因为 类型安全 -》就是类型有错，可能会触发panic，但更安全
*/
func DBPageQueryReturnTypeT[T any](pageNum, pageSize int) ([]T, error) {
	//

	// 1. 校验传参
	// 校验 <=0
	if pageNum <= 0 || pageSize <= 0 {
		// log.Error("DB分页查询失败: pageNum或pageSize参数值错误, 应>0") // 此文件不打日志，错误已经返回给上级
		return nil, errors.New("DB分页查询失败: pageNum或pageSize参数值错误, 应>0")
	}

	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行

	var modelObjs []T
	result := DBComic.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&modelObjs)
	if result.Error != nil {
		// log.Error("分页查询失败: ", result.Error)  // 此文件不打日志，错误已经返回给上级
		return modelObjs, result.Error
	}
	// log.Infof("分页查询成功, 查询到 %d 条记录", len(modelObjs))  // 此文件不打日志，错误已经返回给上级

	// 5. 返回结果
	return modelObjs, result.Error

	// ------------------ 如何调用本方法
	/*
		comics, err := DBPageQuery[*models.Comic](1, 10)
		if err != nil {
			log.Error(err)
		}
		for _, comic := range comics {
			fmt.Println(comic.Title)
		}
	*/
}
