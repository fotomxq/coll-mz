package router

import (
	"html/template"
	"net/http"
)

//设定页面
func pageSetHandle(w http.ResponseWriter, r *http.Request) {
	//如果还未登录，则跳转
	if CheckLoginToURL(w, r) == false {
		return
	}
	//如果已经登录，则进入
	t, err := template.ParseFiles(modGetTempSrc("set.html"))
	if err != nil {
		log.AddErrorLog(err)
	}
	data := map[string]map[string]string{}
	data["collList"] = map[string]string{}
	data["collList"] = collPage.CollList
	t.Execute(w, data)
}
