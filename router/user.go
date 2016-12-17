package router

import "net/http"

//用户管理页面
func PageUser(w http.ResponseWriter, r *http.Request){
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w,r)
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
func ActionUser(w http.ResponseWriter, r *http.Request){
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w,r)
	if userID == ""{
		return
	}
}

