package db

import (
	"fmt"
	"reflect"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ---------------------------- 变量 start ----------------------------
var DBComic *gorm.DB // comic 数据库对象
var once sync.Once   // 使用 sync.Once 确保单例

var dbNameComic = "comic"                 // 数据库名-漫画，用于日志打印
var tableNameWebsiteType = "website_type" // 数据库表名-网站类型，用于日志打印
var tableNameWebsite = "website"          // 数据库表名-网站，用于日志打印
var tableNamePornType = "porntype"        // 数据库表名-色情类型，用于日志打印
var tableNameCountry = "country"          // 数据库表名-国家，用于日志打印
var tableNameType = "type"                // 数据库表名-类型，用于日志打印
var tableNameProcess = "process"          // 数据库表名-进度，用于日志打印
var tableNameAuthor = "author"            // 数据库表名-作者，用于日志打印

// 定义统一的操作接口,方便单元测试的时候调用. 为了把所有表的增删改查都叫Add
// 定义 model 约束
type Model interface {
	*models.Website | *models.Country | *models.PornType | *models.Type | *models.ComicSpider
	// 或者定义通用方法
	// GetID() uint
	// GetNameID() int
}

// type TableOperations[T any] interface { // 定义泛型接口  不能用，方法实现时不兼容
type TableOperations[T Model] interface { // 定义泛型接口  也能用
	Add(modelPointer T) error // 原来的写法 Add(modelPointer interface{})  泛型写法：Add(modelPointer T)
	DeleteById(id uint)
	DeleteByNameId(nameid any)
	DeleteByOther(condition string, other any)
	UpdateById(id uint, updates map[string]interface{})
	UpdateByNameId(nameId int, updates map[string]interface{})
	UpdateByOther(condition string, other any, updates map[string]interface{})
	QueryById(id uint) T
	QueryByNameId(nameId int) T
	QueryByOther(condition string, other any) T

	BatchAdd(modelPointers []T) // 接收特定类型的切片 泛型写法：BatchAdd(modelPointers []T)
	BatchDeleteById(ids []uint)
	BatchDeleteByNameId(nameIds []int)
	BatchDeleteByOther(condition string, others []any)
	BatchUpdateById(updates []map[string]interface{})
	BatchUpdateByNameId(updates []map[string]interface{})
	BatchUpdateByOther(updates []map[string]interface{})

	BatchQueryById(ids []uint) ([]T, error)
	BatchQueryByNameId(nameIds []int) ([]T, error)
	BatchQueryByOther(condition string, others []any, orderby string, sort string) ([]T, error)
}

// 实例化接口操作对象
var WebsiteOps WebsiteOperations // 另一种写法: var WebsiteOps = WebsiteOperations{}
var CountryOps CountryOperations // 另一种写法: var WebsiteOps = WebsiteOperations{}

// ---------------------------- 变量 end ----------------------------

// ---------------------------- 初始化 start ----------------------------
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
		DBComic, err = gorm.Open(dbOpen, &gorm.Config{})
		if err != nil {
			log.Error("单例: 数据库连接失败, 是不是数据库名+密码没配对？ 数据库没创建？ err= ", err)
			panic(err)
		}
		log.Debug("单例: 数据库连接成功")
	})
}

// ---------------------------- 初始化 end ----------------------------

// ---------------------------- 函数 start ----------------------------
// 插入默认数据
/*
思路:
	1. 准备插入默认数据
	2. 插入数据
		- 通过执行事务，保证数据的一致性，都插入或者都回滚 DBConnObj.Transaction(func(tx *gorm.DB) error {
		- 循环调用 BDInsertDefaultData()方法，而不是在 BDInsertDefaultData() 方法里循环执行参数，这样最简单，也不会用到泛型 + 类型
*/
func InsertDefaultData() {

	// v0.3 方式：使用通用 增删改查方法。但是实现方式：用循环实现
	// 1. 准备插入数据
	// 准备默认数据 - website_type ,必须在website 之前插入，否则报错 --
	websiteTypeDefaultNoClass := &models.WebsiteType{Id: 1, Name: "待分类"}
	websiteTypeDefaultComic := &models.WebsiteType{Id: 2, Name: "漫画"}
	websiteTypeDefaultNovel := &models.WebsiteType{Id: 3, Name: "小说"}
	websiteTypeDefaultAudiobook := &models.WebsiteType{Id: 4, Name: "有声书"}
	websiteTypeDefaultVideo := &models.WebsiteType{Id: 5, Name: "视频"}
	websiteTypeDefaultMusic := &models.WebsiteType{Id: 6, Name: "音乐"}
	websiteTypeDefaultCloudDisk := &models.WebsiteType{Id: 7, Name: "网盘"}
	websiteTypeDefaultMuitiEntertainment := &models.WebsiteType{Id: 8, Name: "综合娱乐"} // 多种娱乐. Entertainment -》 娱乐 英文
	defaultDataWebsiteTypeArr := []*models.WebsiteType{websiteTypeDefaultNoClass, websiteTypeDefaultComic,
		websiteTypeDefaultNovel, websiteTypeDefaultAudiobook, websiteTypeDefaultVideo, websiteTypeDefaultMusic,
		websiteTypeDefaultCloudDisk, websiteTypeDefaultMuitiEntertainment}
	websiteTypeUniqueIndexArr := []string{"Name"}
	WebsiteTypeUpdateDBColumnRealNameArr := []string{"name"}

	// 准备默认数据- website --
	websiteDefaultNoClass := &models.Website{Name: "待分类", Id: 1, Domain: "未知", NeedProxy: false, IsHttps: false,
		CoverURLIsNeedHttps: true, ChapterContentURLIsNeedHttps: true,
		CoverURLConcatRule:          "{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}",
		ChapterContentURLConcatRule: "{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}",
		CoverDomain:                 "www.未知.com", ChapterContentDomain: "www.未知.com",
		IsRefer: false, WebsiteTypeId: 1}
	websiteDefaultJ88d := &models.Website{Name: "j88d", Id: 2, Domain: "www.j88d.com", NeedProxy: false, IsHttps: false,
		CoverURLIsNeedHttps: false, ChapterContentURLIsNeedHttps: false,
		CoverURLConcatRule:          "{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}",
		ChapterContentURLConcatRule: "{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}",
		CoverDomain:                 "www.j88d.com", ChapterContentDomain: "www.j88d.com",
		IsRefer: true, WebsiteTypeId: 8}
	websiteDefaultAwsS3 := &models.Website{Name: "aws-s3", Id: 3, Domain: "ap-northeast-2.console.aws.amazon.com/s3/home?region=ap-northeast-2", NeedProxy: false,
		IsHttps: true, CoverURLIsNeedHttps: false, ChapterContentURLIsNeedHttps: false,
		CoverURLConcatRule:          "{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}",
		ChapterContentURLConcatRule: "{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}",
		CoverDomain:                 "www.awsS3.com", ChapterContentDomain: "www.awsS3.com",
		IsRefer: true, WebsiteTypeId: 7}
	websiteDefaultYuliu := &models.Website{Name: "预留", Id: 4, Domain: "www.yuliu.com", NeedProxy: false, IsHttps: false,
		CoverURLIsNeedHttps: false, ChapterContentURLIsNeedHttps: false,
		CoverURLConcatRule:          "{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}",
		ChapterContentURLConcatRule: "{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}",
		CoverDomain:                 "www.预留.com", ChapterContentDomain: "www.预留.com",
		IsRefer: false, WebsiteTypeId: 1} // 预留
	defaultDataWebsiteArr := []*models.Website{websiteDefaultNoClass, websiteDefaultJ88d, websiteDefaultAwsS3, websiteDefaultYuliu} // 要插入数据
	websiteUniqueIndexArr := []string{"Name", "Domain"}                                                                             // 唯一索引
	websiteUpdateDBColumnRealNameArr := []string{"need_proxy", "Is_https", "is_refer",
		"cover_url_is_need_https", "chapter_content_url_is_need_https",
		"cover_url_concat_rule", "chapter_content_url_concat_rule",
		"cover_domain", "chapter_content_domain"} // 要更新的字段

	// 准备默认数据- pornType 色情类型 --
	pornTypeDefaultNoCategory := &models.PornType{Name: "待分类", Id: 1}
	pornTypeDefaultCartoonNormal := &models.PornType{Name: "普通漫画", Id: 2}
	pornTypeDefaultCartoonSex := &models.PornType{Name: "色漫", Id: 3}
	defaultDataPornTypeArr := []*models.PornType{pornTypeDefaultNoCategory, pornTypeDefaultCartoonNormal, pornTypeDefaultCartoonSex}
	pornTypeUniqueIndexArr := []string{"Name"}            // 唯一索引
	pornTypeUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

	// 准备默认数据- country --
	countryDefaultNoType := &models.Country{Name: "待分类", Id: 1}
	countryDefaultChina := &models.Country{Name: "中国", Id: 2}
	countryDefaultKoren := &models.Country{Name: "韩国", Id: 3}
	countryDefaultAmerica := &models.Country{Name: "欧美", Id: 4}
	countryDefaultJapan := &models.Country{Name: "日本", Id: 5}
	defaultDataCountryArr := []*models.Country{countryDefaultNoType, countryDefaultChina, countryDefaultKoren, countryDefaultAmerica, countryDefaultJapan}
	countryUniqueIndexArr := []string{"Name"}            // 唯一索引
	countryUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

	// 准备默认数据-type --
	// 一级分类
	typeDefaultNoTypeLevel1 := &models.Type{Id: 1, Name: "待分类", Level: 1}
	typeDefaultKoren := &models.Type{Id: 2, Name: "韩漫", Level: 1}
	typeDefaultJapan := &models.Type{Id: 3, Name: "日漫", Level: 1}
	typeDefaultRealPerson := &models.Type{Id: 4, Name: "真人漫画", Level: 1}
	typeDefault3D := &models.Type{Id: 5, Name: "3D漫画", Level: 1}
	typeDefaultAmeraica := &models.Type{Id: 6, Name: "欧美漫画", Level: 1}
	typeDefaultSameSex := &models.Type{Id: 7, Name: "同性", Level: 1}
	defaultDataTypeArr := []*models.Type{
		// 一级分类
		typeDefaultNoTypeLevel1, typeDefaultKoren, typeDefaultJapan,
		typeDefaultRealPerson, typeDefault3D, typeDefaultAmeraica,
		typeDefaultSameSex,
	}
	typeUniqueIndexArr := []string{"Name"}                       // 唯一索引
	typeUpdateDBColumnRealNameArr := []string{"level", "parent"} // 要更新的字段

	// 准备默认数据- process --
	processDefaultNoType := &models.Process{Id: 1, Name: "待分类"}
	processDefaultOngoing := &models.Process{Id: 2, Name: "连载"}
	processDefaultCompleted := &models.Process{Id: 3, Name: "完结"}
	defaultDataProcessArr := []*models.Process{processDefaultNoType, processDefaultOngoing, processDefaultCompleted}
	processUniqueIndexArr := []string{"Name"}            // 唯一索引
	processUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

	// 准备默认数据- author --
	authorDefaultNoName := &models.Author{Id: 1, Name: "佚名"}
	defaultDataAuthorArr := []*models.Author{authorDefaultNoName}
	authorUniqueIndexArr := []string{"Name"}            // 唯一索引
	authorUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

	dataObjArr := []any{defaultDataWebsiteTypeArr, defaultDataWebsiteArr, defaultDataPornTypeArr, defaultDataCountryArr,
		defaultDataTypeArr, defaultDataProcessArr, defaultDataAuthorArr} // 插入对象 数组 . 必须website_type 在最前面，否则报错
	indexArr := [][]string{websiteTypeUniqueIndexArr, websiteUniqueIndexArr, pornTypeUniqueIndexArr, countryUniqueIndexArr,
		typeUniqueIndexArr, processUniqueIndexArr, authorUniqueIndexArr} // 唯一索引 数组
	dbColArr := [][]string{WebsiteTypeUpdateDBColumnRealNameArr, websiteUpdateDBColumnRealNameArr, pornTypeUpdateDBColumnRealNameArr,
		countryUpdateDBColumnRealNameArr, typeUpdateDBColumnRealNameArr, processUpdateDBColumnRealNameArr,
		authorUpdateDBColumnRealNameArr} // 要更新的字段 数组
	dbNameArr := []string{dbNameComic, dbNameComic, dbNameComic, dbNameComic, dbNameComic,
		dbNameComic, dbNameComic} // 数据库名称 数组，仅用于日志打印
	tableNameArr := []string{tableNameWebsiteType, tableNameWebsite, tableNamePornType, tableNameCountry,
		tableNameType, tableNameProcess, tableNameAuthor} // 表名称 数组，仅用于日志打印
	// 2. 插入数据
	// -- 校验参数个数是否一致

	// -- 用事务，保证数据的一致性，都插入或者都回滚
	err := DBComic.Transaction(func(tx *gorm.DB) error {
		// -- 循环调用批量插入方法
		for i, dataObj := range dataObjArr {
			// 执行插入操作 --
			err := DBUpsertBatch(tx, dataObj, indexArr[i], dbColArr[i])
			// 打印
			if err != nil {
				log.Errorf("插入默认数据%s-%s 失败, err= %s", dbNameArr[i], tableNameArr[i], err)
				return err
			} else {
				okNum := reflect.ValueOf(dataObj).Len() // 用反射获取到 any类型的长度. 插入成功几个
				log.Infof("插入默认数据%s-%s 成功个数: %v", dbNameArr[i], tableNameArr[i], okNum)
			}
		}
		return nil // 所有事务执行完毕，返回成功。返回给事务的
	})

	if err != nil {
		log.Error("插入默认数据失败, 全部回滚, err= ", err)
		errorutil.ErrorPanic(err, "插入默认数据失败, 全部回滚, err=")
	}

	// v0.2 方式：使用通用 增删改查方法。但是实现方式：每个表都要写一遍，重复代码多。下一步考虑用循环实现
	/*
		// -- 插入默认数据-website
		// 准备插入数据 --
		websiteDefaultNoClass := &models.Website{Name: "待分类", NameId: 0, Domain: "未知", NeedProxy: 0, IsHttps: 0}
		websiteDefaultJ88d := &models.Website{Name: "j88d", NameId: 1, Domain: "http://www.j88d.com", NeedProxy: 0, IsHttps: 0} // 请求domain 时带上http://
		defaultDataWebsiteArr := []*models.Website{websiteDefaultNoClass, websiteDefaultJ88d}                                // 要插入数据

		websiteUniqueIndexArr := []string{"NameId"}                                           // 唯一索引
		websiteUpdateDBColumnRealNameArr := []string{"name", "domain", "need_proxy", "Is_https"} // 要更新的字段

		// 执行插入操作 --
		err := DBUpsertBatch(defaultDataWebsiteArr, websiteUniqueIndexArr, websiteUpdateDBColumnRealNameArr)
		// 打印
		if err != nil {
			log.Errorf("插入默认数据%s-%s 失败, err= %s", dbNameComic, tableNameWebsite, err)
		} else {
			log.Infof("插入默认数据%s-%s 成功个数: %v", dbNameComic, tableNameWebsite, len(defaultDataWebsiteArr))
		}

		// 插入默认数据-category 类别
		// 准备插入数据 --
		pornTypeDefaultNoCategory := &models.Category{Name: "待分类", NameId: 0}
		pornTypeDefaultCartoonNormal := &models.Category{Name: "普通漫画", NameId: 1}
		pornTypeDefaultCartoonSex := &models.Category{Name: "色漫", NameId: 2}
		defaultDataPornTypeArr := []*models.Category{pornTypeDefaultNoCategory, pornTypeDefaultCartoonNormal, pornTypeDefaultCartoonSex}

		pornTypeUniqueIndexArr := []string{"NameId"}          // 唯一索引
		pornTypeUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

		// 执行插入操作 --
		err = DBUpsertBatch(defaultDataPornTypeArr, pornTypeUniqueIndexArr, pornTypeUpdateDBColumnRealNameArr)
		// 打印
		if err != nil {
			log.Errorf("插入默认数据%s-%s 失败, err= %s", dbNameComic, tableNamePornType, err)
		} else {
			log.Infof("插入默认数据%s-%s 成功个数: %v", dbNameComic, tableNamePornType, len(defaultDataPornTypeArr))
		}

		// 插入默认数据-country
		// 准备插入数据 --
		countryDefaultNoType := &models.Country{Name: "待分类", NameId: 0}
		countryDefaultChina := &models.Country{Name: "中国", NameId: 1}
		countryDefaultKoren := &models.Country{Name: "韩国", NameId: 2}
		countryDefaultAmerica := &models.Country{Name: "欧美", NameId: 3}
		countryDefaultJapan := &models.Country{Name: "日本", NameId: 4}
		defaultDataCountryArr := []*models.Country{countryDefaultNoType, countryDefaultChina, countryDefaultKoren, countryDefaultAmerica, countryDefaultJapan}

		countryUniqueIndexArr := []string{"NameId"}          // 唯一索引
		countryUpdateDBColumnRealNameArr := []string{"name"} // 要更新的字段

		// 执行插入操作 --
		err = DBUpsertBatch(defaultDataCountryArr, countryUniqueIndexArr, countryUpdateDBColumnRealNameArr)
		// 打印
		if err != nil {
			log.Errorf("插入默认数据%s-%s 失败, err= %s", dbNameComic, tableNameCountry, err)
		} else {
			log.Infof("插入默认数据%s-%s 成功个数: %v", dbNameComic, tableNameCountry, len(defaultDataCountryArr))
		}

		// -- 插入默认数据-type
		// 一级分类
		typeDefaultNoTypeLevel1 := &models.Type{NameId: 0, Name: "待分类", Level: 1}
		typeDefaultKoren := &models.Type{NameId: 1, Name: "韩漫", Level: 1}
		typeDefaultJapan := &models.Type{NameId: 2, Name: "日漫", Level: 1}
		typeDefaultRealPerson := &models.Type{NameId: 3, Name: "真人漫画", Level: 1}
		typeDefault3D := &models.Type{NameId: 4, Name: "3D漫画", Level: 1}
		typeDefaultAmeraica := &models.Type{NameId: 5, Name: "欧美漫画", Level: 1}
		typeDefaultSameSex := &models.Type{NameId: 6, Name: "同性", Level: 1}

		defaultDataTypeArr := []*models.Type{
			// 一级分类
			typeDefaultNoTypeLevel1, typeDefaultKoren, typeDefaultJapan,
			typeDefaultRealPerson, typeDefault3D, typeDefaultAmeraica,
			typeDefaultSameSex,
		}
		typeUniqueIndexArr := []string{"NameId"}                             // 唯一索引
		typeUpdateDBColumnRealNameArr := []string{"name", "level", "parent"} // 要更新的字段

		// 执行插入操作 --
		err = DBUpsertBatch(defaultDataTypeArr, typeUniqueIndexArr, typeUpdateDBColumnRealNameArr)
		// 打印
		if err != nil {
			log.Errorf("插入默认数据%s-%s 失败, err= %s", dbNameComic, tableNameType, err)
		} else {
			log.Infof("插入默认数据%s-%s 成功个数: %v", dbNameComic, tableNameType, len(defaultDataTypeArr))
		}
	*/

	// v0.1 方式：不使用通用 增删改查方法
	/*
		// 插入默认数据-website
		// 插入默认数据-website
		websiteDefaultNoClass := &models.Website{Name: "待分类", NameId: 0, Domain: "未知", NeedProxy: 0, IsHttps: 0}
		websiteDefaultJ88d := &models.Website{Name: "j88d", NameId: 1, Domain: "http://www.j88d.com", NeedProxy: 0, IsHttps: 0} // 请求domain 时带上http://
		defaultWebsites := []*models.Website{websiteDefaultNoClass, websiteDefaultJ88d}
		WebsiteOps.BatchAdd(defaultWebsites)

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
		CountryOps.BatchAdd(countries)

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
		TypeBatchAdd(defaultTypes)
	*/
}

// TruncateTable 清空指定模型对应的数据表，同时跳过外键检查（适用于 MySQL）
func TruncateTable(db *gorm.DB, model interface{}) error {
	// 优先从 db.Statement 获取表名
	tableName := db.Statement.Table

	// 如果 db.Statement 未设置表名，尝试解析 model
	if tableName == "" && model != nil {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err == nil {
			tableName = stmt.Schema.Table
		}
	}

	// 解析不到table，就return
	if tableName == "" {
		return fmt.Errorf("无法确定表名")
	}

	log.Debug("清空表: ", tableName)

	// 使用事务执行多个 SQL 语句

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

// ---------------------------- 函数 end ----------------------------
