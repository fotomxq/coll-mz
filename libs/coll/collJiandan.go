package coll

import (
	"github.com/PuerkitoBio/goquery"
)

//煎蛋采集脚本
func CollJiandan() (bool, error) {
	// *** 采集首页 ***
	//解析页面内容
	CollPg.SendLog("开始采集煎蛋网首页...")
	doc,err := goquery.NewDocument("http://jandan.net")
	if err != nil{
		CollPg.SendErrorLog(err)
	}
	//获取妹子图ID
	docList := doc.Find("#list-girl").Children().Children().Children().Find(".acv_comment").Find("p").Find("a")
	//遍历节点
	for i := range docList.Nodes{
		node := docList.Eq(i)
		v,b := node.Attr("href")
		if b == false{
			continue
		}
		if v != ""{
			CollPg.SendLog("发现新的文件 : " + v)
			s,err := CollPg.AutoAddData("jiandan",v,"",false)
			if err != nil{
				CollPg.SendErrorLog(err)
				CollPg.SendLog("文件保存失败...")
				continue
			}
			if s == ""{
				CollPg.SendLog("未知错误，文件保存失败...")
				continue
			}
			CollPg.SendLog("文件保存成功，路径：" + s)
		}
	}

	// *** 采集子页面 ***

	//返回成功
	return true, nil
}
