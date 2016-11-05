//网络连接包
//废弃，采集器使用github.com/PuerkitoBio/goquery解决，但保留该模块
package libs

import (
    "io/ioutil"
    "net/http"
    "net/url"
)

//网络通讯类构建
type SimpleHttp struct{
	sendUrl string
	sendParams map[string][]string
	proxyOn bool
}

//设定URL地址
func (this *SimpleHttp) SetSendUrl(sendUrl string){
	this.sendUrl = sendUrl
}

//设定参数
func (this *SimpleHttp) SetSendParams(sendParams map[string][]string) {
	this.sendParams = sendParams
}

//设定是否启动代理
func (this *SimpleHttp) SetProxy(setOn bool){
	this.proxyOn = setOn
}



//Get数据
// url - 网络地址 ; param - 参数 (url.value)
func (this *SimpleHttp) Get() (res []byte, err error){
	var Url *url.URL
	Url,err = url.Parse(this.sendUrl)
	if err != nil{
		return nil, err
	}
	//转换格式
	var urlParams url.Values = this.sendParams
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
func (this *SimpleHttp) Post() (res []byte, err error){
	var urlParams url.Values = this.sendParams
	resp,err := http.PostForm(this.sendUrl, urlParams)
	if err != nil{
		return nil ,err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}