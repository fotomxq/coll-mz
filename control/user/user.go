package user

import (
	"github.com/fotomxq/coll-mz/control/core"
	"database/sql"
	"time"
	"net/http"
	"strconv"
)

//用户处理器包
//可用于用户管理、登录
//支持任意数据库类型，或直接制定单一用户密码
//使用方法：声明User类后初始化，之后选择单一用户还是多用户模式设定即可
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
	//字段列
	fields []string
	//数据库句柄
	db *core.Database
	//默认退出时间
	timeout int
}

//用户字段组
type UserFields struct {
	id int64
	nicename string
	username string
	password string
	last_ip string
	last_time int64
	is_disabled int
}

//初始化
//param session *core.SessionOperate 会话句柄
//param mark string 标识码，用于会话等相关特定处理、密码混合加密
//param timeout int 自动退出时间，秒为单位
func (this *User) init(session *core.SessionOperate,mark string,timeout int){
	this.session = session
	this.oneUserStatus = false
	this.mark = mark
	this.fields = []string{
		"id","nicename","username","password","last_ip","last_time","is_disabled",
	}
	res,b := this.session.SessionGet(this.mark)
	if b == true{
		var loginID int64 = 0
		res["login-id"] = loginID
		var loginTime int64 = 0
		res["login-time"] = loginTime
	}
	this.timeout = timeout
}

//设定数据库
//设定后单一用户开关将关闭，也就是说系统将默认使用数据库方式查询登录
//param db *sql.DB 数据库连接句柄
func (this *User) SetManyUser(db *core.Database){
	this.db = db
	this.oneUserStatus = false
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
	this.oneUserpasswd = this.getPasswdSha1(this.getSha1(passwd))
	this.oneUserStatus = true
}

//获取用户登录ID
//return int64 登录的用户ID
func (this *User) GetLoginStatus() int64{
	//获取值
	var res map[interface{}]interface{}
	var b bool
	res,b = this.session.SessionGet(this.mark)
	if b == false{
		return 0
	}
	//更新登录时间值
	if res["login-id"] > 0{
		var t *time.Time
		t = &time.Now()
		var unixTime int64
		unixTime = t.Unix()
		//超出时间，强行退出
		if this.timeout > unixTime - res["login-time"]{
			var loginID int64 = 0
			res["login-id"] = loginID
			_ = this.session.SessionSet(this.mark,res)
			return false
		}
		res["login-time"] = unixTime
		_ = this.session.SessionSet(this.mark,res)
	}
	//返回
	return res["login-id"]
}

//用户登录
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//param r *http.Request HTTP读取句柄
//return bool 是否登录成功
func (this *User) Login(username string,passwdSha1 string,r *http.Request) bool{
	//初始化变量
	var res map[interface{}]interface{}
	var b bool
	var err error
	res,b = this.session.SessionGet(this.mark)
	if b == false{
		return false
	}
	var loginID int64 = 0
	//是否已经登录，是则返回成功
	if this.GetLoginStatus() > 0{
		return true
	}
	//检查用户名和密码是否合法
	if this.checkUsername(username,passwdSha1) == false{
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//获取IP地址
	var ipAddr string
	ipAddr = r.RemoteAddr
	//获取当前时间
	var t *time.Time
	t = &time.Now()
	var unixTime int64
	unixTime = t.Unix()
	//检查模式
	if this.oneUserStatus == true{
		//如果是单用户模式
		if this.oneUsername == username && passwdSha1Sha1 == this.oneUserpasswd{
			loginID = 1
		}else{
			return false
		}
	}else{
		//如果是多用户模式
		var query string
		query = "select `id` from `user` where `username` = ? and `password` = ? and `is_disabled` = 0"
		var row *sql.Row
		row = this.db.DB.QueryRow(query,username,passwdSha1Sha1)
		var id int64
		err = row.Scan(&id)
		if err != nil{
			core.SendLog(err.Error())
		}
		//用户存在，则修改登录IP和时间
		if id > 0{
			var querySet string
			querySet = "update `user` set `last_ip` = ? , `last_time` = ? where `id` = ?"
			_,err = this.db.DB.Exec(querySet,ipAddr,unixTime,id)
			if err != nil{
				core.SendLog(err.Error())
			}
			loginID = id
		}
	}
	//检查是否验证通过
	if loginID < 1{
		return false
	}
	//输出日志
	core.SendLog("用户" + strconv.FormatInt(loginID,10) + "通过IP地址" + ipAddr + "登录了系统。")
	//修改session
	res["login-id"] = loginID
	res["login-time"] = unixTime
	return this.session.SessionSet(this.mark,res)
}

//用户退出
func (this *User) Logout(){
	var res map[interface{}]interface{}
	var b bool
	res,b = this.session.SessionGet(this.mark)
	if b == true{
		var loginID int64 = 0
		res["login-id"] = loginID
	}
}

//根据ID查询用户信息
//param id int64 用户ID
//return *UserFields,bool 用户信息组，是否成功
func (this *User) GetID(id int64) (*UserFields,bool){
	//初始化变量
	var result UserFields
	var row *sql.Row
	var b bool
	//读取数据
	this.db.Set("user",this.fields)
	row,b = this.db.GetID(id)
	if b == false{
		return &result,false
	}
	row.Scan(&result.id,&result.nicename,&result.username,&result.password,&result.last_ip,&result.last_time,&result.is_disabled)
	//返回数据
	return &result,true
}

//查询用户列表
func (this *User) GetList(searchUser string,page int,max int,sort int,desc bool) ([]map[string]string,bool){

}

//创建新用户
func (this *User) Create(username string,passwd string) (int64){

}

//修改用户名和密码
func (this *User) Edit(id int64,nicename string,username string,passwdSha1 string) bool{
	//检查昵称、用户名和密码是否合法
	if len(nicename) < 1 && len(nicename) > 30 && this.checkUsername(username,passwdSha1) == false{
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//执行修改用户
	var query string
	query = "update `user` set `nicename` = ? , `username` = ? , `password` = ? where `id` = ?"
	var res sql.Result
	var err error
	res,err = this.db.DB.Exec(query,nicename,username,passwdSha1Sha1,id)
	if err != nil{
		core.SendLog(err.Error())
		return false
	}
	var result int64
	result,err = res.RowsAffected()
	if err != nil{
		core.SendLog(err.Error())
		return false
	}
	//返回结果
	return result > 0
}

//删除用户
func (this *User) Delete(id int64) bool{
	//初始化变量
	var row int64
	//执行删除
	this.db.Set("user",this.fields)
	row = this.db.Delete(id)
	//返回结果
	core.SendLog("删除用户" + strconv.FormatInt(id,10))
	return row > 0
}

///////////////////////////////////////////////////////////////////////
//以下是内部函数
//////////////////////////////////////////////////////////////////////
//检查用户名和密码是否正确
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func (this *User) checkUsername(username string,passwdSha1 string) bool{
	return this.matchString.CheckUsername(username,6,20) && len(passwdSha1) == 40
}

//计算字符串的SHA1值
//param str string 字符串
//return string 计算结果，失败返回空字符串
func (this *User) getSha1(str string) string{
	return this.matchString.GetSha1(str)
}

//获取加密密码
//加入mark字符串并再次计算SHA1值
//param passwd string 加密过一次的密码
//return string 加密后的密码，失败返回空字符串
func (this *User) getPasswdSha1(passwd string) string{
	var newPasswd string
	newPasswd = passwd + this.mark
	return this.getSha1(newPasswd)
}