//配置处理模块
//实际是JSON封装处理，可装载配置文件后将配置读出
//格式默认为map模式，如果需要用其他模式，请自行修改代码
//该脚本不依赖任何第三方模块，包括该libs下的其他模块
//注意，读取文件后返回的数据是interface接口类型，所以需要自行对所有配置名称进行接口定义，之后才能使用，强行转换必然失败
//注意，读取的数据需要利用config.Data["name"].(string)转换类型后才能直接使用
package core

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

//配置文件类型
type Config struct {
	//配置文件路径，内部保存
	src string
	//配置数据，可直接修改该数据，但建议先读取出数据
	//也可利用该逻辑，建立新的配置文件
	Data map[string]interface{}
}

//读取配置文件
func (config *Config) LoadFile(src string) (error){
	//保存路径
	config.src = src
	//读取文件
	fd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	//解析JSON
	err = json.Unmarshal(c,&config.Data)
	if err != nil{
		return err
	}
	return nil
}

//将配置写入文件
func (config *Config) SaveFile()(bool,error){
	jsonStr,jsonErr := json.Marshal(config.Data)
	if jsonErr != nil{
		return false,jsonErr
	}
	fileErr := ioutil.WriteFile(config.src, jsonStr, os.ModeAppend)
	if fileErr != nil {
		return false,fileErr
	}
	return true,nil
}
