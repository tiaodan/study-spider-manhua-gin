/**
功能: 分发爬取请求
*/

package spider

import (
	"io"
	"strconv"
	"study-spider-manhua-gin/src/log"
	"study-spider-manhua-gin/src/models"

	"github.com/gin-gonic/gin"
	"github.com/mohae/deepcopy"
	"github.com/tidwall/gjson"
)

// -- 初始化 ------------------------------------------------------------------------------

// -- 初始化 ------------------------------------------- end -----------------------------------

// -- 方法 ------------------------------------------------------------------------------

// -- 分发请求 /spider/oneCategoryByJson。自行判断，该用哪个 表的 ModelMapping。不应该用 _命名方式，但是能看清
/*
作用简单说：
	- 分发请求 /spider/oneCategoryByJson。自行判断，该用哪个 表的 ModelMapping

作用详细说:

核心思路:
	1. 读取 前端JSON里 spiderTag -> website字段
	2. 根据该字段，使用不同的爬虫 ModelMapping映射表
	3. 调用通用 爬取方法

参考通用思路：
	1. 校验传参
	2. 数据清洗
	3. 业务逻辑 需要的数据校验 +清洗
	4. 执行核心逻辑
	5. 返回结果

参数：
	1. context *gin.Context  // 读取 前端JSON里 spiderTag -> website字段，根据该字段，使用不同的爬虫 ModelMapping映射表
	2. xx

返回：

注意：

使用方式：
// gjson 读取 前端JSON里 spiderTag -> website字段 --
website := gjson.Get(string(data), "spiderTag.website").String()
id := gjson.GetBytes(data, "adult.1.id").String()
*/
func DispatchApi_OneCategoryByJSON(c *gin.Context) {
	// 0. 初始化
	okTotal := 0 // 成功条数

	// 1. 校验传参
	// 2. 数据清洗

	// 3. 业务逻辑 需要的数据校验 +清洗
	// -- 找到应该爬哪个网站
	// 读取 JSON Body --
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "func: 通过json爬分类。读取 前端传参 Body 失败"})
		return
	}

	// gjson 读取 前端JSON里 spiderTag -> website字段 --
	website := gjson.Get(string(data), "spiderTag.website").String()
	adultArrGjsonResult := gjson.GetBytes(data, "adult").Array()

	// -- 根据该字段，使用不同的爬虫 ModelMapping映射表
	switch website {
	case "toptoon-tw":
		var gjsonResultArr []map[string]any // 批量插入用的参数。爬取到的 数据表对象 数组

		// -- 通过mapping 循环读取每条内容
		for i := range adultArrGjsonResult {
			// -- 给mapping 赋值
			// ！！！！ 非常重要。临时mapping，如果不每次都用新变量，ComicMappingForSpiderToptoonByJSON 带%d,第二次赋值 不行，会导致后面报错
			mappingTemp := deepcopy.Copy(ComicMappingForSpiderToptoonByJSON).(map[string]models.ModelMapping) // 需要深拷贝写法，并强制转成期望类型。mappingTemp := ComicMappingForSpiderToptoonByJSON 还是浅拷贝写法，还是指针，因为go里map全是指针。
			log.Debug("------- delete, deepCopy的 mapping = ", mappingTemp)
			mapping := mappingAssign(mappingTemp, i) // 返回空，说明有问题。 原来写法：第二次赋值不行，报错。-》 mapping := mappingAssign(ComicMappingForSpiderToptoonByJSON, i)
			log.Debug("------- delete, deepCopy的 mapping,赋值后 = ", mapping)
			if mapping == nil {
				c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeByJson), mapping 赋值失败"}) // 返回错误
				return                                                                                                        // 直接结束                                                                                                      // 直接金额数
			}

			// -- 根据 mapping爬取内容
			oneObjGjsonResult, err := BookTemSpiderTypeByJson(data, mapping) // gin.Context 只能读1次，已经被读取了，所以不能传。因此传的2进制data
			log.Debug("---------- delete oneObjGjsonResult = ", oneObjGjsonResult)
			if err != nil {
				c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeByJson), 爬取失败"}) // 返回错误
				return
			}

			// -- 准备插入db 用的数据
			gjsonResultArr = append(gjsonResultArr, oneObjGjsonResult)
		}

		// -- 批量插入db，循环处理 gjsonResultArr[]
		err := upsertSpiderTableData("comic", gjsonResultArr)
		if err != nil {
			log.Error("func=BookTemSpiderTypeByJson(爬取JSON). 插入db失败, err: ", err)
			c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeByJson), 批量插入db失败"}) // 返回错误
			return                                                                                                    // 返回这个c
		}

		okTotal = len(gjsonResultArr) // 成功条数
		log.Infof("func=DispatchApi_OneCategoryByJSON(分发api: /spider/oneTypeByJson), 爬取成功, 插入%d条数据", okTotal)

	default:
		c.JSON(400, gin.H{"error": "func=DispatchApi_OneCategoryByJSON(分发api- /spider/oneTypeByJson), 没找到到应爬哪个网站. 建议: 排查json参数 apiderTag-website"}) // 返回错误
	}

	// 4. 执行核心逻辑
	// 5. 返回结果
	c.JSON(200, "爬取成功,插入"+strconv.Itoa(okTotal)+"条数据")
}

// -- 方法 ------------------------------------------- end -----------------------------------
