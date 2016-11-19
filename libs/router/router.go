//路由器设定
package router

import (
	"github.com/fotomxq/coll-mz/libs/coll"
	"github.com/fotomxq/coll-mz/libs/core"
	"net/http"
)

//通用错误
var err error

//通用配置
var config core.Config

//日志
var log core.Log

//数据存储位置
var dataSrc string

//收集器
var collPage coll.Coll

//文件操作
var file core.FileOperate

//路径分隔符
var fileSep string

//通用简化http处理器
var simpleHttp core.SimpleHttp

//启动路由器
func Router() {
	//路径分隔符
	fileSep = file.GetPathSep()
	//获取配置数据
	err = config.LoadFile("content" + fileSep + "config.json")
	if err != nil {
		log.AddErrorLog(err)
		return
	}
	//读取配置文件
	dataSrc = config.Data["dataSrc"].(string)
	port := config.Data["serverLocal"].(string)
	//设定日志类
	logDirSrc := dataSrc + fileSep + "log"
	log.SetDirSrc(logDirSrc)
	//创建收集器实例
	b, err := collPage.Create(dataSrc)
	if err != nil {
		log.AddErrorLog(err)
		return
	}
	if b == false {
		log.AddLog("cannot create coll.")
		return
	}
	coll.CollPg = &collPage
	//设定静态绑定
	http.Handle("/assets/", http.FileServer(http.Dir("template")))
	//设定动态绑定
	http.HandleFunc("/", page404Handle)
	http.HandleFunc("/login", pageLoginHandle)
	http.HandleFunc("/action-login", actionLoginHandle)
	http.HandleFunc("/action-logout", actionLogoutHandle)
	http.HandleFunc("/center", pageCenterHandle)
	http.HandleFunc("/set", pageSetHandle)
	http.HandleFunc("/action-coll-run", actionCollRunHandle)
	http.HandleFunc("/action-view", actionViewHandle)
	//启动路由器
	log.AddLog("服务器已经启动，地址：" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.AddErrorLog(err)
		return
	}
}
