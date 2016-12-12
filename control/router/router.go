package router

import (
	"../core"
	"net/http"
)

//路由器设定
//可以在该文件夹下建立多个go文件，将handle和对应的url绑定

//运行服务器
func RunSever(host string){
	//设定静态assets数据绑定
	setAssets()
	//输出日志
	core.SendLog("****** 启动服务器 : " + host + " ******")
	//启动路由器
	http.ListenAndServe(host, nil)
}

//静态路由部分
func setAssets(){
	var src string
	src = "template" + core.PathSeparator + "assets" + core.PathSeparator
	http.Handle("/assets/",http.FileServer(http.Dir(src)))
}