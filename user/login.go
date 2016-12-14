package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"net/http"
	"strconv"
)

//该文件定义用户登录、退出、登录状态操作

//获取用户登录ID
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return string 登录的用户ID，未登录则返回空字符串
func GetLoginStatus(w http.ResponseWriter,r *http.Request) string{
	//获取值
	var res map[interface{}]interface{}
	var b bool
	res,b = SessionOperate.SessionGet(r,Mark,Mark)
	if b == false{
		return ""
	}
	//检查是否存在值
	if res["login-id"] == nil || res["login-time"] == nil{
		res["login-id"] = ""
		res["login-time"] = 0
		_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
		return ""
	}
	//更新登录时间值
	if res["login-id"].(string) != ""{
		var t time.Time
		t = time.Now()
		var unixTime int64
		unixTime = t.Unix()
		//超出时间，强行退出
		if UserLoginTimeoutMinute > unixTime - res["login-time"].(int64){
			var loginID string = ""
			res["login-id"] = loginID
			_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
			return ""
		}
		res["login-time"] = unixTime
		_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
	}
	//返回
	return res["login-id"].(string)
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
	var loginID string
	var loginErr int = 0
	//获取当前时间
	var t time.Time
	t = time.Now()
	var unixTime int64
	unixTime = t.Unix()
	//获取session数据
	res,b = SessionOperate.SessionGet(r,AppName,Mark)
	if b == false{
		return false
	}
	//退出该函数之前确保执行
	defer loginReturn(w,r,loginID,unixTime,loginErr)
	//检查是否超过错误次数5次
	if res["login-error"] != nil{
		loginErr = res["login-error"].(int)
		if loginErr > 5{
			return false
		}
	}
	//是否已经登录，是则返回成功
	if GetLoginStatus(w,r) != ""{
		return true
	}
	//检查用户名和密码是否合法
	if checkUsername(username,passwdSha1) == false{
		loginErr += 1
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = getPasswdSha1(passwdSha1)
	//获取IP地址
	var ipAddr string
	ipAddr = r.RemoteAddr
	//检查模式
	if oneUserStatus == true{
		//如果是单用户模式
		if oneUsername == username && passwdSha1Sha1 == oneUserpasswd{
			loginID = "true"
		}else{
			loginErr += 1
			return false
		}
	}else{
		//如果是多用户模式
		var result UserFields
		err = dbColl.Find(bson.M{"username":username,"password":passwdSha1Sha1}).One(&result)
		if err != nil{
			sendLog("user/login.go",ipAddr,"Login","many-user-find",err.Error())
			loginErr += 1
			return false
		}
		//用户存在，则修改登录IP和时间
		var userID string
		userID = result.ID.Hex()
		if userID != ""{
			err = dbColl.Update(bson.M{"_id":bson.ObjectIdHex(userID)},bson.M{"$set":bson.M{"lastip":ipAddr,"lasttime":unixTime}})
			if err != nil{
				sendLog("user/login.go",ipAddr,"Login","many-user-update",err.Error() + " , user id : "+userID+" , lastip : "+ipAddr+" , lasttime : "+strconv.FormatInt(unixTime,10))
				loginErr += 1
				return false
			}
			loginID = userID
		}
	}
	//检查是否验证通过
	if loginID == ""{
		loginErr += 1
		return false
	}
	//输出日志
	sendLog("user/login.go",ipAddr,"Login","login-success","ID为" + loginID + "的用户成功登录了平台。")
	//修改session并返回
	loginErr = 0
	return true
}

//登录完成后调用该模块
//用户记录Session部分、记录登录失败次数操作
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param loginID string 用户ID
//param unixTime int64 登录unix时间戳
//param loginErr int 失败次数
//return bool 记录是否成功
func loginReturn(w http.ResponseWriter,r *http.Request,loginID string,unixTime int64,loginErr int) bool{
	var res map[interface{}]interface{}
	var b bool
	//为了确保其他存储变量不会被清空，需要先读取再写入
	res,b = SessionOperate.SessionGet(r,AppName,Mark)
	if b == false{
		return false
	}
	res["login-id"] = loginID
	res["login-time"] = unixTime
	res["login-error"] = loginErr
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
	if res["login-id"].(string) == ""{
		return
	}
	var loginID string
	res["login-id"] = loginID
	_ = SessionOperate.SessionSet(w,r,AppName,Mark,res)
}
