package router

import(
	"net/http"
	"html/template"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login/", http.StatusFound)
	}
	t, err := template.ParseFiles("template/404.html")
	if (err != nil) {
		log.AddErrorLog(err)
	}
	t.Execute(w, nil)
}