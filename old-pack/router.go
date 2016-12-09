package controller

import (
	"net/http"
)

//router
type Router struct {
	matchString MatchString
	handle      Handle
}

//Start the server
func (this *Router) RunServer(db *Database) {
	//Initialize the handle
	this.handle.Init(db)
	//Set Static
	sep := GetPathSep()
	templateDir := "." + sep + "template"
	http.Handle("/assets/", http.FileServer(http.Dir(templateDir)))
	//Set dynamic binding
	http.HandleFunc("/", this.handle.page404)
	http.HandleFunc("/favicon.ico", this.handle.pageFavicon)
	http.HandleFunc("/login", this.handle.pageLogin)
	http.HandleFunc("/action-login", this.handle.actionLogin)
	http.HandleFunc("/action-logout", this.handle.actionLogout)
	http.HandleFunc("/set", this.handle.pageSet)
	http.HandleFunc("/action-set", this.handle.actionSet)
	http.HandleFunc("/center", this.handle.pageCenter)
	http.HandleFunc("/action-center", this.handle.actionCenter)
	http.HandleFunc("/action-list", this.handle.actionViewList)
	http.HandleFunc("/action-view", this.handle.actionView)
	http.HandleFunc("/debug", this.handle.actionDebug)
	//Start the server listening
	log.NewLog("Server run : "+configData["server-local"].(string), nil)
	err = http.ListenAndServe(configData["server-local"].(string), nil)
	if err != nil {
		log.NewLog("Unable to connect to the database.", err)
	}
}
