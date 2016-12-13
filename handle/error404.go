package handle

import "net/http"

//404错误页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func Page404(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		showTemplate(w, r, "404.html", nil)
	}else{
		if checkLogged(w,r) == false{
			goURL(w,r,"/login")
		}
	}
}