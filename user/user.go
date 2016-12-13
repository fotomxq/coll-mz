package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"../core"
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

//全局应用名称
var AppName string

//全局标识码
var Mark string

//全局用户自动退出时限
var UserLoginTimeoutMinute int64

//转接日志输出模块
//param message string 日志内容
func sendLog(message string) {
	core.SendLog(message)
}

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
	id bson.ObjectId
	nicename string
	username string
	password string
	last_ip string
	last_time int
	is_disabled bool
}

//初始化
func init() {
	oneUserStatus = false
	table = "user"
	fields = []string{
		"_id","nicename","username","password","last_ip","last_time","is_disabled",
	}
}