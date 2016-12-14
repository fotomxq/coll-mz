package user

import (
	"../core"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//用户处理器包
//可用于用户管理、登录
//支持任意数据库类型，或直接制定单一用户密码
//使用方法：声明User类后初始化，之后选择单一用户还是多用户模式设定即可
//依赖外部包：
// mgo (gopkg.in/mgo.v2 / gopkg.in/mgo.v2/bson)
//依赖本地包：
// core.session-operate.go
// core.match-string.go

//全局DB数据库操作模块
var DB *mgo.Database

//全局验证处理句柄
var MatchString *core.MatchString

//全局session会话操作
var SessionOperate *core.SessionOperate

//日志处理器
var logOperate core.LogOperate

//全局应用名称
var AppName string

//全局标识码
var Mark string

//全局用户自动退出时限
var UserLoginTimeoutMinute int64

//单一用户模式是否启动
var oneUserStatus bool
//单一用户名和密码
var oneUsername string
var oneUserpasswd string
//字段列
var fields []string
//数据表
var table string
//数据库合集
var dbColl *mgo.Collection

//用户字段组
type UserFields struct {
	Id_ bson.ObjectId
	NiceName string
	UserName string
	Password string
	LastIP string
	LastTime int64
	IsDisabled bool
}

//初始化
func Init() {
	oneUserStatus = false
	table = "user"
	fields = []string{
		"_id","nicename","username","password","lastip","lasttime","isdisabled",
	}
	logOperate.Init(DB,"user-log")
}

//日志输出模块
//param fileName string 文件名称
//param ipAddr string IP地址
//param funcName string 函数名称
//param mark string 标记名称
//param message string 消息
func sendLog(fileName string,ipAddr string,funcName string,mark string,message string) {
	logOperate.SetFileName(fileName)
	logOperate.SendLog(ipAddr,funcName,mark,message)
}