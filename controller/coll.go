package controller

import (
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"strconv"
)

//coll struct
type Coll struct {
	db *Database
	dataSrc string
	lang *Language
	collList CollList
	collListK []string
}

//Collector list
type CollList struct {
	local CollChildren
	jiandan CollChildren
	jiandanIndex CollChildren
	xiuren CollChildren
	meizitu CollChildren
}

//Collector list children
type CollChildren struct {
	status bool
	source string
}

//Initialize the collector
func (this *Coll) init(db *Database,dataSrc string){
	this.db = db
	this.dataSrc = dataSrc
	this.collList.local = CollChildren{
		status : false,
		source : "local",
	}
	this.collList.jiandan = CollChildren{
		status : false,
		source : "jiandan",
	}
	this.collList.jiandanIndex = CollChildren{
		status : false,
		source : "jiandan-index",
	}
	this.collList.xiuren = CollChildren{
		status : false,
		source : "xiuren",
	}
	this.collList.meizitu = CollChildren{
		status : false,
		source : "meizitu",
	}
	this.collListK = []string{"local","jiandan","jiandan-index","xiuren","meizitu"}
}

////////////////////////////////////////////////
//External methods
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
	case "":
		//Run all collectors
		log.NewLog(this.lang.Get("coll-run-all"),nil)
		for _,v := range this.collListK{
			this.Run(v)
		}
		break
	default:
		break
	}
}

//Gets the running state
func (this *Coll) GetStatus() (string){
	res := make(map[string]interface{})
	for _,value := range this.collListK{
		valueC := make(map[string]interface{})
		c := this.GetCollChildren(value)
		valueC["status"] = c.status
		valueC["source"] = c.source
		src := this.dataSrc + GetPathSep() + "coll-log" + GetPathSep() + c.source + ".log"
		logContentByte,err := LoadFile(src)
		if err != nil{
			valueC["log"] = ""
			log.NewLog("",err)
		}
		valueC["log"] = string(logContentByte)
		res[value] = valueC
	}
	resJson,err := json.Marshal(res)
	if err != nil{
		log.NewLog("",err)
		return ""
	}
	resJsonStr := string(resJson)
	return resJsonStr
}

////////////////////////////////////////////////
//This part is a variety of Web site data collector
///////////////////////////////////////////////

//Collect jiandan data
func (this *Coll) CollJiandan() {
	//Gets the object
	thisChildren := &this.collList.jiandan
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//start
	nextURL := "http://jandan.net/ooxx"
	var b bool
	errNum := 0
	for{
		//Get the page data
		doc,err := goquery.NewDocument(nextURL)
		if err != nil{
			collOperate.NewLog("",err)
			continue
		}
		//Gets a list of nodes
		nodes := doc.Find(".commentlist").Children()
		for i := range nodes.Nodes {
			aNodes := nodes.Eq(i).Find(".text p").Children().Filter("a")
			for j := range aNodes.Nodes {
				nodeURL, b := aNodes.Eq(j).Attr("href")
				if b == false {
					collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
					errNum += 1
					continue
				}
				newID := collOperate.AutoCollFile(nodeURL,"","",0)
				if newID < 1{
					errNum += 1
					continue
				}
				errNum = 0
			}
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			break
		}
		//Gets the next URL address
		nextURL,b = doc.Find(".previous-comment-page").Eq(0).Next().Attr("href")
		if b == false{
			break
		}
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

//Collect jiandan index data
func (this *Coll) CollJiandanIndex() {
	//Gets the object
	thisChildren := &this.collList.jiandanIndex
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//start
	indexURL := "http://jandan.net"
	//Get the page data
	doc,err := goquery.NewDocument(indexURL)
	if err != nil{
		collOperate.NewLog("",err)
		thisChildren.status = false
		return
	}
	//Gets a list of nodes
	list := doc.Find("#list-girl").Children().Children().Children().Find(".acv_comment").Find("p").Find("a")
	for i := range list.Nodes {
		node := list.Eq(i)
		v, b := node.Attr("href")
		if b == false{
			collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
			continue
		}
		if v == ""{
			collOperate.NewLog(collOperate.lang.Get("coll-error-get-children-empty"),nil)
			continue
		}
		newID := collOperate.AutoCollFile(v,"","",0)
		if newID < 1{
			continue
		}
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

//Collect xiuren data
func (this *Coll) CollXiuren() {
	//Gets the object
	thisChildren := &this.collList.xiuren
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

//Collect Mzitu data
func (this *Coll) CollMeizitu() {
	//Gets the object
	thisChildren := &this.collList.meizitu
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

//Collect local data
func (this *Coll) CollLocal() {
	//Gets the object
	thisChildren := &this.collList.local
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

//Gets the CollListChildren handle
func (this *Coll) GetCollChildren(name string) *CollChildren{
	switch name {
	case "local":
		return &this.collList.local
		break
	case "jiandan" :
		return &this.collList.jiandan
		break
	case "jiandan-index" :
		return &this.collList.jiandanIndex
		break
	case "xiuren":
		return &this.collList.xiuren
		break
	case "meizitu":
		return &this.collList.meizitu
		break
	}
	return &this.collList.local
}

//Create new CollListChildren
func (this *Coll) CreateCollListChildren(name string) CollChildren{
	var thisCollListChildren CollChildren
	thisCollListChildren.source = name
	thisCollListChildren.status = false
	return thisCollListChildren
}

//coll start
func (this *Coll) CollStart(thisChildren *CollChildren,collOperate *CollOperate) bool {
	collOperate.init(this.db,this.dataSrc,thisChildren,this.lang)
	if thisChildren.status == true{
		collOperate.NewLog(collOperate.lang.Get("coll-is-running"),nil)
		return false
	}
	thisChildren.status = true
	collOperate.ClearLog()
	collOperate.NewLog(collOperate.lang.Get("coll-run"),nil)
	return true
}

//coll end
func (this *Coll) CollEnd(thisChildren *CollChildren,collOperate *CollOperate) {
	thisChildren.status = false
	collOperate.NewLog(collOperate.lang.Get("coll-stop"),nil)
	if collOperate.collNum > 0{
		collOperate.NewLog(collOperate.lang.Get("coll-num") + strconv.Itoa(collOperate.collNum),nil)
	}else{
		collOperate.NewLog(collOperate.lang.Get("coll-no"),nil)
	}
}