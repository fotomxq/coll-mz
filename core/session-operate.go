package core

import (
	"github.com/gorilla/sessions"
	"net/http"
	"gopkg.in/mgo.v2"
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
	//采集器存储句柄
	store *sessions.CookieStore
	//会话启动状态
	status bool
	//存储session数据的数据库集合
	dbCollStore *mgo.Collection
	//session是否和IP绑定
	//如果绑定，则cookie必须和IP地址一致
	//潜在可能造成用户特殊环境下，变化IP地址后自动退出的情况
	sessionIPBind bool
}

//创建会话
//必须执行该函数，才能使用其他内部函数
//param name string 标记
//param db *mgo.Database 数据库句柄
//param sessionIPBind bool 是否和IP绑定
func (this *SessionOperate) Create(name string,db *mgo.Database,sessionIPBind bool) {
	this.appName = []byte(name)
	this.store = sessions.NewCookieStore(this.appName)
	this.status = true
	this.dbCollStore = db.C("session-store")
	this.sessionIPBind = sessionIPBind
}

//获取会话标记对象
//param r *http.Request Http读取对象
//param appName string 应用标记
//param mark string 标记
//return map[interface{}]interface{}, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(r *http.Request,appName string,mark string) (map[interface{}]interface{}, bool) {
	//初始化变量
	var res map[interface{}]interface{}
	//确保session已经启动
	if this.status == false{
		this.store = sessions.NewCookieStore(this.appName)
		this.status = true
	}
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
