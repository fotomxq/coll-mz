package main

import(
	"github.com/fotomxq/ftmp-libs"
)

//日志处理
var log ftmplibs.Log
//配置处理
var config ftmplibs.Config
//错误
//某些声明下，直接复用而不是重新声明error变量
var err error

//启动脚本
func main(){
	//开始提示
	log.AddLog("* _ * * _ * 脚本开始运行 * _ * * _ *")
	log.AddLog("初始化参数中...")
	//读取配置信息
	//err = config.LoadFile("config22.json")
	if err != nil{
		log.AddLog("发生一个错误:")
		log.AddErrorLog(err)
	}
	//test
	data := new(GetPageData)
	params
	data.Create(params)
}