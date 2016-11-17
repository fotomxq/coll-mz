package router

import (
	"github.com/gorilla/sessions"
	"net/http"
)

//启动session
var store = sessions.NewCookieStore([]byte("coll-mz-session"))

func LoginIn(w http.ResponseWriter, r *http.Request, user string, passwd string) bool {
	if user == config.Data["user"].(string) && passwd == config.Data["passwd"].(string) {
		return changeLoginSession(w, r, true)
	}
	return false
}

//检查用户是否已经登录，或提供用户和密码检测
func LoginCheck(w http.ResponseWriter, r *http.Request) bool {
	//检查是否已经登陆
	if getLoginSession(w, r) == true {
		return true
	}
	//检测全没通过，则失败
	return false
}

//退出登录
func LoginOut(w http.ResponseWriter, r *http.Request) bool {
	return changeLoginSession(w, r, false)
}

//修改登录状态设定
func changeLoginSession(w http.ResponseWriter, r *http.Request, b bool) bool {
	session, err := store.Get(r, "login")
	if err != nil {
		return false
	}
	bStr := "in"
	if b == false {
		bStr = "out"
	}
	session.Values["logged-ok"] = bStr
	session.Save(r, w)
	return true
}

//获取登录状态
func getLoginSession(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "login")
	if err != nil {
		return false
	}
	if session.Values["logged-ok"] == nil {
		return false
	}
	if session.Values["logged-ok"].(string) == "in" {
		return true
	}
	return false
}
