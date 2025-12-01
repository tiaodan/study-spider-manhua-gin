// db class 相关操作
package db

import (
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	// 导入 clause 包
	"gorm.io/gorm/clause"
)

// 增
func PornTypeAdd(categoryData *models.PornType) error {
	result := DBComic.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "NameId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"name": categoryData.Name}),
	}).Create(categoryData)
	if result.Error != nil {
		log.Error("创建失败:", result.Error)
		return result.Error
	} else {
		log.Info("创建成功:", categoryData, "1")
	}
	return nil
}

// 批量增
func PornTypeBatchAdd(categories []*models.PornType) {
	for i, categoryData := range categories {
		err := PornTypeAdd(categoryData)
		if err == nil {
			log.Debugf("批量创建第%d条成功, category: %v", i+1, categoryData.Name)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删
func PornTypeDelete(id uint) {
	var categoryData models.PornType
	result := DBComic.Delete(&categoryData, id)
	if result.Error != nil {
		log.Error("删除失败:", result.Error)
	} else {
		log.Info("删除成功:", id)
	}
}

// 批量删
func PornTypeBatchDelete(ids []uint) {
	var categories []models.PornType
	result := DBComic.Delete(&categories, ids)
	if result.Error != nil {
		log.Error("批量删除失败:", result.Error)
	} else {
		log.Debug("批量删除成功:", ids)
	}
}

// 改
func PornTypeUpdate(nameId uint, updates map[string]interface{}) {
	var categoryData models.PornType
	// 解决0值不更新
	result := DBComic.Model(&categoryData).Where("name_id = ?", nameId).Select("name").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败:", result.Error)
	} else {
		log.Info("修改成功:", nameId)
	}
}

// 批量改
func PornTypeBatchUpdate(updates map[uint]map[string]interface{}) {
	for nameId, update := range updates {
		var categoryData models.PornType
		// 解决0值不更新
		result := DBComic.Model(&categoryData).Where("name_id = ?", nameId).Select("name").Updates(update)
		if result.Error != nil {
			log.Errorf("更新类型 %d 失败: %v", nameId, result.Error)
		} else {
			log.Debugf("更新类型 %d 成功", nameId)
		}
	}
}

// 查
func PornTypeQueryById(id uint) *models.PornType {
	var categoryData models.PornType
	result := DBComic.First(&categoryData, id)
	if result.Error != nil {
		log.Error("查询失败:", result.Error)
		return nil
	}
	log.Info("查询成功:", categoryData)
	return &categoryData
}

// 批量查
func PornTypeBatchQuery(ids []uint) ([]*models.PornType, error) {
	var categories []*models.PornType
	result := DBComic.Find(&categories, ids)
	if result.Error != nil {
		log.Errorf("批量查询失败: %v", result.Error)
		return categories, result.Error
	}
	log.Infof("批量查询成功, 查询到 %d 条记录", len(categories))
	return categories, nil
}
