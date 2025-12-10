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


# v0.0.0.16
版本总结:
    - 修改数据库基础框架，增加字段

核心改动:
    1. 搭框架相关：
    - doc/json文件，wbsiteId还得改，aws的id=2了，postman保存的数据也要改，先用10吧
    - website还要加几个默认值：aws网站、预留了一个网站 -》 namedId =3
    - website表，
        - 加1个 cover_url_is_need_https 列，表示 cover 是否需要https，因为默认一定是要http请求头的。 // 需要这个，因为图片有的开. 只要是是否的字段，命名时都加个is
        - 加1个 chapter_content_url_is_need_https 列 表示 章节内容url(图片/视频等) 是否需要https
        - 如果加了上述内容，插入默认数据func,修改
    - website表加2列
        - cover_url_concat_rule // 封面ULR拼接规则 string
        - chapter_content_url_concat_rule 章节内容URL拼接规则 string  // concatenate 拼接英文，缩写 concat
        - 如果加了上述内容，插入默认数据func,修改
    - website表还要2加个前缀，
        - cover_domain 封面图片用的域名，也可以用ip比如cover前缀，可能不是域名，比如：https://cdn.mangakakalot.tv/mangakakalot/covers/xxx.jpg
        - chapter_content_domain 章节内容域名，也可以用ip
        - 如果加了上述内容，插入默认数据func,修改
    - 所有的表model，加上check，确保不能输入空字符串
    - comic update字段，改成别的，不能用update关键字，改成 latest_chater。因为用upate报错，占用系统关键字  // 最新章节
    - comic表，需要更新名字：爬取存数据用的，叫comic_spider, comic_my -》我的数据库，业务真实用的comic
    - comic 少字段 cover_save_path_api_path  // 保存路径的api，这个字段需要只在我的数据有吗？ -> 是的
    - comic 少字段 release_date // 发布时间. 时间+日期，如果没有具体日期，按 00:00:00 时间来
        - 增删改查、方法需要同步改
        - mapping是否要改？ 需要，因为toppoon的爬取 By JSON，有这个字段。测试下能否爬到
        - comci_spider, comci_my 都要改
    - website 需要一个字段(isRefer)，是否是参考/参照/refer网站，比如：tooptoon，参考网站，爬取的漫画，需要去tooptoon网站爬取。
        - 插入默认数据，插入默认数据-要更新的列，同步改
    - 加一个网站 娱乐类型分类 website_type表，比如漫画、小说、视频、关联到website表，
        - 插入默认数据，插入默认数据-更新的列, website更新的列, 数据迁移,需要同步修改 
    - 少一个进度表 id process 完结/连载
        - 插入默认数据，插入默认数据-更新的，数据迁移，需要同步修改
        - comic_my + comic_spider 是否得加上 processId外键
        - comic_spider表, 爬取，需要更新字段吗？- 要
            - 那前端传的json也要同步更新，加上 processId字段 -- >doc/json

非核心改动:
    1. 搭框架相关
        - 插入默认数据失败，就panic. 因为默认数据，必须插入成功
        - !!!!!!!! 为什么多了一个字段后，run main.go 不会自动更新列？？？ 因为看错表了，看的是comic,不是comic_spider表!

# v0.0.0.17
版本总结:
    - 1. 打算用gorm，插入comic-author关系表，没成功，因为comic表加了 authroArr,卡在MapByTag。2. 删除name_id字段

核心改动：
    先弄简单的：
    - comic数据库要有个author表，因为一个漫画/有声书/影视，可能有多个作者，所以得有个作者表
        - 有 author 表 
        - 有 comic_author_realation 表，表示 comic 和 author 的关系。多对多 
            - 外键要关联上comic_nameid,author_nameid
            - 建2个表 comci_spider_author_realation, comic_my_author_realation
        - 数据迁移, 需要同步修改 
    - comci表，唯一索引，加上author_concat 字段，叫作者拼接字符串
        - comic_spider + comic_my 表：???
            - comic 表，model, 唯一索引加上 author_concat √
            - comic 表，model, 数据清洗 （判断空，繁体转简体） √
            - 添加 author_concat_type, 作者拼接方式。比如：0 默认，按爬取顺序拼接，1: 按字母升序拼接 2:按我的意愿拼接 3: 参考最权威的网站拼接(b比如有声书，参考喜马拉雅，韩漫参考toptoon，小说参考 起点-建议0 /3) √
                - 修改爬取json/ doc里json √
            - comic 表，爬取操作，update列，加上 author_concat, author_concat_type √
            - mapping 表，update列，唯一索引列, 加上 autho_concat √
            - 爬取插入逻辑，需要修改。还要验证
                - 爬取数据加上  author_concat, author_concat_type √
    - models id类型都换成 int
        - website_type 
        - website
        - comic
        - author
        - process
        - comic_author_realation
        - comic_spider_author_realation
        - comic_my_author_realation

    再弄复杂的：
    - 删除NameId字段，用mysql的id字段作为索引，name_id字段可能多余
    - 解决一个问题，为什么website_type 插入默认数据，不是从1开始，而是从3开始？每次都是这样？
        - 具体原因未知，AI说是，可能为mysql自己优化逻辑导致，因为唯一索引是name
    - 唯一索引换成 非id，比如name √
        - website_type √
        - website √
        - type √
        - process √
        - porntype √
        - country √
        - comic-spider √
        - comic_my √
        - author √
        - comic_author_realation √
        - comic_spider_author_realation √
        - 插入默认数据，索引也要改 √
    - 删除所有表name_id  √
        - website表，删除name_id字段 √
            - 插入默认数据，同步改，简单搜索 用到NameId/name_id的地方，看是否要同步改
            - 报错的地方，同步改
            - comic 表关联外键的地方改 comic_spider, comic_my
            - 如果 运行程序出错，看是不是数据已经有了字段，删除所有表后，重新运行 
            - 插入默认数据，id就不能从0开始了，得从1开始，要不0那条插入不进去。-- 需要重新设计默认数据id
        - website  √
        - type 同上 √
        - process √
        - porntype √
        - country √
        - comic-spider √
        - comic_my √
        - author √
        - comic_author_realation √
        - comic_spider_author_realation √


非核心改动: 


# v0.0.0.18 
版本总结：
    - 解决v0.0.0.17问题: 插入comic+autho+ 加comic_spider_author 关联表

核心改动:
    1. 搭框架相关：
    先弄简单的：

    再弄复杂的：
    - comic author表爬取，插入逻辑修改
        - 插入默认数据，佚名，id=1 √
        - 插入 author表，√
        - 插入 comic_author_realation 表 √
        = 删除 comic_spider_author_relations comic_my_author_relations √
    - 还要爬取作者，有的是2个作者，还得想好怎么存。√
    - comcic还少个作者字段 author,有多个作者怎么办？ 需要加一个作者数据库？多个作者用 作者1&作者2&作者3，拼接成author字段

# v0.0.0.18.1
版本总结: 
    - comic,拆成2个表？？??????????????????? 没完成，计划用cursor实现，怕影响以前代码，先提交
    - 把之前business 里没用的代码，都改成nouse文件了。因为现在都用通用模板方法了
    - 把之前db 里没用的代码，都改成nouse文件了。因为现在都用通用模板方法了

核心改动：
    1. 搭框架相关：
    先弄简单的：
    - comic 少字段 release_date // 发布时间,每一章也要有发布时间 √

    再弄复杂的：

# v0.0.0.19
版本总结:
    - comic基本结构差不多了。comic,拆成2个表,并增加了数据清洗。
        - dbUPsertBatch加了 Omit(clause.Associations) -》 为了不更新关联表，只更新主表
        - dbUPsertBatch加了 Slect(*) -》 搭配Omit,才能 不更新关联表，只更新主表

核心改动:
     1. 搭框架相关：
    先弄简单的：
    - models 里 comic_my表 copy了子表结构，从 comic_spider
        - mian.go 迁移表时，加上了子表
    - comic子表 插入之前，调用下自己实现的数据处理接口：
        - 去除空格
        - 繁体转简体
    - 把 doc/F12找到的json里websiteId 都改为1-未分类，防止因为外键报错，原来是10，导致测试时，每次重新建表后，都要重新加上id=10 的website √

    再弄复杂的：
    - 拆分数据库表，把经常更新的数据，拆成2个表 
        - 拆分comic表，把打分、点击率单独拿出来，因为可能经常更新。这样需要联表查询，这块代码需要改下
            - 插入默认数据要改
            - 插入默认数据-更新的列 需要改 
            - 爬取映射mapping，需要改 
            - 报错的地方，需要改
        - 拆分comic表，就需要考虑 联表操作(增删改查)
        - comic-spider + comic-my都要改
    - comic有2个字段，考虑要不要有，要不要放 经常更新的表里  --> 考虑，cursor推荐。cursor是目前，我用过的最好的AI编程(用过 chatgpt网页版-只能问答、通义灵码、codeGeex-好、github copllot-一般，不能改项目代码) -------------------------？？？？？？？？？？？？？？？
        - 要创建chapter models库。先简单写，后期再补 √ 
            - 迁移表要加上 √ 
        - AI推荐: comic_stats 里放 latest_chapter_id(外键)、total_chapter(总章节数)、latest_chapter_num(序号一章序号)、latest_chapter_name(最后一章名称)。latest_chapter_id(外键) 是冗余的一个，就是为了方便  √
        - AI推荐: comic 里加 latest_chapter_id(外键)，这个外键属于业务场景，不属于统计场景。2个表都更新 √
        - comic_spider_stats, comic_my_stats 表, comic_xx_stats 里加 latest_chapter_id(外键)，在这里属于冗余一个，就是为了查询方便，不用join影响性能  
        - total_chapter // 总章节数量 2个表都更新
            - 可以空？ 不能，让默认为0
                - 还没爬章节场景，可以空，爬完章节，再更新
                - 需要加默认为0吗？默认为0，避免 NULL 值带来的查询和统计复杂度
            - 需要更新：
                - 插入默认数据 x
                - 插入默认数据-更新的列 x
                - 数据迁移 x
                - mapping爬取结构 √ 能爬到
                - 插入stats表数据时的 新增字段赋值 √
                - 插入stats表数据时的 updateCol √ 
            - 2个表都更新 √
        - stats表里 LatestChapter 结构体 考虑删除，想着都是冗余了，如果用不到就先删除，用到再说 √
        - latest_chapter_name // 最后一章名称 
            - 可爬  
            - 插入默认数据 x
            - 插入默认数据-更新的列 x
            - 数据迁移 x
            - mapping爬取结构 √ 
            - 插入stats表数据时的 新增字段赋值 √
            - 插入stats表数据时的 updateCol √
        - comic 少total_chapter字段，不考虑这个字段，最后一章的name_id就是总数，放到 频繁更新的表里, 经常用often 英文 . 考虑做成外键 -》关联chapter_id √
            - 之前是这么考虑的。后来想的是：
                - 还是要total_chapter,为了查询简单，不用join 联表查询了
                - 还要要有外键id chapterId。本来在stats表，这个外键id没用，但是还是加上了，以防万一
        - latest_chapter_id(外键), 外键id考虑可以null,因为爬取的时候，可能没有这个章节。爬完章节再来更新这个字段 √
        - comic_spider、comic_my 表，都要更新 √
        - comic_spider_stats, comic_my_stats 表，都要更新 √
    - comic 少字段 最后更新时间 last_update_date = 最新章节发布时间 ,可以同步最新章节的-发布时间，用于查询：本周有哪些更新
        有2种方式：-- 不用
        - = 最新章节的发布时间，用于查询：对于这本书，官方本周有哪些更新
        - = 表 update_at 字段，表示，我这周更新了哪些章节。用于查询，对于这本书，我 主动 本周有哪些更新
        用以下方式实现：
        - comic 加字段：最后章节更新时间 lastest_chapter_release_date
            - 可爬？是的 √
            - 插入默认数据 x 不用改
            - 插入默认数据-更新的列 x 不用改
            - 数据迁移 x 不用改
            - mapping爬取结构 √ 
            - models 表结构，假如新增字段 (comic_spider comic_my都要改 ) √
            - 插入stats表数据时的 新增字段赋值 √
            - 插入stats表数据时的 updateCol √
    
    2. 校验相关:
    - 插入前数据+业务校验，
        数据校验- 对于通用的东西：
        - string类型 √
            - 去除空格 TrimSpaces()
            - 繁体转简体 Trad2Simple()
        - 注意：插入数据，前调用下 √
        - 注意：comic_spider comic_my都要改 √

        业务校验-对于业务的数据 BusinessDataCleanObj() 
        - int类型
            - 比如star，最高10.0，如果超出，就按0算 √
        - 注意：插入数据，前调用下 √
        - 注意：comic_spider comic_my都要改 √
    - comic 实现一个数据清洗接口，（数据清洗自动实现，去除空格、繁体转简体，自动去除协议头：http/https，超出范围，自动置为某个值）√
        实现思路：√
        - sider.go里加一个统一接口: 数据清洗 DataClean(), 包含 TrimSpaces() Trad2Simple() BusinessDataClean() √
        - comic_spider comic_my 结构体里，加一个数据清洗接口，比如：func (c *Comic) BusinessDataCleanObj() √
        - 注意：comic_spider comic_my都要改  √
    - 数据清洗的时候，如果有https或者http头，自动删除。√
非核心改动:
    -

# v0.0.0.20
版本总结:
    - 打算实现通过html爬comic,l奥报错，先上传
    - 实现 comic_spider 通过html方式爬取 ？？
    - 考虑爬取，插入逻辑。思考：并发爬，插入sql时，只能执行1个。并发爬取结果，存到一个数组里。插入时,从插入库中取，保证数据完整，连续，美观，不混乱
        - 总结就是：爬取并发爬，插入串行插入。没实现？？？
        - 因为我想的gorm 批量插入也是很快的
    - 爬chapter,并完善chapter表, 完善chapter处理逻辑 。没实现？？？？ 
    - chapter也分 chapter_spider chapter_my 2个表。没实现？？？？？

核心改动:
    1. 搭框架相关：
    先弄简单的：
    - 1 先创建chapter表
        - chapter_spider表
            - 迁移表，需要改 √
            - 可爬？是的 
            - 插入默认数据 x 不用改
            - 插入默认数据-更新的列 x 不用改
            - 建立表结构 √ 
            - 实现数据清洗各种接口 √
            - mapping爬取结构 需要改 √
            - models 表结构，假如新增字段 x 不用改,没新增
        - chapter_my表
            - 迁移表，需要改 √
            - 可爬？是的 
            - 插入默认数据 x 不用改
            - 插入默认数据-更新的列 x 不用改
            - 建立表结构 √
            - 实现数据清洗各种接口 √
            - mapping爬取结构 不用，用spider爬到的就行 x 不用改
            - models 表结构，假如新增字段  x 不用改,没新增
    - comci end字段，作为冗余字段，就是纯显示用，，已经关联外键id,- processId。 ？？？？？？？？？？？？
        - end 字段 json数据里爬不到
        - end 类型要 从 bool 改成 int类型.  comic_spider comic_my 2表都要改
        - end 和爬取的processId之间逻辑，还要优化下。-》在业务数据清理的地方. comic_spider comic_my 2表都要改
            processId是人为传的，
            - processId = 1, 表示待分类，end应该 == processId = 1
                - end 爬不到, == 1
                - end 爬到了， == 爬到的值 (2或3)
            - processId = 2, 表示连载，  end应该 == processId = 2
            - processId = 3, 表示完结，  end应该 == processId = 3
        - 应该用哪个字段，作为真实显示？ end √
        - 先这样实现，后续 插入 comic_my数据库，人为严格判断 √

    - website 加几列数据，
        - book_can_spider_type: 可选:json/html/both/bothno  book可以的爬取方式:json还是html还是2者都行,还是都不行，有的网站通过html找不到有用数据
            toptoon-tw: json -> 只有json方式
            - 可爬？x 
            - 插入默认数据 要改 √
            - 插入默认数据-更新的列 要改 √
            - 数据迁移  x 不用改
            - mapping爬取结构  x 不用改
            - models 表结构，假如新增字段 要改 √
        - chapter_can_spider_type: json/html/both/bothno  chapter可以的爬取方式:json还是html还是2者都行,还是都不行，有的网站通过html找不到有用数据
            toptoon-tw: bothno -> 都不行
            - 可爬？x 
            - 插入默认数据 要改 √
            - 插入默认数据-更新的列 要改 √
            - 数据迁移  x 不用改
            - mapping爬取结构  x 不用改
            - models 表结构，假如新增字段 要改 √
        - book_spider_req_body_eg_server_filepath 爬取book时,请求体内容示例,后台服务器路径,不要求必须有值,字符串空也可以
            如：json -> 给出请求体，传什么数据
            如: html -> 给出请求体，传什么数据
            如："爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html"
            - 可爬？x 
            - 插入默认数据 要改 √
            - 插入默认数据-更新的列 要改 √
            - 数据迁移  x 不用改
            - mapping爬取结构  x 不用改
            - models 表结构，假如新增字段 要改 √
        - chapter_spider_req_body_eg_server_filepath 爬取chapter时,请求体内容示例,后台服务器路径,不要求必须有值,字符串空也可以
            如：json -> 给出请求体，传什么数据
            如: html -> 给出请求体，传什么数据
            如："爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html"
            - 可爬？x 
            - 插入默认数据 要改 √
            - 插入默认数据-更新的列 要改 √
            - 数据迁移  x 不用改
            - mapping爬取结构  x 不用改
            - models 表结构，假如新增字段 要改 √
        - 打分方式：star_type string
            - 没有参考的网站，或者网站评分系统无参考价值，就自己主观打分。 
            - my -> 表示我自己打的
            - copy_toptoon -> 表示从toptoon网站爬取的
            - copy_18to -> 表示从18to网站爬取的
            - 可爬？x 
            - 插入默认数据 要改 √
            - 插入默认数据-更新的列 要改 √
            - 数据迁移  x 不用改
            - mapping爬取结构  x 不用改
            - models 表结构，假如新增字段 要改 √

    再弄复杂的

    2. 校验相关:
    3. 爬取相关:
    4. 其他相关

非核心改动:
    -
# v0.0.0.20.2 临时提交
- 简单实现 comic_spider 通过html方式爬取1个页面,多页面无法实现，所以先提交。要改 GetAllObjFromOneHtmlPageUseCollyByMapping() 方法结构 

# v0.0.0.20.3 临时提交
    - 实现chapter_spider 逻辑

# v0.0.0.21
版本总结:
    - 实现 comic_spider 通过html方式爬取 ？？
    - 考虑爬取，插入逻辑。思考：并发爬，插入sql时，只能执行1个。并发爬取结果，存到一个数组里。插入时,从插入库中取，保证数据完整，连续，美观，不混乱
        - 总结就是：爬取并发爬，插入串行插入。没实现？？？
        - 因为我想的gorm 批量插入也是很快的
    - 爬chapter,并完善chapter表, 完善chapter处理逻辑 。没实现？？？？ 
    - chapter也分 chapter_spider chapter_my 2个表。没实现？？？？？

核心改动:
    1. 搭框架相关：
    先弄简单的：
    - 先用kxmanhua，实现爬章节 √
        - 可获取book信息：考虑想更新哪列，就实现更新哪列。比如（cover_url_api_path 已经有了, 能爬到，但是不需要更新，走的是update方法）√
            - 评分 不能爬，不管
            - 作者 √ 用 / 分开 能爬。
            - 最后一章id √ 但需要插入完 chapters 之后，才能获取到，并且更新 book
            - 最后一章name 不能爬，不管
            - 简介-短 x 不能爬，不管
            - 简介-长 √ 能爬
            - 作者拼接 √ 自己判断
            - 作者拼接类型 √ 自己判断

        - 可获取章节信息：(除了id，所有内容都要生成？我的思路，能爬到哪些，就set哪些，爬不到的就默认处理。最后通过DataClean()清洗下) √
            - id x 不用，自行生成
            - chapter_num x 没有，自行生成。必须要，唯一索引 √
            - chapter_sub_num x 没有，自行生成。有没有都行。必须要唯一索引 √
            - chapter_real_sort_num x 没有，自行生成。有没有都行 ？？
            - name 能爬。可能需要截取。必须要 √
            - url_api_path。能爬。需要考虑截取http头+域名。必须要 √
            - release_date。不能爬。就按默认
            - parent_id。不能爬。但是需要，因为是唯一索引    √
        - 休刊公告应该怎么着？要是有多个休刊公告咋办？ - 不要了 √
    - 实现 chapter_spider 通过html方式爬取 √
        - kxmanhua 
            - spider_chapter 需要修改内容：√
                - ChapterRealSortNum: = chapterNum -> 数据清洗的时候弄 √
                - chapterSpider，加上中文字符转英文的 func √ 
                - chapterSpider，mapping 去掉字符串 的 不需要的内容： ♥  ---- 考虑应该在数据清洗弄，还是mappping弄？因为解耦考虑：一个方法就干一个事，mapping就管爬，你让它干那么复杂的事干什么！ √
                - parentId 啥时候加？爬完之后，根据方法上下文传参改，在数据清洗之前 √
                - chapter_spider/ chapter_my 加上字段： √
                    - spiderEndStatus
                    - downloadEndStatus
                    - uploadAwsEndStatus
                    - uploadBaiduEndStatus
    - 问题：
    - 1. 通过html插入，第一次插入 comic + 子表都成功了，第二次更新就不行 √，需要先查，再插入子表 √

    - 实现chapter_content 爬取逻辑
    - book page1 测试完整逻辑，爬取book，爬取chapter，爬取chapter_content，下载chapter_content，上传chapter_content到aws，上传chapter_content到baidu
    - 整个网站，测试完整逻辑，爬取book，爬取chapter，爬取chapter_content，下载chapter_content，上传chapter_content到aws，上传chapter_content到baidu
    - 找一个网站：不带水印，比较好爬 √
        - 鸟鸟韩漫 nnhanman.xyz

    再弄复杂的

    2. 校验相关:
    3. 爬取相关:
    通用性相关：
        - 中文符号 如小括号，中括号，逗号等，转成英文 comic_spider comic_my都要改 √
    4. 其他相关

非核心改动:
    -


# ------------------------------------------ 未解决问题如下： -------------- 钱钱钱，挣钱买摩托
思路：能简单，别复杂。不能老是学自己不会的，要把会的完全应用

要解决：
    - 应该先去爬网站？还是整理基础结构？还是校验数据，确保安全完整不污染
    建议：① 搭框架 → ② 建校验 → ③ 小规模爬取测试 → ④ 扩展到完整爬虫

    1. 搭框架相关：
    先弄简单的：
    ------------------------------- 请求传参 websiteId 等外键id,赋值给comic 对
    - website 加几列数据
        - 这个网站 图片 能存多久/多久就会失效，因为有的网站，1天之后链接就会变，防爬机制
        - chapter_content 是否有水印

    - 加了某些东西之后，插入默认数据，插入默认数据-更新的列，爬取映射mapping, 更新的列，/ 数据迁移。这几个参考点，需要同步修改

    再弄复杂的：

    2. 校验相关:
    - 数据清洗分2个方面：1. 前端传参、方法间传参，数据清洗-属于前端编程人员操作 2. 插入db前数据清洗(这个实现了一部分) -》 属于后端编程人员操作
    - 给所有报错，给出推断原因，让用户自己去简单排查。缩短排查时间

    3. 爬取相关:

    通用性相关:
    - 把mapping映射关系，写出json文件，类似配置文件的方式，通用，以后不用改代码，直接改json文件即可 -》 参考 spider.go -> ComicMappingForSpiderToptoonByJSON 这个变量
    - 实现：通过配置文件，或者键值对变量，控制：爬书的时候处理哪些字段，爬章节的时候，更新哪些父表-book表哪些字段。做成通用的框架
    - 考虑把某个网站的爬取算法，放到一起，比如一个文件，里。方便归类，我喜欢归类清楚的东西
    - 多个项目放到一个项目里，数据库，表、命名，文件结构都容易冲突，考虑如何实现

    爬取后处理 -》 转到我用的表
    - comic_spider -> 能转成 comic_my,不更新 cover_save_path_api_path 字段
        - 考虑cover_save_path_api_path字段,如果已经有了，要是不小心，传空咋办，就会替换成空了.已解决，comic_spider表，不带该字段
        - cover_save_path)api 是最关键的字段，没有它，所有业务都不行
        - !!!重要：comic_spider不带 cover_save_path 字段，comic_my带 cover_save_path 字段，这样comic_spider不管传啥，都不会影响该字段。并且测试，comic_my的 upsert方法，cover_save_path有值，如果传参不带 cover_save_path字段，且update 需要更新此字段，会发生什么琴科给
    - 爬取状态：
        - spider_end_status int 有几种 0 1 2 ，爬取完chapter后，并且爬完chapter_content后，end-》才能变为1
        - download_end_status int 有几种 0 1 2 ，需要spider_end_status=1才能用，把chapterContent下载完后，才能 = 1
        - upload_aws_end_status int 有几种 0 1 2 ，需要download_end_status=1才能用，把chapterContent 所有上传到aws后，才能 = 1

    4. 其他相关
    - 单元测试用例
    - 画完整逻辑图，很长时间后，一看就知道逻辑了
    - 想想自己工作时候，有什么娱乐的东西，无聊的话




------------------------- 解决完再上传



# 待办
- comic估计是没有打分非常好的网站, 就自己打分吧
    - level1: 自己主观打分
    - level2: 根据公式打分
    - level3: 根据公式+用户权重打分
- 梳理项目代码，画图，要很久以后，还一目了然
- 最完美的通用模板，应该是所有用的东西都是配置化的，比如新增列，新加列，都是通过配置实现，根本不用动源代码。现实吗？
- comic表加上时间，创建时间，更新时间，是否删除标志 √ 没测
- 更新至，内容包含 '最终话'、'完结' 认为是完结了 √ 没测
- 人气，字符串里可能没带单位。如何处理？ 如：'人气：5555' √ 没测
- 优化思路：更新过程，应该是有什么字段，只更新你获取到的，没获取到的，不要乱更新！！！！
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
- 考虑，程序里，批量插入，尽量用指针，比如comicArr chapterArr里尽量存指针，因为可能回涉及很多，要修改内容。如果不用指针，就修改不了原数据了