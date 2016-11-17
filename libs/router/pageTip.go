package router

import (
	"html/template"
	"net/http"
)

//提示页面
//可提供一串字符，输出到页面用于提示
func pageTipHandle(w http.ResponseWriter, r *http.Request, data map[string]template.HTML) {
	t, err := template.ParseFiles(modGetTempSrc("tip.html"))
	if err != nil {
		log.AddErrorLog(err)
	}
	t.Execute(w, data)
}
