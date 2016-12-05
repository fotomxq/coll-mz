package controller

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

//Collect feig
func (this *Coll) CollFeig() {
	//Gets the object
	thisChildren := &this.collList.feig
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//start
	var errNum int = 0
	var page int = 1
	//Traverse the page
	for {
		if thisChildren.status == false{
			return
		}
		pageURL := thisChildren.url + strconv.Itoa(page)
		pageDoc,err := goquery.NewDocument(pageURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open") + " , error page doc.",err)
			errNum += 1
			break
		}
		//Traverse the columns in the page
		nodes := pageDoc.Find("#loop-square").Find(".entry-content").Find("img")
		for nodeKey := range nodes.Nodes{
			if thisChildren.status == false{
				return
			}
			node := nodes.Eq(nodeKey)
			nodeTitle,b := node.Attr("alt")
			if b == false{
				errNum += 1
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children-empty") + " , error node title.",nil)
			}
			nodeURL,b := node.Attr("src")
			if b == false{
				errNum += 1
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children-empty") + " , error node src.",nil)
			}
			newID := collOperate.AutoCollFile(nodeURL,nodeTitle,strconv.Itoa(page),0)
			if newID < 1{
				errNum += 1
				collOperate.NewLog(collOperate.lang.Get("coll-error-new-id") + " , error node new id.",nil)
			}
			errNum = 0
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			break
		}
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}