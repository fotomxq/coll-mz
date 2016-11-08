//获取任意站点整站分页数据
//可复用，获取任意站点的，某一类页面的所有页面内容
//注意，某些非通用页面无法使用该模块run功能，但可以使用模块内其他方法辅助实现采集目标
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fotomxq/ftmp-libs"
	"github.com/PuerkitoBio/goquery"
)

//结构
type GetPageData struct {
	//数据库的类型
	dbType string
	//数据库的DNS
	dbDNS string
	//日志对象
	log ftmplibs.Log
	//数据库对象
	db *sql.DB
	//数据库连接失败error
	dbErr error
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
//如果该方法无法使用，请直接构建新的方法，并选择性使用模块内其他方法
func (getPageData *GetPageData) Run() (bool,error){
	//连接数据库
	if getPageData.dbErr != nil{
		getPageData.log.AddLog("发生一个错误:")
		getPageData.log.AddErrorLog(err)
	}
	defer getPageData.db.Close()
	//返回成功
	return true,nil
}

//检查数据是否过时
func (getPageData *GetPageData) check() (bool,error){
	return true,nil
}

//连接数据库
func (getPageData *GetPageData) ConnectDB(){
	getPageData.db, getPageData.dbErr = sql.Open(config.Data["databaseType"].(string),config.Data["databaseDNS"].(string))
}

//获取页面所有子页面链接
func (getPageData *GetPageData) GetPageList(url string,search string)([][]string,error){
	//定义返回清单
	var list [][]string
	//获取页面
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return list,err
	}
	node := doc.Find(search)
	//遍历所有子链接
	for i := range node.Nodes {
		//获取节点数据
		thisNode := node.Eq(i).Children()
		thisHref, hrefExist := thisNode.Attr("href")
		thisTitile,titleExist := thisNode.Attr("title")
		if hrefExist == false || titleExist == false {
			continue
		}
		//检查数据库中，是否存在该数据
		checkBool,checkErr := getPageData.CheckRepeat(thisHref,thisTitile)
		if checkErr != nil{
			return list,checkErr
		}
		if checkBool == true{
			continue
		}
		//如果该链接标题不存在，将链接和标题加入到组
		child := []string{thisHref,thisTitile}
		append(list,child)
	}
	return list,nil
}

//保存某个子页面
func (getPageData *GetPageData) SavePageChildren(url string,title string,search string)(bool,error){
	//读取页面数据
	doc, docErr := goquery.NewDocument(url)
	if docErr != nil {
		return false,docErr
	}
	//检查是否存在对象
	node := doc.Find(search)
	for i := range node.Nodes {

	}
	return true,nil
}

//检查数据是否已采集
func (getPageData *GetPageData) CheckRepeat(url string,title string)(bool,error){
	//检查是否连接到数据库
	if getPageData.dbErr != nil{
		return false,getPageData.dbErr
	}
	//构建查询语句
	var query string = "select id from `loadurl` where `url` = ? and `title` = ?"
	//开始查询
	stmt,stmtErr := getPageData.db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	rows,resErr := stmt.Query(url,title)
	if resErr != nil{
		return false,resErr
	}
	rows.Next()
	var id int
	scanErr := rows.Scan(&id)
	if scanErr != nil{
		return false,nil
	}
	if id > 0{
		return true,nil
	}
	//返回
	return false,nil
}

//向数据库添加新的数据
func (getPageData *GetPageData) AddNewData(url string,title string) (bool,error){
	//检查是否连接到数据库
	if getPageData.dbErr != nil{
		return false,getPageData.dbErr
	}
	//构建添加语句
	var query string = "insert into `loadurl`(`url`,`title`) values(?,?)"
	stmt,stmtErr := getPageData.db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	res,resErr := stmt.Exec(url,title)
	if resErr != nil {
		return false,resErr
	}
	newID,newIDErr := res.LastInsertId()
	if newIDErr != nil {
		return false,newIDErr
	}
	if newID > 0{
		return true,nil
	}
	return false,nil
}

//向数据库添加视频记录
func (getPageData *GetPageData) AddVideoData(url string) (bool,error){
	//检查是否连接到数据库
	if getPageData.dbErr != nil{
		return false,getPageData.dbErr
	}
	var query string = "insert into `video`(`url`) values(?)"
	stmt,stmtErr := getPageData.db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	res,resErr := stmt.Exec(url)
	if resErr != nil {
		return false,resErr
	}
	newID,newIDErr := res.LastInsertId()
	if newIDErr != nil {
		return false,newIDErr
	}
	if newID > 0{
		return true,nil
	}
	return false,nil
}

//将页面存到文件
//该方法可用于保存图像等数据
func (getPageData *GetPageData) SaveUrlToFile(url string,src string) (bool,error){
	http := new(ftmplibs.SimpleHttp)
	http.SetSendUrl(url)
	c,err := http.Get()
	if err != nil{
		return false,err
	}
	t,tErr := ftmplibs.WriteFile(src,c)
	if tErr != nil{
		return false,nil
	}
	if t != false{
		return false,nil
	}
	return true,nil
}