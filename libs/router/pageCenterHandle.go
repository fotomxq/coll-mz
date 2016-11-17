package router

import (
	"html/template"
	"net/http"
)

//中心页面
func pageCenterHandle(w http.ResponseWriter, r *http.Request) {
	//如果还未登录，则跳转
	if CheckLoginToURL(w, r) == false {
		return
	}
	//如果已经登录，则进入
	t, err := template.ParseFiles(modGetTempSrc("center.html"))
	if err != nil {
		log.AddErrorLog(err)
	}
	t.Execute(w, nil)
}
