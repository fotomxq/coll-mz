package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
		sendLog("user/set-many-user.go","0.0.0.0","SetManyUser","database-count",err.Error())
		return
	}
	if num > 0{
		return
	}
	//设定索引
	var index mgo.Index
	index = mgo.Index{
		Key: []string{"-_id"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}
	dbColl.EnsureIndex(index)
	//构建数据
	var passwordSha1Sha1 string
	passwordSha1Sha1 = getPasswdSha1(getSha1(password))
	if passwordSha1Sha1 == ""{
		sendLog("user/set-many-user.go","0.0.0.0","SetManyUser","password-sha1","无法获取密码的SHA1。")
		return
	}
	err = dbColl.Insert(&UserFields{bson.NewObjectId(),username,username,passwordSha1Sha1,"0.0.0.0",0,false})
	if err != nil{
		sendLog("user/set-many-user.go","0.0.0.0","SetManyUser","insert-new-user",err.Error())
		return
	}
	sendLog("user/set-many-user.go","0.0.0.0","SetManyUser","insert-new-user","成功初始化了平台，用户名："+username+"，密码："+password)
}
