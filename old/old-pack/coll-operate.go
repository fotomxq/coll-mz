package controller

import (
	"strconv"
	"time"
	"database/sql"
)

type CollOperate struct {
	//database
	db *Database
	//Database field columns
	fields []string
	//Data folder path
	dataSrc string
	sourceSrc string
	//Authentication module
	matchString MatchString
	//Language data
	lang *Language
	//Acquisition object
	collChildren *CollChildren
	//log
	log Log
	//coll num
	collNum int
	//The number is ignored
	collIgnoreNum int
	//status
	status bool
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
func (this *CollOperate) init(db *Database,dataSrc string,collChildren *CollChildren,lang *Language) {
	this.db = db
	this.dataSrc = dataSrc
	this.sourceSrc = this.dataSrc + GetPathSep() + "coll-file" + GetPathSep() + collChildren.source
	this.fields = []string{"id","parent","star","sha1","src","source","url","name", "file_type","size","coll_time", }
	this.collChildren = collChildren
	logSrc := this.dataSrc + GetPathSep() + "coll-log"
	this.log.init(logSrc,true,false,true,true,true,false)
	this.log.SetOneFileName(collChildren.source)
	this.log.SetIsForward(true)
	this.lang = lang
	this.collNum = 0
	this.status = true
}

//View the data list
func (this *CollOperate) ViewDataList(parent int64,star int,searchTitle string,page int,max int,sort int,desc bool) ([]map[string]string,bool) {
	var result []map[string]string
	if this.status == false{
		return result,false
	}
	query := "select `id`,`star`,`name`,`file_type` from `coll` where `parent` = ?"
	if star > 0{
		query += " and `star` = ?"
	}
	if searchTitle != ""{
		query += " and `name` = ?"
	}
	var sortStr string
	switch(sort){
	case 0:
		sortStr = "id"
		break
	case 2:
		sortStr = "star"
		break
	case 7:
		sortStr = "name"
		break
	case 9:
		sortStr = "size"
		break
	case 10:
		sortStr = "coll_time"
		break
	default:
		sortStr = "id"
		break
	}
	query += " " + this.db.GetPageSortStr(page,max,sortStr,desc)
	var rows *sql.Rows
	if star > 0 && searchTitle != ""{
		rows,err = this.db.db.Query(query,parent,star,searchTitle)
	}
	if star > 0{
		rows,err = this.db.db.Query(query,parent,star)
	}
	if searchTitle != ""{
		rows,err = this.db.db.Query(query,parent,searchTitle)
	}
	if star < 1 && searchTitle == ""{
		rows,err = this.db.db.Query(query,parent)
	}
	if err != nil{
		this.NewLog("",err)
		return result,false
	}
	for {
		b := rows.Next()
		if b == false{
			break
		}
		var thisRes CollFields
		err = rows.Scan(&thisRes.id,&thisRes.star,&thisRes.name,&thisRes.file_type)
		if err != nil{
			log.NewLog("",err)
			return result,false
		}
		idStr := strconv.FormatInt(thisRes.id,10)
		starStr := strconv.Itoa(thisRes.star)
		thisResArr := map[string]string{
			"id" : idStr,
			"star" : starStr,
			"name" : thisRes.name,
			"file-type" : thisRes.file_type,
		}
		result = append(result,thisResArr)
	}
	return result,true
}

//View the data src
func (this *CollOperate) ViewDataSrc(id int64) (string) {
	if this.status == false{
		return ""
	}
	query := "select `src` from `coll` where `id` = ?"
	row := this.db.db.QueryRow(query,id)
	var res CollFields
	err = row.Scan(&res.src)
	if err != nil{
		log.NewLog("",err)
		return ""
	}
	return res.src
}

//Empty a data set
func (this *CollOperate) ClearColl() bool {
	//Close the database and rebuild the database
	err = this.db.Close()
	if err != nil{
		this.NewLog("",err)
		return false
	}
	//Delete all file data
	err = DeleteFile(this.sourceSrc)
	if err != nil{
		this.NewLog(this.lang.Get("coll-err-clear-coll") + this.sourceSrc,err)
		return false
	}
	//Delete the log data
	this.ClearLog()
	//return
	return true
}

//Create an upper level ID
func (this *CollOperate) AutoCollParent(parentTitle string,parentURL string) (int64,string){
	if this.status == false{
		return 0,""
	}
	//Create parent directory data
	parentSha1 := this.matchString.GetSha1(parentTitle + parentURL)
	if parentSha1 == ""{
		this.NewLog(this.lang.Get("coll-error-sha1") + " parent : " + parentTitle + " , url : " + parentURL,nil)
		return 0,""
	}
	//Check parent sha1 if the data already exists
	if this.CheckDataSha1(parentSha1) == true{
		this.NewLog(this.lang.Get("coll-error-repeat-sha1") + parentURL + " , sha1 : " + parentSha1,nil)
		return -1,""
	}
	//Create parent database data
	parentID := this.CreateNewData(0,parentSha1,"",parentURL,parentTitle,"folder","0")
	if parentID > 0{
		this.NewLog(this.lang.Get("coll-new-id") + strconv.FormatInt(parentID,10) + " , URL : " + parentURL,nil)
	}else{
		this.NewLog(this.lang.Get("coll-error-move-file") + parentURL,nil)
	}
	return parentID,parentSha1
}

//Automatically collect collection
func (this *CollOperate) AutoCollParentFiles(parentTitle string,parentURL string,urls []string) int64 {
	if this.status == false{
		return 0
	}
	//Create parent directory data
	parentSha1 := this.matchString.GetSha1(parentTitle + parentURL)
	if parentSha1 == ""{
		this.NewLog(this.lang.Get("coll-error-sha1") + " parent : " + parentTitle + " , url : " + parentURL,nil)
		return 0
	}
	//Check parent sha1 if the data already exists
	if this.CheckDataSha1(parentSha1) == true{
		this.NewLog(this.lang.Get("coll-error-repeat-sha1") + parentURL + " , sha1 : " + parentSha1,nil)
		return -1
	}
	//Create parent database data
	parentID := this.CreateNewData(0,parentSha1,"",parentURL,parentTitle,"folder","0")
	if parentID > 0{
		this.NewLog(this.lang.Get("coll-new-id") + strconv.FormatInt(parentID,10) + " , URL : " + parentURL,nil)
	}else{
		this.NewLog(this.lang.Get("coll-error-move-file") + parentURL,nil)
	}
	//Traverse the generated subfile data
	for _,value := range urls{
		newID := this.AutoCollFile(value,parentTitle,parentSha1,parentID)
		if newID < 1 && newID != -1 {
			continue
		}
	}
	return parentID
}

//Automatically builds new data
// return -1 - The data has been collected
// return 0 - Acquisition failed
// return >0 - Acquires the new ID of the data
func (this *CollOperate) AutoCollFile(url string,name string,parent string,parentID int64) int64 {
	if this.status == false{
		return 0
	}
	//Check url if the data already exists
	if(this.CheckDataURL(url) == true){
		this.NewLog(this.lang.Get("coll-error-repeat-url") + url,nil)
		return -1
	}
	this.NewLog(this.lang.Get("coll-new-url") + url,nil)
	//Buffer files and get information
	cacheFileInfo,b := this.GetCacheFIleInfo(url)
	if b == false{
		this.NewLog(this.lang.Get("coll-error-cache-info"),nil)
		return 0
	}
	if name == ""{
		name = cacheFileInfo["name"]
	}
	//Check sha1 if the data already exists
	if this.CheckDataSha1(cacheFileInfo["sha1"]) == true{
		_ = this.DeleteFile(cacheFileInfo["cache-src"])
		this.NewLog(this.lang.Get("coll-error-repeat-sha1") + url + " , sha1 : " + cacheFileInfo["sha1"],nil)
		return -1
	}
	//Transfer buffer files
	fileSrc := this.SaveCacheToFile(parent,cacheFileInfo)
	if fileSrc == ""{
		this.NewLog(this.lang.Get("coll-error-move-file") + url,nil)
		_ = this.DeleteFile(cacheFileInfo["cache-src"])
		return 0
	}
	//Establish database data
	newID := this.CreateNewData(parentID,cacheFileInfo["sha1"],fileSrc,url,name,cacheFileInfo["type"],cacheFileInfo["size"])
	if newID > 0{
		this.NewLog(this.lang.Get("coll-new-id") + strconv.FormatInt(newID,10) + " , URL : " + url,nil)
	}else{
		this.NewLog(this.lang.Get("coll-error-move-file") + url,nil)
	}
	return newID
}

//Check whether the file is duplicated
func (this *CollOperate) CheckDataSha1(sha1 string) bool {
	query := "select `id` from `coll` where `sha1` = ?"
	row := this.db.db.QueryRow(query,sha1)
	var id int64
	err = row.Scan(&id)
	if err != nil{
		return false
	}
	if id > 0{
		this.collIgnoreNum += 1
		return true
	}
	return false
}

//Check whether the file is duplicated
func (this *CollOperate) CheckDataURL(url string) bool {
	query := "select `id` from `coll` where `url` = ?"
	row := this.db.db.QueryRow(query,url)
	var id int64
	err = row.Scan(&id)
	if err != nil{
		return false
	}
	if id > 0{
		this.collIgnoreNum += 1
		return true
	}
	return false
}

//Create a new data record
func (this *CollOperate) CreateNewData(parent int64,sha1 string,src string,url string,name string,fileType string,size string) int64 {
	query := "insert into `coll`(" + this.db.GetFieldsToStr(this.fields) + ") values(null,?,0,?,?,?,?,?,?,?,?)"
	stmt,err := this.db.db.Exec(query,parent,sha1,src,this.collChildren.source,url,name,fileType,size,this.GetNowTimeUnix())
	if err != nil{
		this.NewLog("",err)
		return 0
	}
	newID,err := stmt.LastInsertId()
	if err != nil{
		this.NewLog("",err)
		return 0
	}
	this.collNum += 1
	return newID
}

//Update the data record information
func (this *CollOperate) UpdateData(id int64,parent int64,star int,name string) bool {
	query := "update `coll` set `parent` = ? , `star` = ? , `name` = ? where `id` = ?"
	stmt,err := this.db.db.Exec(query,parent,star,name,id)
	if err != nil{
		log.NewLog("",err)
		return false
	}
	row,err := stmt.RowsAffected()
	if err != nil{
		log.NewLog("",err)
		return false
	}
	return row > 0
}

//Delete a collection of data
func (this *CollOperate) DeleteData(id int64) bool {
	b,err := this.db.Delete("coll",id)
	if err != nil{
		log.NewLog("",err)
		return false
	}
	return b > 0
}

//Gets the str SHA1 value
func (this *CollOperate) GetStrSha1(str string) string {
	return this.matchString.GetSha1(str)
}

//Gets the file SHA1 value
func (this *CollOperate) GetFileSha1(src string) string {
	res,err := GetFileSha1(src)
	if err != nil{
		log.NewLog("",err)
		return ""
	}
	return res
}

//send log
// name string - Log name, corresponding to language configuration name
// append string - Additional content
// err error - error
func (this *CollOperate) NewLog(msg string,err error) {
	if msg != ""{
		msg = this.collChildren.source + " ~ " + msg + "<br />"
	}
	if err != nil && msg == ""{
		msg = this.collChildren.source + " Error ~ " + msg + "<br />"
	}
	if IsFile(this.log.lastSrc) == true{
		//if log file size > 20KB
		if GetFileSize(this.log.lastSrc) > 20480 {
			content,err := LoadFile(this.log.lastSrc)
			if err != nil{
				log.NewLog("",err)
			}
			contentStr := string(content)
			newContent := this.matchString.SubStr(contentStr,0,len(contentStr) / 2)
			err = WriteFile(this.log.lastSrc,[]byte(newContent))
			if err != nil{
				log.NewLog("",err)
			}
		}
	}
	this.log.NewLog(msg,err)
}

//Get the latest log content
func (this *CollOperate) GetLog(name string) string{
	if IsFile(this.log.lastSrc) == false{
		return ""
	}
	content,err := LoadFile(this.log.lastSrc)
	if err != nil{
		log.NewLog("",err)
	}
	return string(content)
}

//clear log
func (this *CollOperate) ClearLog(){
	src := this.dataSrc + GetPathSep() + "coll-log" + GetPathSep() + this.collChildren.source + ".log"
	cByte := []byte("")
	err = WriteFile(src,cByte)
	if err != nil{
		log.NewLog("",err)
	}
}

//Gets the buffer file information
func (this *CollOperate) GetCacheFIleInfo(url string) (map[string]string,bool) {
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
		this.NewLog("",err)
		return result,false
	}
	urls := GetURLNameType(url)
	if urls["full-name"] == "" || urls["only-name"] == "" || urls["type"] == ""{
		this.NewLog("Error ~ full-name : " + urls["full-name"] + " ; only-name : " + urls["only-name"] + " ; type : " + urls["type"],nil)
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
		this.NewLog("",err)
		return result,false
	}
	err = WriteFile(cacheSrc,cByte)
	if err != nil{
		this.NewLog("",err)
		return result,false
	}
	result["cache-src"] = cacheSrc
	//Gets additional file information
	result["size"] = strconv.FormatInt(GetFileSize(cacheSrc),10)
	result["sha1"],err = GetFileSha1(cacheSrc)
	if err != nil{
		this.NewLog("",err)
		return result,false
	}
	return result,true
}

//Save the buffer file to the file database
// parentName string - The name of the previous archive
// cacheFileInfo map[string]string - File information
// source string - data source
func (this *CollOperate) SaveCacheToFile(parentName string,cacheFileInfo map[string]string) string {
	//Create a directory path
	t := time.Now()
	dirSrc := this.sourceSrc + GetPathSep() + t.Format("200601")
	if parentName != ""{
		dirSrc += GetPathSep() + parentName
	}
	err = CreateDir(dirSrc)
	if err != nil{
		this.NewLog("",err)
		return ""
	}
	//Creates a file path
	fileSrc := dirSrc + GetPathSep() + cacheFileInfo["full-name"]
	//Transfer files
	err = CutFile(cacheFileInfo["cache-src"],fileSrc)
	if err != nil{
		this.NewLog("",err)
		return ""
	}
	return fileSrc
}

//Delete a file
func (this *CollOperate) DeleteFile(src string) bool {
	err = DeleteFile(src)
	if err != nil{
		this.NewLog("",err)
		return false
	}
	return true
}

//get now time unix
func (this *CollOperate) GetNowTimeUnix() int64{
	t := time.Now()
	return t.Unix()
}