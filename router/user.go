package router

import (
	"net/http"
	"strings"
)

//管理用户界面、基本操作处理，以及其他一些模块。

//用户管理页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageUser(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w, r, "user")
	if userID == "" {
		userID = checkIPAndLogged(w, r, "user-self")
		if userID != ""{
			goURL(w,r,"/user-self")
		}
		return
	}
	//初始化
	var data map[string]interface{} = map[string]interface{}{
		"refCSS": []string{
			"theme", "user",
		},
		"refJS": []string{
			"user", "sha1", "message",
		},
	}
	//输出页面
	showTemplate(w, r, "user.html", data)
}

//用户管理界面动作处理
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func ActionUser(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w, r, "user")
	if userID == "" {
		return
	}
	//检查post
	if checkPost(r) == false {
		return
	}
	//初始化
	var data map[string]interface{} = map[string]interface{}{}
	var b bool
	//post action
	var postAction string
	postAction = r.FormValue("action")
	switch postAction {
	case "permissions":
		data["permissions"] = glob.UserOperate.PermissionsData
		b = true
	case "list":
		//搜索提交
		var search string
		search = glob.MatchString.CheckFilterStr(r.FormValue("search"),1,50)
		var pages PageMaxSortDesc
		pages = getPageMaxSortDesc(r)
		//获取用户列表
		var userData *[]map[string]interface{}
		userData, b = glob.UserOperate.List(search, pages.page, pages.max, pages.sort, pages.desc)
		data["list"] = userData
		data["page"] = pages.page
		data["max"] = pages.max
		data["sort"] = pages.sort
		data["desc"] = pages.desc
		data["search"] = search
	case "create":
		var niceName string
		niceName = glob.MatchString.CheckFilterStr(r.FormValue("nicename"),2,30)
		if niceName == "" {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","add-check-nicename","昵称存在错误。")
			break
		}
		var userName string
		userName = r.FormValue("username")
		if glob.MatchString.CheckUsername(userName) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","add-check-username","用户名存在错误。")
			break
		}
		var password string
		password = r.FormValue("password")
		if glob.MatchString.CheckHexSha1(password) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","add-check-password","密码存在错误。")
			break
		}
		var permissionsStr string
		permissionsStr = r.FormValue("permissions")
		if permissionsStr == ""{
			sendLog("router/user.go",getIPAddr(r),"ActionUser","add-check-permissions","权限存在错误。")
			break
		}
		var permissions []string
		permissions = strings.Split(permissionsStr, "|")
		var newUserID string
		newUserID = glob.UserOperate.Create(niceName, userName, password, permissions)
		b = newUserID != ""
		if b == false{
			sendLog("router/user.go",getIPAddr(r),"ActionUser","add-result","添加失败。")
			break
		}
	case "edit":
		var postUserID string
		postUserID = r.FormValue("id")
		if glob.MatchString.CheckHexSha1(postUserID) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-id","用户ID存在错误。")
			break
		}
		var niceName string
		niceName = glob.MatchString.CheckFilterStr(r.FormValue("nicename"),2,30)
		if niceName == "" {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-nicename","昵称存在错误。")
			break
		}
		var userName string
		userName = r.FormValue("username")
		if glob.MatchString.CheckUsername(userName) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-username","用户名存在错误。")
			break
		}
		var password string
		password = r.FormValue("password")
		if password != "" && glob.MatchString.CheckHexSha1(password) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-password","密码存在错误。")
			break
		}
		var permissionsStr string
		permissionsStr = r.FormValue("permissions")
		if permissionsStr == ""{
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-permissions","权限存在错误。")
			break
		}
		var permissions []string
		permissions = strings.Split(permissionsStr, "|")
		var isDisabledStr string
		var isDisabled bool = false
		isDisabledStr = r.FormValue("isdisabled")
		if isDisabledStr == "true"{
			isDisabled = true
		}
		if glob.MatchString.CheckUsername(userName) == false {
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-check-username","用户名存在错误。")
			break
		}
		b = glob.UserOperate.Edit(postUserID, niceName, userName, password, permissions,isDisabled)
		if b == false{
			sendLog("router/user.go",getIPAddr(r),"ActionUser","edit-result","修改失败。")
			break
		}
	case "delete":
		var postUserID string
		postUserID = r.FormValue("id")
		if glob.MatchString.CheckHexSha1(postUserID) == false {
			break
		}
		if postUserID == userID {
			break
		}
		b = glob.UserOperate.Delete(postUserID)
	}
	postJSONData(w, r, data, b, userID != "")
}

//////////////////////////////////////////////////////////////////////////////////////
//用户相关的通用模块
//////////////////////////////////////////////////////////////////////////////////////

//检查是否已经登录
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
//param page string 当前页面名称
//return string 用户ID
func checkIPAndLogged(w http.ResponseWriter, r *http.Request, page string) string {
	//检查IP是否可访问
	if checkIP(r) == false {
		return ""
	}
	//检查是否已经登录了
	var userID string
	userID = userCheckLogged(w, r)
	if userID == "" {
		goURL(w, r, "/login")
		return ""
	}
	//检查用户权限，是否足够访问该页面？
	if glob.UserOperate.CheckUserVisitPage(userID, page) == false {
		return ""
	}
	//返回
	return userID
}
