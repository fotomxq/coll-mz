//收集数据处理器
package coll

import (
	"../core"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//通用错误
var err error

//收集器类
type Coll struct {
	//数据存储路径
	dataSrc string
	//文件句柄
	file core.FileOperate
	//数据库句柄
	db *sql.DB
	//数据库错误句柄
	dbErr error
	//日志句柄
	log core.Log
	//收集路径
	dataCollSrc string
	//收集缓冲路径
	dataCollCacheSrc string
}

func (c *Coll) Create(dataSrc string)(bool,error){
	//检查目录是否存在，不存在则创建
	c.dataSrc = dataSrc
	b,err := c.file.CreateDir(dataSrc)
	if err != nil || b == false{
		return b,err
	}
	//构建完成后，在数据目录下创建log、coll等子目录及文件
	dataLogSrc := dataSrc + "/log"
	c.dataCollSrc = dataSrc + "/coll"
	c.dataCollCacheSrc = dataSrc + "/coll-cache"
	dataDatabaseDirSrc := dataSrc + "/database"
	dataDatabaseSrc := dataDatabaseDirSrc + "/database.sqlite"
	b,err = c.file.CreateDir(dataLogSrc)
	if err != nil || b == false{
		return b,err
	}
	b,err = c.file.CreateDir(c.dataCollSrc)
	if err != nil || b == false{
		return b,err
	}
	b,err = c.file.CreateDir(c.dataCollCacheSrc)
	if err != nil || b == false{
		return b,err
	}
	b,err = c.file.CreateDir(dataDatabaseDirSrc)
	if err != nil || b == false{
		return b,err
	}
	b,err = c.file.CopyFile("./content/database-default.sqlite",dataDatabaseSrc)
	//创建日志结构
	c.log.SetDirSrc(dataLogSrc)
	//连接到数据库，直接返回结果
	return c.connectDB("sqlite3",dataDatabaseSrc)
}

//连接到数据库
func (c *Coll) connectDB(t string,dns string)(bool,error){
	c.db,c.dbErr = sql.Open(t,dns)
	if c.dbErr != nil{
		return false,c.dbErr
	}
	return true,nil
}

//关闭数据库
func (c *Coll) CloseDB(){
	err = c.db.Close()
	if err != nil{
		c.log.AddErrorLog(err)
	}
}

//自动将文件存入数据集
//整合检查、添加、保存功能
func (c *Coll) AutoAddData(source string,url string,name string) (bool,error){

}

//将数据采集到数据库
func (c *Coll) AddData(sha1 string,source string,url string,name string,t string,size int)(bool,error){
	//检查数据库是否连接
	if c.dbErr != nil{
		return false,c.dbErr
	}
	//获取当前时间
	datatimeObj := time.Now()
	//构建sql
	query := "insert into `coll`(`sha1`,`source`,`url`,`name`,`type`,`size`,`coll_time`) values(?,?,?,?,?,?,?)"
	stmt,stmtErr := c.db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	res,resErr := stmt.Exec(sha1,source,url,name,t,size,datatimeObj.Format("2006-01-02 15:04:05"))
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

//检查数据是否存在
func (c *Coll) CheckData(sha1 string) (bool,error){
	//检查数据库是否连接
	if c.dbErr != nil{
		return false,c.dbErr
	}
	//构建查询语句
	var query string = "select id from `coll` where `sha1` = ?"
	//开始查询
	stmt,stmtErr := c.db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	rows,resErr := stmt.Query(sha1)
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

//保存URL到文件
func (c *Coll) SaveUrl(url string,src string) (bool,error){
	http := new(core.SimpleHttp)
	http.SetSendUrl(url)
	return http.Save(src)
}