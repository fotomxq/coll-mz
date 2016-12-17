package router

import (
	"net/http"
)

//中心页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageCenter(w http.ResponseWriter, r *http.Request) {
	//检查IP是否可访问
	if checkIP(r) == false{
		return
	}
	//检查是否已经登录
	var userID string
	userID = userCheckLogged(w, r)
	if userID == "" {
		goURL(w, r, "/login")
		return
	}else{
		var data map[string]interface{} = map[string]interface{}{}
		showTemplate(w,r,"center.html",data)
		return
	}
}
