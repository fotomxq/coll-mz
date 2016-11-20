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
	//确保post/get正常
	err = r.ParseForm()
	if err != nil {
		return
	}
	//获取提交动作类型
	postAction := r.FormValue("action")
	//判断动作类型，并反馈结果
	data := map[string]template.HTML{}
	switch postAction {
	case "coll-all":
		//并发启动搜集模块
		//直接返回调用成功的提示
		data = PageTipHandleData("启动采集", "开始采集", "已经启动了采集程序，请稍等。", "set")
		collPage.ClearLogContent()
		go collPage.RunAll()
		modOutputSimpleHtml(w, r, "coll-run-ok")
		break
	case "get-log":
		logSrc := collPage.GetLogSrc()
		t, err := template.ParseFiles(logSrc)
		if err != nil {
			log.AddErrorLog(err)
		}
		t.Execute(w, nil)
		break
	case "backup-database":
		break
	case "return-database":
		break
	default:
		data = PageTipHandleData("错误", "非法参数", "您提交了一个错误的指令。", "set")
		pageTipHandle(w, r, data)
		break
	}
}
