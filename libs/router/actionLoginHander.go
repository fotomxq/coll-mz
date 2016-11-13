package router

import (
	"net/http"
	"github.com/fotomxq/coll-mz/libs/core"
)

func actionLoginHander(w http.ResponseWriter, r *http.Request){
	err = r.PostForm()
	if err != nil{
		log.AddErrorLog(err)
	}
	postUser := r.FormValue("email")
	postPasswd := r.FormValue("password")
	//验证数据
	ms := new(core.MatchString)
	if ms.CheckUsername(postUser) == false || ms.CheckPassword(postPasswd) == false{
		http.Redirect(w, r, "/login/", http.StatusFound)
	}else{
		http.Redirect(w, r, "/index/", http.StatusFound)
	}
}
