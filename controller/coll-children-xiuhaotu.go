package controller

import (
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

//Collect xiuhaotu data
func (this *Coll) CollXiuhaotu() {
	//Gets the object
	thisChildren := &this.collList.xiuhaotu
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	defer this.CollEnd(thisChildren,&collOperate)
	//start
	var page int = 0
	var errNum int = 0
	//Traverse the page count
	for{
		//Forced interrupt handling
		if thisChildren.status == false{
			return
		}
		//Get the page
		pageURL := thisChildren.url + strconv.Itoa(page)
		pageDoc,err := goquery.NewDocument(pageURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-err-next"),nil)
			break
		}
		this.DebugErrorHTMLDoc("coll-xiuhaotu-doc-page",pageDoc)
		//$('.pad-content-listing').eq(0).find('img')
		nodes := pageDoc.Find(".pad-content-listing").Eq(0).Find("img")
		this.DebugErrorHTMLNode("coll-xiuhaotu-nodes",nodes)
		//Traverse all nodes
		for nodeKey := range nodes.Nodes{
			//Forced interrupt handling
			if thisChildren.status == false{
				return
			}
			//get node data
			node := nodes.Eq(nodeKey)
			nodeTitle,b := node.Attr("alt")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
				errNum += 1
				continue
			}
			nodeURL,b := node.Attr("src")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-large-img"),nil)
				errNum += 1
				continue
			}
			nodeURL = strings.Replace(nodeURL,".md","",-1)
			newID := collOperate.AutoCollFile(nodeURL,nodeTitle,strconv.Itoa(page),0)
			if newID < 1 && newID != -1{
				errNum += 1
				continue
			}
		}
		//More than 10 times the error is to exit
		if errNum > 10{
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			return
		}
		//Gets the next page
		page += 1
	}
}

