package db

import (
	"fmt"
	"study-spider-manhua-gin/log"
	"study-spider-manhua-gin/models"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once // 使用 sync.Once 确保单例

// 初始化数据库连接
/*
参数：
	dbType string 数据库类型 如 mysql、sqlite3、postgres 等
	dbName string 数据库名
	dbUser string 数据库用户名
	dbPass string 数据库密码
*/
func InitDB(dbType, dbName, dbUser, dbPass string) {
	once.Do(func() { // 使用 sync.Once 确保只执行一次
		// dsn := "root:password@tcp(127.0.0.1:3306)/pdd_order?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		var err error

		var dbOpen gorm.Dialector // 用什么数据库打开
		if dbType == "mysql" {
			dbOpen = mysql.Open(dsn)
		}
		// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		DB, err = gorm.Open(dbOpen, &gorm.Config{})
		if err != nil {
			log.Error("单例: 数据库连接失败, 是不是数据库名+密码没配对？ 数据库没创建？ err= ", err)
			panic(err)
		}
		log.Debug("单例: 数据库连接成功")
	})
}

// 插入默认数据
/*
思路:
1. 插入默认数据-1
2. 插入默认数据-2
*/
func InsertDefaultData() {
	// 插入默认数据-website
	// 插入默认数据-website
	websiteDefaultNoClass := &models.Website{Name: "待分类", NameId: 0, Url: "未知", NeedProxy: 0, IsHttps: 0}
	websiteDefaultJ88d := &models.Website{Name: "j88d", NameId: 1, Url: "http://www.j88d.com", NeedProxy: 0, IsHttps: 0} // 请求url 时带上http://
	defaultWebsites := []*models.Website{websiteDefaultNoClass, websiteDefaultJ88d}
	WebsiteBatchAdd(defaultWebsites)

	// 插入默认数据-category 类别
	classDefaultNoCategory := &models.Category{Name: "待分类", NameId: 0}
	classDefaultCartoonNormal := &models.Category{Name: "普通漫画", NameId: 1}
	classDefaultCartoonSex := &models.Category{Name: "色漫", NameId: 2}
	classes := []*models.Category{classDefaultNoCategory, classDefaultCartoonNormal, classDefaultCartoonSex}
	CategoriesBatchAdd(classes)

	// 插入默认数据-country
	countryDefaultNoType := &models.Country{Name: "待分类", NameId: 0}
	countryDefaultChina := &models.Country{Name: "中国", NameId: 1}
	countryDefaultKoren := &models.Country{Name: "韩国", NameId: 2}
	countryDefaultAmerica := &models.Country{Name: "欧美", NameId: 3}
	countryDefaultJapan := &models.Country{Name: "日本", NameId: 4}
	countries := []*models.Country{countryDefaultNoType, countryDefaultChina, countryDefaultKoren, countryDefaultAmerica, countryDefaultJapan}
	CountriesBatchAdd(countries)

	// 插入默认数据-type
	// 一级分类
	typeDefaultNoTypeLevel1 := &models.Type{NameId: 0, Name: "待分类", Level: 1}
	typeDefaultKoren := &models.Type{NameId: 1, Name: "韩漫", Level: 1}
	typeDefaultJapan := &models.Type{NameId: 2, Name: "日漫", Level: 1}
	typeDefaultRealPerson := &models.Type{NameId: 3, Name: "真人漫画", Level: 1}
	typeDefault3D := &models.Type{NameId: 4, Name: "3D漫画", Level: 1}
	typeDefaultAmeraica := &models.Type{NameId: 5, Name: "欧美漫画", Level: 1}
	typeDefaultSameSex := &models.Type{NameId: 6, Name: "同性", Level: 1}

	defaultTypes := []*models.Type{
		// 一级分类
		typeDefaultNoTypeLevel1, typeDefaultKoren, typeDefaultJapan,
		typeDefaultRealPerson, typeDefault3D, typeDefaultAmeraica,
		typeDefaultSameSex,
	}
	TypesBatchAdd(defaultTypes)

}

// TruncateTable 清空指定模型对应的数据表，同时跳过外键检查（适用于 MySQL）
func TruncateTable(db *gorm.DB, model interface{}) error {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return err
	}
	tableName := stmt.Schema.Table

	// 注释部分执行报错，gorm不让一次执行多个语句
	// sql := fmt.Sprintf(`
	// 	SET FOREIGN_KEY_CHECKS = 0;
	// 	TRUNCATE TABLE %s;
	// 	SET FOREIGN_KEY_CHECKS = 1;
	// `, tableName)

	// return db.Exec(sql).Error

	// 分开执行每个 SQL 语句
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1;") // 确保最终恢复外键检查

	return db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", tableName)).Error
}

// 删除表所有数据，不用truncate, 生产环境避免用truncate
func DeleteTableAllData(db *gorm.DB, model interface{}) error {
	result := db.Where("1 = 1").Delete(&model)
	if result.Error != nil {
		log.Error("删除表所有数据失败:", result.Error)
		return result.Error
	}
	return nil
}
