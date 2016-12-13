package user

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

//该文件定义修改用户数据

//创建新用户
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return int 新的用户ID，失败返回0，发现用户名存在返回-1
func Create(nicename string,username string,passwdSha1 string) (int){
	//检查昵称、用户名、密码是否合法
	if checkNicename(nicename) == false || checkUsername(username,passwdSha1) == false{
		return 0
	}
	//检查用户是否存在
	if checkUsernameIsExisit(username) == true{
		return -1
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = getPasswdSha1(passwdSha1)
	//执行创建用户
	var err error
	var data UserFields
	data.nicename = nicename
	data.username = username
	data.password = passwdSha1Sha1
	data.last_ip = "0.0.0.0"
	data.last_time = 0
	data.is_disabled = false
	err = dbColl.Insert(&data)
	if err != nil{
		return 0
	}
	//查询新的ID
	var newData UserFields
	err = dbColl.Find(bson.M{"username":username}).One(&newData)
	if err != nil{
		sendLog(err.Error())
		return 0
	}
	var newID int
	newID,err = strconv.Atoi(newData.id.String())
	return newID
}

//修改用户名和密码
//param id int 要编辑的用户ID
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func Edit(id int,nicename string,username string,passwdSha1 string) bool{
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
//param id int 用户ID
//return bool是否成功
func Delete(id int) bool{
	var err error
	err = dbColl.RemoveId(id)
	return err != nil
}
