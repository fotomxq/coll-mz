package router

import "net/http"

//管理用户界面、基本操作处理，以及其他一些模块。

//用户管理页面
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func PageUserSelf(w http.ResponseWriter, r *http.Request) {
	//检查是否已经登录
	var userID string
	userID = checkIPAndLogged(w, r, "user-self")
	if userID == "" {
		return
	}
	//获取当前用户信息
	userInfo,b := glob.UserOperate.GetIDFields(userID)
	if b == false{
		return
	}
	if checkPost(r) == false{
		return
	}
	//如果存在提交，则先处理后输出页面
	var postStatus string = "no"
	//获取数据
	var postNicename string
	postNicename = r.FormValue("nicename")
	var postPassword string
	postPassword = r.FormValue("password")
	//如果存在数据
	if postNicename != "" && postPassword != ""{
		postStatus = "has"
		postNicename = glob.MatchString.CheckFilterStr(postNicename,2,30)
		if postNicename == "" {
			sendLog("router/user-self.go",getIPAddr(r),"PageUserSelf","check-nicename","昵称存在错误。")
			postStatus = "error"
		}
		if glob.MatchString.CheckHexSha1(postPassword) == false {
			sendLog("router/user-self.go",getIPAddr(r),"PageUserSelf","check-password","密码存在错误。")
			postStatus = "error"
		}
		if postStatus != "error"{
			b = glob.UserOperate.Edit(userInfo.ID.Hex(),postNicename,userInfo.UserName,postPassword,userInfo.Permissions,userInfo.IsDisabled)
			if b == true{
				postStatus = "ok"
				//重新获取当前用户信息
				userInfo,b = glob.UserOperate.GetIDFields(userID)
				if b == false{
					return
				}
			}
		}
	}
	//初始化
	var data map[string]interface{} = map[string]interface{}{
		"refCSS": []string{
			"theme",
		},
		"refJS": []string{
			"user-self", "sha1","message",
		},
		"userNicename" : userInfo.NiceName,
		"postStatus" : postStatus,
	}
	//输出页面
	showTemplate(w, r, "user-self.html", data)
}

