//获取任意站点整站分页数据
//可复用，获取任意站点的，某一类页面的所有页面内容
//注意，某些非通用页面无法使用该模块run功能，但可以使用模块内其他方法辅助实现采集目标
package collmzLibs

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fotomxq/ftmp-libs"
)

//结构
type GetPageData struct {
	//数据库的类型
	dbType string
	//数据库的DNS
	dbDNS string
	//数据库对象
	db *sql.DB
	//数据库连接失败error
	dbErr error
}

//建立后初始化该结构
//configData 从[config].json文件内读出的配置数据，必须确保采用data-default.json的格式
func (getPageData *GetPageData) Create(configData map[string]interface{}) (bool,error){
	//保存基本配置数据
	getPageData.dbType = configData["databaseType"].(string)
	getPageData.dbDNS = configData["databaseDNS"].(string)
	//检查并返回结果
	return getPageData.check()
}

//检查数据是否过时、存在错误
func (getPageData *GetPageData) check() (bool,error){
	//尝试连接数据库
	getPageData.ConnectDB()
	if getPageData.dbErr != nil{
		return false,getPageData.dbErr
	}
	return true,nil
}

//连接数据库
func (getPageData *GetPageData) ConnectDB(){
	getPageData.db, getPageData.dbErr = sql.Open(getPageData.dbType,getPageData.dbDNS)
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
	return http.Save(src)
}

//关闭数据库
func (getPageData *GetPageData) CloseDB(){
	getPageData.db.Close()
}