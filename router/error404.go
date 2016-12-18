package router

import "net/http"

//404错误页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func Page404(w http.ResponseWriter, r *http.Request){
	//检查IP是否可访问
	if checkIP(r) == false{
		return
	}
	//判断URL
	if r.URL.Path != "/" {
		var data map[string]interface{} = map[string]interface{}{
			"addTitle" : "找不到页面",
			"refCSS" : []string{
				"404",
			},
		}
		showTemplate(w, r, "404.html", data)
	}else{
		goURL(w,r,"/login")
	}
}