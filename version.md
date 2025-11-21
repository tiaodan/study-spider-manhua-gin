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

# v0.0.0.8
- 能实现漫画的爬取，章节还不行，增加几个字段，下个版本封装爬取

# v0.0.0.9 2025.5.7
- 并发爬取某个网站的漫画，但未完全实现，待更新

# v0.0.0.10 2025.5.7v2
- 改了些,没改完

# v0.0.0.10 2025.5.8
- 改了些,没改完

# v0.0.0.10 2025.5.8v2
- 改了些单元测试报错,没改完。最新问题：查询结果 不是按照 name_id 升序排列，而是按 id（主键）排序，这种情况怎么处理？如何按指定排序

# v0.0.0.11 2025.5.8v3
- 上一次提交，好像需要合并，用服务器提交

# v0.0.0.12 2025.5.9
- website 单元测试改完,封装一点点

# v0.0.0.12 2025.5.12
- website 单元测试封装完, 减少200行代码 950 -> 750

# v0.0.0.12 2025.5.13
- 重新封装website,初始化db时，用公用的db，不自己再新建连接

# v0.0.0.12 2025.5.13v2
- 移除git.zip 

# v0.0.0.12 2025.5.13v3
- 封装db单元测试，没改完

#  v0.0.0.13 2025.5.14
- 修改完错误，go test通过

#  v0.0.0.13 2025.5.15
- 加异常单元测试，没改完

# v0.0.0.13 2025.5.25
- 改了一部分，没改完

# v0.0.0.13 2025.5.25v2
- 改了一部分db，没改完，重新提交

# v0.0.0.13 2025.5.26
- 阶段一：db的测试函数，每个增删改成都占用一个函数。相当于每个测试用例都占用一个函数
- 阶段二：db的测试函数，封装一个通用函数，把增删改查、批量增删改查的操作，都放到一个函数里，然后给函数传测试用例。没写完

# v0.0.0.13 2025.5.27
- 单元测试，封装成一个最简单通用的函数？如何实现？接口+map？没写完

# v0.0.0.13 2025.5.27v2
- 单元测试，封装成一个最简单通用的函数？如何实现？接口+map？没写完
- 整合完一部分website_normal代码，待完善

# v0.0.13 2025.5.28
- 把checkHasId 和NoId 整合成一个方法。自动判断有没有id。但是单元测试运行到TestCommon 不清空表

# v0.0.13 2025.5.28v2
- 把checkHasId 和NoId 整合成一个方法。自动判断有没有id。但是单元测试运行到TestCommon 可以清空表

# v0.0.13 2025.5.30
- 思考测试用例 思维导图怎么设计。修改test.go文件序号，方便go test按顺序执行

# v0.0.13 2025.5.30v2
- 改了改测试用例思维导图

# v0.0.13 2025.6.3
- 实现增、批量增测试用例，测试没问题

# v0.0.13 2025.6.4
- 实现增、批量增测试用例，测试没问题。把测试用例使用变量，改为不用指针方式

# v0.0.13 2025.6.5
- website  normal 增删改查弄完了，测试通过
- returnObjZeroOneNegate updates 多一个

# v0.0.13 2025.6.6
- 修复问题: returnObjZeroOneNegate updates 多一个(v0.0.13 2025.6.5)

# v0.0.0.14
版本总结：
    - 封装db 通用增删改查模板
    - /doc 目录, 画了一些逻辑图

核心改动：
    - db操作封装成通用模板，所有表增删改成操作，都用 dbtempalte.go
    - 把源码都放到src目录,包括main.go
待解决问题：！！！！！！
- dbtemplate.go 模板方法，只能连comic数据库,因为用的comic数据库的DB对象。得实现一个通用方式：能连各种数据库，各种操作
    - 思路，只要给方法传一个 ，数据库链接对象的，指针就行
    - 实现：dbtempalte.go 里批量upsert方法里，加了 数据库链接对象
非核心改动：
- 学会:用Swark VSCode插件一键生成项目流程图，它可以读取项目。 Ctrl+Shift+R 快捷键选中项目后，会在选中目录生成md文件，点击Mermaid Live Editor: View。但是swark只能显示目录结构，函数调用关系，无法看到
- 把 model/category.go 模型改为 -》PornType

# v0.0.0.15
版本总结：
    - 封装 通用的爬取漫画 方法 (爬JSON方式)
核心改动：
    - 封装 通用的爬取漫画 方法
    - sex_types表改为 porn_types表,表示是否是色情内容
    - 封装 处理对象 string字段，前后空格的 工具方法
        一般有2种实现方式：
        1. 用反射实现，缺点：反射比较耗时，写起来麻烦，不容易理解 -- 弃用
        2. 用接口实现，util里 定义一个结果，所有models里的 struct 都实现这个接口 -- 推荐，现在用的这种方法，参考stringutils.go
    - comic 表 download_end / upload_end / upload_baidu_end 最好改为 download_end—_status 因为可能有好几个数，比如0 1 2 3 各戴白不同状态，如果叫end，好像是只能true/false了
	- 更新comic表结构，把needTcp/coverNeedTcp 改成bool，不用bigint
    - 修改db.DB -> 改为 db.DBComic。表示comic数据库 的对象 -> 为了以后能适配各种数据库，适配一个数据库，加一个对象
    - 实现，tooptoon.net-台湾的网站,能根据JSON爬取某一类的所有book。JSON数据 --> https://tooptoon.net/api/v1/comics?category=1&limit=10&page=1 -> F12 -> Fetcg.XHR -> 找第一个json。比如：https://d1dkh1tjti8mih.cloudfront.net/www_v1/jsonComic/tw/complete/54bc1d768276cf834a8276670d457013ce2ec7d078b4fbef707b16fa30508a9d.json
    - comic表，need_tcp用不到，删除， coverNeedTcp 删除。链接是否用到https，放到website表里

非核心改动:
    - 把website model 的 url字段改成domain，表示域名，url容易误解
    - 把website model 的 needProxy isHttps 都改为bool类型
    - 用接口，实现model，只要是string，都转成简体中文 !!!!
    - 把comics 表改成 comcic_url -> comic_url_api_path, comic_url 容易误解  
    - 把comics 表改成 cover_url -> cover_url_api_path, comic_url 容易误解
    - 把请求的json，放到1个文件保存起来 -》 放到doc/F12找到的json里了


------------------------------------------ 未解决问题如下：
要解决：
    - website表，加1个coverNeedHttps 列，表示 cover 是否需要https ，因为默认一定是要http请求头的。 ，// 需要这个，因为图片有的开
    - 插入前数据校验，你如star，最高10.0，如果超出，就按0算
    - 数据清洗的时候，如果有https或者http头，自动删除。comic 实现一个数据清洗接口，（数据清洗自动实现，去除空格、繁体转简体，自动去除协议头：http/https，超出范围，自动置为某个值）
    - 数据清洗分2个方面：1. 前端传参、方法间传参，数据清洗-属于前端编程人员操作 2. 插入db前数据清洗 -》 属于后端编程人员操作
    - website还要加个前缀，img_prefix 图片前缀，比如cover前缀，可能不是域名，比如：https://cdn.mangakakalot.tv/mangakakalot/covers/xxx.jpg
    - website表加2列，章节url拼接规则，举例。图片拼接规则，举例
    - 拆分comic表，把打分、点击率单独拿出来，因为可能经常更新。这样需要联表查询，这块代码需要改下
    - 给所有报错，给出推断原因，让用户自己去简单排查。缩短排查时间
    - 把mapping映射关系，写出json文件，类似配置文件的方式，通用，以后不用改代码，直接改json文件即可 -》 参考 spider.go -> ComicMappingForSpiderToptoonByJSON 这个变量

------------------------- 解决完再上传

# 待办
- comic表加上时间，创建时间，更新时间，是否删除标志 √ 没测
- 更新至，内容包含 '最终话'、'完结' 认为是完结了 √ 没测
- 人气，字符串里可能没带单位。如何处理？ 如：'人气：5555' √ 没测
- 漫画爬取过程去重 √ 没测
- 把队列池循环20改成1 √
- 爬取的漫画，关联外键type表 √ 没测
- 爬取的漫画，关联外键website表 √ 没测
- 爬取请求，主动区分完结+未完结。如果传的参数区分完结或者未完结，程序就不判断了 √ 没测
- 漫画添加是否删除标志位 √ 没测
- comic唯一索引，应该改为组合形式 网站+ 名字  √ 没测
- 删除传参，要传deleted 字段 √ 没测
- 单元测试
- 图片加一个是否重复字段
- 可以人为控制爬一个章节、一个图片、上传一个章节、上传一个图片，爬取上传、手动上传
- log大小限制
- 并发数，写成配置文件
- 找一个官方网站，类似喜马拉雅，判断是完结还是连载
- 写android程序，测试漫画显示
- 单元测试t.Log不打日志，fmt.Println打日志
- 单元测试 go run ./db 不打日志，执行不到test文件, cd db && go test 可以打日志
- 增加其它项目db 增删改查方法，可以通过id操作，也可以通过其它字段操作，增加灵活性 study-db-gorm study-gin-canuse
- 修改db 文件增删改查byid byother 批量byid byother
- 单元测批量操作
- website 单元测试写完，再写其它的
- db封装 byid byNameId byOther 增删改查操作都封装这几个方法
- 其它项目model 有全部大写字母的内容，改成首字母大写。study-db-gorm study-gin-canuse 如 website里的model ID URL
- 安卓、平板，实现自动翻页，自动滚屏，定时滚屏
- 单元测试，查询结果 不是按照 name_id 升序排列，而是按 id（主键）排序，这种情况怎么处理？如何按指定排序
- 完善通用场景测试用例，完善xmind+代码。单元测试、接口测试、功能测试使用的用例还是稍微有区别的
- 单元测试，会在目录生成app.log，如何生成到项目根目录
- 单元测试，封装成一个最简单通用的函数？如何实现？接口+map？
- 单元测试，通用方法，不传model.Website 指针，通过用户数据的表名，自动判断用哪个dbOps
- 把checkHasId 和NoId 整合成一个方法。自动判断有没有id
- 测试清空函数
- 写测试用例的时候，每一个功能点弄一组测试用例。日志打印，打每组有多少用例
- 写测试用例的时候，每个用例对应一个level，方便以后只测试指定level用例