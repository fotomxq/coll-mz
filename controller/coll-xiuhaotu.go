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
	//start
	nextPageURL := thisChildren.url
	page := 1
	var errNum int = 0
	var runStop = true
	//Traverse the page count
	for{
		//Forced interrupt handling
		if thisChildren.status == false{
			return
		}
		//Gets the current page URL
		thisURL := nextPageURL + strconv.Itoa(page)
		doc,err := goquery.NewDocument(thisURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open"),err)
			errNum += 1
			break
		}
		nodes := doc.Find(".image-container")
		//Traverse nodes
		for nodeKey := range nodes.Nodes{
			//Forced interrupt handling
			if thisChildren.status == false{
				return
			}
			runStop = false
			//Gets the image URL address
			nodeImgURL,b := nodes.Eq(nodeKey).Children().Attr("src")
			if b == false || nodeImgURL == ""{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				break
			}
			//Analysis URL address, access to large picture URL
			largeImgURL := strings.Replace(nodeImgURL,".md","",-1)
			if largeImgURL == ""{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-large-img") + nodeImgURL,nil)
				errNum += 1
				continue
			}
			//save large picture url
			newID := collOperate.AutoCollFile(largeImgURL,"","",0)
			if newID < 1 && newID != -1{
				errNum += 1
				continue
			}
			//The number of error records is 0
			errNum = 0
		}
		//Forced interrupt handling
		if errNum > 10{
			collOperate.NewLog(collOperate.lang.Get("coll-error-too-many"),nil)
			break
		}
		if runStop == true{
			collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
			break
		}
		//Increase the number of pages
		page += 1
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

