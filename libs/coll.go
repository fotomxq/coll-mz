//收集数据处理器
package collmzLibs

//通用错误
var err error
//通用页面操作对象
var collPage GetPageData
//文件存储路径
var collFileSrc string

//使用各模块将采用自动化处理，请不要直接调用对应模块，在配置文件中设定即可
func Coll(obj string)(bool,error){
	switch obj{
		case "xiuren":
			return CollXiuren()
			break
		case "jiandan":
			return CollJiandan()
			break
	}
	return false,nil
}