package router

import (
	"html/template"
	"net/http"
)

//检查用户是否登陆并跳转
func CheckLoginToURL(w http.ResponseWriter, r *http.Request) bool {
	if LoginCheck(w, r) == false {
		data := map[string]template.HTML{
			"title":        template.HTML("您尚未登录"),
			"contentTitle": template.HTML("请先登录"),
			"content":      template.HTML("您尚未登录平台，请先进入登录页面，使用用户和密码登录后再进行访问。"),
			"gotoURL":      template.HTML("login"),
		}
		pageTipHandle(w, r, data)
		return false
	} else {
		return true
	}
}
