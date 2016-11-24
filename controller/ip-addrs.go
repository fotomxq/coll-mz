package controller

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

//Obtain an IP address from the network
func IPAddrsGetExternal() string {
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
	html = strings.Replace(html, " ", "", -1)
	html = strings.Replace(html, "\n", "", -1)
	return html
}

//Get the local IP address
func IPAddrsGetInternal() string {
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
