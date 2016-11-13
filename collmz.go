package main

import(
	"github.com/fotomxq/coll-mz/libs/core"
	"github.com/fotomxq/coll-mz/libs/router"
)

//日志处理
var log core.Log
//错误
var err error

//启动脚本
func main(){
	//设定错误前缀
	log.SetErrorPrefix("发生一个错误 : ")
	//获取配置数据
	config := new(core.Config)
	err = config.LoadFile("content/config/config.json")
	if err != nil{
		log.AddErrorLog(err)
		return
	}
	//激活服务器
	router.Router()
}