package handle

import (
	"gopkg.in/mgo.v2"
	"../core"
	"../user"
	"net/http"
)

//路由句柄处理器
//路由URL对应的句柄函数在这里声明

//全局DB数据库操作模块
var DB *mgo.Database

//全局Session
var SessionOperate *core.SessionOperate

//全局APP名称
var AppName string

//全局系统path符号
var PathSep string

//转接日志输出模块
//param message string 日志内容
func sendLog(message string) {
	core.SendLog(message)
}

//转接用户登录检查模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//return int 用户ID
func userCheckLogged(w http.ResponseWriter,r *http.Request) int{
	return user.GetLoginStatus(w,r)
}

//转接用户登录模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否登录成功
func userLogin(w http.ResponseWriter,r *http.Request,username string,passwdSha1 string) bool{
	return user.Login(w,r,username,passwdSha1)
}

//转接用户退出模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func userLogout(w http.ResponseWriter,r *http.Request){
	user.Logout(w,r)
}