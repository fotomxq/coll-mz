package controller

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

//Collect jiandan.net page data
func (this *Coll) CollJiandan() {
	//Gets the object
	thisChildren := &this.collList.jiandan
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	defer this.CollEnd(thisChildren,&collOperate)
	//start
	nextURL := thisChildren.url
	var b bool
	errNum := 0
	var parent string
	for{
		if thisChildren.status == false{
			return
		}
		//Get the page data
		doc,err := goquery.NewDocument(nextURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open"),err)
			continue
		}
		//Get the number of pages
		var nowPage string
		nowPage,err = doc.Find(".current-comment-page").Eq(0).Html()
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-now-page"),err)
			break
		}
		parent = strings.Replace(nowPage,"[","",-1)
		parent = strings.Replace(parent,"]","",-1)
		if parent == ""{
			errNum += 1
			collOperate.NewLog(collOperate.lang.Get("coll-error-now-page"),nil)
			continue
		}
		//Gets a list of nodes
		nodes := doc.Find(".commentlist").Children()
		for i := range nodes.Nodes {
			if thisChildren.status == false{
				return
			}
			aNodes := nodes.Eq(i).Find(".text p").Children().Filter("a")
			for j := range aNodes.Nodes {
				if thisChildren.status == false{
					return
				}
				nodeURL, b := aNodes.Eq(j).Attr("href")
				if b == false {
					collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
					errNum += 1
					continue
				}
				newID := collOperate.AutoCollFile(nodeURL,"",parent,0)
				if newID < 1 && newID != -1{
					errNum += 1
					continue
				}
				errNum = 0
			}
		}
		//More than 10 times the error is to exit
		if errNum > 10 {
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			return
		}
		//Gets the next URL address
		nextURL,b = doc.Find(".previous-comment-page").Eq(0).Attr("href")
		if b == false{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next"),nil)
			break
		}
	}
}
