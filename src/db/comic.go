// db comic 相关操作
package db

import (
	// 导入 clause 包
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

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
思路:
   1. 插入或更新1条数据
	- 不存在 唯一索引，插入新数据
	- 存在 唯一索引，更新数据
*/
func ComicUpsert(comic *models.Comic) error {
	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "Name"}}, // 判断唯一索引: Name
		DoUpdates: clause.Assignments(map[string]interface{}{
			"country_id":       comic.CountryId,
			"website_id":       comic.WebsiteId,
			"sex_type_id":      comic.SexTypeId,
			"type_id":          comic.TypeId,
			"update":           comic.Update,
			"hits":             comic.Hits,
			"comic_url":        comic.ComicUrl,
			"cover_url":        comic.CoverUrl,
			"brief_short":      comic.BriefShort,
			"brief_long":       comic.BriefLong,
			"end":              comic.End,
			"star":             comic.Star,
			"need_tcp":         comic.NeedTcp,
			"cover_need_tcp":   comic.CoverNeedTcp,
			"spider_end":       comic.SpiderEnd,
			"download_end":     comic.DownloadEnd,
			"upload_aws_end":   comic.UploadAwsEnd,
			"upload_baidu_end": comic.UploadBaiduEnd,
		}),
	}).Create(comic)
	if result.Error != nil {
		log.Error("创建失败: ", result.Error)
		return result.Error
	} else {
		log.Info("创建成功: ", comic)
	}
	return nil
}

// 批量增
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

// 删
/*
作用简单说：
  - 删除1条数据

作用详细说:

思路:
  1. 准备删除参数
	- 从方法参数拿id
  2. 调用删除方法
	- 根据id删除 对应表数据 (这里有2种写法，因为DB.Delete要传一个数据表对象)
	  写法1: 直接传表对象的指针，这样写的代码更少，更简洁。推荐！！
	  写法2：先var一个表 空对象，再用空对象作为参数
  3. 返回错误信息
*/
func ComicDelete(id uint) error {
	// 1. 准备删除参数
	log.Debug("删除漫画, 参数id= ", id)

	// 2. 调用删除方法
	// -- 写法2：先var一个表 空对象，再用空对象作为参数 --》 不推荐
	// var comic models.Comic
	// result := DB.Delete(&comic, id)

	// -- 写法1: 直接传表对象的指针，这样写的代码更少，更简洁。推荐！！
	result := DB.Delete(&models.Comic{}, id)
	if result.Error != nil {
		log.Error("删除失败: ", result.Error)
		return result.Error
	} else {
		log.Info("删除成功: ", id)
	}

	// 3. 返回错误信息
	return nil
}

// 批量删
/*
现在实现方式：
	- 循环调用单一删除。实现简单，但不是业界推荐
	- 业界推荐：1条sql一次性删除
	- 考虑软删除+硬删除

作用简单说：
  - 删除1条数据

作用详细说:

*/
func ComicsBatchDelete(ids []uint) {
	var comics []models.Comic
	result := DB.Delete(&comics, ids)
	if result.Error != nil {
		log.Error("批量删除失败: ", result.Error)
	} else {
		log.Debug("批量删除成功: ", ids)
	}
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

思路:
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

// id 可以int,可以string。go默认定义的 any = interface{},忘了写这个注释啥意思
*/
func ComicUpdateByIdOmitIndex(comicId any, comic *models.Comic) error {
	// 1. 准备要用的参数

	// 2. 调用DB方法
	// 方式4，不推荐
	// result := DB.Model(&comic).Where("id = ?", comicId).Omit("name").Updates(comic)
	// 方式2，也不推荐。安全，写法没问题，就是乱
	// result := DB.Model(&comic).Where("id = ?", comicId).Select("country_id", "website_id",
	// 	"category_id", "type_id", "update", "hits", "comic_url",
	// 	"cover_url", "brief_short", "brief_long", "end", "need_tcp", "cover_need_tcp",
	// 	"spider_end", "download_end", "upload_aws_end", "upload_baidu_end").Updates(comic)

	// 方式3：只调Updates()方法，传入 map[string]interface{} -》 推荐！！
	// 更新参数
	updateDataMap := map[string]any{
		"country_id":       comic.CountryId,
		"website_id":       comic.WebsiteId,
		"sex_type_id":      comic.SexTypeId,
		"type_id":          comic.TypeId,
		"update":           comic.Update,
		"hits":             comic.Hits,
		"comic_url":        comic.ComicUrl,
		"cover_url":        comic.CoverUrl,
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
	result := DB.Model(&comic).Where("id = ?", comicId).Updates(updateDataMap)
	if result.Error != nil {
		log.Error("修改失败: ", result.Error)
		return result.Error
	} else {
		log.Info("修改成功: ", comicId)
	}

	// 3. 返回错误信息
	return nil
}

// 批量改
func ComicsBatchUpdate(updates map[uint]map[string]interface{}) {
	for comicId, update := range updates {
		var comic models.Comic
		result := DB.Model(&comic).Where("id = ?", comicId).Updates(update)
		if result.Error != nil {
			log.Errorf("更新漫画 %d 失败: %v", comicId, result.Error)
		} else {
			log.Debugf("更新漫画 %d 成功", comicId)
		}
	}
}

// 查
func ComicQueryById(id uint) *models.Comic {
	var comic models.Comic
	result := DB.First(&comic, id)
	if result.Error != nil {
		log.Error("查询失败: ", result.Error)
		return nil
	}
	log.Info("查询成功: ", comic)
	return &comic
}

// 批量查
func ComicsBatchQuery(ids []uint) ([]*models.Comic, error) {
	var comics []*models.Comic
	result := DB.Find(&comics, ids)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return comics, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(comics))
	return comics, nil
}

// 查所有
func ComicsQueryAll() ([]*models.Comic, error) {
	var comics []*models.Comic
	result := DB.Find(&comics)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return comics, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(comics))
	return comics, nil
}

// 查数据总数
func ComicsTotal() (int64, error) {
	var count int64
	result := DB.Model(&models.Comic{}).Count(&count)
	if result.Error != nil {
		log.Error("查询数据总数失败: ", result.Error)
		return 0, result.Error
	}
	log.Infof("查询数据总数成功, 总数为 %d", count)
	return count, nil
}

// 分页查询
func ComicsPageQuery(pageNum, pageSize int) ([]*models.Comic, error) {
	var comics []*models.Comic
	result := DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comics)
	if result.Error != nil {
		log.Error("分页查询失败: ", result.Error)
		return comics, result.Error
	}
	log.Infof("分页查询成功, 查询到 %d 条记录", len(comics))
	return comics, result.Error
}
