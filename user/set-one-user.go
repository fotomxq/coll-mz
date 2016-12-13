package user

//设定为单一用户模式
//该模式下指定特定的用户名和密码，即可实现登录和退出效果
//如果不需要用户名，直接给定空字符串即可实现
//但其他获取用户列表、信息等信息无法使用
//启动后无法关闭
//param username string 用户名
//param passwd string 密码
func SetOneUser(username string,password string){
	oneUsername = username
	oneUserpasswd = getPasswdSha1(getSha1(password))
	oneUserStatus = true
}
