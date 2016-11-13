package router

import(
	"net/http"
	"html/template"
)

//登录页面
func loginHandler(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("template/login.html")
	if (err != nil) {
		log.AddErrorLog(err)
	}
	t.Execute(w, nil)
}

