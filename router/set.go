package router

import "net/http"

//设定页面

//显示所有可修改的config.json选项
//只有管理员可进入

//设定页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageSet(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w, r, "admin")
	if userID == "" {
		return
	}
	//初始化
	var data map[string]interface{} = map[string]interface{}{
		"refCSS": []string{
			"theme",
		},
		"refJS": []string{
			"set", "sha1","message",
		},
	}
	//输出页面
	showTemplate(w, r, "set.html", data)
}

//设定动作处理
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionSet(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w, r, "admin")
	if userID == "" {
		return
	}
	//检查post
	if checkPost(r) == false {
		return
	}
	//初始化
	var data map[string]interface{} = map[string]interface{}{}
	var b bool
	//post action
	var postAction string
	postAction = r.FormValue("action")
	switch postAction {
	}
	postJSONData(w, r, data, b, userID != "")
}