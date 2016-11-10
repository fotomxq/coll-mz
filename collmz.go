package main

import(
	"github.com/fotomxq/ftmp-libs"
	"github.com/fotomxq/coll-mz/libs"
)

//日志处理
var log ftmplibs.Log
//错误
//某些声明下，直接复用而不是重新声明error变量
var err error

//启动脚本
func main(){
	//开始提示
	log.AddLog("* _ * * _ * 脚本开始运行 * _ * * _ *")
	log.AddLog("初始化参数中...")
	//配置基本参数
	log.SetErrorPrefix("发生一个错误 : ")
	//xiuren
	collmzLibs.CollXiuren()
}