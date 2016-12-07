package controller

import "github.com/PuerkitoBio/goquery"

//Collect xiuren data
func (this *Coll) CollXiuren() {
	//Gets the object
	thisChildren := &this.collList.xiuren
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	defer this.CollEnd(thisChildren,&collOperate)
	//start
	var errNum int = 0
	pageURL := thisChildren.url
	for{
		//Forced interrupt handling
		if thisChildren.status == false{
			return
		}
		//get page doc
		pageDoc,err := goquery.NewDocument(pageURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-err-next") + " ~ page doc.",nil)
			break
		}
		//get nodes
		pageNodes := pageDoc.Find("#main").Find("a")
		for pageNodeKey := range pageNodes.Nodes{
			//Forced interrupt handling
			if thisChildren.status == false{
				return
			}
			//get node
			pageNode := pageNodes.Eq(pageNodeKey)
			//get next page
			if pageNodeKey == len(pageNodes.Nodes) - 1{
				var b bool
				pageURL,b = pageNode.Attr("href")
				if b == false{
					collOperate.NewLog(collOperate.lang.Get("coll-err-next") + " ~ get next page",nil)
					return
				}
				continue
			}
			//parent and childrens
			var childrenURLs []string
			//Get the children title and url
			parentTitle,b := pageNode.Children().Attr("alt")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children") + " ~ c 1.",nil)
				errNum += 1
				continue
			}
			parentURL,b := pageNode.Attr("href")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children") + " ~ c 2.",nil)
				errNum += 1
				continue
			}
			//check parent sha1
			parentSha1 := collOperate.matchString.GetSha1(parentTitle + parentURL)
			if parentSha1 == ""{
				collOperate.NewLog(this.lang.Get("coll-error-sha1") + " parent : " + parentTitle + " , url : " + parentURL,nil)
				return
			}
			//Get the children page
			childrenDoc,err := goquery.NewDocument(parentURL)
			if err != nil{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children") + " ~ c 3.",nil)
				errNum += 1
				continue
			}
			childrenNodes := childrenDoc.Find(".photoThum").Find("a")
			//Gets all the picture nodes of the subpage
			for childrenKey := range childrenNodes.Nodes{
				childrenNode := childrenNodes.Eq(childrenKey)
				childrenURL,b := childrenNode.Attr("href")
				if b == false{
					collOperate.NewLog(collOperate.lang.Get("coll-error-get-children") + " ~ c 4",nil)
					errNum += 1
					continue
				}
				childrenURLs = append(childrenURLs,childrenURL)
			}
			//If the collected image data is too small
			if len(childrenURLs) < 2{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-large-img") + " ~ c 5.",nil)
				break
			}
			//create database id
			newID := collOperate.AutoCollParentFiles(parentTitle,parentURL,childrenURLs)
			if newID < 0 && newID != -1{
				continue
			}
			if newID == -1{
				continue
			}
			if newID > 0{
				errNum = 0
			}
			//too many error
			if errNum > 10{
				collOperate.NewLog(collOperate.lang.Get("coll-error-too-many") + " ~ c 6.",nil)
				return
			}
		}
	}
}
