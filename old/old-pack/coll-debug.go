package controller

import (
	"github.com/PuerkitoBio/goquery"
)

//Collect debug
func (this *Coll) CollDebug() () {
	//Gets the object
	thisChildren := &this.collList.meizitu
	var collOperate CollOperate
	if this.CollStart(thisChildren, &collOperate) == false {
		return
	}
	//test content src
	collErrorSrc := this.dataSrc + GetPathSep() + "coll-error"
	if IsFolder(collErrorSrc) == false {
		err = CreateDir(collErrorSrc)
		if err != nil {
			collOperate.NewLog("", err)
			return
		}
	}
	//test
	doc, err := goquery.NewDocument("http://www.mzitu.com/80446")
	if err != nil {
		collOperate.NewLog("test get url : " + "http://www.mzitu.com/80446", err)
	} else {
		html, _ := doc.Html()
		_ = WriteFile(collErrorSrc, []byte(html))
		return
	}
	if thisChildren.status == false {
		return
	}
	//finish
	this.CollEnd(thisChildren, &collOperate)
	return
}

//Sending an HTML to the error file was originally text for debugging.
func (this *Coll) DebugErrorHTMLNode(name string,html *goquery.Selection) {
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
func (this *Coll) DebugErrorHTMLDoc(name string,doc *goquery.Document) {
	src := this.collErrSrc + GetPathSep() + name + ".html"
	c,err := doc.Html()
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
func (this *Coll) DebugErrorHTMLStr(name string,html string) {
	src := this.collErrSrc + GetPathSep() + name + ".html"
	err = WriteFile(src,[]byte(html))
	if err != nil{
		log.NewLog("",err)
	}
}