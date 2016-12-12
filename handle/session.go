package handle

import (
	"net/http"
)

//启动会话模块

//启动会话
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func startSession(w http.ResponseWriter, r *http.Request){
	if SessionOperate.Status == false{
		SessionOperate.Create(AppName,w,r)
	}
}
