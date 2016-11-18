package router

import (
	"net/http"
)

//退出登录
func actionLogoutHandle(w http.ResponseWriter, r *http.Request) {
	b := LoginOut(w, r)
	if b == false {
		log.AddLog("user logout.")
	}
	data := PageTipHandleData("退出COLL-MZ平台","退出登录","您已经成功退出了COLL-MZ平台。","login")
	pageTipHandle(w, r, data)
}
