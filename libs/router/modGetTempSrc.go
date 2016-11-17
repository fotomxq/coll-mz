package router

//获取模板路径
func modGetTempSrc(name string) string {
	return "template" + fileSep + name
}
