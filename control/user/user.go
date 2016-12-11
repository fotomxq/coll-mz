package user

import "github.com/fotomxq/coll-mz/control/core"

//用户处理器包
//可用于用户管理、登录
//支持任意数据库类型，或直接制定单一用户密码
//使用方法：声明User类后初始化即可使用
//依赖外部包：
//依赖本地包：
// core.session-operate.go
// core.id-addrs.go
// core.match-string.go
// core.language.go
// core.database.go

//用户类
type User struct {
	//session会话操作
	session *core.SessionOperate
	//单一用户模式是否启动
	oneUserStatus bool
	//单一用户名和密码
	oneUsername string
	oneUserpasswd string
	//标识码
	mark string
	//验证句柄
	matchString core.MatchString
}

//初始化
//param session *core.SessionOperate 会话句柄
//param mark string 标识码，用于会话等相关特定处理
func (this *User) init(session *core.SessionOperate,mark string){
	this.session = session
	this.oneUserStatus = false
	this.mark = mark
}

//设定为单一用户模式
//该模式下指定特定的用户名和密码，即可实现登录和退出效果
//如果不需要用户名，直接给定空字符串即可实现
//但其他获取用户列表、信息等信息无法使用
//启动后无法关闭
//param username string 用户名
//param passwd string 密码
func (this *User) SetOneUser(username string,passwd string){
	this.oneUsername = username
	this.oneUserpasswd = passwd
	this.oneUserStatus = true
}

//用户登录
func (this *User) Login(username string,passwdSha1 string){
	if this.oneUserStatus == true{
		if this.oneUsername
	}
}

//用户退出
func (this *User) Logout(){

}

//根据ID查询用户信息
func (this *User) GetID(){

}

//查询用户列表
func (this *User) GetList(searchUser string,page int,max int,sort int,desc bool) ([]map[string]string,bool){

}

//创建新用户
func (this *User) Create(username string,passwd string) (int64){

}

//修改用户名和密码
func (this *User) Edit(id int64,username string,passwd string) bool{

}

//删除用户
func (this *User) Delete(id int64) bool{

}

///////////////////////////////////////////////////////////////////////
//以下是内部函数
//////////////////////////////////////////////////////////////////////
//检查用户名和密码是否正确
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func (this *User) checkUsername(username string,passwdSha1 string) bool{

}