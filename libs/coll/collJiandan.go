package coll

import (
	"github.com/PuerkitoBio/goquery"
)

//煎蛋采集脚本
func CollJiandan() (bool, error) {
	// *** 采集首页 ***
	//解析页面内容
	CollPg.SendLog("开始采集煎蛋网首页...")
	doc, err := goquery.NewDocument("http://jandan.net")
	if err != nil {
		CollPg.SendErrorLog(err)
	}
	//获取妹子图ID
	docList := doc.Find("#list-girl").Children().Children().Children().Find(".acv_comment").Find("p").Find("a")
	//遍历节点
	for i := range docList.Nodes {
		node := docList.Eq(i)
		v, b := node.Attr("href")
		if b == false {
			continue
		}
		if v != "" {
			vNames := CollPg.simhttp.GetURLNameType(v)
			b,err := CollPg.AutoAddData("jiandan", v, vNames[1], false)
			if err != nil{
				continue
			}
			if b == ""{
				continue
			}
		}
	}

	// *** 采集子页面 ***

	//返回成功
	return true, nil
}
