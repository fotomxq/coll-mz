package handle

import (
	"net/http"
)

//用户登录、退出检查部分

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageLogin(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	if checkLogged(w, r) == true {
		goURL(w, r, "/center")
	}else{
		showTemplate(w,r,"login.html",nil)
	}
}

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionLogin(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	if checkLogged(w, r) == true {
		goURL(w,r,"/center")
	}else{
		//检查post提交
		if checkPost(r) == false{
			goURL(w,r,"/login")
		}
		//获取登录用户名和密码
		var username string
		var passwdSha1 string
		username = r.FormValue("username")
		passwdSha1 = r.FormValue("passwd")
		if len(username) < 4  && len(passwdSha1) < 10 {
			goURL(w,r,"/login")
		}
		//提交给登录模块
		var b bool
		b = userLogin(w,r,username,passwdSha1)
		if b == true{
			goURL(w,r,"/center")
		}
		goURL(w,r,"/login")
	}
}

//退出操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionLogout(w http.ResponseWriter, r *http.Request){
	if checkLogged(w,r) == true{
		userLogout(w,r)
	}
	goURL(w,r,"/login")
}

//检查是否已登录
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//return bool 是否登录
func checkLogged(w http.ResponseWriter, r *http.Request) bool {
	//确保启动会话
	startSession()
	//返回是否已经登录
	return userCheckLogged(w,r) > 0
}
