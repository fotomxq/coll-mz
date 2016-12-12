package core

import (
	"github.com/gorilla/sessions"
	"net/http"
)

//操作Session会话
//依赖内部库：
// core.SendLog()
//依赖外部库：
// github.com/gorilla/sessions

//Session操作类型
type SessionOperate struct {
	//采集器存储句柄
	store *sessions.CookieStore
	//会话启动状态
	Status bool
	//接收和反馈操作句柄
	w http.ResponseWriter
	r *http.Request
}

//创建会话
//必须执行该函数，才能使用其他内部函数
//param name string 标记
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
func (this *SessionOperate) Create(name string,w http.ResponseWriter, r *http.Request) {
	this.store = sessions.NewCookieStore([]byte(name))
	this.w = w
	this.r = r
	this.Status = true
}

//获取会话标记对象
//param mark string 标记
//return map[interface{}]interface{}, bool 会话变量组合，是否失败
func (this *SessionOperate) SessionGet(mark string) (map[interface{}]interface{}, bool) {
	var res map[interface{}]interface{}
	if this.Status == false{
		return res,false
	}
	s, err := this.store.Get(this.r, mark)
	if err != nil {
		SendLog(err.Error())
		return res, false
	}
	return s.Values, true
}

//写入会话数据
//param mark string 会话标记
//param data map[interface{}]interface{} 会话变量组合
//return bool 是否失败
func (this *SessionOperate) SessionSet(mark string, data map[interface{}]interface{}) bool {
	if this.Status == false{
		return false
	}
	s, err := this.store.Get(this.r, mark)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	s.Values = data
	err = s.Save(this.r, this.w)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	return true
}
