package controller

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
		"jiandan" : "煎蛋网",
		"xiuren" : "秀人网",
		"meizitu" : "妹子图",
		"local" : "本地文件",
	}
	this.collStatus = map[string]bool{
		"jiandan" : false,
		"xiuren" : false,
		"meizitu" : false,
		"local" : false,
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
		if this.CheckCollStatus(name) == true{
			this.NewLog("coll-is-run",nil)
			return
		}
		//Set the operating status
		this.collStatus[name] = true
		//Output a script that prompts you to run the Collector
		this.NewLog("coll-run-" + name,nil)
	}
	//Select the corresponding content
	switch name {
	case "jiandan" :
		go this.CollJiandan()
		break
	case "xiuren":
		go this.CollXiuren()
		break
	case "meizitu":
		go this.CollMeizitu()
		break
	case "local":
		go this.CollLocal()
		break
	case "":
		//Run all collectors
		for i := range this.collNames {
			this.Run(this.collNames[i])
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
		this.NewLog("",err)
	}
	return string(content)
}

//
func (this *Coll) ViewDataList(parent int64,star int,searchTitle string) {

}

//
func (this *Coll) ViewData(id int64) (CollFields,bool) {
	query := "select `id` from `coll` where `id` = ?"
	row := this.db.db.QueryRow(query,id)
	var res CollFields
	err = row.Scan(&res.id,&res.sha1,&res.coll_time,&res.file_type,&res.name,&res.parent,&res.size,&res.source,&res.src,&res.star,&res.url)
	if err != nil{
		this.NewLog("",err)
		return res,false
	}
	return res,true
}

////////////////////////////////////////////////
//This part is a variety of Web site data collector
///////////////////////////////////////////////

//Collect jiandan data
func (this *Coll) CollJiandan() {
	//ok.
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

//Check whether the file is duplicated
func (this *Coll) CheckDataSha1(sha1 string) bool {
	query := "select `id` from `coll` where `sha1` = ?"
	row := this.db.db.QueryRow(query,sha1)
	var id int64
	err = row.Scan(&id)
	if err != nil{
		this.NewLog("",err)
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

//
func (this *Coll) CreateNewData(parent int64,sha1 string,src string,source string,url string,name string,fileType string,size string) {

}

//
func (this *Coll) UpdateData(parent int64,) {
	//...
}

//Delete a collection of data
func (this *Coll) DeleteData(id int64) bool {
	b,err := this.db.Delete("coll",id)
	if err != nil{
		this.NewLog("",err)
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
		return ""
	}
	return res
}

//send log
func (this *Coll) NewLog(name string,err error) {
	if name != ""{
		content := this.lang.Get(name) + "<br />"
		this.log.NewLog(content,nil)
	}
	if err != nil{
		contentErr := err.Error() + "<br />"
		this.log.NewLog(contentErr,nil)
	}
}