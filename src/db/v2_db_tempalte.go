/*
作用：db操作 通用模板V2。实现方式：插件函数工厂/函数选项模式
简单理解：比如查询函数，where的所有条件，都变成一个函数，用的时候加函数就行
目标：1套代码，解决所有项目 增删改查操作

数据库操作详细日志：不在此go文件打。原因：
  - 此文件如果出错，已经返给上级错误原因了

排名,方案名称,通用程度,代码复杂度,扩展性,推荐指数,备注
1,Option/Func 选项模式,★★★★★,中,★★★★★,最高,当前最推荐的方式
2,Builder 模式（链式）,★★★★☆,中高,★★★★☆,很高,GORM 自己就是这种风格
3,扩展参数结构体,★★★★,低~中,★★★★,高,比较务实，容易理解
4,map + 特殊约定key,★★★,低,★★☆☆☆,中低,容易失控，不建议长期用
5,完全不同的新方法,★★☆☆☆,-,-,低,违背“改一个最通用方法”的初衷


最推荐写法：使用 Option/Func 模式（函数选项模式）

使用示例：
使用示例（你要的场景：某个 parent_id 下 sort_num 最大的章节）
chapter, err := FindOne[ComicChapter](
    WithWhere("parent_id = ?", parentID),
    WithOrder("chapter_real_sort_num DESC"),
    WithLimit(1),
)

// 或者用快捷方式
chapter, err := FindOne[ComicChapter](
    WithWhere("parent_id = ?", parentID),
    WithMaxSortNum(),
)
*/

package db

import "gorm.io/gorm"

// 查询选项函数类型
type FindOption func(*gorm.DB) *gorm.DB

// 核心通用查询方法
func DBFindOneV2[T any](dbConn *gorm.DB, opts ...FindOption) (*T, error) {
	var result T

	db := dbConn // 假设这是你的全局 db 实例

	// 依次应用所有选项
	for _, opt := range opts {
		db = opt(db)
	}

	err := db.First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// 统计数量
func DBCountV2[T any](dbConn *gorm.DB, opts ...FindOption) (int, error) {
	var count int64

	db := dbConn // 假设这是你的全局 db 实例

	// 依次应用所有选项
	for _, opt := range opts {
		db = opt(db)
	}

	err := db.Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// 通用执行选项 函数
func applyOptions(dbConn *gorm.DB, opts ...FindOption) *gorm.DB {
	db := dbConn
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// 通用 提取某一列 Pluck
/*
参数:
	column string, // 要 pluck 的列名(数据库小写形式)，例如 "id"
	dest *[]R, // 接收结果的切片指针，例如 *[]int64
*/
func DBPluckV2[T any, R any](dbConn *gorm.DB, column string, dest *[]R, opts ...FindOption) error {
	db := applyOptions(dbConn.Model(new(T)), opts...)
	return db.Pluck(column, dest).Error
}

// ------------------ 常用的选项工厂函数 ------------------

// Where 条件
func WithWhere(query any, args ...any) FindOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

// 排序
func WithOrder(value string) FindOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

// Limit
func WithLimit(limit int) FindOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// Offset（配合分页用）
func WithOffset(offset int) FindOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

// 预加载关联
func WithPreload(query string, args ...any) FindOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}
