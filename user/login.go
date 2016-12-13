package user

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"net/http"
)

//该文件定义用户登录、退出、登录状态操作

//获取用户登录ID
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return int 登录的用户ID
func GetLoginStatus(w http.ResponseWriter,r *http.Request) int{
	//获取值
	var res map[interface{}]interface{}
	var b bool
	res,b = SessionOperate.SessionGet(r,Mark,Mark)
	if b == false{
		return 0
	}
	//检查是否存在值
	if res["login-id"] == nil || res["login-time"] == nil{
		res["login-id"] = 0
		res["login-time"] = 0
		_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
		return 0
	}
	//更新登录时间值
	if res["login-id"].(int) > 0{
		var t time.Time
		t = time.Now()
		var unixTime int64
		unixTime = t.Unix()
		//超出时间，强行退出
		if UserLoginTimeoutMinute > unixTime - res["login-time"].(int64){
			var loginID int = 0
			res["login-id"] = loginID
			_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
			return 0
		}
		res["login-time"] = unixTime
		_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
	}
	//返回
	return res["login-id"].(int)
}

//用户登录
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//param r *http.Request HTTP读取句柄
//return bool 是否登录成功
func Login(w http.ResponseWriter,r *http.Request,username string,passwdSha1 string) bool{
	//初始化变量
	var res map[interface{}]interface{}
	var b bool
	var err error
	res,b = SessionOperate.SessionGet(r,AppName,Mark)
	if b == false{
		return false
	}
	var loginID int = 0
	//是否已经登录，是则返回成功
	if GetLoginStatus(w,r) > 0{
		return true
	}
	//检查用户名和密码是否合法
	if checkUsername(username,passwdSha1) == false{
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = getPasswdSha1(passwdSha1)
	//获取IP地址
	var ipAddr string
	ipAddr = r.RemoteAddr
	//获取当前时间
	var t time.Time
	t = time.Now()
	var unixTime int64
	unixTime = t.Unix()
	//检查模式
	if oneUserStatus == true{
		//如果是单用户模式
		if oneUsername == username && passwdSha1Sha1 == oneUserpasswd{
			loginID = 1
		}else{
			return false
		}
	}else{
		//如果是多用户模式
		var result UserFields
		err = dbColl.Find(bson.M{"username":username,"password":passwdSha1Sha1}).One(&result)
		if err != nil{
			sendLog(err.Error())
			return false
		}
		//用户存在，则修改登录IP和时间
		var userID int
		userID,err = strconv.Atoi(result.id.String())
		if err != nil{
			sendLog(err.Error())
			return false
		}
		if userID > 0{
			err = dbColl.UpdateId(userID,bson.M{"last_ip":ipAddr,"last_time":unixTime})
			if err != nil{
				sendLog(err.Error())
				return false
			}
			loginID = userID
		}
	}
	//检查是否验证通过
	if loginID < 1{
		return false
	}
	//输出日志
	sendLog("用户" + strconv.Itoa(loginID) + "通过IP地址" + ipAddr + "登录了系统。")
	//修改session
	res["login-id"] = loginID
	res["login-time"] = unixTime
	return SessionOperate.SessionSet(w,r,AppName,Mark,res)
}

//用户退出
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
func Logout(w http.ResponseWriter,r *http.Request){
	var res map[interface{}]interface{}
	var b bool
	res,b = SessionOperate.SessionGet(r,AppName,Mark)
	if b == false{
		return
	}
	if res["login-id"].(int) < 1{
		return
	}
	var loginID int = 0
	res["login-id"] = loginID
	_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
}
