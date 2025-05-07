// db comic 相关操作
package db

import (
	// 导入 clause 包
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"

	"gorm.io/gorm/clause"
)

// 增
func ComicAdd(comic *models.Comic) error {
	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "Name"}}, // 判断唯一索引: Name
		DoUpdates: clause.Assignments(map[string]interface{}{
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
		err := ComicAdd(comic)
		if err == nil {
			log.Debugf("批量创建第%d条成功, comic: %v", i+1, &comic)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删
func ComicDelete(id uint) error {
	log.Debug("删除漫画, 参数id= ", id)
	var comic models.Comic
	result := DB.Delete(&comic, id)
	if result.Error != nil {
		log.Error("删除失败: ", result.Error)
		return result.Error
	} else {
		log.Info("删除成功: ", id)
	}
	return nil
}

// 批量删
func ComicsBatchDelete(ids []uint) {
	var comics []models.Comic
	result := DB.Delete(&comics, ids)
	if result.Error != nil {
		log.Error("批量删除失败: ", result.Error)
	} else {
		log.Debug("批量删除成功: ", ids)
	}
}

// 改 - 参数用结构体, 0值不更新
// id 可以int,可以string。go默认定义的 any = interface{}
func ComicUpdate(comicId any, comic *models.Comic) error {
	// 解决0值不更新 -> 指定更新字段
	result := DB.Model(&comic).Where("id = ?", comicId).Select("update", "hits", "comic_url",
		"cover_url", "brief_short", "brief_long", "end", "need_tcp", "cover_need_tcp",
		"spider_end", "download_end", "upload_aws_end", "upload_baidu_end").Updates(comic)
	if result.Error != nil {
		log.Error("修改失败: ", result.Error)
		return result.Error
	} else {
		log.Info("修改成功: ", comicId)
	}

	return nil
}

// 改 - 根据id, 排除唯一索引 参数用结构体
// id 可以int,可以string。go默认定义的 any = interface{}
func ComicUpdateByIdOmitIndex(comicId any, comic *models.Comic) error {
	// result := DB.Model(&comic).Where("id = ?", comicId).Omit("name").Updates(comic)
	result := DB.Model(&comic).Where("id = ?", comicId).Select("update", "hits", "comic_url",
		"cover_url", "brief_short", "brief_long", "end", "need_tcp", "cover_need_tcp",
		"spider_end", "download_end", "upload_aws_end", "upload_baidu_end").Updates(comic)
	if result.Error != nil {
		log.Error("修改失败: ", result.Error)
		return result.Error
	} else {
		log.Info("修改成功: ", comicId)
	}

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
