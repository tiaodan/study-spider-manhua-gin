// db website 相关操作
package db

import (
	"strings"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"
	"study-spider-manhua-gin/util/stringutil"

	// 导入 clause 包
	"gorm.io/gorm/clause"
)

// 增
func WebsiteAdd(website *models.Website) error {
	// 预处理空格
	website.Name = strings.TrimSpace(website.Name)
	website.Url = strings.TrimSpace(website.Url)

	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "NameId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":       website.Name,
			"url":        website.Url,
			"need_proxy": website.NeedProxy,
			"is_https":   website.IsHttps,
		}),
	}).Create(website)
	if result.Error != nil {
		log.Error("创建失败:", result.Error)
		return result.Error
	} else {
		log.Info("创建成功:", website)
	}
	return nil
}

// 批量增
func WebsiteBatchAdd(websites []*models.Website) {
	for i, website := range websites {
		err := WebsiteAdd(website)
		if err == nil {
			log.Debugf("批量创建第%d条成功, website: %v", i+1, website.Name)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删-通过id
func WebsiteDeleteById(id uint) {
	var website models.Website
	result := DB.Delete(&website, id)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", id)
	}
}

// 删-通过 nameId
func WebsiteDeleteByNameId(nameId any) {
	var website models.Website
	result := DB.Where("name_id = ?", nameId).Delete(&website)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", nameId)
	}
}

// 删-通过其它
func WebsiteDeleteByOther(condition string, other any) {
	var website models.Website
	result := DB.Where(condition+" = ?", other).Delete(&website)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", other)
	}
}

// 批量删-通过id
func WebsitesBatchDeleteById(ids []uint) {
	var websites []models.Website
	result := DB.Delete(&websites, ids)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", ids)
	}
}

// 批量删-通过nameIds
func WebsitesBatchDeleteByNameId(nameIds []int) {
	var websites []models.Website
	result := DB.Where("name_id IN ?", nameIds).Delete(&websites)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", nameIds)
	}
}

// 批量删-通过other
func WebsitesBatchDeleteByOther(condition string, others []any) {
	var websites []models.Website
	// result := DB.Where(condition+" IN ?", others).Delete(&websites) // other这样写错？
	result := DB.Where(condition+" IN ?", others).Delete(&websites)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", others)
	}
}

// 改 - by Id
func WebsiteUpdateById(id uint, updates map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(updates)

	var website models.Website
	// 解决0值不更新
	result := DB.Model(&website).Where("id = ?", id).Select("name", "url", "need_proxy", "is_https").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", id)
	}
}

// 改 - by nameId
func WebsiteUpdateByNameId(nameId int, updates map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(updates)

	var website models.Website
	// 解决0值不更新
	result := DB.Model(&website).Where("name_id = ?", nameId).Select("name", "url", "need_proxy", "is_https").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", nameId)
	}
}

// 改 - by other
func WebsiteUpdateByOther(condition string, other any, updates map[string]interface{}) {
	// 预处理：去除字符串字段的首尾空格
	stringutil.TrimSpaceMap(updates)

	var website models.Website
	// 解决0值不更新
	// result := DB.Model(&website).Where("name_id = ?", nameId).Select("name", "url", "need_proxy", "is_https").Updates(updates)  // 之前写法
	result := DB.Model(&website).Where(condition+" = ?", other).Select("name", "url", "need_proxy", "is_https").Updates(updates) // 之前写法
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", other)
	}
}

// 改 - 批量 by id
func WebsitesBatchUpdateById(updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where("id = ?", update["Id"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by nameId
func WebsitesBatchUpdateByNameId(updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where("name_id = ?", update["NameId"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by other
func WebsitesBatchUpdateByOther(updates []map[string]interface{}) {
	for _, update := range updates {
		// 预处理：去除字符串字段的首尾空格
		stringutil.TrimSpaceMap(update)

		var website models.Website
		// 解决0值不更新
		result := DB.Model(&website).Where("name_id = ?", update["NameId"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 查 - by id
func WebsiteQueryById(id uint) *models.Website {
	var website models.Website
	result := DB.First(&website, id)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 查 - by nameId
func WebsiteQueryByNameId(nameId int) *models.Website {
	var website models.Website
	result := DB.Where("name_id = ?", nameId).First(&website)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 查 - by other
func WebsiteQueryByOther(condition string, other any) *models.Website {
	var website models.Website
	result := DB.Where(condition+" = ?", other).First(&website)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", website)
	return &website
}

// 批量查 - by ids
func WebsitesBatchQueryById(ids []uint) ([]*models.Website, error) {
	var websites []*models.Website
	result := DB.Find(&websites, ids)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// 批量查 - by nameIds
func WebsitesBatchQueryByNameId(nameIds []int) ([]*models.Website, error) {
	var websites []*models.Website
	result := DB.Where("name_id IN ?", nameIds).Order("name_id").Find(&websites) // 默认升序
	// result := DB.Where("name_id IN ?", nameIds).Order("name_id DESC")Find(&websites) // 倒序排列
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}

// 批量查 - by others
// 参数：orderby 排序字符串 如: name_id   sort 排序方式，ASC 为正序，DESC 为倒序
func WebsitesBatchQueryByOther(condition string, others []any, orderby string, sort string) ([]*models.Website, error) {
	var websites []*models.Website
	// result := DB.Where(condition+" IN ?", others).Find(&websites)  // other这样写错？
	// result := DB.Where(condition+" IN ?", others).Order("name_id DESC").Find(&websites)
	result := DB.Where(condition+" IN ?", others).Order(orderby + " " + sort).Find(&websites)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return websites, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(websites)) // 原 log.Debug无需修改
	return websites, nil
}
