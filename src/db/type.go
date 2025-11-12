// db comicType 相关操作
package db

import (
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	// 导入 clause 包
	"gorm.io/gorm/clause"
)

// 增
func TypeAdd(comicType *models.Type) error {
	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "NameId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name": comicType.Name, "level": comicType.Level,
			"parent": comicType.Parent,
		}),
	}).Create(comicType)
	if result.Error != nil {
		log.Error("创建失败:", result.Error)
		return result.Error
	} else {
		log.Info("创建成功:", comicType)
	}
	return nil
}

// 批量增
func TypeBatchAdd(comicTypes []*models.Type) {
	for i, comicType := range comicTypes {
		err := TypeAdd(comicType)
		if err == nil {
			log.Debugf("批量创建第%d条成功, comicType: %v", i+1, comicType.Name)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删-通过id
func TypeDeleteById(id uint) {
	var comicType models.Type
	result := DB.Delete(&comicType, id)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", id)
	}
}

// 删-通过 nameId
func TypeDeleteByNameId(nameId any) {
	var comicType models.Type
	result := DB.Where("name_id = ?", nameId).Delete(&comicType)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", nameId)
	}
}

// 删-通过其它
func TypeDeleteByOther(condition string, other any) {
	var comicType models.Type
	result := DB.Where(condition+" = ?", other).Delete(&comicType)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", other)
	}
}

// 批量删-通过id
func TypesBatchDeleteById(ids []uint) {
	var comicTypes []models.Type
	result := DB.Delete(&comicTypes, ids)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", ids)
	}
}

// 批量删-通过nameIds
func TypesBatchDeleteByNameId(nameIds []int) {
	var comicTypes []models.Type
	result := DB.Where("name_id IN ?", nameIds).Delete(&comicTypes)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", nameIds)
	}
}

// 批量删-通过other
func TypesBatchDeleteByOther(condition string, others []any) {
	var comicTypes []models.Type
	// result := DB.Where(condition+" IN ?", others).Delete(&comicTypes) // other这样写错？
	result := DB.Where(condition+" IN ?", others).Delete(&comicTypes)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", others)
	}
}

// 改 - by Id
func TypeUpdateById(id uint, updates map[string]interface{}) {
	var comicType models.Type
	// 解决0值不更新
	result := DB.Model(&comicType).Where("id = ?", id).Select("name", "level", "parent").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", id)
	}
}

// 改 - by nameId
func TypeUpdateByNameId(nameId int, updates map[string]interface{}) {
	var comicType models.Type
	// 解决0值不更新
	result := DB.Model(&comicType).Where("name_id = ?", nameId).Select("name", "level", "parent").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", nameId)
	}
}

// 改 - by other
func TypeUpdateByOther(condition string, other any, updates map[string]interface{}) {
	var comicType models.Type
	// 解决0值不更新
	// result := DB.Model(&comicType).Where("name_id = ?", nameId).Select("name", "url", "need_proxy", "is_https").Updates(updates)  // 之前写法
	result := DB.Model(&comicType).Where(condition+" = ?", other).Select("name", "level", "parent").Updates(updates) // 之前写法
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", other)
	}
}

// 改 - 批量 by id
func TypesBatchUpdateById(updates []map[string]interface{}) {
	for _, update := range updates {
		var comicType models.Type
		// 解决0值不更新
		result := DB.Model(&comicType).Where("id = ?", update["Id"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by nameId
func TypesBatchUpdateByNameId(updates []map[string]interface{}) {
	for _, update := range updates {
		var comicType models.Type
		// 解决0值不更新
		result := DB.Model(&comicType).Where("name_id = ?", update["NameId"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 改 - 批量 by other
func TypesBatchUpdateByOther(updates []map[string]interface{}) {
	for _, update := range updates {
		var comicType models.Type
		// 解决0值不更新
		result := DB.Model(&comicType).Where("name_id = ?", update["NameId"]).Select("name", "url", "need_proxy", "is_https").Updates(update)
		if result.Error != nil {
			log.Errorf("更新网站 %d 失败: %v", update["Id"], result.Error)
		} else {
			log.Debugf("更新网站 %d 成功", update["Id"])
		}
	}
}

// 查 - by id
func TypeQueryById(id uint) *models.Type {
	var comicType models.Type
	result := DB.First(&comicType, id)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", comicType)
	return &comicType
}

// 查 - by nameId
func TypeQueryByNameId(nameId int) *models.Type {
	var comicType models.Type
	result := DB.Where("name_id = ?", nameId).First(&comicType)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", comicType)
	return &comicType
}

// 查 - by other
func TypeQueryByOther(condition string, other any) *models.Type {
	var comicType models.Type
	result := DB.Where(condition+" = ?", other).First(&comicType)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", comicType)
	return &comicType
}

// 批量查 - by ids
func TypesBatchQueryById(ids []uint) ([]*models.Type, error) {
	var comicTypes []*models.Type
	result := DB.Find(&comicTypes, ids)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return comicTypes, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(comicTypes)) // 原 log.Debug无需修改
	return comicTypes, nil
}

// 批量查 - by nameIds
func TypesBatchQueryByNameId(nameIds []int) ([]*models.Type, error) {
	var comicTypes []*models.Type
	result := DB.Where("name_id IN ?", nameIds).Order("name_id").Find(&comicTypes) // 默认升序
	// result := DB.Where("name_id IN ?", nameIds).Order("name_id DESC")Find(&comicTypes) // 倒序排列
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return comicTypes, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(comicTypes)) // 原 log.Debug无需修改
	return comicTypes, nil
}

// 批量查 - by others
// 参数：orderby 排序字符串 如: name_id   sort 排序方式，ASC 为正序，DESC 为倒序
func TypesBatchQueryByOther(condition string, others []any, orderby string, sort string) ([]*models.Type, error) {
	var comicTypes []*models.Type
	// result := DB.Where(condition+" IN ?", others).Find(&comicTypes)  // other这样写错？
	// result := DB.Where(condition+" IN ?", others).Order("name_id DESC").Find(&comicTypes)
	result := DB.Where(condition+" IN ?", others).Order(orderby + " " + sort).Find(&comicTypes)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return comicTypes, result.Error
	}
	log.Debugf("批量查询成功, 查询到 %d 条记录", len(comicTypes)) // 原 log.Debug无需修改
	return comicTypes, nil
}
