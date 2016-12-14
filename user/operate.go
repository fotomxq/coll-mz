package user

import (
	"gopkg.in/mgo.v2/bson"
)

//该文件定义修改用户数据

//创建新用户
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return string 新的用户ID，失败返回空字符串，发现用户名存在返回"user-exist"字符串
func Create(nicename string,username string,passwdSha1 string) string{
	//检查昵称、用户名、密码是否合法
	if checkNicename(nicename) == false || checkUsername(username,passwdSha1) == false{
		sendLog("user/operate.go","0.0.0.0","Create","check-user-password","尝试创建新的用户，但用户名和密码不正确。")
		return ""
	}
	//检查用户是否存在
	if checkUsernameIsExisit(username) == true{
		sendLog("user/operate.go","0.0.0.0","Create","user-exisit","创建新的用户，但用户已经存在了。")
		return "user-exist"
	}
	//计算密码
	var passwordSha1Sha1 string
	passwordSha1Sha1 = getPasswdSha1(passwdSha1)
	//执行创建用户
	var err error
	err = dbColl.Insert(&UserFields{bson.NewObjectId(),username,username,passwordSha1Sha1,"0.0.0.0",0,false})
	if err != nil{
		sendLog("user/operate.go","0.0.0.0","Create","insert-new-user",err.Error())
		return ""
	}
	//查询新的ID
	var newData UserFields
	err = dbColl.Find(bson.M{"username":username}).One(&newData)
	if err != nil{
		sendLog("user/operate.go","0.0.0.0","Create","insert-after-find",err.Error())
		return ""
	}
	return newData.Id_.String()
}

//修改用户名和密码
//param id string 要编辑的用户ID
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func Edit(id string,nicename string,username string,passwdSha1 string) bool{
	//检查昵称、用户名和密码是否合法
	if checkNicename(nicename) == false || checkUsername(username,passwdSha1) == false{
		return false
	}
	//检查用户是否存在
	if checkUsernameIsExisit(username) == true{
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = getPasswdSha1(passwdSha1)
	//执行修改用户
	var err error
	err = dbColl.UpdateId(id,bson.M{"nicename":nicename,"username":username,"password":passwdSha1Sha1})
	return err != nil
}

//删除用户
//param id string 用户ID
//return bool是否成功
func Delete(id string) bool{
	var err error
	err = dbColl.RemoveId(id)
	return err != nil
}
