package core

//该文件只声明该模块类型

//通用错误
var err error

//通用系统文件路径
var PathSeparator string

//日志操作句柄
var Log *LogOperate

//核心类
type ControlCore struct {

}

//初始化该模块
func init(){
	PathSeparator = GetPathSep()
}