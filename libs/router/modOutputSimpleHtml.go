package router

import (
	"net/http"
)

//向浏览器输出一个简单的HTML数据
func modOutputSimpleHtml(w http.ResponseWriter, r *http.Request, content string) {
	err = simpleHttp.PostText(w,r,content)
	if err != nil{
		log.AddErrorLog(err)
	}
}