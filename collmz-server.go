package main

import (
	"./core"
	"./router"
	"./user"
	"gopkg.in/mgo.v2"
	"strconv"
	"./handle"
)

//全局APP名称
var AppName string

//全局DB数据库操作模块
var DB *mgo.Database

//全局User操作模块
var UserOperate user.User

//全局Session
var SessionOperate core.SessionOperate

//控制器主程序
//该函数用于启动整个项目
func main(){
	//读取配置文件信息
	var configSrc string
	configSrc = "config" + core.PathSeparator + "config.json"
	var configData map[string]interface{}
	var b bool
	configData,b = core.LoadConfig(configSrc)
	if b == false{
		core.SendLog("无法读取config.json配置数据。")
		return
	}
	//读取APP名称
	AppName = configData["app-name"].(string)
	//连接数据库
	var session *mgo.Session
	var err error
	session,err = mgo.Dial(configData["mgo-host"].(string))
	if err != nil{
		core.SendLog("无法连接到数据库，错误 : "+err.Error())
		return
	}
	core.SendLog("数据库连接成功 : " + configData["mgo-host"].(string))
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	DB = session.DB(configData["mgo-db"].(string))
	//构建用户处理器
	var userLoginTimeoutMinute int64
	userLoginTimeoutMinute,err = strconv.ParseInt(configData["user-login-timeout-minute"].(string),10,64)
	if err != nil{
		core.SendLog(err.Error())
		return
	}
	UserOperate.Init(&SessionOperate,AppName,userLoginTimeoutMinute)
	UserOperate.SetManyUser(DB,configData["user-username"].(string),configData["user-password"].(string))
	//将全局变量赋予路由绑定内部
	handle.AppName = AppName
	handle.DB = DB
	handle.UserOperate = &UserOperate
	handle.SessionOperate = &SessionOperate
	//启动服务器
	router.RunSever(configData["server-host"].(string))
}