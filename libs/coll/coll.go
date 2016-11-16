//收集数据处理器
package coll

import (
	"github.com/fotomxq/coll-mz/libs/core"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
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
	//全局列表设定
	CollList map[string]string
	//http操作句柄
	simhttp core.SimpleHttp
}

//初始化结构
func (c *Coll) Create(dataSrc string)(bool,error){
	//初始化列表
	c.CollList = map[string]string{
		"jiandan" : "煎蛋网",
		"xiuren" : "秀人网",
	}
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
	if c.file.IsFile(dataDatabaseSrc) == false{
		b,err = c.file.CopyFile("./content/database-default.sqlite",dataDatabaseSrc)
		if err != nil || b == false{
			return b,err
		}
	}
	//创建日志结构
	c.log.SetDirSrc(dataLogSrc)
	//连接到数据库，直接返回结果
	return c.connectDB("sqlite3",dataDatabaseSrc)
}

//启动所有脚本
func (c *Coll) RunAll()(bool,error){
	c.SendLog("启动了所有采集程序，可能时间较长，请耐心等待...")
	for k,_ := range c.CollList{
		b,err := c.Run(k)
		if err != nil || b == false{
			return b,err
		}
	}
	c.SendLog("所有采集程序运行结束，现在可以快乐的访问内容了。")
	return true,nil
}

//启动某个脚本
func (c *Coll) Run(wt string) (bool,error){
	if c.CollList[wt] != nil{
		c.SendLog("开始启动" + c.CollList[wt] + "采集程序...")
	}
	var b bool
	switch wt{
		case "jiandan":
			b,err = CollJiandan()
			break
		case "xiuren":
			b,err = CollXiuren()
			break
		default:
			c.SendLog("指定的采集程序不存在。")
			return false,nil
			break
	}
	c.SendLog("该采集程序运行结束。")
	return b,err
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
	return false,nil
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

//发送日志
func (c *Coll) SendLog(str string){
	c.log.AddLog(str)
}

//发送错误日志
func (c *Coll) SendErrorLog(err error){
	c.log.AddErrorLog(err)
}

//根据URL构建文件路径
//过程中会自动创建需要的目录
//src - 父级目录
//url - URL地址
//name - 指定文件名称
//返回 文件路径,错误
func (c *Coll) CreateFileSrc(src string,url string,name string) (string,error){
	//尝试解析URL文件名称
	urls := c.simhttp.GetURLNameType(url)
	if urls == nil{
		return "",nil
	}
	if urls[2] == ""{
		return "",nil
	}
	//构建存储目录
	dirSrc,err := c.CreateDirSrc(src)
	if err != nil || dirSrc == ""{
		return "",err
	}
	//路径分隔符
	sep := c.file.GetPathSep()
	//根据目录路径生成文件路径
	fileSrc := dirSrc + sep + name + urls[2]
	return fileSrc
}

//创建新的目录
//src - 父级目录
func (c *Coll) CreateDirSrc(src string)(string,error){
	//路径分隔符
	sep := c.file.GetPathSep()
	//构建子目录路径
	var dataSrc string
	var dataSrcF string = src + sep + c.GetNowDateYM() + sep + c.GetNowDateD() + sep
	for i := 1 ; i < 100 ; i ++{
		//构建路径
		dataSrc = dataSrcF + strconv.Itoa(i)
		//判断该目录是否存在
		if c.file.IsFolder(dataSrc) {
			//如果存在，则判断该文件夹下文件数量是否超过100，超过则进入下一个循环
			max,err := c.file.GetFileListCount(dataSrc)
			if err != nil{
				return "",err
			}
			if max > 100{
				dataSrc = ""
				continue
			}else{
				return dataSrc,nil
			}
		}else{
			b,err := c.file.CreateDir(dataSrc)
			if b == false || err != nil{
				return "",err
			}
			return dataSrc,nil
		}
	}
	//如果今天1-100个目录全满，则创建返回创建失败
	if dataSrc == ""{
		return "",nil
	}
	//其他逻辑错误，到达这里，直接返回失败
	return "",nil
}

//获取今天年份
func (c *Coll) GetNowDateYM() string {
	t := time.Now()
	return t.Format("200601")
}

//获取今天月日
func (c *Coll) GetNowDateD() string {
	t := time.Now()
	return t.Format("02")
}