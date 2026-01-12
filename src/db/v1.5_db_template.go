/*
作用：db操作 通用模板，针对v1.5 请求
目标：1套代码，解决所有项目 增删改查操作

数据库操作详细日志：不在此go文件打。原因：
  - 此文件如果出错，已经返给上级错误原因了
*/
package db

import "gorm.io/gorm"

// ------------------------------------------- 方法 -------------------------------------------

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
func DBFindOneByFieldV1_5[T any](DBConnObj *gorm.DB, field string, value any) (*T, error) {
	// 1. 校验传参
	// 2. 数据清洗
	// 3. 准备数据库执行，需要的参数
	// 4. 数据库执行
	var result T
	db := DBConnObj.Where(field+" = ?", value).First(&result)
	if db.Error != nil {
		// log.Error("查询失败: ", db.Error)  // 此文件不打日志，错误已经返回给上级
		return nil, db.Error
	}
	// log.Infof("查询成功, 查询到 %d 条记录", 1)  // 此文件不打日志，错误已经返回给上级

	// 5. 返回结果
	return &result, nil
}
