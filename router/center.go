package router

import (
	"net/http"
)

//中心页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageCenter(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w,r,"center")
	if userID == ""{
		return
	}
	//输出页面
	var data map[string]interface{} = map[string]interface{}{
		"refLocalCss" : []string{
			"center",
		},
		"refLocalJs" : []string{
			"center",
		},
	}
	showTemplate(w,r,"center.html",data)
}
