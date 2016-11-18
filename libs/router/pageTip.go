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

//获取提示页面所需的data数据
//返回格式 template.HTML
func PageTipHandleData(title string,title2 string,content string,toURL string) map[string]template.HTML{
	data := map[string]template.HTML{
		"title":        template.HTML(title),
		"contentTitle": template.HTML(title2),
		"content":      template.HTML(content),
		"gotoURL":      template.HTML(toURL),
	}
	return data
}