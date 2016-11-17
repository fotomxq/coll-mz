package router

import (
	"html/template"
	"net/http"
)

//强行启动采集器
func actionCollRunHandle(w http.ResponseWriter, r *http.Request) {
	//如果还未登录，则跳转
	if LoginCheck(w, r) == false {
		return
	}
	//根据采集需求，启动采集脚本

	//提示已经开启采集数据
	data := map[string]template.HTML{
		"title":        template.HTML("启动采集"),
		"contentTitle": template.HTML("开始采集"),
		"content":      template.HTML("已经启动了采集程序，请稍等。"),
		"gotoURL":      template.HTML("center"),
	}
	pageTipHandle(w, r, data)
}
