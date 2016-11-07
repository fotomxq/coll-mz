//获取任意站点整站分页数据
//可复用，获取任意站点的，某一类页面的所有页面内容
//注意，某些极特殊页面无法使用该模块，使用前需要给定参数测试一下
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fotomxq/ftmp-libs"
)

//结构
type GetPageData struct {
	//页码页面的URL前缀部分
	pageURLPrefix string
	//后缀部分
	pageURLSuffix string
	//在页码页面上查询出所有可选的子页面，jquery.find字符串形式
	pageFindStr string
	//以页码查询子节点为基准，查找标题
	//如果该部分留空，则忽略
	//必须确保childrenTitleFindStr有内容，否则数据无法采集
	pageTitleFindStr string
	//是否直接从该页面上采集数据
	//true - 不递归查询子页面，直接将找出来的数据存档
	//false - 继续递归查询子页面内容，根据指定jquery.find内容保存
	isOnlyPage bool
	//子页面查询内容，jquery.find字符串形式
	childrenFindStr string
	//子页面的标题，jquery.find字符串形式
	childrenTitleFindStr string
	//在页码中查询最大页数的查询路径，jquery.find形式
	//注意，必须确保可通过该对象的jquery.html()获取数据
	//如果留空，则从最开始一直读取到没法获取页面为止
	pageMaxFindStr string
	//从第几页开始
	pageStart int
	//第几页强行停止
	pageEnd int
	//开始和结束是否为正序
	isPageSortAsc bool
	//数据库的类型
	dbType string
	//数据库的DNS
	dbDNS string
	//日志对象
	log ftmplibs.Log
	//数据库对象
	db *sql.DB
}

//建立后初始化该结构
//log 日志对象
//configData 从[config].json文件内读出的配置数据，必须确保采用data-default.json的格式
func (getPageData *GetPageData) Create(log ftmplibs.Log,configData map[string]interface{}) bool{
	//获取日志对象
	getPageData.log = log
	//返回建立成功
	return true
}

//开始运行采集
func (getPageData *GetPageData) Run() (bool,error){
	defer db.Close()
	return true,nil
}

//检查数据是否过时
func (getPageData *GetPageData) check() (bool,error){
	return true,nil
}

//连接数据库
func (getPageData *GetPageData) connectDB() error{
	getPageData.db, err = sql.Open(config.Data["databaseType"].(string),config.Data["databaseDNS"].(string))
	if err != nil{
		getPageData.log.AddLog("发生一个错误:")
		getPageData.log.AddErrorLog(err)
	}
	return err
}