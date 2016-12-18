package router

import (
	"net/http"
)

//用户登录、退出检查部分

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageLogin(w http.ResponseWriter, r *http.Request) {
	//检查IP是否可访问
	if checkIP(r) == false{
		return
	}
	//检查是否已经登录
	if userCheckLogged(w, r) != "" {
		goURL(w, r, "/center")
	}else{
		if r.URL.Path != "/login" {
			sendLog("router/login.go",getIPAddr(r),"PageLogin","user-no-logged","用户尚未登录，但访问了内部页面，URL："+r.URL.Path)
		}
		var data map[string]interface{} = map[string]interface{}{
			"addTitle" : "登录",
			"refCSS" : []string{
				"login",
			},
			"refJS" : []string{
				"login","sha1",
			},
		}
		showTemplate(w,r,"login.html",data)
	}
}

//登录操作
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionLogin(w http.ResponseWriter, r *http.Request) {
	//检查IP是否可访问
	if checkIP(r) == false{
		return
	}
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
	//检查IP是否可访问
	if checkIP(r) == false{
		return
	}
	if userCheckLogged(w,r) != ""{
		userLogout(w,r)
	}
	goURL(w,r,"/login")
}