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