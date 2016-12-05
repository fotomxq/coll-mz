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
	//start
	page := 1
	errNum := 0
	//Traverse all page numbers
	for {
		if thisChildren.status == false{
			return
		}
		thisPageURL := thisChildren.url + strconv.Itoa(page)
		doc,err := goquery.NewDocument(thisPageURL)
		if err != nil{
			collOperate.NewLog(collOperate.lang.Get("coll-error-next-open"),err)
			break
		}
		//collOperate.NewLog("Get URL : " + thisPageURL,nil)
		//Gets and traverses the nodes of all page pages
		nodes := doc.Find("#pins").Children()
		for nodeKey := range nodes.Nodes{
			//collOperate.NewLog("Get the page sub-node , key : " + strconv.Itoa(nodeKey),nil)
			if thisChildren.status == false{
				return
			}
			node := nodes.Eq(nodeKey).Children()
			//Gets the page URL and title
			thisCURL,b := node.Eq(0).Attr("href")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			//collOperate.NewLog("Get the page sub-node a href , href url : " + thisCURL,nil)
			thisCTitle,err := node.Eq(1).Children().Html()
			if err != nil{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),err)
				errNum += 1
				continue
			}
			//Thumbnail analysis of all internal image URL address
			//eg : http://i.meizitu.net/thumbs/2016/11/80538_29a51_236.jpg
			thisCThumbsURL,b := node.Eq(0).Children().Attr("data-original")
			if b == false{
				collOperate.NewLog(collOperate.lang.Get("coll-error-get-children"),nil)
				errNum += 1
				continue
			}
			//collOperate.NewLog("Get thisCThumbsURL : " + thisCThumbsURL,nil)
			//eg : http://i.meizitu.net/2016/11/80538_29a51_236.jpg
			thisCThumbsURL = strings.Replace(thisCThumbsURL,"thumbs/","",-1)
			//eg : {"http:","","i.meizitu.net","2016","11","80538_29a51_236.jpg"}
			thisCThumbsURLs := strings.Split(thisCThumbsURL,"/")
			//If urls is not equal to 7, the file structure changes.
			if len(thisCThumbsURLs) != 6{
				collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
				errNum = 11
				break
			}
			//Disassemble the end name part
			//eg : {"80538_29a51_236","jpg"}
			thisCThumbsURLNames := strings.Split(thisCThumbsURLs[len(thisCThumbsURLs) - 1],".")
			if len(thisCThumbsURLNames) != 2{
				collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
				errNum = 11
				break
			}
			//eg : {"80538","29a51","236"}
			thisCThumbsURLNames2 := strings.Split(thisCThumbsURLNames[0],"_")
			if len(thisCThumbsURLNames2) != 3{
				collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
				errNum = 11
				break
			}
			//Analyze the name section
			thisCThumbsURLNamesSeps := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u"}
			sep := ""
			for _,v := range thisCThumbsURLNamesSeps{
				if strings.Index(thisCThumbsURLNames2[1],v) > 0{
					sep = v
					break
				}
			}
			if sep == ""{
				collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
				errNum = 11
				break
			}
			//eg : {"29","51"}
			thisCThumbsURLNames3 := strings.Split(thisCThumbsURLNames2[1],sep)
			if len(thisCThumbsURLNames3) != 2{
				collOperate.NewLog(collOperate.lang.Get("coll-error-doc"),nil)
				errNum = 11
				break
			}
			//Create an upper level ID
			parentID,parentSha1 := collOperate.AutoCollParent(thisCTitle,thisCURL)
			if parentID < 1 && parentID != -1{
				errNum += 1
				continue
			}
			//Synthesize the final URL address
			imgURL := thisCThumbsURLs[0] +"/"+ thisCThumbsURLs[1] +"/"+ thisCThumbsURLs[2] +"/"+ thisCThumbsURLs[3] +"/"+ thisCThumbsURLs[4] + "/" + thisCThumbsURLNames3[0] + sep
			nextNum := 1
			for {
				nextNumStr := strconv.Itoa(nextNum)
				if nextNum < 10 {
					nextNumStr = "0" + nextNumStr
				}
				nextImgURL := imgURL + nextNumStr + "." + thisCThumbsURLNames[1]
				newID := collOperate.AutoCollFile(nextImgURL,thisCTitle,parentSha1,parentID)
				if newID < 1{
					if nextNum < 2 {
						_ = collOperate.DeleteData(parentID)
						errNum = 11
					}
					break
				}
				errNum = 0
				nextNum += 1
			}
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