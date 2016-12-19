package core

//版本号及最后一次修改日期：
// 2016.12.15
//该模块是core核心，用于绝大部分内部处理
//该文件用于声明模块内全局变量

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
func init() {
	PathSeparator = GetPathSep()
}
