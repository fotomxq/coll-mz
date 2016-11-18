package router

import (
	"net/http"
)

//检查用户是否登陆并跳转
func CheckLoginToURL(w http.ResponseWriter, r *http.Request) bool {
	if LoginCheck(w, r) == false {
		data := PageTipHandleData("您尚未登录","请先登录","您尚未登录平台，请先进入登录页面，使用用户和密码登录后再进行访问。","login")
		pageTipHandle(w, r, data)
		return false
	} else {
		return true
	}
}
