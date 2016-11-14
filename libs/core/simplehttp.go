//网络连接包
//该模块主要用于非解析页面
//需要解析页面，建议使用github.com/PuerkitoBio/goquery
package core

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//网络通讯类构建
type SimpleHttp struct{
	sendUrl string
	sendParams map[string][]string
	proxyOn bool
	f FileOperate
}

//设定URL地址
func (simpleHttp *SimpleHttp) SetSendUrl(sendUrl string){
	simpleHttp.sendUrl = sendUrl
}

//设定参数
func (simpleHttp *SimpleHttp) SetSendParams(sendParams map[string][]string) {
	simpleHttp.sendParams = sendParams
}

//设定是否启动代理
func (simpleHttp *SimpleHttp) SetProxy(setOn bool){
	simpleHttp.proxyOn = setOn
}

//Get数据
// url - 网络地址 ; param - 参数 (url.value)
func (simpleHttp *SimpleHttp) Get() ([]byte,error){
	var Url *url.URL
	var err error
	Url,err = url.Parse(simpleHttp.sendUrl)
	if err != nil{
		return nil, err
	}
	//转换格式
	var urlParams url.Values = simpleHttp.sendParams
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = urlParams.Encode()
	resp,err := http.Get(Url.String())
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//post数据
func (simpleHttp *SimpleHttp) Post() ([]byte,error){
	var urlParams url.Values = simpleHttp.sendParams
	resp,err := http.PostForm(simpleHttp.sendUrl, urlParams)
	if err != nil{
		return nil ,err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//将URL保存到文件
func (simpleHttp *SimpleHttp) Save(fileSrc string) (bool,error){
	//获取URL内容
	urlC,urlErr := simpleHttp.Get()
	if urlErr != nil{
		return false,urlErr
	}
	//将数据写入文件
	writeFileBool,writeErr := simpleHttp.f.WriteFile(fileSrc,urlC)
	return writeFileBool,writeErr
}

//获取URL末尾文件名称和类型
//注意某些文件无法获取，如经过特殊处理的URL路径
//返回结构 []string key : 0-全名 1-仅名称 2-仅类型
func (simpleHttp *SimpleHttp) GetURLNameType(url string) []string{
	var res []string
	urls := strings.Split(url,"/")
	if len(urls) < 1 {
		return res
	}
	res[0] = urls[len(urls) - 1]
	if res[0] == ""{
		return res
	}
	names := strings.Split(res[0],".")
	res[1] = names[0]
	res[2] = names[1]
	return res
}