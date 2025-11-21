// db country 相关操作
package db

import (
	"study-spider-manhua-gin/src/models"

	"study-spider-manhua-gin/src/log"
	// 导入 clause 包
	"gorm.io/gorm/clause"
)

// 初始化
// 为 Website 表实现完整的操作接口
type CountryOperations struct{}

// 增
func (c CountryOperations) Add(country *models.Country) error {
	result := DBComic.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "NameId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"name": country.Name}),
	}).Create(country)
	if result.Error != nil {
		log.Error("创建失败:", result.Error)
		return result.Error
	} else {
		log.Info("创建成功:", country)
	}
	return nil
}

// 批量增
func (c CountryOperations) BatchAdd(countries []*models.Country) {
	for i, country := range countries {
		err := c.Add(country)
		if err == nil {
			log.Debugf("批量创建第%d条成功, country: %v", i+1, country.Name)
		} else {
			log.Errorf("批量创建第%d条失败, err: %v", i+1, err)
		}
	}
}

// 删
func (c CountryOperations) Delete(id uint) {
	var country models.Country
	result := DBComic.Delete(&country, id)
	if result.Error != nil {
		log.Error("删除失败: ", result.Error)
	} else {
		log.Info("删除成功: ", id)
	}
}

// 批量删
func (c CountryOperations) BatchDelete(ids []uint) {
	var countries []models.Country
	result := DBComic.Delete(&countries, ids)
	if result.Error != nil {
		log.Error("批量删除失败: ", result.Error)
	} else {
		log.Debug("批量删除成功: ", ids)
	}
}

// 改
func (c CountryOperations) Update(nameId uint, updates map[string]interface{}) {
	var country models.Country
	// 解决0值不更新
	result := DBComic.Model(&country).Where("name_id = ?", nameId).Select("name").Updates(updates)
	if result.Error != nil {
		log.Error("修改失败: ", result.Error)
	} else {
		log.Info("修改成功: ", nameId)
	}
}

// 批量改
func (c CountryOperations) BatchUpdate(updates map[uint]map[string]interface{}) {
	for nameId, update := range updates {
		var country models.Country
		// 解决0值不更新
		result := DBComic.Model(&country).Where("name_id = ?", nameId).Select("name").Updates(update)
		if result.Error != nil {
			log.Errorf("更新国家 %d 失败: %v", nameId, result.Error)
		} else {
			log.Debugf("更新国家 %d 成功", nameId)
		}
	}
}

// 查
func (c CountryOperations) QueryById(id uint) *models.Country {
	var country models.Country
	result := DBComic.First(&country, id)
	if result.Error != nil {
		log.Error("查询失败: ", result.Error)
		return nil
	}
	log.Info("查询成功: ", country)
	return &country
}

// 批量查
func (c CountryOperations) BatchQuery(ids []uint) ([]*models.Country, error) {
	var countries []*models.Country
	result := DBComic.Find(&countries, ids)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return countries, result.Error
	}
	log.Infof("批量查询成功, 查询到 %d 条记录", len(countries))
	return countries, nil
}

// 批量查 - by NameIds
func (c CountryOperations) BatchQueryByNameId(nameIds []int) ([]*models.Country, error) {
	var countries []*models.Country
	result := DBComic.Where("name_id IN ?", nameIds).Find(&countries)
	if result.Error != nil {
		log.Error("批量查询失败: ", result.Error)
		return countries, result.Error
	}
	log.Infof("批量查询成功, 查询到 %d 条记录", len(countries))
	return countries, nil
}
