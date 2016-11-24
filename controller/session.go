package controller

import (
	"github.com/gorilla/sessions"
	"net/http"
)

//session store
var store = sessions.NewCookieStore([]byte("coll-mz-session"))

//Get the session data
func SessionGet(w http.ResponseWriter, r *http.Request,sessionMark string) (map[interface{}]interface{},bool){
	s,err := store.Get(r, sessionMark)
	var res map[interface{}]interface{}
	if err != nil{
		log.NewLog("",err)
		return res,false
	}
	return s.Values,true
}

//Write session data
func SessionSet(w http.ResponseWriter, r *http.Request,sessionMark string,data map[interface{}]interface{}) (bool){
	s,err := store.Get(r, sessionMark)
	if err != nil{
		log.NewLog("",err)
		return false
	}
	s.Values = data
	err = s.Save(r,w)
	if err != nil{
		log.NewLog("",err)
		return false
	}
	return true
}