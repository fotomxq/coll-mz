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
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open"),err)
			errNum += 1
			break
		}
		//Traverse the columns in the page
		colums := pageDoc.Find("#loop-square").Children()
		for columKey := range colums.Nodes{
			if thisChildren.status == false{
				return
			}
			listNodes := colums.Eq(columKey).Children()
			listPage := 1
			for _ = range listNodes.Nodes{
				if thisChildren.status == false{
					return
				}
				node := listNodes.Eq(listPage).Children().Eq(1).Children().Children().Children()
				nodeTitle,b := node.Attr("alt")
				if b == false{
					errNum += 1
					collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
					this.SendErrorHTML("coll-feig-node",node)
					continue
				}
				nodeURL,b := node.Attr("src")
				if b == false{
					errNum += 1
					collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
					this.SendErrorHTML("coll-feig-node",node)
					continue
				}
				newID := collOperate.AutoCollFile(nodeURL,nodeTitle,strconv.Itoa(page),0)
				if newID < 1{
					errNum += 1
					continue
				}
				listPage += 1
				errNum = 0
			}
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			this.SendErrorHTML("coll-feig-colums",colums)
			break
		}
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}