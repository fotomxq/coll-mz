package controller

import (
	"strconv"
)

//coll struct
type Coll struct {
	db *Database
	dataSrc string
	lang *Language
	collList CollList
	collListV map[string]*CollChildren
}

//Collector list
type CollList struct {
	local CollChildren
	jiandan CollChildren
	jiandanIndex CollChildren
	xiuren CollChildren
	meizitu CollChildren
	xiuhaotu CollChildren
}

//Collector list children
type CollChildren struct {
	status bool
	source string
	url string
}

//Initialize the collector
func (this *Coll) init(db *Database,dataSrc string){
	this.db = db
	this.dataSrc = dataSrc
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
		url : "",
	}
	this.collList.xiuhaotu = CollChildren{
		status : false,
		source : "xiuhaotu",
		url : "http://showhaotu.xyz/explore/?list=images&sort=date_desc&page=",
	}
	this.collListV = map[string]*CollChildren{
		"local" : &this.collList.local,
		"jiandan" : &this.collList.jiandan,
		"jiandan-index" : &this.collList.jiandanIndex,
		"xiuren" : &this.collList.xiuren,
		"meizitu" : &this.collList.meizitu,
		"xiuhaotu" : &this.collList.xiuhaotu,
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
		collOperate.NewLog(collOperate.lang.Get("coll-num") + strconv.Itoa(collOperate.collNum) + collOperate.lang.Get("coll-ig-num") + strconv.Itoa(collOperate.collIgnoreNum),nil)
	}else{
		collOperate.NewLog(collOperate.lang.Get("coll-no"),nil)
	}
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