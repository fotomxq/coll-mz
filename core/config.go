package core

import "encoding/json"

//该模块用于读取和写入配置文件信息
//配置文件采用JSON格式
//使用方式：直接调用函数即可使用
//依赖内部模块：core.SendLog()
//依赖外部库：无

//读取配置文件
//param src string 配置文件路径
//return map[string]interface{},bool 配置信息，是否成功读取
func LoadConfig(src string) (map[string]interface{}, bool) {
	var res map[string]interface{}
	c, b := LoadFile(src)
	if b == false {
		return res, false
	}
	err = json.Unmarshal(c, &res)
	if err != nil {
		SendLog(err.Error())
		return res, false
	}
	return res, true
}

//保存配置文件
//param src string 配置文件路径
//param data map[string]interface{} 配置信息
//return bool 是否写入成功
func SaveConfigFile(src string, data map[string]interface{}) bool {
	dataJson, err := json.Marshal(data)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	return WriteFile(src, dataJson)
}
