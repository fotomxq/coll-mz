package router

import(
	"net/http"
	"html/template"
)

//登录页面
func pageLoginHandle(w http.ResponseWriter, r *http.Request){
	if LoginCheck(w,r) == true{
		//如果已经登录，则跳转
		data := map[string]template.HTML{
			"title" : template.HTML("您已登录"),
			"contentTitle" : template.HTML("您已经登录过了"),
			"content" : template.HTML("您已经登录过该平台，5秒后自动跳转到首页。"),
			"gotoURL" : template.HTML("center"),
		}
		pageTipHandle(w,r,data)
	}else{
		//如果未登录，则展开界面
		t, err := template.ParseFiles(modGetTempSrc("login.html"))
		if (err != nil) {
			log.AddErrorLog(err)
		}
		t.Execute(w, nil)
	}
}

