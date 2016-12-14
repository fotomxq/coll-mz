package router

import (
	"../core"
	"net/http"
	"gopkg.in/mgo.v2"
)

//路由器设定

//系统path符号
var pathSep string

//全局对接句柄
var glob *GlobOperate

//全局对接类型
type GlobOperate struct{
	//数据库操作模块
	DB *mgo.Database
	//Session句柄
	SessionOperate *core.SessionOperate
	//APP名称
	AppName string
	//用户句柄
	UserOperate *core.User
}

//初始化
func Init(params *GlobOperate){
	pathSep = core.PathSeparator
	glob = params
}

//运行服务器
func RunSever(host string){
	//绑定静态assets数据
	http.Handle("/assets/",http.StripPrefix("/assets/",http.FileServer(http.Dir(getTemplateSrc("assets")+core.PathSeparator))))
	http.HandleFunc("/favicon.ico", FileFavicon)
	//绑定错误页面
	http.HandleFunc("/",Page404)
	//绑定登录和退出页面
	http.HandleFunc("/login",PageLogin)
	http.HandleFunc("/action-login",ActionLogin)
	http.HandleFunc("/action-logout",ActionLogout)
	//绑定中心页面
	http.HandleFunc("/center",PageCenter)
	//输出日志
	core.SendLog("****** 启动服务器 : " + host + " ******")
	//启动路由器
	var err error
	err = http.ListenAndServe(host, nil)
	if err != nil{
		core.SendLog(err.Error())
		return
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// 转接外部方法
/////////////////////////////////////////////////////////////////////////////////////////////////////////


//转接日志输出模块
//param message string 日志内容
func sendLog(message string) {
	core.SendLog(message)
}

//转接用户登录检查模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//return string 用户ID
func userCheckLogged(w http.ResponseWriter,r *http.Request) string{
	//确保启动会话
	startSession()
	//返回登录用户ID，无登录或失败则返回空字符串
	return glob.UserOperate.GetLoginStatus(w,r)
}

//转接用户登录模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否登录成功
func userLogin(w http.ResponseWriter,r *http.Request,username string,passwdSha1 string) bool{
	return glob.UserOperate.Login(w,r,username,passwdSha1)
}

//转接用户退出模块
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func userLogout(w http.ResponseWriter,r *http.Request){
	glob.UserOperate.Logout(w,r)
}