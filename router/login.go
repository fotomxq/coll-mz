package router

import (
	"net/http"
)

//用户登录、退出检查部分

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageLogin(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	if userCheckLogged(w, r) != "" {
		goURL(w, r, "/center")
	}else{
		showTemplate(w,r,"login.html",nil)
	}
}

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionLogin(w http.ResponseWriter, r *http.Request) {
	//初始化变量
	var data string = "no-login"
	defer postJSONData(w,r,&data,true)
	//检查是否已经登录
	if userCheckLogged(w, r) != "" {
		data = "logged"
		return
	}else{
		//检查post提交
		if checkPost(r) == false{
			data = "error"
			return
		}
		//获取登录用户名和密码
		var username string
		var passwdSha1 string
		username = r.FormValue("username")
		passwdSha1 = r.FormValue("password")
		if len(username) < 4  && len(passwdSha1) < 10 {
			data = "error-username-or-passwd"
			return
		}
		//提交给登录模块
		var b bool
		b = userLogin(w,r,username,passwdSha1)
		if b == true{
			data = "success"
		}else{
			data = "error-login"
		}
		return
	}
}

//退出操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionLogout(w http.ResponseWriter, r *http.Request){
	if userCheckLogged(w,r) != ""{
		userLogout(w,r)
	}
	goURL(w,r,"/login")
}