package core

import (
	"net/http"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

//操作Session会话
// 内部数据必须统一为string
//警告：
// 使用该模块操作cookie，会降低系统整体运行效率，但会提升整体安全性。
// 可有效避免跨站请求、伪造cookie漏洞
//依赖内部库：
// core.LogOperate
//依赖外部库：
// github.com/gorilla/sessions

//Session操作类型
type SessionOperate struct {
	//应用名称
	appName string
	//cookie句柄
	store http.Cookie
	//验证处理句柄
	MatchString *MatchString
	//存储session数据的数据库集合
	dbCollStore *mgo.Collection
	//session是否和IP绑定
	//如果绑定，则cookie值必须和IP地址一致
	sessionIPBind bool
	//session超时时间
	sessionTimeout int
	sessionTimeout64 int64
}

//Session数据字段
type SessionFields struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//标记名称
	Name string
	//创建Unix时间戳
	CreateTime int64
	//IP地址
	IP string
	//数据值
	Value string
}

//创建会话
//必须执行该函数，才能使用其他内部函数
//param appName string 应用标记
//param db *mgo.Database 数据库句柄
//param sessionIPBind bool 是否和IP绑定
//param sessionTimeout int cookie过期时间
//param MatchString *MatchString 验证模块
func (this *SessionOperate) Create(appName string,db *mgo.Database,sessionIPBind bool,sessionTimeout int,MatchString *MatchString) {
	//保存变量
	this.appName = appName
	this.sessionIPBind = sessionIPBind
	this.sessionTimeout = sessionTimeout
	this.sessionTimeout64 = int64(sessionTimeout)
	this.MatchString = MatchString
	this.dbCollStore = db.C("session-store")
}

//获取会话标记对象
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param name string 标记
//return map[string]string, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(w http.ResponseWriter,r *http.Request,name string) (map[string]string, bool) {
	//检查所有超过时间期限的数据，自动删除cookie记录
	var timeOut int64
	timeOut = time.Now().Unix() - this.sessionTimeout64
	err = this.dbCollStore.Remove(bson.M{"createtime":bson.M{"$lt":timeOut}})
	if err == nil{
		//Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionGet","clear-timeout-cookie","清理了过期的Cookie记录。")
	}
	//初始化变量
	var result map[string]map[string]string
	var res SessionFields
	//获取标识信息
	var mark string
	mark = this.getCookieValue(w,r)
	if mark == ""{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionGet","get-mark","无法获取cookie标识码。")
		return map[string]string{},false
	}
	//在数据库查找该值，不存在则返回
	err = this.dbCollStore.Find(bson.M{"name":mark}).One(&res)
	if err != nil{
		//Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionGet","get-database",err.Error())
		return map[string]string{},false
	}
	//如果当前IP地址和结果IP不一致，则重新构建cookie
	if this.sessionIPBind == true{
		if res.IP != IPAddrsGetRequest(r){
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionGet","ip-no-bind","客户端IP地址和Cookie记录IP地址不符，数据集合内的IP地址是：" + res.IP)
			return map[string]string{},false
		}
	}
	//如果不存在数据，则返回
	if res.Value == ""{
		return map[string]string{},true
	}
	//解析数据
	err = json.Unmarshal([]byte(res.Value),&result)
	if err != nil{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionGet","un-json",err.Error())
	}
	//如果不存在数据，则返回
	if result[name] == nil{
		return map[string]string{},true
	}
	//返回数据集合
	return result[name],true
}

//写入会话数据
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param name string 标记
//param data map[string]string 会话变量组合
//return bool 是否失败
func (this *SessionOperate) SessionSet(w http.ResponseWriter, r *http.Request,name string, data map[string]string) bool {
	//初始化变量
	var res SessionFields
	var result map[string]map[string]string
	//获取标识信息
	var mark string
	mark = this.getCookieValue(w,r)
	//如果标识不存在，则创建
	if mark == ""{
		result = map[string]map[string]string{
			name : data,
		}
		var resJsonByte []byte
		resJsonByte,err = json.Marshal(result)
		if err != nil{
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","no-mark-en-json",err.Error())
			return false
		}
		var resJson string
		resJson = string(resJsonByte)
		mark = this.createCookie(w,r,resJson)
		if mark == ""{
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","get-mark",err.Error())
			return false
		}
		//直接构建了新的数据，返回成功
		return true
	}
	//在数据库查找该值，不存在则重建
	err = this.dbCollStore.Find(bson.M{"name":mark}).One(&res)
	if err != nil{
		result = map[string]map[string]string{
			name : data,
		}
		var resJsonByte []byte
		resJsonByte,err = json.Marshal(result)
		if err != nil{
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","no-find-en-json",err.Error())
			return false
		}
		var resJson string
		resJson = string(resJsonByte)
		mark = this.createCookie(w,r,resJson)
		if mark == ""{
			return false
		}
		err = this.dbCollStore.Find(bson.M{"name":mark}).One(&res)
		if err != nil{
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","get-database",err.Error())
			return false
		}
		//直接构建了新的数据，返回成功
		return true
	}
	//如果当前IP地址和结果IP不一致，则重新构建cookie
	if this.sessionIPBind == true{
		if res.IP != IPAddrsGetRequest(r){
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","ip-no-bind","客户端IP地址和Cookie记录IP地址不符，数据集合内的IP地址是：" + res.IP)
			return false
		}
	}
	//解析数据，并覆盖数据
	if res.Value != ""{
		err = json.Unmarshal([]byte(res.Value),&result)
		if err != nil{
			Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","un-json",err.Error())
			return false
		}
		result[name] = data
	}else{
		result = map[string]map[string]string{
			name : data,
		}
	}
	//JSON
	var resJsonByte []byte
	resJsonByte,err = json.Marshal(result)
	if err != nil{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","en-json",err.Error())
		return false
	}
	var resJson string
	resJson = string(resJsonByte)
	//写入数据库
	err = this.dbCollStore.UpdateId(bson.ObjectIdHex(res.ID.Hex()),bson.M{"$set":bson.M{"value":resJson}})
	if err != nil{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.SessionSet","update-db",err.Error())
		return false
	}
	//返回
	return true
}

//删除该用户客户端cookie
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return bool 是否成功
func (this *SessionOperate) RemoveCookie(w http.ResponseWriter, r *http.Request) bool {
	var cookieValue *http.Cookie
	//查找cookie
	cookieValue,err = r.Cookie(this.appName)
	//如果获取失败、空值、非40位的SHA1，则重新创建cookie
	if err != nil || cookieValue.Value == "" || len(cookieValue.Value) != 40{
		return true
	}
	var mark string = cookieValue.Value
	err = this.dbCollStore.Remove(bson.M{"name":mark})
	if err != nil {
		Log.SendLog("core/session-operate.go", IPAddrsGetRequest(r), "SessionOperate.removeCookie", "remove-db", "无法删除cookie的值。")
		return false
	}
	cookieValue = &http.Cookie{
		Name:   this.appName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge: 1,
	}
	http.SetCookie(w,cookieValue)
	return true
}

//获取cookie标识码
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return string mark标记值
func (this *SessionOperate) getCookieValue(w http.ResponseWriter, r *http.Request) string{
	//初始化cookie值
	var cookieValue *http.Cookie
	//查找cookie
	cookieValue,err = r.Cookie(this.appName)
	//如果获取失败、空值、非40位的SHA1，则重新创建cookie
	if err != nil || cookieValue.Value == "" || len(cookieValue.Value) != 40{
		return ""
	}
	//存在，则获取值
	return cookieValue.Value
}

//创建cookie
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param data string 要写入的数据
//return string cookie标识码
func (this *SessionOperate) createCookie(w http.ResponseWriter, r *http.Request,data string) string{
	//初始化cookie值
	var cookieValue *http.Cookie
	//如果不存在，则创建
	var mark string
	mark = this.getCookieMark(r)
	if mark == ""{
		return ""
	}
	//设定cookie
	cookieValue = &http.Cookie{
		Name:   this.appName,
		Value:    mark,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   this.sessionTimeout,
	}
	http.SetCookie(w,cookieValue)
	//获取当前时间戳
	var unixTime int64
	unixTime = time.Now().Unix()
	//将数据保存到数据库中
	err = this.dbCollStore.Insert(&SessionFields{bson.NewObjectId(),mark,unixTime,IPAddrsGetRequest(r),data})
	if err != nil{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.createCookie","insert-db",err.Error())
		return ""
	}
	//解析返回数据，直接将mark作为返回值
	return mark
}

//获取cookie mark
//param r *http.Request HTTP读句柄
//return string mark值，失败则返回空字符串
func (this *SessionOperate) getCookieMark(r *http.Request) string{
	//mark通过IP地址+应用标记+当前时间+随机数
	var mark string
	var t time.Time
	t = time.Now()
	mark = IPAddrsGetRequest(r) + this.appName + "cookie" + t.String() + this.MatchString.GetRandStr(99999)
	mark = this.MatchString.GetSha1(mark)
	if mark == ""{
		Log.SendLog("core/session-operate.go",IPAddrsGetRequest(r),"SessionOperate.getCookieMark","get-sha1","无法获取cookie的sha1值。")
		return ""
	}
	return mark
}