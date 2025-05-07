# v0.0.0.0 
- 整合配置文件, from github项目: study-config-viper
- 整合日志文件, from github项目: study-log-go-original
- 整合错误处理, lian, from github项目: study-error-go-original
- 整合数据库  , from github项目: study-db-gorm
- 参考项目: study-restful-api-gin
- 做成一个gin服务,前后端分离

# 待办
- 插入默认数据country , category, website, √
- log打印，全换成 logger。config.go 不能换，因为先读取了配置, 日志级别才能生效 √
- gorm 项目，把addType 都改为 TypeAdd √
- 插入默认数据 √
- 插入默认数据types(types还不知道加哪些,先简单按国家分) √
- 删除order数据库，保留,别的项目参考 √
- 解决日志打印xx.go 文件名不对, 都是logger.go文件问题 √
- 处理gin日志和自己写的日志冲突 √ gin设置成release模式即可
- 日志分Debug,和Debugf,为了对应 log.Println 和log.Printf，后来整合成一个Debug √
- 如何通过lock的mod文件, 引入新项目, 防止版本冲突。用的时候拷贝go.mod+go.sum就行 √
- 更新study-log-go-original 代码 √
- 做一个整合好所有项目的项目, 保证拿来就能用 √

- gorm 漫画 项目，没配置外键关联
- http请求传输，考虑关键字防屏蔽。如改为 混乱的pinyin,或者混乱的英文
- 做一个待办列表项目
- 调研网站, 哪个适合爬,分类好

# v0.0.0.1
- 做好一个整合好所有项目的项目, 保证拿来就能用 

# v0.0.0.2
- 解决logger.Debug(xxxx)不打印问题

# v0.0.0.3
- main.go 去除defer app.Close(),因为写不进去文件

# v0.0.0.4
- 爬不到内容, 待修改

# v0.0.0.5
- 实现漫画的增删改查，ComicUpdate() 参数可以是int/string 都能插入成功

# v0.0.0.6
- 爬取漫画操作, 去除前后空格

# v0.0.0.7
- fix-bug: 爬取漫画,插入报错:1054 (42S22): Unknown column 'comicUrl' in 'field list'。原因comic.go ComicAdd()方法，参数用_下划线方式
- 增加更新方法: 排除唯一索引字段 方式更新
- fix-bug：解决，db.Updates() 0值不更新问题