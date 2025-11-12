// 功能: 封装restfult api - comic模块
package comic

import (
	"strconv"
	"study-spider-manhua-gin/src/db"
	"study-spider-manhua-gin/src/errorutil"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	"github.com/gin-gonic/gin"
)

// 增
/*
作用简单说：
	- 接口: 新增1条数据

作用详细说:

思路:
	1. 解析前端传参
	2. 调用数据库操作
	3. 返回结果

参考通用思路：
	1. 校验传参
	2. 解析前端传参。 参数 -》 转成 对象
	3. 数据清洗
	4. 执行核心逻辑
		1. 准备数据库执行，需要的参数
		2. 数据库执行
	5. 返回结果
*/
func ComicAdd(c *gin.Context) {
	log.Debug("增加漫画")

	// 1. 校验传参

	// 2. 解析前端传参。 参数 -》 转成 对象
	// 参数 -》 转成 对象
	var comic models.Comic
	if err := c.ShouldBindJSON(&comic); err != nil {
		log.Error("解析请求体失败, err: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return // 必须保留 return，确保绑定失败时提前退出
	}

	// 3. 数据清洗

	// 4. 调用数据库操作
	// v1.0 方式，调用comic的 增删改查方法，弃用。因为有了通用方法模板
	// err := db.ComicUpsert(&comic)

	// v2.0 方式，调用通用的增删改查方法模板。 推荐!!!
	// -- 准备 唯一索引字段名
	unquieIndexArr := []string{"Name"}
	// -- 准备 要更新的字段 数据 map
	// v0.1 写法。db.DBUpsert 第3个参数，还要手动传值。弃用
	/*
		updateColumnsMap := map[string]any{
			"country_id":       comic.CountryId,
			"website_id":       comic.WebsiteId,
			"category_id":      comic.CategoryId,
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
		}
		err := db.DBUpsert(&comic, unquieIndexArr, updateColumnsMap)
	*/

	// v0.2 写法。db.DBUpsert 第3个参数，直接传 数据库真实列名。 推荐!!!
	updateDBColumnRealNameArr := []string{
		"country_id",
		"website_id",
		"category_id",
		"type_id",
		"update",
		"hits",
		"comic_url",
		"cover_url",
		"brief_short",
		"brief_long",
		"end",
		"star",
		"need_tcp",
		"cover_need_tcp",
		"spider_end",
		"download_end",
		"upload_aws_end",
		"upload_baidu_end",
	}
	err := db.DBUpsert(&comic, unquieIndexArr, updateDBColumnRealNameArr)

	if err != nil {
		log.Error("增加漫画失败, err: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回结果
	c.JSON(200, "添加成功")
}

// 删
/*
作用简单说：
  - 接口: 删除1条数据

作用详细说:

思路:
	1. 解析前端传参
	2. 调用数据库操作
	3. 返回结果

参考通用思路：
	1. 校验传参
	2. 解析前端传参。 参数 -》 转成 ?
	3. 数据清洗
	4. 执行核心逻辑
		1. 准备数据库执行，需要的参数
		2. 数据库执行
	5. 返回结果

注意：
	DB.Delete(&models.Comic{}, id) 这种传id写法，gorm会不同处理。建议校验id >0。因为为负数，会执行sql，以防万一
	参数      是否执行SQL      SQL内容
	- 0      否             无操作
	- 1      是             DELETE FROM `comics` WHERE `id` = 1
	- -1     是             DELETE FROM `comics` WHERE `id` = -1
*/
func ComicDelete(c *gin.Context) {
	// 1. 校验传参

	// 2. 解析前端传参

	// 提取前端传递的 id 参数
	idStr := c.Param("id")
	log.Debug("删除漫画, 参数= ", idStr)
	// id, err := strconv.ParseUint(idStr, 10, 64) // 转换为 十进制 64 位无符号整数。// v0.1 没用 DB通用方法之前，设置 id类型是 unit，所以要进行转换处理
	id, err := strconv.Atoi(idStr) // v0.2 用DB通用方法写法。不用转成Uint了 Atoi，默认是32位int
	if err != nil {
		log.Error("删除漫画, 参数错误, ParseInt出现错误")
		c.JSON(400, gin.H{"error": "删除漫画, 参数错误, ParseInt出现错误"})
		return
	}
	// 参数校验
	if id <= 0 {
		log.Error("删除漫画, 参数错误, id不合法, 应>0")
		c.JSON(400, gin.H{"error": "删除漫画, 参数错误, id不合法, 应>0"})
		return
	}

	// 3. 数据清洗
	// 4. 调用数据库操作

	// err = db.ComicDelete(uint(id)) // 弃用. v0.1 没用 DB通用方法之前，设置 id类型是 unit，所以要进行转换处理
	err = db.DBDeleteById(&models.Comic{}, id) // v0.2 用DB通用方法写法。不用转成Uint了
	if err != nil {
		log.Error("删除漫画失败, err: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回结果
	c.JSON(200, "删除成功")
}

// 改
func ComicUpdate(c *gin.Context) {
	log.Debug("修改漫画")
	// bodyBytes, _ := io.ReadAll(c.Request.Body)  // 测试用-可以删
	// log.Debug("请求内容c.request.body= ", string(bodyBytes))  // 测试用-可以删，这段代码影响c.ShouldBindJson

	// 绑定前端数据
	var comic models.Comic
	if err := c.ShouldBindJSON(&comic); err != nil {
		log.Error("解析请求体失败, err: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return // 必须保留 return，确保绑定失败时提前退出
	}
	err := db.ComicUpdateByIdOmitIndex(comic.Id, &comic)

	if err != nil {
		log.Error("修改漫画失败, err: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "修改成功")
}

// 改 - 根据id, 排除唯一索引
/*
疑问: 为什么要排除唯一索引字段?
答: 唯一索引很关键,作用比id还重要。防止误更新 唯一索引字段

作用简单说：
  - 接口: 更新
  	- 只更新 指定字段，如DB.Model().Select(指定字段).Update()，中Select()中的字段
	- 不更新 唯一索引字段。如唯一索引叫 name, 写代码的时候要排除它

作用详细说:

思路:
   1. 解析前端传参
   2. 调用数据库操作
   3. 返回结果

参考通用思路：
	1. 校验传参
	2. 解析前端传参。 参数 -》 转成 ?
	3. 数据清洗
	4. 执行核心逻辑
		1. 准备数据库执行，需要的参数
		2. 数据库执行
	5. 返回结果
*/
func ComicUpdateByIdOmitIndex(c *gin.Context) {
	log.Debug("修改漫画")

	// 1. 校验传参

	// 2. 解析前端传参
	// 前端数据 -》 转成对象
	var comic models.Comic
	if err := c.ShouldBindJSON(&comic); err != nil {
		log.Error("解析请求体失败, err: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return // 必须保留 return，确保绑定失败时提前退出
	}
	log.Debug("修改漫画, 参数.needTcp= ", comic.NeedTcp)

	// 校验参数
	if comic.Id <= 0 {
		log.Error("修改漫画, 参数错误, id不合法, 应>0")
		c.JSON(400, gin.H{"error": "修改漫画, 参数错误, id不合法, 应>0"})
		return
	}

	// 3. 数据清洗
	// 4. 执行核心逻辑
	// -- 调用数据库操作
	// 准备数据库执行，需要的参数 --
	// 要更新的列
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

	// err := db.ComicUpdateByIdOmitIndex(comic.Id, &comic) // 弃用. v0.1 没用 DB通用方法之前，设置 id类型是 unit，所以要进行转换处理
	err := db.DBUpdateByIdOmitIndex(&comic, comic.Id, updateDataMap) // v0.2 用DB通用方法写法
	if err != nil {
		log.Error("修改漫画失败, err: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回结果
	c.JSON(200, "修改成功")
}

// 查
/*
返回: json对象
{
	"total": 0,
	"data": []
}
*/
func ComicsQuery(c *gin.Context) {
	log.Debug("查询所有漫画")
	total, err := db.ComicsTotal() // 补充总数获取
	errorutil.ErrorPrint(err, "查询漫画总数失败")
	comics, _ := db.ComicsQueryAll()

	c.JSON(200, gin.H{
		"total": total,
		"data":  comics,
	})
}

// 查-分页
/*


作用简单说：
  - 接口: 更新
  	- 只更新 指定字段，如DB.Model().Select(指定字段).Update()，中Select()中的字段
	- 不更新 唯一索引字段。如唯一索引叫 name, 写代码的时候要排除它

作用详细说:

核心思路:
	1. 获取前端传参,并做校验。没传page和size, 不处理, 返回
	2. 参数缺失校验
	3. 参数类型校验
	4. 业务逻辑

参考通用思路：
	1. 校验传参
	2. 解析前端传参。 参数 -》 转成 ?
	3. 数据清洗
	4. 执行核心逻辑
		1. 准备数据库执行，需要的参数
		2. 数据库执行
	5. 返回结果

返回: json对象
{
	"total": 0,
	"data": []
}
*/
func ComicsQueryByPage(c *gin.Context) {
	log.Debug("分页查询漫画")

	// 1. 校验传参

	// 强校验参数类型
	pageStr := c.DefaultQuery("page", "") // 之前写法默认为 1, pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "") // 之前写法默认为 10 ,所以不存在类型不是string类型, sizeStr := c.DefaultQuery("size", "10")
	log.Debugf("前端传参, page=%v, size=%v", pageStr, sizeStr)

	// 参数缺失校验
	if pageStr == "" || sizeStr == "" {
		c.JSON(400, gin.H{"error": "参数缺失"})
		return
	}

	// 参数类型校验
	if _, err := strconv.Atoi(pageStr); err != nil {
		c.JSON(400, gin.H{"error": "page参数类型错误"})
		return
	}
	if _, err := strconv.Atoi(sizeStr); err != nil {
		c.JSON(400, gin.H{"error": "size参数类型错误"})
		return
	}

	// 2. 解析前端传参。 参数 -》 转成 ?
	// 3. 数据清洗
	// 4. 执行核心逻辑
	// -- 准备数据库执行，需要的参数
	total, err := db.ComicsTotal() // 总数
	errorutil.ErrorPrint(err, "查询漫画总数失败")

	page, _ := strconv.Atoi(pageStr) // 因为默认都是数字str了，所以不存在报错情况
	size, _ := strconv.Atoi(sizeStr) // 因为默认都是数字str了，所以不存在报错情况
	// 校验 <=0
	if page <= 0 || size <= 0 {
		c.JSON(400, gin.H{"error": "page或size参数值错误, 应>0"})
		return
	}

	// -- 数据库执行
	// comics, _ := db.ComicsPageQuery(page, size)  // v1.0 方式，弃用。因为有了通用方法模板
	comics, err := db.DBPageQueryReturnTypeT[*models.Comic](page, size) // v2.0 方式，调用通用的增删改查方法模板。 推荐!!!
	if err != nil {
		log.Error("分页查询漫画失败, 执行sql失败, err: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回结果
	// 构造指定的返回结构
	c.JSON(200, gin.H{
		"total": total,
		"data":  comics,
	})
}
