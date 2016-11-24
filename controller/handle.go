package controller

import(
	"net/http"
	"html/template"
)

//The page handle
type Handle struct {
	//User processor
	user User
	//language configuration processor
	lang Language
	//Database processor
	db *Database
}

func (this *Handle) Init(db *Database){
	//Initialize the language configuration processor
	this.lang.Init(configData["language"].(string))
	//Save the database processor
	this.db = db
	//Initialize the user processor
	this.user.Init(db,3600)
}

/////////////////////////////////////
//This part is a generic module
/////////////////////////////////////

//Get the template file path
func (this *Handle) GetTempSrc(name string) string{
	return "template" + GetPathSep() + name
}

//Output text directly to the browser
func (this *Handle) PostText(w http.ResponseWriter, r *http.Request, content string){
	var contentByte []byte = []byte(content)
	_,err := w.Write(contentByte)
	if err != nil{
		log.NewLog("You can not directly output string data.",err)
		return
	}
}

//Jump to URL
func (this *Handle) ToURL(w http.ResponseWriter, r *http.Request,urlName string) {
	http.Redirect(w, r, urlName, http.StatusFound)
}

//Output template
func (this *Handle) ShowTemplate(w http.ResponseWriter, r *http.Request,templateFileName string,data interface{}){
	t, err := template.ParseFiles(this.GetTempSrc(templateFileName))
	if err != nil {
		log.NewLog("The template does not output properly,template file name : " + templateFileName,err)
		return
	}
	t.Execute(w, data)
}

//Output the prompt page
func (this *Handle) showTip(w http.ResponseWriter, r *http.Request,title string,contentTitle string,content string,gotoURL string){
	data := map[string]string{
		"title" : title,
		"contentTitle" : contentTitle,
		"content" : content,
		"gotoURL" : gotoURL,
	}
	this.ShowTemplate(w,r,"tip.html",data)
}

//Check that you are logged in
func (this *Handle) CheckLogin(w http.ResponseWriter, r *http.Request) bool{
	if this.user.CheckLogin(w, r) == false {
		log.NewLog("User has not logged in, but visited the home page.",nil)
		this.ToURL(w,r,"/login")
		return false
	}
	return true
}

//Check the post data
func (this *Handle) CheckURLPost(r *http.Request) bool{
	err = r.ParseForm()
	if err != nil {
		log.NewLog("Failed to get get / post data.",err)
		return false
	}
	return true
}

/////////////////////////////////////
//This section is the page
/////////////////////////////////////

//404 error handling
func (this *Handle) page404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if this.CheckLogin(w,r) == false{
			return
		}else{
			this.ToURL(w,r,"/center")
		}
	}else{
		log.NewLog("The page can not be found,url path : " + r.URL.Path,nil)
		this.ShowTemplate(w,r,"404.html",nil)
	}
}

//Resolve the login page
func (this *Handle) pageLogin(w http.ResponseWriter, r *http.Request){
	if this.user.CheckLogin(w, r) == true {
		this.ToURL(w,r,"/center")
		return
	} else {
		this.ShowTemplate(w,r,"login.html",nil)
		return
	}
}

//Get the site icon file
func (this *Handle) pageFavicon(w http.ResponseWriter, r *http.Request) {
	this.ToURL(w,r,"/assets/favicon.ico")
}

//Output the set page
func (this *Handle) pageSet(w http.ResponseWriter, r *http.Request) {
	if this.CheckLogin(w,r) == false{
		return
	}
	this.ShowTemplate(w,r,"set.html",nil)
}

//Output the center page
func (this *Handle) pageCenter(w http.ResponseWriter, r *http.Request) {
	if this.CheckLogin(w,r) == false{
		return
	}
	this.ShowTemplate(w,r,"center.html",nil)
}

/////////////////////////////////////
//This section is the feedback page
/////////////////////////////////////

//Submit data Try to log in
func (this *Handle) actionLogin(w http.ResponseWriter, r *http.Request) {
	postUser := r.FormValue("email")
	postPasswd := r.FormValue("password")
	b := this.user.LoginIn(w,r,postUser,postPasswd)
	if b == false{
		this.ToURL(w,r,"/login")
		return
	}else{
		this.ToURL(w,r,"/center")
	}
}

//sign out
func (this *Handle) actionLogout(w http.ResponseWriter, r *http.Request) {
	if this.user.CheckLogin(w,r) == false{
		this.ToURL(w,r,"/login")
		return
	}
	this.showTip(w,r,this.lang.Get("handle-logout-title"),this.lang.Get("handle-logout-contentTitle"),this.lang.Get("handle-logout-content"),"/login")
}

//Resolution settings page
func (this *Handle) actionSet(w http.ResponseWriter, r *http.Request) {
	//If not, jump
	if this.CheckLogin(w,r) == false{
		return
	}
	//Make sure that post / get is fine
	b := this.CheckURLPost(r)
	if b == false{
		return
	}
	//Gets the submit action type
	postAction := r.FormValue("action")
	switch postAction {
	case "coll-all":
		this.PostText(w, r, "coll-run-ok")
		break
	case "get-log":
		break
	default:
		this.page404(w,r)
		return
		break
	}
}

//Feedback center action
func (this *Handle) actionCenter(w http.ResponseWriter, r *http.Request){
	if this.CheckLogin(w,r) == false{
		return
	}
}

//Feedback center view content action
func (this *Handle) actionView(w http.ResponseWriter, r *http.Request) {
	if this.CheckLogin(w,r) == false{
		return
	}
}