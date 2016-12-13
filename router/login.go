package router

import (
	"net/http"
	"../handle"
)

//设定登录部分
func setLogin(){
	http.HandleFunc("/login",handle.PageLogin)
	http.HandleFunc("/action-login",handle.ActionLogin)
	http.HandleFunc("/action-logout",handle.ActionLogout)
}