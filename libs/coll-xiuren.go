//收集秀人网数据模块
package collmzLibs

import (
	"github.com/fotomxq/ftmp-libs"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

//通用错误
var err error
//通用页面操作对象
var collPage GetPageData
//文件存储路径
var collFileSrc string

//运行的脚本
func CollXiuren() (bool,error){
	//读取配置文件
	config := new(ftmplibs.Config)
	err = config.LoadFile("content/config/coll-xiuren.json")
	if err != nil{
		return false,err
	}
	//保存参数
	collFileSrc = config.Data["collFileSrc"].(string)
	//连接到数据库
	collBool,err := collPage.Create(config.Data)
	if err != nil || collBool == false{
		return collBool,err
	}
	defer collPage.CloseDB()
	//获取首页
	firstURL := config.Data["parentURL"].(string) + "1" + config.Data["parentURLAppend"].(string)
	firstDoc,err := goquery.NewDocument(firstURL)
	if err != nil{
		return false,err
	}
	firstNodePageMax,err := firstDoc.Find(".info").Html()
	if err != nil{
		return false,err
	}
	pageInfos := strings.Split(firstNodePageMax,"/")
	pageMax,err := strconv.Atoi(pageInfos[1])
	if err != nil{
		return false,err
	}
	//根据页码长度，建立循环
	for page := 1 ; page <= pageMax ; page++ {
		pageURL := config.Data["parentURL"].(string) + strconv.Itoa(page) + config.Data["parentURLAppend"].(string)
		pageBool,err := collXPage(pageURL)
		if err != nil || pageBool == false{
			return pageBool,err
		}
	}
	return true,nil
}

//解析某分页数据，并收集子页面数据
func collXPage(pageURL string) (bool,error){
	doc,err := goquery.NewDocument(pageURL)
	if err != nil{
		return false,err
	}
	nodes := doc.Find(".content")
	for i := range nodes.Nodes{
		node := nodes.Eq(i).Children()
		href,hrefExist := node.Attr("href")
		title,titleExist := node.Attr("title")
		if hrefExist == false || titleExist == false{
			return false,nil
		}
		_,err = collChildPage(href,title)
		if err != nil{
			return false,err
		}
	}
	return true,nil
}

//收集子页面数据
func collChildPage(url string,title string) (bool,error){
	//获取子页面内容
	doc,err := goquery.NewDocument(url)
	if err != nil{
		return false,err
	}
	//查找是否存在视频节点，存在则存入数据库
	videoHtml,err := doc.Find("#jwplayer_1").Html()
	if err != nil {
		return false,err
	}
	if videoHtml != ""{
		collPage.AddVideoData(url)
	}
	//如果不存在照片，则返回
	_,exist := doc.Find(".photoThum").Children().Attr("href")
	if exist == false{
		collPage.AddNewData(url,title)
		return true,nil
	}
	//构建存储路径
	fileSrc := collFileSrc + "/" + title
	//查找所有图片，进行遍历归档操作
	nodes := doc.Find(".photoThum")
	for i := range nodes.Nodes {
		imgURL,exist := nodes.Eq(i).Children().Attr("href")
		if exist == false{
			return false,nil
		}
		_,err = collPage.SaveUrlToFile(imgURL,fileSrc)
		if err != nil{
			return false,err
		}
	}
	//完成操作，将数据存入数据库
	collPage.AddNewData(url,title)
	return true,nil
}