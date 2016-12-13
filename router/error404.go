package router

import (
	"net/http"
	"../handle"
)

//设定错误页面
func set404(){
	http.HandleFunc("/",handle.Page404)
}
