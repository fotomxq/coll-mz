package user

import (
	"gopkg.in/mgo.v2/bson"
)

//该文件定义内部模块

//检查用户名和密码是否正确
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func checkUsername(username string,passwdSha1 string) bool{
	return MatchString.CheckUsername(username) && len(passwdSha1) == 40
}

//检查昵称是否合法
//param nicename string 昵称
//return bool 是否成功
func checkNicename(nicename string) bool{
	if len(nicename) < 1 || len(nicename) > 30 || MatchString.CheckUsername(nicename){
		return false
	}
	return true
}

//计算字符串的SHA1值
//param str string 字符串
//return string 计算结果，失败返回空字符串
func getSha1(str string) string{
	return MatchString.GetSha1(str)
}

//获取加密密码
//加入mark字符串并再次计算SHA1值
//param passwd string 加密过一次的密码
//return string 加密后的密码，失败返回空字符串
func getPasswdSha1(passwd string) string{
	var newPasswd string
	newPasswd = passwd + Mark
	return getSha1(newPasswd)
}

//检查用户名是否存在
func checkUsernameIsExisit(username string) bool{
	var result UserFields
	var err error
	err = dbColl.Find(bson.M{"username":username}).One(&result)
	if err != nil{
		return false
	}
	var userID int
	userID = result.Id_.String()
	return userID != ""
}
