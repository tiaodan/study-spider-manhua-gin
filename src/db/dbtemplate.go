/*
作用：db操作 通用模板
目标：1套代码，解决所有项目 增删改查操作

数据库操作详细日志：不在此go文件打。原因：
  - 此文件如果出错，已经返给上级错误原因了
*/
package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"unicode"

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
func DBUpsert(DBConnObj *gorm.DB, modelObj any, uniqueIndexArr []string, updateDBColumnRealNameArr []string) error { // 写法2 : 更新内容用 数据库真实字段名
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
	result := DBConnObj.Clauses(clause.OnConflict{
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
	3 uniqueIndexArr []string 类型 // 用Model里定义的字段，用数据库真实列名也行，首字母大写也行，首字母小写也行。 唯一索引字段,可以是多个 如 []string{"Name", "Id"}
		注意：
		- 首字母用大写，小写均可以。建议用首字母大写，因为 大写更适合 struct 结构定义，显得更规范
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
func DBUpdateByIdOmitIndex_nouse_bymap(model any, id int, updateDataMap map[string]any) error {
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

// 根据指定字段查询 - 通用，使用于任何数据表
/*
作用简单说：
	- 查询 1条数据

作用详细说:
	-

核心思路:
	1.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1 field string 用数据库小写列名，如 "chapter_num"
*/
func DBFindOneByField[T any](field string, value any) (*T, error) {
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行
	var result T
	db := DBComic.Where(field+" = ?", value).First(&result)
	if db.Error != nil {
		// log.Error("查询失败: ", db.Error)  // 此文件不打日志，错误已经返回给上级
		return nil, db.Error
	}
	// log.Infof("查询成功, 查询到 %d 条记录", 1)  // 此文件不打日志，错误已经返回给上级

	// 5. 返回结果
	return &result, nil
}

// 根据map条件查询 - 通用，使用于任何数据表
/*
作用简单说：
	- 使用map条件查询1条数据

作用详细说:
	-

核心思路:
	1.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1 condition map[string]interface{} 查询条件，如 map[string]interface{}{"name": "xxx", "type_id": 1}
*/
func DBFindOneByMapCondition[T any](condition map[string]any) (*T, error) {
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行
	var result T
	db := DBComic.Where(condition).First(&result)
	if db.Error != nil {
		// log.Error("查询失败: ", db.Error)  // 此文件不打日志，错误已经返回给上级
		return nil, db.Error
	}
	// log.Infof("查询成功, 查询到 %d 条记录", 1)  // 此文件不打日志，错误已经返回给上级

	// 5. 返回结果
	return &result, nil
}

// 根据对象的唯一索引字段查询 - 通用，使用于任何数据表
/*
作用简单说：
	- 根据obj 和提供的索引，生成一个只有索引的 map，并根据此map，查询1条数据
*/
func DBFindOneByUniqueIndexMapCondition[T any](obj *T, uniqueIndexArr []string) (*T, error) {
	/* 之前代码 v0.1 不用
	var existingComic models.ComicSpider
	condition := map[string]interface{}{
		"name":          onePageBookArr[i].Name,
		"country_id":    onePageBookArr[i].CountryId,
		"website_id":    onePageBookArr[i].WebsiteId,
		"porn_type_id":  onePageBookArr[i].PornTypeId,
		"type_id":       onePageBookArr[i].TypeId,
		"author_concat": onePageBookArr[i].AuthorConcat,
	}
	result := db.DBComic.Where(condition).First(&existingComic)
	if result.Error == nil {
		// 更新对象的ID为数据库中的实际ID
		onePageBookArr[i].Id = existingComic.Id
		log.Debugf("更新comic ID: %s -> %d", onePageBookArr[i].Name, existingComic.Id)
	} else {
		log.Errorf("查询comic失败: %s, err: %v", onePageBookArr[i].Name, result.Error)
	}
	*/

	/* v0.2 tongyilingma AI 给鸡巴改坏了
	// 1. 获取obj 的反射信息
	val := reflect.ValueOf(obj).Elem() // 反射值
	// typ := val.Type()                  // 反射类型,暂时没用着

	// 先测试生成查询 condition map
	log.Debug("--------- delete 生成的对象 obj = ", obj)
	log.Debug("--------- delete 生成的索引 uniqueIndexArr = ", uniqueIndexArr)
	// 2. 构建查询条件 map
	whereConditions := make(map[string]interface{})
	// 遍历 判断 uniqueIndexArr 中的值. 是由在 obj的key里。通过能不能获取到 obj key的值,来判断
	for _, uniqueIndex := range uniqueIndexArr {
		fieldValue := val.FieldByName(uniqueIndex)
		if fieldValue.IsValid() {
			whereConditions[uniqueIndex] = fieldValue.Interface()
		} else {
			log.Debugf("通过唯一索引字段+obj, 查找1个obj失败, 通过索引=%s 获取不到 obj的值", uniqueIndex)
			return nil, errors.New("字段 " + uniqueIndex + " 不存在")
		}
	}
	log.Debug("------------ delete , whereconditon = ", whereConditions)

	// 构建查询条件后，执行查询
	var result T
	db := DBComic.Where(whereConditions).First(&result)
	if db.Error != nil {
		return nil, db.Error
	}
	return &result, nil
	*/

	//  v0.3 修改 tongyilingma AI 改错的，垃圾
	// 1. 获取obj 的反射信息
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	// 2. 构建查询条件 map
	whereConditions := make(map[string]interface{})
	for _, uniqueIndex := range uniqueIndexArr {
		field, found := typ.FieldByName(uniqueIndex)
		if !found {
			log.Debugf("通过唯一索引字段+obj, 查找1个obj失败, 通过索引=%s 获取不到 obj的值", uniqueIndex)
			return nil, errors.New("字段 " + uniqueIndex + " 不存在")
		}

		fieldValue := val.FieldByName(uniqueIndex)
		if fieldValue.IsValid() {
			// 使用 getColumnName 转换为数据库列名
			columnName := getColumnName(field)
			whereConditions[columnName] = fieldValue.Interface()
		} else {
			log.Debugf("通过唯一索引字段+obj, 查找1个obj失败, 通过索引=%s 获取不到 obj的值", uniqueIndex)
			return nil, errors.New("字段 " + uniqueIndex + " 不存在")
		}
	}

	log.Debug("------------ where condition = ", whereConditions)

	// 构建查询条件后，执行查询
	var result T
	db := DBComic.Where(whereConditions).First(&result)
	if db.Error != nil {
		return nil, db.Error
	}
	return &result, nil
}

// 根据对象查询 - 通用，使用于任何数据表
/*
作用简单说：
	- 使用对象的非零值字段作为条件查询1条数据

作用详细说:
	- 使用反射获取对象的字段和值，忽略零值字段，将非零值字段作为查询条件

核心思路:
	1. 使用反射获取对象字段
	2. 过滤掉零值字段
	3. 构建查询条件

参数：
	1 obj T 查询对象，其非零值字段将作为查询条件
*/
func DBFindOneByStruct[T any](obj *T) (*T, error) {
	// 1. 使用反射获取对象字段
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()

	// 2. 构建查询条件 map
	whereConditions := make(map[string]interface{})

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		// 获取数据库列名
		columnName := getColumnName(fieldType)

		// 检查字段是否为零值，如果是则跳过
		if !isZeroValue(field) {
			whereConditions[columnName] = field.Interface()
		}
	}

	// 3. 使用构建的条件进行查询
	var result T
	db := DBComic.Where(whereConditions).First(&result)
	if db.Error != nil {
		// log.Error("查询失败: ", db.Error)  // 此文件不打日志，错误已经返回给上级
		return nil, db.Error
	}
	// log.Infof("查询成功, 查询到 %d 条记录", 1)  // 此文件不打日志，错误已经返回给上级

	// 4. 返回结果
	return &result, nil
}

// isZeroValue 检查反射值是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		// 对于结构体类型，我们可能需要递归检查，这里简单返回false
		return false
	default:
		return v.IsZero()
	}
}

// update - 通用方法
/*
uniqueIndexUpperCaseArr 大写写法- strct 里的 字段
*/
func DBUpdate(DBConnObj *gorm.DB, modelObj any, uniqueIndexUpperCaseArr []string, updateDBColumnRealNameArr []string) error {
	// 1. 校验传参
	if len(uniqueIndexUpperCaseArr) == 0 {
		return errors.New("DB更新失败: uniqueIndexUpperCaseArr 不能为空")
	}
	if len(updateDBColumnRealNameArr) == 0 {
		return errors.New("DB更新失败: updateDBColumnRealNameArr 不能为空")
	}

	// 2. 使用反射从 modelObj 中提取 uniqueIndexUpperCaseArr 对应字段的值作为查询条件
	modelValue := reflect.ValueOf(modelObj)
	modelType := reflect.TypeOf(modelObj)

	// 处理指针类型
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	// 构建查询条件 map（只包含 uniqueIndexArr 指定的字段）
	whereConditions := make(map[string]interface{})
	for _, fieldName := range uniqueIndexUpperCaseArr {
		// 查找结构体字段（不区分大小写）
		field, found := modelType.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, fieldName)
		})
		if !found {
			return errors.New("DB更新失败: 在模型中找不到字段 " + fieldName)
		}

		fieldValue := modelValue.FieldByIndex(field.Index)

		// 获取数据库列名
		columnName := getColumnName(field)
		whereConditions[columnName] = fieldValue.Interface()
	}

	// 3. 执行更新
	result := DBConnObj.Model(modelObj).Where(whereConditions).Select(updateDBColumnRealNameArr).Updates(modelObj)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DBUpdateById 通过主键id更新记录
/*
基于DBUpdate()函数实现，专门用于通过主键id更新
DBConnObj 数据库连接对象
modelObj 模型对象指针
id 主键id值
updateDBColumnRealNameArr 要更新的数据库字段名数组（小写蛇形命名）
*/
func DBUpdateById(DBConnObj *gorm.DB, modelObj any, id int, updateDBColumnRealNameArr []string) error {
	// 1. 校验传参
	if id <= 0 {
		return errors.New("DB更新失败: id必须大于0")
	}
	if len(updateDBColumnRealNameArr) == 0 {
		return errors.New("DB更新失败: updateDBColumnRealNameArr 不能为空")
	}

	// 2. 执行更新
	// WHERE id = ? AND 只更新指定字段
	result := DBConnObj.Model(modelObj).Where("id = ?", id).Select(updateDBColumnRealNameArr).Updates(modelObj)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// getColumnName 获取结构体字段对应的数据库列名
func getColumnName(field reflect.StructField) string {
	// 检查gorm标签中的column
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		// 解析gorm标签，查找column
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// 如果没有指定column，使用默认的snake_case转换
	fieldName := field.Name
	var result []rune
	for i, r := range fieldName {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// 获取总数，通过条件 : map条件参数。现在没用这，用的时候再说
func DBCountByMapConditon(db *gorm.DB, model interface{}, conds ...interface{}) int {
	return 0
}

// 获取总数，通过某个列名 : 比如 parent_id ，这里参数要用 小写_,用数据库真实列名形式
func DBCountByField[T any](DBConnObj *gorm.DB, field string, value any) (int, error) {
	// v0.2写法，err判断简洁
	var count int64
	err := DBConnObj.Model(new(T)).Where(field+" = ?", value).Count(&count).Error

	return int(count), err

	/* v0.1写法，err判断不够简洁
	var count int64
	db := DBConnObj.Model(new(T)).Where(field+" = ?", value).Count(&count)

	if db.Error != nil {
		return 0, db.Error
	}
	// 5. 返回结果
	return int(count), nil
	*/
}

// 批量更新
func DBUpdateBatchByIdArr[T any](DBConnObj *gorm.DB, idArr []int, updates map[string]any) error {
	if len(idArr) == 0 { // 没有数据
		return nil
	}

	return DBConnObj.Model(new(T)).Where("id in (?)", idArr).Updates(updates).Error
}

// 根据指定字段查询 多个 - 通用，使用于任何数据表
/*
作用简单说：
	- 查询 1条数据

作用详细说:
	-

核心思路:
	1.

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 准备数据库执行，需要的参数
	4. 数据库执行
	5. 返回结果

参数：
	1 field string 用数据库小写列名，如 "chapter_num"
*/
func DBFindManyByField[T any](field string, value any) ([]T, error) {
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行

	// 1. 参数校验
	if field == "" {
		return nil, fmt.Errorf("field is empty")
	}

	// 2. 查询
	var result []T
	db := DBComic.Where(field+" = ?", value).Find(&result)

	if db.Error != nil {
		return nil, db.Error
	}

	// 5. 返回结果
	return result, nil
}
