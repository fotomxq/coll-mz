package controller

//coll struct
type Coll struct {
	db *Database
	fields []string
	dataSrc string
	matchString MatchString
	collNames map[string]string
	collStatus map[string]bool
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
}

//Runs a collector
func (this *Coll) Run(name string) {
	if this.CheckAllCollStatus() == true{
		this.NewLog("The collector is already running, do not repeat it.",nil)
		return
	}
	this.collStatus[name] = true
	switch name {
	case "jiandan" :
		go this.CollJiandan()
		break
	case "xiuren":
		go this.CollXiuren()
		break
	case "meizitu":
		go this.CollMzitu()
		break
	case "local":
		go this.CollLocal()
		break
	case "":
		for i := range this.collNames {
			this.Run(this.collNames[i])
		}
		break
	default:
		break
	}
}

//Collect jiandan data
func (this *Coll) CollJiandan() {

}

//Collect xiuren data
func (this *Coll) CollXiuren() {

}

//Collect Mzitu data
func (this *Coll) CollMzitu() {

}

//Collect local data
func (this *Coll) CollLocal() {

}

//
func (this *Coll) ViewDataList(parent int64,star int,searchTitle string) {

}

//
func (this *Coll) ViewData(id int64) CollFields {
	query := "select `id` from `coll` where `id` = ?"
	row := this.db.db.QueryRow(query,id)
	var res CollFields
	err = row.Scan(&res.id,&res.sha1,&res.coll_time,&res.file_type,&res.name,&res.parent,&res.size,&res.source,&res.src,&res.star,&res.url)
	if err != nil{
		this.NewLog("",err)
		return nil
	}
	return res
}

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
func (this *Coll) CheckAllCollStatus() bool{
	for i := range this.collStatus {
		if this.collStatus[i] == true{
			return true
		}
	}
	return false
}

//
func (this *Coll) CreateNewData(parent int64,sha1 string,src string,source string,url string,name string,fileType string,size string) {

}

//
func (this *Coll) UpdateData(parent int64,) {
	//...
}

//
func (this *Coll) DeleteData(id int64) bool {

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
func (this *Coll) NewLog(content string,err error) {
	log.NewLog(content,err)
}