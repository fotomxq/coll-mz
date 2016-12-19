package main

import (
	"./core"
	"./router"
	"gopkg.in/mgo.v2"
	"strconv"
)

//全局APP名称
var AppName string

//全局APP标识
var AppMark string

//全局数据库操作模块
var DB *mgo.Database

//全局Session句柄
var SessionOperate core.SessionOperate

//全局验证处理句柄
var MatchString core.MatchString

//全局日志数据库操作
var LogOperate core.LogOperate

//全局用户操作句柄
var UserOperate core.User

//IP黑名单处理器
var IPAddrOperate core.IPAddrBan

//控制器主程序
//该函数用于启动整个项目
func main() {
	//读取配置文件信息
	var configSrc string
	configSrc = "config" + core.PathSeparator + "config.json"
	var configData map[string]interface{}
	var b bool
	configData, b = core.LoadConfig(configSrc)
	if b == false {
		core.SendLog("无法读取config.json配置数据。")
		return
	}

	//读取APP名称
	AppName = configData["app-name"].(string)
	AppMark = configData["app-mark"].(string)

	//连接数据库
	var session *mgo.Session
	var err error
	session, err = mgo.Dial(configData["mgo-host"].(string))
	if err != nil {
		core.SendLog("无法连接到数据库，错误 : " + err.Error())
		return
	}
	core.SendLog("数据库连接成功 : " + configData["mgo-host"].(string))
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	DB = session.DB(configData["mgo-db"].(string))

	//初始化日志操作句柄
	LogOperate.Init(DB, AppMark)

	//将日志句柄赋予给core
	core.Log = &LogOperate

	//创建SESSION
	var sessionIPBind bool
	sessionIPBind = configData["session-ip-bind"].(string) == "true"
	var sessionTimeout int
	sessionTimeout, err = strconv.Atoi(configData["session-timeout"].(string))
	if err != nil {
		core.SendLog(err.Error())
		return
	}
	SessionOperate.Create(AppMark, DB, sessionIPBind, sessionTimeout, &MatchString)

	//构建用户处理器var userLoginTimeout int64
	var userLoginTimeout int64
	userLoginTimeout, err = strconv.ParseInt(configData["user-login-timeout"].(string), 10, 64)
	if err != nil {
		core.SendLog(err.Error())
		return
	}
	var userOneStatus bool
	userOneStatus = configData["user-one"].(string) == "true"
	UserOperate.Init(&core.UserParams{
		DB, &MatchString,
		&SessionOperate,
		&LogOperate,
		AppMark,
		AppMark,
		userLoginTimeout,
		userOneStatus,
		configData["user-username"].(string),
		configData["user-password"].(string),
		[]string{"admin", "normal"},
		map[string]map[string]interface{}{
			"admin": {
				"name": "管理员",
				"page": []string{"*"},
			},
			"normal": {
				"name": "普通用户",
				"page": []string{"center"},
			},
		},
	})

	//创建IP名单处理器
	var ipBanOn bool
	ipBanOn = configData["ip-ban-on"].(string) == "true"
	var ipWriteOn bool
	ipWriteOn = configData["ip-white-on"].(string) == "true"
	IPAddrOperate.Init(DB, ipBanOn, ipWriteOn)

	//将本机IP添加到白名单内，用于开发工作
	//IPAddrOperate.SaveToList("[::1]",false)

	//初始化路由
	var debug bool
	debug = configData["debug"].(string) == "true"
	router.Init(&router.GlobOperate{
		debug,
		DB,
		&SessionOperate,
		&LogOperate,
		&IPAddrOperate,
		configData["app-name"].(string),
		configData["app-des"].(string),
		configData["app-copyright"].(string),
		&UserOperate,
		&MatchString,
	})
	//启动服务器
	router.RunSever(configData["server-host"].(string))
}
