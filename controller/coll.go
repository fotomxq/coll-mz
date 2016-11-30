package controller

import (
	"strconv"
	"time"
	"github.com/PuerkitoBio/goquery"
)

//coll struct
type Coll struct {
	//database
	db *Database
	//Database field columns
	fields []string
	//Data folder path
	dataSrc string
	//Authentication module
	matchString MatchString
	//A list of collector names
	collNames map[string]string
	//Collector operating status list
	collStatus map[string]bool
	//Language data
	lang *Language
	//Internal log module
	log Log
}

//Collector Type
type CollFields struct {
	id int64
	parent int64
	star int
	sha1 string
	src string
	source string
	url string
	name string
	file_type string
	size int64
	coll_time string
}

//Initialize the collector
func (this *Coll) init(db *Database,dataSrc string){
	this.db = db
	this.dataSrc = dataSrc
	this.fields = []string{
		"id",
		"parent",
		"star",
		"sha1",
		"src",
		"source",
		"url",
		"name",
		"file_type",
		"size",
		"coll_time",
	}
	this.collNames = map[string]string{
		"local" : "local.file",
		"jiandan" : "jandan.net page",
		"jiandan-index" : "jandan.net index",
		"xiuren" : "xiuren.org",
		"meizitu" : "mzitu.com",
	}
	this.collStatus = map[string]bool{
		"local" : false,
		"jiandan" : false,
		"jiandan-index" : false,
		"xiuren" : false,
		"meizitu" : false,
	}
	logSrc := dataSrc + GetPathSep() + "coll-log"
	this.log.init(logSrc,true,false,true,true,true,false)
	this.log.isOneFile = true
	this.log.isForward = true
}

////////////////////////////////////////////////
//External methods
///////////////////////////////////////////////

//Runs a collector
func (this *Coll) Run(name string) {
	if name != ""{
		//If it is running, it returns
		if this.collStatus[name] == true{
			this.NewLog("coll-is-run",this.collNames[name],nil)
			return
		}
		//Set the operating status
		this.collStatus[name] = true
		//Output a script that prompts you to run the Collector
		this.NewLog("coll-run",this.collNames[name],nil)
	}
	//Select the corresponding content
	switch name {
	case "local":
		go this.CollLocal()
		break
	case "jiandan" :
		go this.CollJiandan()
		break
	case "jiandan-index" :
		go this.CollJiandanIndex()
		break
	case "xiuren":
		go this.CollXiuren()
		break
	case "meizitu":
		go this.CollMeizitu()
		break
	case "":
		//Run all collectors
		for k := range this.collNames {
			this.Run(k)
		}
		break
	default:
		break
	}
}

//Get the latest log content
func (this *Coll) GetLog() string{
	content,err := LoadFile(this.log.logDirSrc + GetPathSep() + "log.log")
	if err != nil{
		log.NewLog("",err)
	}
	return string(content)
}

//
func (this *Coll) ViewDataList(parent int64,star int,searchTitle string,page int,max int,sort int,desc bool) {

}

//
func (this *Coll) ViewData(id int64) (CollFields,bool) {
	query := "select `id` from `coll` where `id` = ?"
	row := this.db.db.QueryRow(query,id)
	var res CollFields
	err = row.Scan(&res.id,&res.sha1,&res.coll_time,&res.file_type,&res.name,&res.parent,&res.size,&res.source,&res.src,&res.star,&res.url)
	if err != nil{
		this.NewLog("","",err)
		return res,false
	}
	return res,true
}

//clear log
func (this *Coll) ClearLog(){
	src := this.dataSrc + GetPathSep() + "coll-log" + GetPathSep() + "log.log"
	err = WriteFile(src,[]byte(""))
	if err != nil{
		log.NewLog("",err)
	}
}

//Gets the running state
func (this *Coll) GetStatus(){
	var res map[string]string
	for i := range this.collStatus{
		if this.collStatus[i] == true{
			res[i] = this.collNames[i] + this.lang.Get("coll-is-running")
		}else{
			res[i] = this.collNames[i] + this.lang.Get("coll-is-stopped")
		}
	}
}

////////////////////////////////////////////////
//This part is a variety of Web site data collector
///////////////////////////////////////////////

//Collect jiandan data
func (this *Coll) CollJiandan() {
	nextURL := "http://jandan.net/ooxx"
	source := "jiandan"
	var b bool
	errNum := 0
	for{
		//Get the page data
		doc,err := goquery.NewDocument(nextURL)
		if err != nil{
			this.NewLog("coll-jiandan-err-a","",err)
			continue
		}
		//Gets a list of nodes
		nodes := doc.Find(".commentlist").Children()
		for i := range nodes.Nodes {
			aNodes := nodes.Eq(i).Find(".text p").Children().Filter("a")
			for j := range aNodes.Nodes {
				nodeURL, b := aNodes.Eq(j).Attr("href")
				if b == false {
					this.NewLog("coll-jiandan-err-b","",nil)
					errNum += 1
					continue
				}
				newID := this.AutoCollFile(nodeURL,source,"","",0)
				if newID < 1{
					errNum += 1
					continue
				}
				errNum = 0
			}
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			this.NewLog("coll-jiandan-err-c","",nil)
			break
		}
		//Gets the next URL address
		nextURL,b = doc.Find(".previous-comment-page").Eq(0).Next().Attr("href")
		if b == false{
			break
		}
	}
	//finish
	this.collStatus["jiandan"] = false
}

//Collect jiandan index data
func (this *Coll) CollJiandanIndex() {
	indexURL := "http://jandan.net"
	source := "jiandan-index"
	//Get the page data
	doc,err := goquery.NewDocument(indexURL)
	if err != nil{
		this.NewLog("","",err)
		this.collStatus["jiandan"] = false
		return
	}
	//Gets a list of nodes
	list := doc.Find("#list-girl").Children().Children().Children().Find(".acv_comment").Find("p").Find("a")
	for i := range list.Nodes {
		node := list.Eq(i)
		v, b := node.Attr("href")
		if b == false{
			this.NewLog("coll-jiandan-index-err-a","",nil)
			continue
		}
		if v == ""{
			this.NewLog("coll-jiandan-index-err-b","",nil)
			continue
		}
		this.NewLog("coll-jiandan-index-find-new-url",v,nil)
		newID := this.AutoCollFile(v,source,"","",0)
		if newID < 1{
			continue
		}
	}
	//finish
	this.collStatus["jiandan"] = false
}

//Collect xiuren data
func (this *Coll) CollXiuren() {
	//ok.
	this.collStatus["xiuren"] = false
}

//Collect Mzitu data
func (this *Coll) CollMeizitu() {
	//ok.
	this.collStatus["meizitu"] = false
}

//Collect local data
func (this *Coll) CollLocal() {
	//ok.
	this.collStatus["local"] = false
}

////////////////////////////////////////////////
//Relevant function modules used in the collector
///////////////////////////////////////////////

//Automatically builds new data
func (this *Coll) AutoCollFile(url string,source string,name string,parent string,parentID int64) int64 {
	//Buffer files and get information
	cacheFileInfo,b := this.GetCacheFIleInfo(url)
	if b == false{
		this.NewLog("coll-add-err-file-info",url,nil)
		return 0
	}
	if name == ""{
		name = cacheFileInfo["name"]
	}
	//Check if the data already exists
	if this.CheckDataSha1(cacheFileInfo["sha1"]) == true{
		_ = this.DeleteFile(cacheFileInfo["cache-src"])
		return 0
	}
	//Transfer buffer files
	fileSrc := this.SaveCacheToFile(parent,cacheFileInfo,source)
	if fileSrc == ""{
		this.NewLog("coll-add-new-data-file-src-failed",url,nil)
		_ = this.DeleteFile(cacheFileInfo["cache-src"])
		return 0
	}
	//Establish database data
	newID := this.CreateNewData(parentID,cacheFileInfo["sha1"],fileSrc,source,url,name,cacheFileInfo["type"],cacheFileInfo["size"])
	if newID > 0{
		this.NewLog("coll-add-new-data",strconv.FormatInt(newID,10),nil)
	}else{
		this.NewLog("coll-add-new-data-failed",url,nil)
	}
	return newID
}

//Check whether the file is duplicated
func (this *Coll) CheckDataSha1(sha1 string) bool {
	query := "select `id` from `coll` where `sha1` = ?"
	row := this.db.db.QueryRow(query,sha1)
	var id int64
	err = row.Scan(&id)
	if err != nil{
		//this.NewLog("","",err)
		return false
	}
	if id > 0{
		return true
	}
	return false
}

//Check the collector's current operating status
// true - running
// false - Not running
func (this *Coll) CheckCollStatus(collName string) bool{
	if collName == ""{
		for i := range this.collStatus {
			if this.collStatus[i] == true{
				return true
			}
		}
		return false
	}else{
		return this.collStatus[collName]
	}
}

//Create a new data record
func (this *Coll) CreateNewData(parent int64,sha1 string,src string,source string,url string,name string,fileType string,size string) int64 {
	query := "insert into `coll`(" + this.db.GetFieldsToStr(this.fields) + ") values(null,?,0,?,?,?,?,?,?,?,now())"
	stmt,err := this.db.db.Exec(query,parent,sha1,src,source,url,name,fileType,size)
	if err != nil{
		this.NewLog("","",err)
		return 0
	}
	newID,err := stmt.LastInsertId()
	if err != nil{
		this.NewLog("","",err)
		return 0
	}
	return newID
}

//Update the data record information
func (this *Coll) UpdateData(id int64,parent int64,star int,name string) bool {
	query := "update `coll` set `parent` = ? , `star` = ? , `name` = ? where `id` = ?"
	stmt,err := this.db.db.Exec(query,parent,star,name,id)
	if err != nil{
		this.NewLog("","",err)
		return false
	}
	row,err := stmt.RowsAffected()
	if err != nil{
		this.NewLog("","",err)
		return false
	}
	return row > 0
}

//Delete a collection of data
func (this *Coll) DeleteData(id int64) bool {
	b,err := this.db.Delete("coll",id)
	if err != nil{
		this.NewLog("","",err)
		return false
	}
	return b > 0
}

//Gets the str SHA1 value
func (this *Coll) GetStrSha1(str string) string {
	return this.matchString.GetSha1(str)
}

//Gets the file SHA1 value
func (this *Coll) GetFileSha1(src string) string {
	res,err := GetFileSha1(src)
	if err != nil{
		this.NewLog("","",err)
		return ""
	}
	return res
}

//send log
// name string - Log name, corresponding to language configuration name
// append string - Additional content
// err error - error
func (this *Coll) NewLog(name string,append string,err error) {
	if name != ""{
		content := this.lang.Get(name) + append + "<br />"
		this.log.NewLog(content,nil)
	}
	if err != nil{
		contentErr := err.Error() + "<br />"
		this.log.NewLog(contentErr,nil)
	}
}

//Gets the buffer file information
func (this *Coll) GetCacheFIleInfo(url string) (map[string]string,bool) {
	//Build a basic information framework
	result := map[string]string{
		"full-name" : "",
		"name" : "",
		"type" : "",
		"size" : "",
		"sha1" : "",
		"cacheSrc" : "",
	}
	//Build the cache path
	cacheDir := this.dataSrc + GetPathSep() + "cache"
	err = CreateDir(cacheDir)
	if err != nil{
		this.NewLog("","",err)
		return result,false
	}
	urls := GetURLNameType(url)
	if urls["full-name"] == "" || urls["only-name"] == "" || urls["type"] == ""{
		this.NewLog("coll-cache-err-urls","full-name : " + urls["full-name"] + " ; only-name : " + urls["only-name"] + " ; type : " + urls["type"],err)
		return result,false
	}
	result["full-name"] = urls["full-name"]
	result["name"] = urls["only-name"]
	result["type"] = urls["type"]
	var cacheSrc string
	for {
		cacheSrc = cacheDir + GetPathSep() + this.matchString.GetRandStr(10000) + "." + urls["type"]
		if IsFile(cacheSrc) == false{
			break
		}
	}
	//Download the file from the URL
	var params map[string][]string
	cByte,err := SimpleHttpGet(url,params)
	if err != nil{
		this.NewLog("","",err)
		return result,false
	}
	err = WriteFile(cacheSrc,cByte)
	if err != nil{
		this.NewLog("","",err)
		return result,false
	}
	result["cache-src"] = cacheSrc
	//Gets additional file information
	result["size"] = strconv.FormatInt(GetFileSize(cacheSrc),10)
	result["sha1"],err = GetFileSha1(cacheSrc)
	if err != nil{
		this.NewLog("","",err)
		return result,false
	}
	return result,true
}

//Save the buffer file to the file database
// parentName string - The name of the previous archive
// cacheFileInfo map[string]string - File information
// source string - data source
func (this *Coll) SaveCacheToFile(parentName string,cacheFileInfo map[string]string,source string) string {
	//Create a directory path
	t := time.Now()
	dirSrc := this.dataSrc + GetPathSep() + "coll-file" + GetPathSep() + source + GetPathSep() + t.Format("200601")
	if parentName != ""{
		dirSrc += GetPathSep() + parentName
	}
	err = CreateDir(dirSrc)
	if err != nil{
		this.NewLog("","",err)
		return ""
	}
	//Creates a file path
	fileSrc := dirSrc + GetPathSep() + cacheFileInfo["full-name"]
	//Transfer files
	err = CutFile(cacheFileInfo["cache-src"],fileSrc)
	if err != nil{
		this.NewLog("","",err)
		return ""
	}
	return fileSrc
}

//Delete a file
func (this *Coll) DeleteFile(src string) bool {
	err = DeleteFile(src)
	if err != nil{
		this.NewLog("","",err)
		return false
	}
	return true
}