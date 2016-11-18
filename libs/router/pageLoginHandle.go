package router

import (
	"html/template"
	"net/http"
)

//登录页面
func pageLoginHandle(w http.ResponseWriter, r *http.Request) {
	if LoginCheck(w, r) == true {
		//如果已经登录，则跳转
		data := PageTipHandleData("您已登录","您已经登录过了","您已经登录过该平台，5秒后自动跳转到首页。","center")
		pageTipHandle(w, r, data)
	} else {
		//如果未登录，则展开界面
		t, err := template.ParseFiles(modGetTempSrc("login.html"))
		if err != nil {
			log.AddErrorLog(err)
		}
		t.Execute(w, nil)
	}
}
