package main

import(
	"reptilexiuren/libs"
	"github.com/PuerkitoBio/goquery"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"crypto/sha1"
	"fmt"
)

//配置信息
var config libs.Config

func main(){
	//初始化日志模块
	log := new(libs.Log)
	log.SetNewLogType(2)
	log.SetDirType(0)
	//读取配置文件
	configErr := config.LoadFile("content/data-db/config.json")
	if configErr != nil{
		log.AddErrorLog(configErr)
		return
	}
	//开始执行脚本
	var startUrl string = "http://www.xiuren.org/page-"
	//如果发现首页无变化，则说明整体无新增内容，直接跳出
	firstUrl := startUrl + "1.html"
	checkFirstChange := checkFirstPageChange(firstUrl)
	if checkFirstChange == false {
		//return
		//由于刚刚运行，为了避免后续无法有效获取，所以暂时不中断脚本
		log.AddLog("Check first page no update.")
	}
	//遍历总页面1-100页数
	pageStart,pageStartErr := strconv.Atoi(config.Data["pageStart"].(string))
	pageEnd,pageEndErr := strconv.Atoi(config.Data["pageEnd"].(string))
	if pageStartErr != nil{
		log.AddErrorLog(pageStartErr)
		return
	}
	if pageEndErr != nil{
		log.AddErrorLog(pageEndErr)
		return
	}
	for pageI := pageStart ; pageI <= pageEnd ; pageI ++ {
		//输出日志
		log.AddLog("Page total : 100 , Get page : " + strconv.Itoa(pageI))
		//生成URL地址
		pageUrl := startUrl + strconv.Itoa(pageI) + ".html"
		pageList,pageListErr := getPageList(pageUrl)
		if pageListErr != nil{
			continue
		}
		//遍历获取的页面
		for listI := range pageList{
			//输出该页面信息
			log.AddLog("Get Page : " + pageList[listI][1] + " , url : " + pageList[listI][0])
			//获取页面内容并保存到文件
			saveBool, saveErr := saveChildrenPage(pageList[listI])
			if saveErr != nil{
				log.AddLog("Save page failed.")
				continue
			}
			if saveBool == false{
				log.AddLog("Save page false.")
			}else{
				log.AddLog("Save page ok.")
			}
		}
	}
}

//获取当前页面所有子页面连接
func getPageList(url string) (urlList [][]string,err error){
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return urlList,err
	}
	node := doc.Find(".content")
	//遍历所有子链接
	for i := range node.Nodes {
		thisNode := node.Eq(i).Children()
		thisHref, hrefExist := thisNode.Attr("href")
		thisTitile,titleExist := thisNode.Attr("title")
		if hrefExist == true && titleExist == true {
			//检查链接标题是否存在，如果存在则跳过添加
			checkBool,checkErr := checkRepeat(thisHref,thisTitile)
			if checkErr != nil{
				return urlList,checkErr
			}
			if checkBool == true{
				continue
			}
			//如果该链接标题不存在，则将链接和标题添加到组
			urlList = append(urlList,[]string{thisHref,thisTitile})
		}
	}
	return urlList,nil
}

//将指定页面保存
func saveChildrenPage(urls []string) (bool,error){
	//拆解变量
	var url string = urls[0]
	var title string = urls[1]
	//声明新的日志
	log := new(libs.Log)
	//读取页面数据
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return false,err
	}
	//检查是否存在视频
	videoHtml, videoErr := doc.Find("#jwplayer_1").Html()
	if videoErr != nil{
		log.AddErrorLog(videoErr)
	}
	if videoHtml != "" {
		//如果存在则向数据库添加记录
		videoBool,videoErr := addVideoData(url)
		if videoErr != nil{
			log.AddErrorLog(videoErr)
		}
		if videoBool != false{
			log.AddLog("Add new video page url.")
		}
	}
	//查询是否存在图像
	_,exists := doc.Find(".photoThum").Children().Attr("href")
	if exists != true{
		//如果同时不存在照片和视频，则返回false
		if videoHtml == ""{
			log.AddLog("This page no photo and video.")
			return false, nil
		}
	}else{
		//根据标题建立目录
		mediaSrc := "content/data-media/" + title
		file := new(libs.FileOperate)
		createDirBool := file.CreateDir(mediaSrc)
		if createDirBool != true {
			return false,nil
		}
		//遍历数据，存储入文件
		var hrefUrl string
		var exisitHref bool
		imageMax := doc.Find(".photoThum").Length()
		fmt.Println("Image total : " + strconv.Itoa(imageMax) + ", Save image file : ")
		doc.Find(".photoThum").Each(func(i int, s *goquery.Selection) {
			//遍历获取图像URL地址
			hrefUrl,exisitHref = s.Children().Attr("href")
			if exisitHref == false {
				return
			}
			//将该URL保存到对应目录
			fileSrc := mediaSrc + "/" + strconv.Itoa(i) + ".jpg"
			saveImgToFile(hrefUrl,fileSrc)
			fmt.Print(strconv.Itoa(i) + ",")
		})
	}
	//给数据库添加标识记录
	addBool,addErr := addUrlData(url,title)
	if addErr != nil{
		return false,addErr
	}
	if addBool != true{
		return false,nil
	}
	//返回
	return true,nil
}

//将页面存到文件
//该方法可用于保存图像等数据
func saveImgToFile(url string,src string) (res bool,err error){
	http := new(libs.SimpleHttp)
	http.SetSendUrl(url)
	c,err := http.Get()
	if err != nil{
		return false,err
	}
	f := new(libs.FileOperate)
	t := f.WriteFile(src,c)
	if t != false{
		return false,nil
	}
	return true,nil
}

//检查首页面变化情况
//firstUrl - 首页面URL地址
//return bool - false 不通过 ; true 检查通过
func checkFirstPageChange(firstUrl string) bool{
	log := new(libs.Log)
	http := new(libs.SimpleHttp)
	http.SetSendUrl(firstUrl)
	firstHtml,firstErr := http.Get()
	//如果获取失败，则说明整体逻辑存在问题或网络故障，即退出执行
	if firstErr != nil{
		return false
	}
	//计算在线页面的SHA1值
	sha := sha1.New()
	sha.Write(firstHtml)
	sha1C := sha.Sum(nil)
	sha1Str := string(sha1C)
	//声明文件模块
	file := new(libs.FileOperate)
	sha1FileSrc := "content/data-db/first-html-sha1.txt"
	sha1FileC := file.ReadFile(sha1FileSrc)
	if sha1FileC != nil{
		sha1FileCStr := string(sha1FileC)
		//如果发现两个值相同，则中断脚本
		if sha1Str == sha1FileCStr{
			log.AddLog("No update.")
			return false
		}
	}
	//将在线页面的SHA1值写入文件
	writeBool := file.WriteFile(sha1FileSrc,sha1C)
	if writeBool == false{
		log.AddLog("Update check , Write file faild.")
		return false
	}
	return true
}

//检查是否重复摘取
//如果存在于数据库，则返回false
func checkRepeat(url string,title string) (bool,error) {
	db,dbErr := connectDb()
	if dbErr != nil{
		return false,dbErr
	}
	defer db.Close()
	var query string = "select id from `loadurl` where `url` = ? and `title` = ?"
	stmt,stmtErr := db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	rows,resErr := stmt.Query(url,title)
	if resErr != nil{
		return false,resErr
	}
	rows.Next()
	var id int
	scanErr := rows.Scan(&id)
	if scanErr != nil{
		return false,nil
	}
	if id > 0{
		return true,nil
	}
	return false,nil
}

//向数据库添加新的数据
func addUrlData(url string,title string) (bool,error){
	db,dbErr := connectDb()
	if dbErr != nil{
		return false,dbErr
	}
	var query string = "insert into `loadurl`(`url`,`title`) values(?,?)"
	stmt,stmtErr := db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	res,resErr := stmt.Exec(url,title)
	if resErr != nil {
		return false,resErr
	}
	newID,newIDErr := res.LastInsertId()
	if newIDErr != nil {
		return false,newIDErr
	}
	if newID > 0{
		return true,nil
	}
	return false,nil
}

//向数据库添加视频记录
func addVideoData(url string) (bool,error){
	db,dbErr := connectDb()
	if dbErr != nil{
		return false,dbErr
	}
	var query string = "insert into `video`(`url`) values(?)"
	stmt,stmtErr := db.Prepare(query)
	if stmtErr != nil{
		return false,stmtErr
	}
	res,resErr := stmt.Exec(url)
	if resErr != nil {
		return false,resErr
	}
	newID,newIDErr := res.LastInsertId()
	if newIDErr != nil {
		return false,newIDErr
	}
	if newID > 0{
		return true,nil
	}
	return false,nil
}

//连接数据库
//采用了Golang diver/sql库，所以未来可替换为其他数据库
func connectDb() (*sql.DB,error){
	dbType := config.Data["databaseType"].(string)
	dbDNS := config.Data["databaseDNS"].(string)
	db,err := sql.Open(dbType,dbDNS)
	if err != nil{
		return db,err
	}
	return db,err
}