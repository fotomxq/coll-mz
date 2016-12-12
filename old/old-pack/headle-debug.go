package controller

import "net/http"

//debug
func (this *Handle) actionDebug(w http.ResponseWriter, r *http.Request) {
	if configData["debug"] != "true"{
		this.page404(w,r)
		return
	}
	//If not, jump
	if this.CheckLogin(w, r) == false {
		return
	}
	//Make sure that post / get is fine
	b := this.CheckURLPost(r)
	if b == false {
		return
	}
	//Gets the submit action type
	_,_ = coll.GetStatus()
	postAction := r.FormValue("action")
	postType := r.FormValue("type")
	var resStr string
	switch postAction {
	case "coll":
		go coll.CollDebug()
		break
	default:
		break
	}
	//show data
	switch postType {
	case "json":
		this.postJSONData(w, r, resStr, b)
		return
		break
	case "html":
		this.PostText(w,r,resStr)
		return
		break
	default:
		data := map[string]interface{}{
			"debug" : configData["debug"].(string),
			"html" : resStr,
		}
		this.ShowTemplate(w,r,"debug.html",data)
		return
		break
	}
}
