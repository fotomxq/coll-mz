package router

import (
	"html/template"
	"net/http"
)

//向浏览器输出一个简单的HTML数据
func modOutputSimpleHtml(w http.ResponseWriter, r *http.Request,content string){
	t := template.New("modOutputSimpleHtml.html")
	if err != nil{
		return
	}
	values := map[string]template.HTML{"html": template.HTML(content)}
	err = t.Execute(w,values)
	if err != nil{
		return
	}
}