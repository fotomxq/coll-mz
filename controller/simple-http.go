package controller

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Gets the URL data get mode
func SimpleHttpGet(sendURL string, params map[string][]string) ([]byte, error) {
	var urlU *url.URL
	var err error
	urlU, err = url.Parse(sendURL)
	if err != nil {
		return nil, err
	}
	var urlParams url.Values = params
	//If the parameter has Chinese parameters, this method will be URLEncode.
	urlU.RawQuery = urlParams.Encode()
	resp, err := http.Get(urlU.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

//Gets the URL data post mode
func SimpleHttpPost(sendURL string, params map[string][]string) ([]byte, error) {
	var urlParams url.Values = params
	resp, err := http.PostForm(sendURL, urlParams)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

//Get the file name and type by URL
func GetURLNameType(sendURL string) map[string]string {
	res := map[string]string{
		"full-name": "",
		"only-name": "",
		"type": "",
	}
	urls := strings.Split(sendURL, "/")
	if len(urls) < 1 {
		return res
	}
	res["full-name"] = urls[len(urls)-1]
	if res["full-name"] == "" {
		res["only-name"] = res["full-name"]
		return res
	}
	names := strings.Split(res["full-name"], ".")
	if len(names) < 2 {
		return res
	}
	res["type"] = names[len(names)-1]
	for i := 0 ; i <= len(names) ; i ++{
		if i == len(names) - 1{
			break
		}
		if res["only-name"] == ""{
			res["only-name"] = names[i]
		}else{
			res["only-name"] += "." + names[i]
		}
	}
	return res
}
