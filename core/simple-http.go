package core

import (
	"net/http"
	"io/ioutil"
	"net/url"
)

//通过GET获取URL数据
//param sendURL string URL地址
//param params map[string][]string 参数
//return []byte, bool 数据，是否成功
func SimpleHttpGet(sendURL string, params map[string][]string) ([]byte, bool) {
	var urlU *url.URL
	urlU, err = url.Parse(sendURL)
	if err != nil {
		Log.SendLog("core/simple-http.go","0.0.0.0","SimpleHttpGet","url-parse",err.Error())
		return nil, false
	}
	var urlParams url.Values = params
	//对URL进行编码
	urlU.RawQuery = urlParams.Encode()
	resp, err := http.Get(urlU.String())
	if err != nil {
		Log.SendLog("core/simple-http.go","0.0.0.0","SimpleHttpGet","http-getr",err.Error())
		return nil, false
	}
	defer resp.Body.Close()
	res,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		Log.SendLog("core/simple-http.go","0.0.0.0","SimpleHttpGet","ioutil-read-all",err.Error())
		return res,false
	}
	return res,true
}

//通过POST获取URL数据
//param sendURL string URL地址
//param params map[string][]string 参数
//return []byte, bool 数据，是否成功
func SimpleHttpPost(sendURL string, params map[string][]string) ([]byte, bool) {
	var urlParams url.Values = params
	resp, err := http.PostForm(sendURL, urlParams)
	if err != nil {
		Log.SendLog("core/simple-http.go","0.0.0.0","SimpleHttpPost","http-form",err.Error())
		return nil, false
	}
	defer resp.Body.Close()
	res,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		Log.SendLog("core/simple-http.go","0.0.0.0","SimpleHttpPost","ioutil-read-all",err.Error())
		return res,false
	}
	return res,true
}
