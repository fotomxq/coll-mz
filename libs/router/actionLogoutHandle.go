package router

import (
	"html/template"
	"net/http"
)

//退出登录
func actionLogoutHandle(w http.ResponseWriter, r *http.Request) {
	b := LoginOut(w, r)
	if b == false {
		log.AddLog("user logout.")
	}
	data := map[string]template.HTML{
		"title":        template.HTML("退出COLL-MZ平台"),
		"contentTitle": template.HTML("退出登录"),
		"content":      template.HTML("您已经成功退出了COLL-MZ平台。"),
		"gotoURL":      template.HTML("login"),
	}
	pageTipHandle(w, r, data)
}
