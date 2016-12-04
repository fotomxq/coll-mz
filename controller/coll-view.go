package controller

//Gets the list to the JSON data format
func (this *Coll) ViewList(collName string,parent int64,star int,searchTitle string,page int,max int,sort int,desc bool) ([]map[string]string,bool) {
	var result []map[string]string
	//Gets the object
	thisChildren := this.GetCollChildren(collName)
	var collOperate CollOperate
	if this.CollConnectDB(thisChildren,&collOperate) == false{
		return result,false
	}
	//get data
	result,b := collOperate.ViewDataList(parent,star,searchTitle,page,max,sort,desc)
	if b == false{
		log.NewLog("Failed to get database list data.",nil)
		return result,false
	}
	//close db
	this.CollCloseDB(thisChildren,&collOperate)
	//return
	return result,true
}
