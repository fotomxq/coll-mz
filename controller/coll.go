package controller

import (
	"strconv"
	"time"
	"github.com/PuerkitoBio/goquery"
)

//coll struct
type Coll struct {
	db *Database
	dataSrc string
	lang *Language
	collList CollList
	collListV map[string]*CollChildren
	collDatabaseTemplateSrc string
	collErrSrc string
}

//Collector list
type CollList struct {
	local CollChildren
	jiandan CollChildren
	jiandanIndex CollChildren
	xiuren CollChildren
	meizitu CollChildren
	xiuhaotu CollChildren
	feig CollChildren
}

//Collector list children
type CollChildren struct {
	status bool
	source string
	url string
	db Database
	dbStatus bool
}

//Initialize the collector
func (this *Coll) init(db *Database,dataSrc string,collDatabaseTemplateSrc string){
	this.db = db
	this.dataSrc = dataSrc
	this.collDatabaseTemplateSrc = collDatabaseTemplateSrc
	if IsFile(this.dataSrc) == false{
		err = CreateDir(this.dataSrc + GetPathSep() + "database")
		if err != nil{
			return
		}
	}
	this.collErrSrc = this.dataSrc + GetPathSep() + "coll-error"
	if IsFolder(this.collErrSrc) == false {
		err = CreateDir(this.collErrSrc)
		if err != nil {
			log.NewLog("",err)
			return
		}
	}
	this.collList.local = CollChildren{
		status : false,
		source : "local",
		url : "",
	}
	this.collList.jiandan = CollChildren{
		status : false,
		source : "jiandan",
		url : "http://jandan.net/ooxx",
	}
	this.collList.jiandanIndex = CollChildren{
		status : false,
		source : "jiandan-index",
		url : "http://jandan.net",
	}
	this.collList.xiuren = CollChildren{
		status : false,
		source : "xiuren",
		url : "",
	}
	this.collList.meizitu = CollChildren{
		status : false,
		source : "meizitu",
		url : "http://www.mzitu.com/page/",
	}
	this.collList.xiuhaotu = CollChildren{
		status : false,
		source : "xiuhaotu",
		url : "http://showhaotu.xyz/explore/?list=images&sort=date_desc&page=",
	}
	this.collList.feig = CollChildren{
		status : false,
		source : "feig",
		url : "http://www.girl13.com/page/",
	}
	this.collListV = map[string]*CollChildren{
		"local" : &this.collList.local,
		"jiandan" : &this.collList.jiandan,
		"jiandan-index" : &this.collList.jiandanIndex,
		"xiuren" : &this.collList.xiuren,
		"meizitu" : &this.collList.meizitu,
		"xiuhaotu" : &this.collList.xiuhaotu,
		"feig" : &this.collList.feig,
	}
}

////////////////////////////////////////////////
//This part is the direct access method of the router
///////////////////////////////////////////////

//Runs a collector
func (this *Coll) Run(name string) {
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
	case "xiuhaotu":
		go this.CollXiuhaotu()
		break
	case "feig":
		go this.CollFeig()
		break
	case "":
		//Run all collectors
		log.NewLog(this.lang.Get("coll-run-all"),nil)
		for key := range this.collListV{
			this.Run(key)
		}
		break
	default:
		break
	}
}

//Gets the running state
func (this *Coll) GetStatus() (map[string]interface{},bool){
	res := make(map[string]interface{})
	var b bool = false
	for key := range this.collListV{
		valueC := make(map[string]interface{})
		c := this.GetCollChildren(key)
		valueC["status"] = c.status
		valueC["source"] = c.source
		valueC["url"] = c.url
		src := this.dataSrc + GetPathSep() + "coll-log" + GetPathSep() + c.source + ".log"
		if IsFile(src) == true{
			logContentByte,err := LoadFile(src)
			if err != nil{
				valueC["log"] = ""
				log.NewLog("",err)
			}else{
				valueC["log"] = string(logContentByte)
			}
		}else{
			valueC["log"] = ""
		}
		res[key] = valueC
	}
	b = true
	return res,b
}

//Empty a data set
func (this *Coll) ClearColl(name string) bool {
	exisit := this.CheckCollExisit(name)
	if exisit == false{
		return false
	}
	thisChildren := this.GetCollChildren(name)
	var collOperate CollOperate
	collOperate.init(this.db,this.dataSrc,thisChildren,this.lang)
	b := collOperate.ClearColl()
	return b
}

//Change the collector operating status
func (this *Coll) ChangeStatus(name string,b bool) bool {
	exisit := this.CheckCollExisit(name)
	if exisit == false{
		return false
	}
	thisChildren := this.GetCollChildren(name)
	thisChildren.status = b
	return true
}

//auto coll task
func (this *Coll) AutoTask(){
	task := time.NewTicker(time.Minute * 120)
	for _ = range task.C{
		for _,collValue := range this.collListV{
			if collValue.status == true{
				continue
			}
			this.Run(collValue.source)
		}
	}
}

////////////////////////////////////////////////
//This section is the method used internally by the Collector
///////////////////////////////////////////////

//Create new CollListChildren
func (this *Coll) CreateCollListChildren(name string) CollChildren{
	var thisCollListChildren CollChildren
	thisCollListChildren.source = name
	thisCollListChildren.status = false
	return thisCollListChildren
}

//coll start
func (this *Coll) CollStart(thisChildren *CollChildren,collOperate *CollOperate) bool {
	if this.CollConnectDB(thisChildren,collOperate) == false{
		return false
	}
	collOperate.init(&thisChildren.db, this.dataSrc, thisChildren, this.lang)
	if thisChildren.status == true {
		collOperate.NewLog(collOperate.lang.Get("coll-is-running"), nil)
		return false
	}
	thisChildren.status = true
	collOperate.ClearLog()
	collOperate.NewLog(collOperate.lang.Get("coll-run"),nil)
	return true
}

//coll connect db
func (this *Coll) CollConnectDB(thisChildren *CollChildren,collOperate *CollOperate) bool{
	if thisChildren.dbStatus == true {
		return true
	}
	dbSrc := this.dataSrc + GetPathSep() + "database" + GetPathSep() + thisChildren.source + ".sqlite"
	if IsFile(dbSrc) == false {
		b, err := CopyFile(this.collDatabaseTemplateSrc, dbSrc)
		if err != nil || b == false {
			collOperate.NewLog(collOperate.lang.Get("coll-error-database-create"), err)
			return false
		}
	}
	err = thisChildren.db.Connect("sqlite3", dbSrc)
	if err != nil {
		collOperate.NewLog(collOperate.lang.Get("coll-error-database-connect"), err)
		return false
	}
	thisChildren.dbStatus = true
	return true
}

//coll end
func (this *Coll) CollEnd(thisChildren *CollChildren,collOperate *CollOperate) {
	this.CollCloseDB(thisChildren,collOperate)
	thisChildren.status = false
	collOperate.NewLog(collOperate.lang.Get("coll-stop"),nil)
	if collOperate.collNum > 0{
		collOperate.NewLog(collOperate.lang.Get("coll-num") + strconv.Itoa(collOperate.collNum) + collOperate.lang.Get("coll-ig-num") + strconv.Itoa(collOperate.collIgnoreNum),nil)
	}else{
		collOperate.NewLog(collOperate.lang.Get("coll-no"),nil)
	}
}

//coll close database
func (this *Coll) CollCloseDB(thisChildren *CollChildren,collOperate *CollOperate){
	err = thisChildren.db.Close()
	if err != nil{
		collOperate.NewLog(collOperate.lang.Get("coll-error-database-close"),err)
	}
	thisChildren.dbStatus = false
}

//Check that the Collector is present
func (this *Coll) CheckCollExisit(name string) bool{
	for key := range this.collListV{
		if name == key{
			return true
		}
	}
	return false
}

//Gets the CollListChildren handle
func (this *Coll) GetCollChildren(name string) *CollChildren{
	return this.collListV[name]
}

//Sending an HTML to the error file was originally text for debugging.
func (this *Coll) SendErrorHTML(name string,html *goquery.Selection) {
	src := this.collErrSrc + GetPathSep() + name + ".html"
	c,err := html.Html()
	if err != nil{
		log.NewLog("",err)
		return
	}
	err = WriteFile(src,[]byte(c))
	if err != nil{
		log.NewLog("",err)
	}
}

//Sending an HTML to the error file was originally text for debugging.
func (this *Coll) SendErrorHTMLStr(name string,html string) {
	src := this.collErrSrc + GetPathSep() + name + ".html"
	err = WriteFile(src,[]byte(html))
	if err != nil{
		log.NewLog("",err)
	}
}