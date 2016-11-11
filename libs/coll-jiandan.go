//煎蛋网妹子收集模块
package collmzLibs

import "github.com/PuerkitoBio/goquery"

//运行收集模块
func CollJiandan() (bool,error){
	//获取首页信息

	//返回数据
	return true,err
}

//获取指定页面信息，并返回下一页URL
func getPageIndex(url string)(string,error){
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
	return collPage.SaveUrlToFile(url,src)
}