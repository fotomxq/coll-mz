package core

import (
	"net/http"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

//操作Session会话
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
	//会造成用户特殊环境下，变化IP地址后自动退出的情况，常见于移动端访问。
	sessionIPBind bool
	//session超时时间
	sessionTimeout int
}

//Session数据字段
type SessionFields struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	IP string
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
	this.MatchString = MatchString
	this.dbCollStore = db.C("session-store")
}

//获取会话标记对象
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param name string 标记
//return map[string]interface{}, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(w http.ResponseWriter,r *http.Request,name string) (map[string]interface{}, bool) {
	//获取数据
	var res map[string]map[string]interface{}
	var result *SessionFields
	var b bool
	result,b = this.getData(w,r)
	if b == false{
		return map[string]interface{}{},false
	}
	//如果不存在数据，则返回
	if result.Value == ""{
		return map[string]interface{}{},true
	}
	//解析数据
	err = json.Unmarshal([]byte(result.Value),&res)
	if err != nil{
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionGet","un-json",err.Error())
	}
	//如果不存在数据，则返回
	if res[name] == nil{
		return map[string]interface{}{},true
	}
	//返回数据集合
	return res[name],true
}

//写入会话数据
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param name string 标记
//param data map[string]interface{} 会话变量组合
//return bool 是否失败
func (this *SessionOperate) SessionSet(w http.ResponseWriter, r *http.Request,name string, data map[string]interface{}) bool {
	//获取数据
	var result *SessionFields
	var b bool
	result,b = this.getData(w,r)
	if b == false{
		return false
	}
	//解析数据，并覆盖数据
	var res map[string]map[string]interface{}
	if result.Value != ""{
		err = json.Unmarshal([]byte(result.Value),&res)
		if err != nil{
			Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","un-json",err.Error())
			return false
		}
		res[name] = data
	}else{
		res = map[string]map[string]interface{}{
			name : data,
		}
	}
	//JSON
	var resJsonByte []byte
	resJsonByte,err = json.Marshal(res)
	if err != nil{
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","en-json",err.Error())
		return false
	}
	var resJson string
	resJson = string(resJsonByte)
	//写入数据库
	err = this.dbCollStore.UpdateId(bson.ObjectIdHex(result.ID.Hex()),bson.M{"$set":bson.M{"value":resJson}})
	if err != nil{
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","update-db",err.Error())
		return false
	}
	//返回
	return true
}

//确保cookie值存在，并获取值
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return string mark标记值
func (this *SessionOperate) getCookieValue(w http.ResponseWriter, r *http.Request) string{
	//初始化cookie值
	var cookieValue *http.Cookie
	//查找cookie
	cookieValue,err = r.Cookie(this.appName)
	if err != nil{
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
		//将数据保存到数据库中
		err = this.dbCollStore.Insert(&SessionFields{bson.NewObjectId(),mark,r.RemoteAddr,""})
		if err != nil{
			Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.getCookieValue","insert-db",err.Error())
			return ""
		}
		//解析返回数据，直接将mark作为返回值
		return mark
	}
	//存在，则获取值
	return cookieValue.Value
}

//获取cookie mark
//param r *http.Request HTTP读句柄
//return string mark值，失败则返回空字符串
func (this *SessionOperate) getCookieMark(r *http.Request) string{
	//mark通过IP地址+应用标记+当前时间+随机数
	var mark string
	var t time.Time
	t = time.Now()
	mark = r.RemoteAddr + this.appName + "cookie" + t.String() + this.MatchString.GetRandStr(99999)
	mark = this.MatchString.GetSha1(mark)
	if mark == ""{
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.getCookieMark","get-sha1","无法获取cookie的sha1值。")
		return ""
	}
	return mark
}

//获取数据集合数据
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return *SessionFields,bool 数据集合，是否成功
func (this *SessionOperate) getData(w http.ResponseWriter, r *http.Request) (*SessionFields,bool){
	//初始化变量
	var res SessionFields
	//获取标识信息
	var mark string
	mark = this.getCookieValue(w,r)
	if mark == ""{
		return &res,false
	}
	//在数据库查找该值
	err = this.dbCollStore.Find(bson.M{"name":mark}).One(&res)
	if err != nil{
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.getData","get-database",err.Error())
		return &res,false
	}
	//如果当前IP地址和结果IP不一致，则返回
	if this.sessionIPBind == true{
		if res.IP != r.RemoteAddr{
			Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.getData","ip-no-bind","客户端IP地址和Cookie记录IP地址不符。")
			return &res,false
		}
	}
	return &res,true
}