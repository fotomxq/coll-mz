package controller

import (
	"strconv"
	"github.com/PuerkitoBio/goquery"
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
		nodes := doc.Find("#pins li a")
		for nodeKey := range nodes.Nodes{
			//Resolve to get subpages
			thisChildrenPageURL,b := nodes.Eq(nodeKey).Attr("href")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			childrenDoc,err := goquery.NewDocument(thisChildrenPageURL)
			if err != nil{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			//Loop traversal gets all subpage nodes
			///////////////
		}
		page += 1
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}