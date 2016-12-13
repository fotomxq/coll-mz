package router

import (
	"net/http"
	"../core"
)

//路由器设定
//可以在该文件夹下建立多个go文件，将handle和对应的url绑定

//运行服务器
func RunSever(host string){
	//设定错误页面
	set404()
	//设定静态assets数据绑定
	setAssets()
	//设定登录部分
	setLogin()
	//输出日志
	core.SendLog("****** 启动服务器 : " + host + " ******")
	//启动路由器
	var err error
	err = http.ListenAndServe(host, nil)
	if err != nil{
		core.SendLog(err.Error())
		return
	}
}