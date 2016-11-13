//获取IP地址模块
//该包直接调用函数即可
package core

import (
	"net"
	"io/ioutil"
	"net/http"
)

//通过网络获取IP地址
func GetExternal() string {
	var url string = "http://myexternalip.com/raw"
	resp, err := http.Get(url)
	if err != nil {
		return "0.0.0.0"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "0.0.0.0"
	}
	html := string(body)
	if err != nil {
		return "0.0.0.0"
	}
	return html
}

//获取本机IP地址
func GetInternal() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}
