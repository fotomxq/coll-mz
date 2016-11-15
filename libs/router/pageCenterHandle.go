package router

import(
	"net/http"
	"html/template"
)

//登录页面
func pageCenterHandle(w http.ResponseWriter, r *http.Request){
	//如果还未登录，则跳转
	if LoginCheck(w,r) == false{
		data := map[string]template.HTML{
			"title" : template.HTML("您尚未登录"),
			"contentTitle" : template.HTML("请先登录"),
			"content" : template.HTML("您尚未登录平台，请先进入登录页面，使用用户和密码登录后再进行访问。"),
			"gotoURL" : template.HTML("login"),
		}
		pageTipHandle(w,r,data)
	}else{
		//如果已经登录，则进入
		t, err := template.ParseFiles("template/center.html")
		if (err != nil) {
			log.AddErrorLog(err)
		}
		t.Execute(w, nil)
	}
}
