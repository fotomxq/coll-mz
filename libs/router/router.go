//路由器设定
package router

import(
	"net/http"
	"github.com/fotomxq/coll-mz/libs/core"
)

//通用错误
var err error
//通用配置
var config core.Config
//日志
var log core.Log
//数据存储位置
var dataSrc string

//启动路由器
func Router(){
	//获取配置数据
	err = config.LoadFile("./content/config.json")
	if err != nil{
		log.AddErrorLog(err)
		return
	}
	//读取配置文件
	dataSrc = config.Data["dataSrc"].(string)
	port := config.Data["serverLocal"].(string)
	//设定日志类
	logDirSrc := dataSrc + "/log"
	log.SetDirSrc(logDirSrc)
	//设定静态绑定
	http.Handle("/assets/",http.FileServer(http.Dir("template")))
	//设定动态绑定
	http.HandleFunc("/",page404Handle)
	http.HandleFunc("/login",pageLoginHandle)
	http.HandleFunc("/action-login",actionLoginHandle)
	http.HandleFunc("/action-logout",actionLogoutHandle)
	http.HandleFunc("/center",pageCenterHandle)
	//启动路由器
	log.AddLog("服务器已经启动，地址：" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil{
		log.AddErrorLog(err)
	}
}