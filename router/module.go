package router

import "../core"

//该文件用于定义通用内部函数部分

//获取template路径
//param name string 路径末尾文件名称
//return string 路径
func getTemplateSrc(name string) string{
	return "." + core.PathSeparator + "template" + core.PathSeparator + name
}
