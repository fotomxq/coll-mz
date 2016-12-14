package core

import (
	"github.com/gorilla/sessions"
	"net/http"
)

//操作Session会话
//依赖内部库：
// core.LogOperate
//依赖外部库：
// github.com/gorilla/sessions

//Session操作类型
type SessionOperate struct {
	//采集器存储句柄
	store *sessions.CookieStore
	//会话启动状态
	Status bool
}

//创建会话
//必须执行该函数，才能使用其他内部函数
//param name string 标记
func (this *SessionOperate) Create(name string) {
	this.store = sessions.NewCookieStore([]byte(name))
	this.Status = true
}

//获取会话标记对象
//param r *http.Request Http读取对象
//param appName string 应用标记
//param mark string 标记
//return map[interface{}]interface{}, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(r *http.Request,appName string,mark string) (map[interface{}]interface{}, bool) {
	var res map[interface{}]interface{}
	if this.Status == false{
		this.Create(appName)
	}
	s, err := this.store.Get(r, mark)
	if err != nil {
		LogOperate.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionGet","store-get",err.Error())
		return res, false
	}
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
	if this.Status == false{
		this.Create(appName)
	}
	s, err := this.store.Get(r, mark)
	if err != nil {
		LogOperate.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","store-get",err.Error())
		return false
	}
	s.Values = data
	err = s.Save(r, w)
	if err != nil {
		LogOperate.SendLog("core/session-operate.go",r.RemoteAddr,"SessionOperate.SessionSet","store-save",err.Error())
		return false
	}
	return true
}
