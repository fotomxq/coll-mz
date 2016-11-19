//收集数据处理器
package coll

import (
	"database/sql"
	"github.com/fotomxq/coll-mz/libs/core"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

//通用错误
var err error

//通用coll类句柄
var CollPg *Coll

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
	//可读取日志存储路径
	logReadSrc string
	//全局列表设定
	CollList map[string]string
	//http操作句柄
	simhttp core.SimpleHttp
	//过滤器句柄
	ms core.MatchString
}

//初始化结构
func (c *Coll) Create(dataSrc string) (bool, error) {
	//初始化列表
	c.CollList = map[string]string{
		"jiandan": "煎蛋网",
		"xiuren":  "秀人网",
	}
	//检查目录是否存在，不存在则创建
	c.dataSrc = dataSrc
	b, err := c.file.CreateDir(dataSrc)
	if err != nil || b == false {
		return b, err
	}
	//构建完成后，在数据目录下创建log、coll等子目录及文件
	sep := c.file.GetPathSep()
	dataLogSrc := dataSrc + sep + "log"
	c.dataCollSrc = dataSrc + sep + "coll"
	c.dataCollCacheSrc = dataSrc + sep + "coll-cache"
	dataDatabaseDirSrc := dataSrc + sep + "database"
	dataDatabaseSrc := dataDatabaseDirSrc + sep + "database.sqlite"
	c.logReadSrc = dataDatabaseDirSrc + sep + "log-read.log"
	b, err = c.file.CreateDir(dataLogSrc)
	if err != nil || b == false {
		return b, err
	}
	b, err = c.file.CreateDir(c.dataCollSrc)
	if err != nil || b == false {
		return b, err
	}
	b, err = c.file.CreateDir(c.dataCollCacheSrc)
	if err != nil || b == false {
		return b, err
	}
	b, err = c.file.CreateDir(dataDatabaseDirSrc)
	if err != nil || b == false {
		return b, err
	}
	if c.file.IsFile(dataDatabaseSrc) == false {
		b, err = c.file.CopyFile("."+sep+"content"+sep+"database-default.sqlite", dataDatabaseSrc)
		if err != nil || b == false {
			return b, err
		}
	}
	if c.file.IsFile(c.logReadSrc) == false {
		b, err = c.file.CopyFile("."+sep+"content"+sep+"log-read.log", c.logReadSrc)
		if err != nil || b == false {
			return b, err
		}
	}
	//创建日志结构
	c.log.SetDirSrc(dataLogSrc)
	//连接到数据库，直接返回结果
	return c.connectDB("sqlite3", dataDatabaseSrc)
}

//启动所有脚本
func (c *Coll) RunAll() (bool, error) {
	c.SendLog("启动了所有采集程序，可能时间较长，请耐心等待...")
	for k, v := range c.CollList {
		b, err := c.Run(k)
		if err != nil {
			c.SendErrorLog(err)
			c.SendLog("发生一个错误，导致" + v + "采集器无法进行。")
			continue
		}
		if b == false {
			c.SendLog("因为某些未知的原因，导致" + v + "采集器无法进行。")
			continue
		}
	}
	c.SendLog("所有采集程序运行结束，现在可以快乐的访问内容了。")
	return true, nil
}

//启动某个脚本
func (c *Coll) Run(wt string) (bool, error) {
	if c.CollList[wt] != "" {
		c.SendLog("开始启动" + c.CollList[wt] + "采集程序...")
	}
	var b bool
	switch wt {
	case "jiandan":
		b, err = CollJiandan()
		break
	case "xiuren":
		b, err = CollXiuren()
		break
	default:
		c.SendLog("指定的采集程序不存在。")
		return false, nil
		break
	}
	c.SendLog("该采集程序运行结束。")
	return b, err
}

//连接到数据库
func (c *Coll) connectDB(t string, dns string) (bool, error) {
	c.db, c.dbErr = sql.Open(t, dns)
	if c.dbErr != nil {
		return false, c.dbErr
	}
	return true, nil
}

//关闭数据库
func (c *Coll) CloseDB() {
	err = c.db.Close()
	if err != nil {
		c.log.AddErrorLog(err)
	}
}

//自动将文件存入数据集
//整合检查、添加、保存功能
//source string 存储点识别码，如jiandan
//url string 文件URL地址
//name string 识别标题
//urlIsParent bool URL是否为合集，如果是合集，则需要自行将文件保存，但这里会做标记到数据库，以避免重复提交
//return string 反馈保存的文件路径
//return error 错误
func (c *Coll) AutoAddData(source string, url string, name string,urlIsParent bool) (string, error) {
	//根据url和name构建sha1值
	sha1 := c.ms.GetSha1(url+name)
	if sha1 == ""{
		c.SendLog("无法生成SHA1匹配码。")
		return "",nil
	}
	//检查数据是否已经存在
	checkBool,err := c.CheckData(sha1)
	if err != nil{
		c.SendLog("检查数据过程失败。")
		c.SendErrorLog(err)
		return "",err
	}
	//如果存在则返回
	if checkBool == true{
		return "",nil
	}
	//文件父级目录路径
	parentSrc := c.dataCollSrc + c.file.GetPathSep() + source
	//准备文件基本数据
	var fileSize int64
	var fileType string
	var newFileSrc string
	//如果URL不是合集，则尝试构建缓冲文件，并读取相关数据
	if urlIsParent == false{
		//将文件下载缓冲
		cacheName := c.NewCacheFile(url)
		if cacheName == ""{
			return "",nil
		}
		//获取文件大小
		fileSize = c.file.GetFileSize(cacheName)
		//获取文件格式
		fileNames,err := c.file.GetFileNames(cacheName)
		if err != nil{
			c.SendErrorLog(err)
			return "",err
		}
		fileType = fileNames["type"]
		//构建文件路径
		newFileSrc,err = c.CreateFileSrc(parentSrc,url,name)
		if err != nil{
			c.SendErrorLog(err)
			return "",err
		}
		//将缓冲文件保存到指定文件
		b,err := c.MoveCacheFile(cacheName,newFileSrc)
		if err != nil{
			c.SendErrorLog(err)
			return "",err
		}
		if b == false{
			c.SendLog("无法保存该文件。")
			return "",nil
		}
	}else{
		//如果URL是合集，则建立目录
		fileSize = 0
		fileType = "folder"
		newFileSrc,err = c.CreateFileSrc(parentSrc,"",name)
		if err != nil{
			c.SendErrorLog(err)
			c.SendLog("构建文件路径失败。")
			return "",err
		}
		//根据文件路径，创建文件夹
		b,err := c.file.CreateDir(newFileSrc)
		if err != nil{
			c.SendErrorLog(err)
			c.SendLog("创建文件夹失败。")
			return "",err
		}
		if b == false{
			c.SendLog("未知原因，创建文件夹失败。")
			return "",nil
		}
	}
	//向数据库添加新的数据
	addBool,err := c.AddData(sha1,newFileSrc,source,url,name,fileType,fileSize)
	if err != nil{
		c.SendLog("添加新的数据失败。")
		c.SendErrorLog(err)
		return "",err
	}
	if addBool == false{
		c.SendLog("因为未知原因，无法添加新的数据。")
		return "",nil
	}
	//返回结果
	return newFileSrc, nil
}

//建立缓冲数据
func (c *Coll) NewCacheFile(url string) string{
	//获取文件名称
	fileNames := c.simhttp.GetURLNameType(url)
	fileName := fileNames[0]
	//保存文件
	src := c.dataCollCacheSrc + c.file.GetPathSep() + fileName
	b,err := c.SaveUrl(url,src)
	if err != nil{
		c.SendErrorLog(err)
		c.SendLog("下载到缓冲文件时，发生了错误。")
		return ""
	}
	if b == false{
		c.SendLog("未知原因，缓冲失败。")
		return ""
	}
	return src
}

//移动缓冲文件
func (c *Coll) MoveCacheFile(name string,src string) (bool,error){
	cacheSrc := c.dataCollCacheSrc + c.file.GetPathSep() + name
	if c.file.IsFile(cacheSrc) == false{
		return false,nil
	}
	return c.file.EditFileName(cacheSrc,src)

}

//删除缓冲数据
func (c *Coll) DeleteCacheFile(name string) bool{
	src := c.dataCollCacheSrc + c.file.GetPathSep() + name
	return c.file.DeleteFile(src)
}

//清空缓冲数据
func (c *Coll) DeleteAllCacheFile(){
	files,err := c.file.GetFileList(c.dataCollCacheSrc)
	if err != nil{
		c.SendErrorLog(err)
	}
	for i := range files{
		fiSrc := c.dataCollCacheSrc + c.file.GetPathSep() + files[i]
		b := c.file.DeleteFile(fiSrc)
		if b == false{
			c.SendLog("无法删除缓冲文件。")
		}
	}
}

//将数据采集到数据库
func (c *Coll) AddData(sha1 string,src string, source string, url string, name string, t string, size int64) (bool, error) {
	//检查数据库是否连接
	if c.dbErr != nil {
		return false, c.dbErr
	}
	//获取当前时间
	datatimeObj := time.Now()
	//构建sql
	query := "insert into `coll`(`sha1`,`src`,`source`,`url`,`name`,`type`,`size`,`coll_time`) values(?,?,?,?,?,?,?)"
	stmt, stmtErr := c.db.Prepare(query)
	if stmtErr != nil {
		return false, stmtErr
	}
	res, resErr := stmt.Exec(sha1,src, source, url, name, t, size, datatimeObj.Format("2006-01-02 15:04:05"))
	if resErr != nil {
		return false, resErr
	}
	newID, newIDErr := res.LastInsertId()
	if newIDErr != nil {
		return false, newIDErr
	}
	if newID > 0 {
		return true, nil
	}
	return false, nil
}

//检查数据是否存在
func (c *Coll) CheckData(sha1 string) (bool, error) {
	//检查数据库是否连接
	if c.dbErr != nil {
		return false, c.dbErr
	}
	//构建查询语句
	var query string = "select id from `coll` where `sha1` = ?"
	//开始查询
	stmt, stmtErr := c.db.Prepare(query)
	if stmtErr != nil {
		return false, stmtErr
	}
	rows, resErr := stmt.Query(sha1)
	if resErr != nil {
		return false, resErr
	}
	rows.Next()
	var id int
	scanErr := rows.Scan(&id)
	if scanErr != nil {
		return false, nil
	}
	if id > 0 {
		return true, nil
	}
	//返回
	return false, nil
}

//保存URL到文件
func (c *Coll) SaveUrl(url string, src string) (bool, error) {
	http := new(core.SimpleHttp)
	http.SetSendUrl(url)
	return http.Save(src)
}

//发送日志
func (c *Coll) SendLog(str string) {
	//向日志句柄发送日志
	c.log.AddLog(str)
	//判断日志是否超出范围，超出则清空
	//默认范围为文件大小20kb
	var maxFileSize int64 = 20 * 1024
	fileSize := c.file.GetFileSize(c.logReadSrc)
	if fileSize > maxFileSize {
		var newC []byte
		_, _ = c.file.WriteFile(c.logReadSrc, newC)
	}
	//向可读取日志发送日志
	newStr := str + "<br/ >"
	strByte := []byte(newStr)
	_, _ = c.file.WriteFileForward(c.logReadSrc, strByte)
}

//发送错误日志
func (c *Coll) SendErrorLog(err error) {
	c.log.AddErrorLog(err)
}

//根据URL构建文件路径
//过程中会自动创建需要的目录
//src - 父级目录
//url - URL地址
//name - 指定文件名称
//返回 文件路径,错误
func (c *Coll) CreateFileSrc(src string, url string, name string) (string, error) {
	//文件类型
	var fileType string
	//如果URL给定空，则类型也为空
	if url == ""{
		fileType = ""
	}else{
		//尝试解析URL文件名称
		urls := c.simhttp.GetURLNameType(url)
		if urls == nil {
			return "", nil
		}
		if urls[2] == "" {
			return "", nil
		}
		fileType = "." + urls[2]
	}
	//构建存储目录
	dirSrc, err := c.CreateDirSrc(src)
	if err != nil || dirSrc == "" {
		return "", err
	}
	//根据目录路径生成文件路径
	fileSrc := dirSrc + c.file.GetPathSep() + name + fileType
	return fileSrc, nil
}

//创建新的目录
//src - 父级目录
func (c *Coll) CreateDirSrc(src string) (string, error) {
	//路径分隔符
	sep := c.file.GetPathSep()
	//构建子目录路径
	var dataSrc string
	var dataSrcF string = src + sep + c.GetNowDateYM() + sep + c.GetNowDateD() + sep
	for i := 1; i < 100; i++ {
		//构建路径
		dataSrc = dataSrcF + strconv.Itoa(i)
		//判断该目录是否存在
		if c.file.IsFolder(dataSrc) {
			//如果存在，则判断该文件夹下文件数量是否超过100，超过则进入下一个循环
			max, err := c.file.GetFileListCount(dataSrc)
			if err != nil {
				return "", err
			}
			if max > 100 {
				dataSrc = ""
				continue
			} else {
				return dataSrc, nil
			}
		} else {
			b, err := c.file.CreateDir(dataSrc)
			if b == false || err != nil {
				return "", err
			}
			return dataSrc, nil
		}
	}
	//如果今天1-100个目录全满，则创建返回创建失败
	if dataSrc == "" {
		return "", nil
	}
	//其他逻辑错误，到达这里，直接返回失败
	return "", nil
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

//获取日志内容
func (c *Coll) GetLog() (string, error) {
	contentByte, err := c.file.ReadFile(c.logReadSrc)
	content := string(contentByte)
	return content, err
}

//获取日志文件路径
//用于输出该文件
func (c *Coll) GetLogSrc() string {
	return c.logReadSrc
}

//清空日志内容
func (c *Coll) ClearLogContent() {
	var newC []byte
	_, _ = c.file.WriteFile(c.logReadSrc, newC)
}
