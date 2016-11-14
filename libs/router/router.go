//路由器设定
package router

import(
	"net/http"
	"../core"
)

//通用错误
var err error
//通用配置
var config core.Config
//日志
var log core.Log

//启动路由器
func Router(l core.Log,local string){
	log = &l
	//读取配置文件
	err = config.LoadFile("content/config/server.json")
	if err != nil{
		log.AddErrorLog(err)
	}
	port := config.Data["port"].(string)
	//设定静态绑定
	http.Handle("/assets",http.FileServer(http.Dir("template")))
	//设定动态绑定
	http.HandleFunc("/",notFoundHandler)
	http.HandleFunc("/login",loginHandler)
	http.HandleFunc("/action-login",actionLoginHander)
	//启动路由器
	log.AddLog("服务器已经启动，端口：" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil{
		log.AddErrorLog(err)
	}
}