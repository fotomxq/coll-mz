package controller

import (
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

//Collect Mzitu data
func (this *Coll) CollMeizitu() {
	//Gets the object
	thisChildren := &this.collList.meizitu
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//
	if thisChildren.status == false{
		return
	}
	//start
	page := 1
	errNum := 0
	//Traverse all page numbers
	for {
		thisPageURL := thisChildren.url + strconv.Itoa(page)
		doc,err := goquery.NewDocument(thisPageURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open"),err)
			break
		}
		//Gets and traverses the nodes of all page pages
		nodes := doc.Find("#pins li")
		for nodeKey := range nodes.Nodes{
			//Gets the page URL and title
			thisCURL,b := nodes.Eq(nodeKey).Children().Eq(0).Attr("href")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			thisCTitle,b := nodes.Eq(nodeKey).Children().Eq(1).Html()
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			//Get the page doc
			childrenDoc,err := goquery.NewDocument(thisCURL)
			if err != nil{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			//Loop traversal gets all subpage nodes
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			break
		}
		//Gets the next page
		page += 1
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}