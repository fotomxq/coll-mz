package router

import(
	"net/http"
	"html/template"
)

//404错误页面
func page404Handle(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	t, err := template.ParseFiles(modGetTempSrc("404.html"))
	if (err != nil) {
		log.AddErrorLog(err)
	}
	t.Execute(w, nil)
}