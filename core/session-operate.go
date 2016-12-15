package core

import (
	"github.com/gorilla/sessions"
	"net/http"
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

//操作Session会话
//依赖内部库：
// core.LogOperate
//依赖外部库：
// github.com/gorilla/sessions

//Session操作类型
type SessionOperate struct {
	//应用名称
	appName []byte
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
	Value map[string]interface{}
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
	this.appName = []byte(appName)
	this.sessionIPBind = sessionIPBind
	this.sessionTimeout = sessionTimeout
	this.MatchString = MatchString
	this.dbCollStore = db.C("session-store")
}

//获取会话标记对象
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param name string 变量信息
//return map[interface{}]interface{}, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(w http.ResponseWriter,r *http.Request,name string) (map[interface{}]interface{}, bool) {
	//初始化变量
	var res map[interface{}]interface{}
	//获取标识信息
	var mark string
	mark = this.getCookieValue(w,r)
	if mark == ""{
		return res,false
	}
	//在数据库查找该值
	this.dbCollStore.Find(bson.M{"_id"})

	//获取session数据
	s, err := this.store.Get(r, mark)
	if err != nil {
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionGet","store-get",err.Error())
		return res, false
	}
	//返回
	return s.Values, true
}

//写入会话数据
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param appName string 应用标记
//param mark string 会话标记
//param data map[interface{}]interface{} 会话变量组合
//return bool 是否失败
func (this *SessionOperate) SessionSet(w http.ResponseWriter, r *http.Request,appName string,mark string, data map[interface{}]interface{}) bool {
	//确保session启动
	if this.status == false{
		this.store = sessions.NewCookieStore(this.appName)
		this.status = true
	}
	//获取session值
	s, err := this.store.Get(r, mark)
	if err != nil {
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","store-get",err.Error())
		return false
	}
	//保存到session
	s.Values = data
	err = s.Save(r, w)
	if err != nil {
		Log.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","store-save",err.Error())
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
	var res string
	//查找cookie
	cookieValue,err = r.Cookie(this.appName)
	if err != nil{
		//如果不存在，则创建
		var mark string
		mark = this.getCookieMark(r)
		if mark == ""{
			return res
		}
		//设定cookie
		cookieValue = *http.Cookie{
			Name:   this.appName,
			Value:    mark,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   this.sessionTimeout,
		}
		http.SetCookie(w,cookieValue)
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