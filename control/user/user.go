package user

import (
	"../core"
	"time"
	"net/http"
	"strconv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//用户处理器包
//可用于用户管理、登录
//支持任意数据库类型，或直接制定单一用户密码
//使用方法：声明User类后初始化，之后选择单一用户还是多用户模式设定即可
//依赖外部包：
// mgo (gopkg.in/mgo.v2)
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
	//数据表
	table string
	//数据库合集
	dbColl *mgo.Collection
	//默认退出时间
	timeout int64
}

//用户字段组
type UserFields struct {
	id bson.ObjectId
	nicename string
	username string
	password string
	last_ip string
	last_time int
	is_disabled bool
}

//初始化
//param session *core.SessionOperate 会话句柄
//param mark string 标识码，用于会话等相关特定处理、密码混合加密
//param timeout int64 自动退出时间，秒为单位
func (this *User) Init(session *core.SessionOperate,mark string,timeout int64){
	this.session = session
	this.oneUserStatus = false
	this.mark = mark
	this.table = "user"
	this.fields = []string{
		"_id","nicename","username","password","last_ip","last_time","is_disabled",
	}
	res,b := this.session.SessionGet(this.mark)
	if b == true{
		var loginID int = 0
		res["login-id"] = loginID
		var loginTime int = 0
		res["login-time"] = loginTime
	}
	this.timeout = timeout
}

//设定数据库
//设定后单一用户开关将关闭，也就是说系统将默认使用数据库方式查询登录
//param db *sql.DB 数据库连接句柄
//param username string 初始化用户
//param password string 初始化密码
func (this *User) SetManyUser(db *mgo.Database,username string,password string){
	//导入到内部
	this.dbColl = db.C(this.table)
	this.oneUserStatus = false
	//如果数据集合中，不存在数据，则自动创建用户
	var num int
	var err error
	num,err = this.dbColl.Count()
	if err != nil{
		core.SendLog(err.Error())
		return
	}
	if num > 0{
		return
	}
	//构建数据
	var data UserFields
	data.nicename = username
	data.username = username
	data.password = this.getPasswdSha1(this.getSha1(password))
	data.last_ip = "0.0.0.0"
	data.last_time = 0
	data.is_disabled = false
	err = this.dbColl.Insert(&data)
	if err != nil{
		core.SendLog(err.Error())
		return
	}
	core.SendLog("初始化用户成功，新建立的用户名 : " + username + " , 密码 : " + password + "。")
}

//设定为单一用户模式
//该模式下指定特定的用户名和密码，即可实现登录和退出效果
//如果不需要用户名，直接给定空字符串即可实现
//但其他获取用户列表、信息等信息无法使用
//启动后无法关闭
//param username string 用户名
//param passwd string 密码
func (this *User) SetOneUser(username string,password string){
	this.oneUsername = username
	this.oneUserpasswd = this.getPasswdSha1(this.getSha1(password))
	this.oneUserStatus = true
}

//获取用户登录ID
//return int 登录的用户ID
func (this *User) GetLoginStatus() int{
	//获取值
	var res map[interface{}]interface{}
	var b bool
	res,b = this.session.SessionGet(this.mark)
	if b == false{
		return 0
	}
	//更新登录时间值
	if res["login-id"].(int) > 0{
		var t time.Time
		t = time.Now()
		var unixTime int64
		unixTime = t.Unix()
		//超出时间，强行退出
		if this.timeout > unixTime - res["login-time"].(int64){
			var loginID int = 0
			res["login-id"] = loginID
			_ = this.session.SessionSet(this.mark,res)
			return 0
		}
		res["login-time"] = unixTime
		_ = this.session.SessionSet(this.mark,res)
	}
	//返回
	return res["login-id"].(int)
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
	var loginID int = 0
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
	var t time.Time
	t = time.Now()
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
		var result UserFields
		err = this.dbColl.Find(bson.M{"username":username,"password":passwdSha1Sha1}).One(&result)
		if err != nil{
			core.SendLog(err.Error())
			return false
		}
		//用户存在，则修改登录IP和时间
		var userID int
		userID,err = strconv.Atoi(result.id.String())
		if err != nil{
			core.SendLog(err.Error())
			return false
		}
		if userID > 0{
			err = this.dbColl.UpdateId(userID,bson.M{"last_ip":ipAddr,"last_time":unixTime})
			if err != nil{
				core.SendLog(err.Error())
				return false
			}
			loginID = userID
		}
	}
	//检查是否验证通过
	if loginID < 1{
		return false
	}
	//输出日志
	core.SendLog("用户" + strconv.Itoa(loginID) + "通过IP地址" + ipAddr + "登录了系统。")
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
	if b == false{
		return
	}
	if res["login-id"].(int) < 1{
		return
	}
	var loginID int = 0
	res["login-id"] = loginID
}

//根据ID查询用户信息
//param id int 用户ID
//return *UserFields,bool 用户信息组，是否成功
func (this *User) GetID(id int) (*UserFields,bool){
	//初始化变量
	var result UserFields
	//获取数据
	var err error
	err = this.dbColl.FindId(id).One(&result)
	if err != nil{
		core.SendLog(err.Error())
		return &result,false
	}
	//返回数据
	return &result,true
}

//查询用户列表
//param search string 搜索昵称或用户名
//param page int 页数
//param max int 页码
//param sort int 排序字段键值
//param desc bool 是否倒序
//return []UserFields,bool 数据结果，是否成功
func (this *User) GetList(search string,page int,max int,sort int,desc bool) (*[]UserFields,bool){
	//分析sortStr
	var sortStr string
	if this.fields[sort] != ""{
		sortStr = this.fields[sort]
	}else{
		sortStr = "_id"
	}
	//分析desc
	if desc == true{
		sortStr = "-" + sortStr
	}
	//获取数据
	var result []UserFields
	var err error
	var skip int
	skip = (page - 1) * max
	if search == ""{
		err = this.dbColl.Find(nil).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	}else{
		err = this.dbColl.Find(bson.M{"$or":bson.M{"nicename":search,"username":search}}).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	}
	if err != nil{
		core.SendLog(err.Error())
		return &result,false
	}
	//返回结果
	return &result,true
}

//创建新用户
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return int 新的用户ID，失败返回0，发现用户名存在返回-1
func (this *User) Create(nicename string,username string,passwdSha1 string) (int){
	//检查昵称、用户名、密码是否合法
	if this.checkNicename(nicename) == false || this.checkUsername(username,passwdSha1) == false{
		return 0
	}
	//检查用户是否存在
	if this.checkUsernameIsExisit(username) == true{
		return -1
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//执行创建用户
	var err error
	var data UserFields
	data.nicename = nicename
	data.username = username
	data.password = passwdSha1Sha1
	data.last_ip = "0.0.0.0"
	data.last_time = 0
	data.is_disabled = false
	err = this.dbColl.Insert(&data)
	if err != nil{
		return 0
	}
	//查询新的ID
	var newData UserFields
	err = this.dbColl.Find(bson.M{"username":username}).One(&newData)
	if err != nil{
		core.SendLog(err.Error())
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
func (this *User) Edit(id int,nicename string,username string,passwdSha1 string) bool{
	//检查昵称、用户名和密码是否合法
	if this.checkNicename(nicename) == false || this.checkUsername(username,passwdSha1) == false{
		return false
	}
	//检查用户是否存在
	if this.checkUsernameIsExisit(username) == true{
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//执行修改用户
	var err error
	err = this.dbColl.UpdateId(id,bson.M{"nicename":nicename,"username":username,"password":passwdSha1Sha1})
	return err != nil
}

//删除用户
//param id int 用户ID
//return bool是否成功
func (this *User) Delete(id int) bool{
	var err error
	err = this.dbColl.RemoveId(id)
	return err != nil
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

//检查昵称是否合法
//param nicename string 昵称
//return bool 是否成功
func (this *User) checkNicename(nicename string) bool{
	if len(nicename) < 1 || len(nicename) > 30{
		return false
	}
	return true
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

//检查用户名是否存在
func (this *User) checkUsernameIsExisit(username string) bool{
	var result UserFields
	var err error
	err = this.dbColl.Find(bson.M{"username":username}).One(&result)
	if err != nil{
		return false
	}
	var userID int
	userID,err = strconv.Atoi(result.id.String())
	if err != nil{
		core.SendLog(err.Error())
		return false
	}
	return userID > 0
}