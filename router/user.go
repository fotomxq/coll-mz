package router

import "net/http"

//管理用户界面、基本操作处理，以及其他一些模块。

//用户管理页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageUser(w http.ResponseWriter, r *http.Request){
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w,r,"user")
	if userID == ""{
		return
	}
	//输出页面
	var data map[string]interface{} = map[string]interface{}{
		"refLocalCss" : []string{
			"user",
		},
		"refLocalJs" : []string{
			"user",
		},
	}
	showTemplate(w,r,"user.html",data)
}

//用户管理界面动作处理
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionUser(w http.ResponseWriter, r *http.Request){
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w,r,"user")
	if userID == ""{
		return
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//用户相关的通用模块
//////////////////////////////////////////////////////////////////////////////////////

//检查是否已经登录
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//param page string 当前页面名称
//return string 用户ID
func checkIPAndLogged(w http.ResponseWriter, r *http.Request,page string) string{
	//检查IP是否可访问
	if checkIP(r) == false{
		return ""
	}
	//检查是否已经登录了
	var userID string
	userID = userCheckLogged(w, r)
	if userID == "" {
		goURL(w, r, "/login")
		return ""
	}
	//检查用户权限，是否足够访问该页面？
	if glob.UserOperate.CheckUserVisitPage(userID,page) == false{
		return ""
	}
	//返回
	return userID
}
