package core

import (
	"bytes"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"time"
)

//日志处理模块
//使用方式：
// 1、如果只使用文件存储日志，则直接调用SendLog函数即可
// 可预先配置好LogOperateSrc变量，指定其他存储日志的路径
// 2、如果是数据库类型，则可根据需要声明LogOperate类，初始化后即可使用内部SendLog函数
//依赖内部模块：无
//依赖外部库：无

////////////////////////////////////////////////////////////////////////////////////////////////////////
//以下部分是文件存储方式
////////////////////////////////////////////////////////////////////////////////////////////////////////

//日志存储文件夹路径
//默认存储到程序所在目录的log文件夹下
var LogOperateSrc = "log"

//是否开启向文件存储
var LogOperateFileBool = false

//向控制台发送一个日志
//自动向控制台和日志文件输出信息
//日志同时将存储到"LogOperateSrc/200601/02/2006010215.log"文件内
//param message string 日志信息
func SendLog(message string) {
	//加入时间
	var t time.Time
	t = time.Now()
	message = t.Format("2006-01-02 15:04:05.999999999") + " " + message + "\n"
	//向控制台输出日志
	fmt.Print(message)
	if LogOperateFileBool == false {
		return
	}
	//生成日志文件路径
	var logDir string
	logDir = LogOperateSrc + PathSeparator + t.Format("200601") + PathSeparator + t.Format("02")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " Failed to create the folder : " + logDir)
		return
	}
	var logSrc string
	logSrc = logDir + PathSeparator + t.Format("2006010215") + ".log"
	//向文件输出日志
	var messageByte []byte
	messageByte = []byte(message)
	//如果文件不存在，则直接创建写入日志
	_, err = os.Stat(logSrc)
	if err != nil || os.IsNotExist(err) == true {
		err = ioutil.WriteFile(logSrc, messageByte, os.ModeAppend)
		if err != nil {
			fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		}
		return
	}
	//如果文件存在，则读出文件内容，添加新日志
	var fd *os.File
	fd, err = os.Open(logSrc)
	if err != nil {
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
	defer fd.Close()
	var c []byte
	c, err = ioutil.ReadAll(fd)
	if err != nil {
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
	var s [][]byte
	s = [][]byte{
		c,
		messageByte,
	}
	var sep []byte
	messageByte = bytes.Join(s, sep)
	err = ioutil.WriteFile(logSrc, messageByte, os.ModeAppend)
	if err != nil {
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//以下部分是数据库存储方式
////////////////////////////////////////////////////////////////////////////////////////////////////////

//日志操作类
//该模块同样是日志操作，但依赖于mongodb数据库
type LogOperate struct {
	//全局DB数据库操作模块
	db *mgo.Database
	//数据表名称
	table string
	//数据库合集
	dbColl *mgo.Collection
}

//数据表结构
type LogOperateFields struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//创建时间
	CreateTime string
	//IP地址
	IpAddr string
	//文件名称
	FileName string
	//函数名称
	FuncName string
	//标记
	Mark string
	//消息
	Message string
}

//初始化模块
//使用之前必须确保进行该步骤
//param db *mgo.Database 数据库句柄
//param appName 应用名称
func (this *LogOperate) Init(db *mgo.Database, appName string) {
	//保存数据库连接
	this.db = db
	//构建数据集合
	this.table = "log"
	this.dbColl = db.C(this.table)
}

//发送新的日志
//param fileName string 文件名称
//param ipAddr string IP地址
//param funcName string 函数名称
//param mark string 标记名称
//param message string 消息
func (this *LogOperate) SendLog(fileName string, ipAddr string, funcName string, mark string, message string) {
	//向数据库添加日志
	err = this.dbColl.Insert(&LogOperateFields{bson.NewObjectId(), time.Now().Format("2006-01-02 15:04:05.999999999"), ipAddr, fileName, funcName, mark, message})
	if err != nil {
		fmt.Println("无法向数据库添加日志数据。")
	}
}

//查看日志数据
//强制按照ID倒序排序
//param page int 页数
//param max int 页长
//return []map[string]string,bool 日志数据组，是否获取成功
func (this *LogOperate) View(page int, max int) ([]map[string]string, bool) {
	var result []LogOperateFields
	var res []map[string]string = []map[string]string{}
	var skip int
	skip = (page - 1) * max
	err = this.dbColl.Find(nil).Sort("-_id").Skip(skip).Limit(max).All(&result)
	if err != nil {
		return res, false
	}
	for key := range result {
		res = append(res, map[string]string{
			"ID":         result[key].ID.Hex(),
			"CreateTime": result[key].CreateTime,
			"IpAddr":     result[key].IpAddr,
			"FileName":   result[key].FileName,
			"FuncName":   result[key].FuncName,
			"Mark":       result[key].Mark,
			"Message":    result[key].Message,
		})
	}
	return res, true
}

//清理日志
//仅用于debug模式，其他模式下请勿使用该模块
//return bool 是否成功
func (this *LogOperate) Clear() bool {
	_, err = this.dbColl.RemoveAll(nil)
	return err == nil
}
