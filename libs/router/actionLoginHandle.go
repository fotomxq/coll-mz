package router

import (
	"github.com/fotomxq/coll-mz/libs/core"
	"html/template"
	"net/http"
)

//登录页面动作
func actionLoginHandle(w http.ResponseWriter, r *http.Request) {
	//获取数据
	postUser := r.FormValue("email")
	postPasswd := r.FormValue("password")
	//验证数据
	if LoginIn(w, r, postUser, postPasswd) == true {
		http.Redirect(w, r, "/center", http.StatusFound)
	} else {
		data := map[string]template.HTML{
			"title":        template.HTML("登录失败"),
			"contentTitle": template.HTML("登录失败"),
			"content":      template.HTML("您输入的用户名或密码存在错误，请重新输入。"),
			"gotoURL":      template.HTML("login"),
		}
		pageTipHandle(w, r, data)
	}
}

//检查用户名和密码是否合法
func checkUserPasswd(user string, passwd string) bool {
	ms := new(core.MatchString)
	if ms.CheckEmail(user) == true && ms.CheckPassword(passwd) == true {
		return true
	}
	return false
}
