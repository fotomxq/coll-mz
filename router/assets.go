package router

import (
	"net/http"
	"../core"
	"../handle"
)

//静态路由部分
func setAssets(){
	http.Handle("/assets/",http.StripPrefix("/assets/",http.FileServer(http.Dir(getTemplateSrc("assets")+core.PathSeparator))))
	http.HandleFunc("/favicon.ico", handle.FileFavicon)
}
