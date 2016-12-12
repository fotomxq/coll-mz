package router

import (
	"net/http"
	"../handle"
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

//设定登录部分
func setLogin(){
	http.HandleFunc("/login",handle.PageLogin)
	http.HandleFunc("/action-login",handle.Login)
	http.HandleFunc("/action-logout",handle.Logout)
}

//设定错误页面
func set404(){
	http.HandleFunc("/",handle.Page404)
}

//静态路由部分
func setAssets(){
	http.Handle("/assets/",http.StripPrefix("/assets/",http.FileServer(http.Dir(getTemplateSrc("assets")+core.PathSeparator))))
	http.HandleFunc("/favicon.ico", handle.FileFavicon)
}

//获取template路径
//param name string 路径末尾文件名称
//return string 路径
func getTemplateSrc(name string) string{
	return "." + core.PathSeparator + "template" + core.PathSeparator + name
}