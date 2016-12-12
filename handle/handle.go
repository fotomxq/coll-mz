package handle

import (
	"gopkg.in/mgo.v2"
	"../core"
	"../user"
)

//路由句柄处理器
//路由URL对应的句柄函数在这里声明

//全局DB数据库操作模块
var DB *mgo.Database

//全局User操作模块
var UserOperate *user.User

//全局Session
var SessionOperate *core.SessionOperate

//全局APP名称
var AppName string