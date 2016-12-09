package core

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

//该模块用于通过外网、内部获取IP地址
//使用IPAddrsGetExternal()前请确保联网
//使用方式：
// 直接调用函数即可获取IP地址，失败将返回0.0.0.0
//依赖内部模块：core.SendLog()
//依赖外部库：无

//通过外部网络获取IP地址
//return string IP地址
func IPAddrsGetExternal() string {
	var url string = "http://myexternalip.com/raw"
	resp, err := http.Get(url)
	if err != nil {
		SendLog(err.Error())
		return "0.0.0.0"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendLog(err.Error())
		return "0.0.0.0"
	}
	html := string(body)
	if err != nil {
		SendLog(err.Error())
		return "0.0.0.0"
	}
	html = strings.Replace(html, " ", "", -1)
	html = strings.Replace(html, "\n", "", -1)
	return html
}

//通过内部获取IP地址
//return string IP地址
func IPAddrsGetInternal() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		SendLog(err.Error())
		return "0.0.0.0"
	}
	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}
