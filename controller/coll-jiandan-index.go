package controller

import "github.com/PuerkitoBio/goquery"

//Collect jiandan.net index data
func (this *Coll) CollJiandanIndex() {
	//Gets the object
	thisChildren := &this.collList.jiandanIndex
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//start
	indexURL := thisChildren.url
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
		if thisChildren.status == false{
			return
		}
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
