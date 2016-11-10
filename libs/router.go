//路由器设定
package collmzLibs

import(
	"net/http"
)

func Router(){
	http.Handle("/assets/",http.FileServer(http.Dir("template")))
	http.ListenAndServe(":8888", nil)
}