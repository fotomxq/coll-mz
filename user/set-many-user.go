package user

//设定数据库
//设定后单一用户开关将关闭，也就是说系统将默认使用数据库方式查询登录
//在使用前，务必转移好User全局数据库句柄
//param username string 初始化用户
//param password string 初始化密码
func SetManyUser(username string,password string){
	//导入到内部
	dbColl = DB.C(table)
	oneUserStatus = false
	//如果数据集合中，不存在数据，则自动创建用户
	var num int
	var err error
	num,err = dbColl.Count()
	if err != nil{
		sendLog(err.Error())
		return
	}
	if num > 0{
		return
	}
	//构建数据
	var data UserFields
	data.nicename = username
	data.username = username
	data.password = getPasswdSha1(getSha1(password))
	data.last_ip = "0.0.0.0"
	data.last_time = 0
	data.is_disabled = false
	err = dbColl.Insert(&data)
	if err != nil{
		sendLog(err.Error())
		return
	}
	sendLog("初始化用户成功，新建立的用户名 : " + username + " , 密码 : " + password + "。")
}
