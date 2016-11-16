//煎蛋网妹子收集模块
package coll

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/fotomxq/ftmp-libs"
)

//运行收集模块
func CollJiandan() (bool,error){
	//读取配置文件
	config := new(ftmplibs.Config)
	err = config.LoadFile("content/config/coll-jiandan.json")
	if err != nil{
		return false,err
	}
	nextURL := config.Data["parentURL"].(string)
	dataSrc := config.Data["collFileSrc"].(string)
	//循环获取页面，首先给予首页
	for i:=0;i<=9000;i++{
		nextURL,err = getPages(nextURL,dataSrc)
		if nextURL == "" || err != nil{
			break
		}
	}
	//返回数据
	return true,err
}

//获取指定页面信息，并返回下一页URL
func getPages(url string,dataSrc string)(string,error){
	//获取页面内容
	doc,err := goquery.NewDocument(url)
	if err != nil{
		return "",err
	}
	//根据页面内容，获取所有子节点
	nodes := doc.Find(".commentlist").Children()
	for i := range nodes.Nodes{
		//获取子节点
		node := nodes.Eq(i)
		//获取子节点的所有img原图URL
		imgNodes := node.Children().Children().Find(".text").Children().Eq(2).Find("a")
		//获取节点的子img标签，之后将文件保存
		for j := range imgNodes.Nodes{
			imgURL,exist := imgNodes.Eq(j).Attr("href")
			if exist == false{
				continue
			}
			names := collPage.GetURLNameType(imgURL)
			src := dataSrc + "/" + string(i) + "-" + string(j) + "." + names[2]
			saveBool,err := saveImgUrl(imgURL,src)
			if err != nil || saveBool == false{
				continue
			}
		}
	}
	//获取下一页URL
	nextPageURL,exist := doc.Find(".cp-pagenavi").Eq(0).Children().Eq(1).Attr("href")
	if exist == false{
		return "",nil
	}
	//返回下一页URL
	return nextPageURL,nil
}

//将URL保存到文件
func saveImgUrl(url string,src string)(bool,error){
	//解析出名称
	names := collPage.GetURLNameType(url)
	name := names[0]
	//在数据库查询是否存在
	checkBool,err := collPage.CheckRepeat(url,name)
	if err != nil || checkBool == true{
		return false,err
	}
	//保存到文件
	saveBool,err := collPage.SaveUrlToFile(url,src)
	if saveBool == false || err != nil{
		return saveBool,err
	}
	//建立数据
	return collPage.AddNewData(url,name)
}