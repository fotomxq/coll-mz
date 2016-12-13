package handle

//启动会话模块

//启动会话
func startSession(){
	if SessionOperate.Status == false{
		SessionOperate.Create(AppName)
	}
}
