package router

import "net/http"

//debug模式处理器

//debug模式
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageDebug(w http.ResponseWriter, r *http.Request) {
	//debug模式下，IP锁定、用户是否登录都无效
	//发布环境下请勿启动该模式
	//如果没有指定任何URL动作，则按照输出日志数据处理
	//检查post有效性
	var action string
	if checkPost(r) == false {
		action = "log"
	} else {
		action = r.FormValue("action")
		if action == "" {
			action = "log"
		}
	}
	//初始化变量
	var data map[string]interface{} = map[string]interface{}{
		"refCSS": []string{"debug", "theme"},
		"refJS":  []string{"debug"},
	}
	//判断动作类型
	switch action {
	case "clear-log":
		//清理日志
		glob.LogOperate.Clear()
		fallthrough
	case "log":
		//查看日志
		data["actionLog"] = true
		var logData []map[string]string
		var b bool
		logData, b = glob.LogOperate.View(1, 9999)
		if b == true {
			data["logData"] = logData
		}
	case "clear-user":
		_ = glob.UserOperate.DeleteAll()
	}
	//输出变量及类型
	showTemplate(w, r, "debug.html", data)
}
