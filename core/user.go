package core

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

//用户处理器包
//可用于用户管理、登录
//支持mongo数据库，或直接指定用户名和密码
//使用方法：声明User类后初始化，之后选择单一用户还是多用户模式设定即可
//依赖外部包：
// mgo (gopkg.in/mgo.v2 / gopkg.in/mgo.v2/bson)
//依赖本地包：
// core.LogOperate
// core.session-operate.go
// core.match-string.go

//用户类
type User struct {
	//数据库操作模块
	db *mgo.Database
	//验证处理句柄
	matchString *MatchString
	//session会话操作
	sessionOperate *SessionOperate
	//日志处理器
	logOperate *LogOperate
	//应用名称，用于session大集合
	appName string
	//标识码，用于session小集合
	mark string
	//用户自动退出时限，单位：秒
	userLoginTimeout int64
	//单一用户模式是否启动
	OneUserStatus bool
	//初始化或单一用户名和密码
	oneUsername string
	onePassword string
	//字段列
	fields []string
	//数据表
	table string
	//数据库合集
	dbColl *mgo.Collection
	//权限列表
	Permissions []string
	//权限对应数据
	PermissionsData map[string]map[string]interface{}
}

//初始化用户类需要的参数
type UserParams struct {
	//数据库操作模块
	Db *mgo.Database
	//验证处理句柄
	MatchString *MatchString
	//session会话操作
	SessionOperate *SessionOperate
	//日志处理器
	LogOperate *LogOperate
	//应用名称，用于session大集合
	AppName string
	//标识码，用于session小集合
	Mark string
	//用户自动退出时限，单位：秒
	UserLoginTimeout int64
	//单一用户模式是否启动
	OneUserStatus bool
	//初始化或单一用户名和密码
	OneUsername string
	OnePassword string
	//权限列表
	Permissions []string
	//权限对应数据
	PermissionsData map[string]map[string]interface{}
}

//用户字段组
type UserFields struct {
	//索引
	ID bson.ObjectId `bson:"_id"`
	//昵称
	NiceName string
	//用户名
	UserName string
	//密码
	Password string
	//上一次登录IP
	LastIP string
	//上一次登录时间
	LastTime int64
	//是否禁用
	IsDisabled bool
	//权限标识
	Permissions []string
}

//Session结构
type UserSession struct {
	//当前登录的用户ID
	UserID string
	//最后一次活动时间，unix时间戳
	LastTime int64
	//登录时发生错误次数
	LoginErrorNum int
}

//初始化
func (this *User) Init(params *UserParams) {
	//设定基本值
	this.OneUserStatus = false
	this.fields = []string{
		"_id", "nicename", "username", "password", "lastip", "lasttime", "isdisabled", "permissions",
	}
	this.table = "user"
	//设定相关参数
	this.db = params.Db
	this.matchString = params.MatchString
	this.sessionOperate = params.SessionOperate
	this.logOperate = params.LogOperate
	this.appName = params.AppName
	this.mark = params.Mark
	this.userLoginTimeout = params.UserLoginTimeout
	this.OneUserStatus = params.OneUserStatus
	this.oneUsername = params.OneUsername
	this.onePassword = this.getPasswdSha1(this.getSha1(params.OnePassword))
	this.Permissions = params.Permissions
	this.PermissionsData = params.PermissionsData
	//如果设定的结构为数据库存储用户，则执行连接数据库集合等操作
	if params.OneUserStatus == false {
		//导入到内部
		this.dbColl = this.db.C(this.table)
		this.OneUserStatus = false
		//如果数据集合中，不存在数据，则自动创建用户
		var num int
		var err error
		num, err = this.dbColl.Count()
		if err != nil {
			this.sendLog("0.0.0.0", "SetManyUser", "database-count", err.Error())
			return
		}
		if num > 0 {
			return
		}
		//设定索引
		var index mgo.Index
		index = mgo.Index{
			Key:        []string{"-_id"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}
		this.dbColl.EnsureIndex(index)
		//构建数据
		var passwordSha1Sha1 string
		passwordSha1Sha1 = this.getPasswdSha1(this.getSha1(params.OnePassword))
		if passwordSha1Sha1 == "" {
			this.sendLog("0.0.0.0", "SetManyUser", "password-sha1", "无法获取密码的SHA1。")
			return
		}
		var permissions []string = []string{"admin"}
		err = this.dbColl.Insert(&UserFields{bson.NewObjectId(), params.OneUsername, params.OneUsername, passwordSha1Sha1, "0.0.0.0", 0, false, permissions})
		if err != nil {
			this.sendLog("0.0.0.0", "SetManyUser", "insert-new-user", err.Error())
			return
		}
		this.sendLog("0.0.0.0", "SetManyUser", "insert-new-user", "成功初始化了平台，用户名："+params.OneUsername+"，密码："+params.OnePassword)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// 用户状态、登录、退出
////////////////////////////////////////////////////////////////////////////////////////////////////////

//获取用户登录ID
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return string 登录的用户ID，未登录则返回空字符串
func (this *User) GetLoginStatus(w http.ResponseWriter, r *http.Request) string {
	//获取session
	var res *UserSession
	var b bool
	res, b = this.getSession(w, r)
	if b == false {
		return ""
	}
	if res.UserID == "" {
		return ""
	}
	if res.UserID == "one-true" {
		if this.OneUserStatus == false {
			//如果是单用户记录，但系统未启动单用户模式，则强制清空cookie并记录到log中
			this.sendLog(IPAddrsGetRequest(r), "User.GetLoginStatus", "one-user-status-is-false", "安全问题，系统未开启单用户模式，但该用户"+res.UserID+"尝试单用户模式登录。")
			_ = this.removeCookie(w, r)
		}
		return "one-true"
	} else {
		if this.OneUserStatus == true {
			//如果是多用户记录，但系统启动单用户模式，则强制清空cookie并记录到log中
			this.sendLog(IPAddrsGetRequest(r), "User.GetLoginStatus", "one-user-status-is-false", "安全问题，系统开启单用户模式，但该用户"+res.UserID+"尝试多用户模式登录。")
			_ = this.removeCookie(w, r)
		}
	}
	//更新登录时间值
	var unixTime int64
	unixTime = time.Now().Unix()
	//超出时间，强行退出
	if this.userLoginTimeout < unixTime-res.LastTime {
		res.UserID = ""
		_ = this.setSession(w, r, res)
		this.sendLog(IPAddrsGetRequest(r), "User.GetLoginStatus", "user-login-timeout-minute", res.UserID+"用户登录超时，自动退出。")
	}
	res.LastTime = unixTime
	_ = this.setSession(w, r, res)
	//返回
	return res.UserID
}

//用户登录
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//param r *http.Request HTTP读取句柄
//return bool 是否登录成功
func (this *User) Login(w http.ResponseWriter, r *http.Request, username string, passwdSha1 string) bool {
	//初始化变量
	var res *UserSession
	var b bool
	var err error
	//获取当前时间
	var unixTime int64
	unixTime = time.Now().Unix()
	//获取session数据
	res, b = this.getSession(w, r)
	if b == false {
		return false
	}
	//检查是否超过错误次数5次
	if res.LoginErrorNum > 5 {
		//永久性无法登录，除非cookie失效
		//_ = this.setSession(w,r,res)
		return false
	}
	//是否已经登录，是则返回成功
	if this.GetLoginStatus(w, r) != "" {
		return true
	}
	//检查用户名和密码是否合法
	if this.checkUsername(username, passwdSha1) == false {
		res.LoginErrorNum += 1
		_ = this.setSession(w, r, res)
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//获取IP地址
	var ipAddr string
	ipAddr = IPAddrsGetRequest(r)
	//检查模式
	if this.OneUserStatus == true {
		//如果是单用户模式
		if this.oneUsername == username && passwdSha1Sha1 == this.onePassword {
			res.UserID = "one-true"
		} else {
			res.LoginErrorNum += 1
			_ = this.setSession(w, r, res)
			return false
		}
	} else {
		//如果是多用户模式
		var result UserFields
		err = this.dbColl.Find(bson.M{"username": username, "password": passwdSha1Sha1, "isdisabled": false}).One(&result)
		if err != nil {
			this.sendLog(ipAddr, "Login", "many-user-find", err.Error())
			res.LoginErrorNum += 1
			_ = this.setSession(w, r, res)
			return false
		}
		//用户存在，则修改登录IP和时间
		var userID string
		userID = result.ID.Hex()
		if userID != "" {
			err = this.dbColl.Update(bson.M{"_id": bson.ObjectIdHex(userID)}, bson.M{"$set": bson.M{"lastip": ipAddr, "lasttime": unixTime}})
			if err != nil {
				this.sendLog(ipAddr, "Login", "many-user-update", err.Error()+" , user id : "+userID+" , lastip : "+ipAddr+" , lasttime : "+strconv.FormatInt(unixTime, 10))
				res.LoginErrorNum += 1
				this.setSession(w, r, res)
				return false
			}
			res.UserID = userID
		} else {
			res.LoginErrorNum += 1
			_ = this.setSession(w, r, res)
			return false
		}
	}
	//输出日志
	this.sendLog(ipAddr, "Login", "login-success", "ID为"+res.UserID+"的用户成功登录了平台。")
	//修改session并返回
	res.LastTime = unixTime
	res.LoginErrorNum = 0
	return this.setSession(w, r, res)
}

//用户退出
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
func (this *User) Logout(w http.ResponseWriter, r *http.Request) {
	var res UserSession
	res.UserID = ""
	_ = this.setSession(w, r, &res)
}

//检查用户是否可访问该页面
//param userID string 用户ID
//param page string 页面名称
//return bool 是否成功
func (this *User) CheckUserVisitPage(userID string, page string) bool {
	//如果是单用户，则直接返回
	if this.OneUserStatus == true || userID == "one-true" {
		return true
	}
	//初始化
	var userInfo *UserFields
	//获取数据
	var err error
	err = this.dbColl.Find(bson.M{"_id": bson.ObjectIdHex(userID), "isdisabled": false}).One(&userInfo)
	if err != nil {
		this.sendLog("0.0.0.0", "User.CheckUserVisitPage", "no-user", "检查用户访问权限的时候，无法获取用户信息。")
		return false
	}
	for _, v := range userInfo.Permissions {
		var permissionData map[string]interface{} = map[string]interface{}{}
		permissionData = this.PermissionsData[v]
		var pages []string
		pages = permissionData["page"].([]string)
		for _, v2 := range pages {
			if v2 == "*" || page == v2 {
				return true
			}
		}
	}
	this.sendLog("0.0.0.0", "User.CheckUserVisitPage", "no-permission", "用户"+userID+"尝试访问其无权访问的页面。")
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// 查看用户、用户列表
////////////////////////////////////////////////////////////////////////////////////////////////////////

//根据ID查询用户信息
//param id string 用户ID
//return *map[string]interface{},bool 用户信息组，是否成功
func (this *User) GetID(id string) (*map[string]interface{}, bool) {
	//初始化变量
	var res map[string]interface{} = map[string]interface{}{}
	var result UserFields
	//获取数据
	var err error
	err = this.dbColl.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return &res, false
	}
	//返回数据
	res = this.fieldsToMap(&result)
	return &res, true
}

//查询用户列表
//param search string 搜索昵称或用户名
//param page int 页数
//param max int 页码
//param sort int 排序字段键值
//param desc bool 是否倒序
//return *[]map[string]interface{},bool 数据结果，是否成功
func (this *User) List(search string, page int, max int, sort int, desc bool) (*[]map[string]interface{}, bool) {
	//分析sortStr
	var sortStr string
	if this.fields[sort] != "" {
		sortStr = this.fields[sort]
	} else {
		sortStr = this.fields[0]
	}
	//分析desc
	if desc == true {
		sortStr = "-" + sortStr
	}
	//限制max最大值和最小值
	if max > 100 {
		max = 100
	}
	if max < 1 {
		max = 1
	}
	//获取数据
	var result []UserFields
	var err error
	var skip int
	skip = (page - 1) * max
	if search == "" {
		err = this.dbColl.Find(nil).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	} else {
		var conditions []bson.M = []bson.M{
			{"nicename": bson.M{"$regex": bson.RegEx{search, "i"}}},
			{"username": bson.M{"$regex": bson.RegEx{search, "i"}}},
		}
		err = this.dbColl.Find(bson.M{"$or": conditions}).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	}
	var res []map[string]interface{} = []map[string]interface{}{}
	if err != nil {
		return &res, true
	}
	for _, value := range result {
		res = append(res, this.fieldsToMap(&value))
	}
	//返回结果
	return &res, true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// 创建、修改、删除用户
////////////////////////////////////////////////////////////////////////////////////////////////////////

//创建新用户
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//param permissions []string 权限列表，admin管理员权限，其他可自定义
//return string 新的用户ID，失败返回空字符串，发现用户名存在返回"user-exist"字符串
func (this *User) Create(nicename string, username string, passwdSha1 string, permissions []string) string {
	//检查用户是否存在
	if this.checkUsernameIsExisit(username) == true {
		this.sendLog("0.0.0.0", "Create", "user-exisit", "创建新的用户，但用户已经存在了。")
		return "user-exist"
	}
	//检查权限是否均有效
	if this.checkPermissions(permissions) == false {
		this.sendLog("0.0.0.0", "User.Create", "user-permissions", "该用户尝试添加一个不存在的权限。")
		return "user-permission"
	}
	//计算密码
	var passwordSha1Sha1 string
	passwordSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//执行创建用户
	var err error
	err = this.dbColl.Insert(&UserFields{bson.NewObjectId(), username, username, passwordSha1Sha1, "0.0.0.0", 0, false, permissions})
	if err != nil {
		this.sendLog("0.0.0.0", "Create", "insert-new-user", err.Error())
		return ""
	}
	//查询新的ID
	var newData UserFields
	err = this.dbColl.Find(bson.M{"username": username}).One(&newData)
	if err != nil {
		this.sendLog("0.0.0.0", "Create", "insert-after-find", err.Error())
		return ""
	}
	return newData.ID.String()
}

//修改用户名和密码
//param id string 要编辑的用户ID
//param nicename string 昵称
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//param permissions []string 权限列表，admin管理员权限，其他可自定义
//return bool 是否成功
func (this *User) Edit(id string, nicename string, username string, passwdSha1 string, permissions []string) bool {
	//获取该ID数据
	//初始化变量
	var result UserFields
	//获取数据
	var err error
	err = this.dbColl.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return false
	}
	//如果修改了用户名，则检查用户是否存在
	if result.UserName != username{
		if this.checkUsernameIsExisit(username) == true {
			return false
		}
	}
	//检查权限是否均有效
	if this.checkPermissions(permissions) == false {
		this.sendLog("0.0.0.0", "User.Edit", "user-permissions", "该用户尝试添加一个不存在的权限。")
		return false
	}
	//计算密码
	var passwdSha1Sha1 string
	passwdSha1Sha1 = this.getPasswdSha1(passwdSha1)
	//执行修改用户
	err = this.dbColl.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"nicename": nicename, "username": username, "password": passwdSha1Sha1, "permissions": permissions}})
	return err == nil
}

//删除用户
//param userID string 用户ID
//return bool是否成功
func (this *User) Delete(userID string) bool {
	var err error
	err = this.dbColl.RemoveId(bson.ObjectIdHex(userID))
	return err == nil
}

//删除所有用户
//该函数只用于debug模式
//return bool是否成功
func (this *User) DeleteAll() bool{
	var err error
	_,err = this.dbColl.RemoveAll(nil)
	return err == nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// 内部函数
////////////////////////////////////////////////////////////////////////////////////////////////////////

//检查用户名和密码是否正确
//param username string 用户名
//param passwdSha1 string 密码SHA1值
//return bool 是否成功
func (this *User) checkUsername(username string, passwdSha1 string) bool {
	return this.matchString.CheckUsername(username) && len(passwdSha1) == 40
}

//检查昵称是否合法
//param nicename string 昵称
//return bool 是否成功
func (this *User) checkNicename(nicename string) bool {
	if len(nicename) < 1 || len(nicename) > 30 || this.matchString.CheckUsername(nicename) {
		return false
	}
	return true
}

//计算字符串的SHA1值
//param str string 字符串
//return string 计算结果，失败返回空字符串
func (this *User) getSha1(str string) string {
	return this.matchString.GetSha1(str)
}

//获取加密密码
//加入mark字符串并再次计算SHA1值
//param passwd string 加密过一次的密码
//return string 加密后的密码，失败返回空字符串
func (this *User) getPasswdSha1(passwd string) string {
	var newPasswd string
	newPasswd = passwd + this.mark
	return this.getSha1(newPasswd)
}

//检查用户名是否存在
func (this *User) checkUsernameIsExisit(username string) bool {
	var result UserFields
	var err error
	err = this.dbColl.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return false
	}
	var userID string
	userID = result.ID.String()
	return userID != ""
}

//日志输出模块
//param ipAddr string IP地址
//param funcName string 函数名称
//param mark string 标记名称
//param message string 消息
func (this *User) sendLog(ipAddr string, funcName string, mark string, message string) {
	this.logOperate.SendLog("core/user.go", ipAddr, funcName, mark, message)
}

//获取session值
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return UserSession,bool Session信息组，是否成功
func (this *User) getSession(w http.ResponseWriter, r *http.Request) (*UserSession, bool) {
	var result UserSession
	var res map[string]string
	var b bool
	res, b = this.sessionOperate.SessionGet(w, r, "login")
	if b == false || res["login-user-id"] == "" || res["login-last-time"] == "" || res["login-error-num"] == "" {
		res["login-user-id"] = ""
		res["login-last-time"] = "0"
		res["login-error-num"] = "0"
		b = this.sessionOperate.SessionSet(w, r, "login", res)
		if b == false {
			this.sendLog(IPAddrsGetRequest(r), "User.getSession", "set-session-res", "无法设定session数据。")
		}
	}
	result.UserID = res["login-user-id"]
	result.LastTime, err = strconv.ParseInt(res["login-last-time"], 10, 64)
	if err != nil {
		this.sendLog(IPAddrsGetRequest(r), "User.getSession", "get-lasttime-int64", err.Error())
	}
	result.LoginErrorNum, err = strconv.Atoi(res["login-error-num"])
	if err != nil {
		result.LoginErrorNum = 0
		this.sendLog(IPAddrsGetRequest(r), "User.getSession", "get-login-error-num-int", err.Error())
	}
	return &result, true
}

//设定session
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//param data UserSession Session信息组
//return bool 是否成功
func (this *User) setSession(w http.ResponseWriter, r *http.Request, data *UserSession) bool {
	var res map[string]string
	var b bool
	res, b = this.sessionOperate.SessionGet(w, r, "login")
	if b == false {
		this.sendLog(IPAddrsGetRequest(r), "User.setSession", "get-session-res", "无法获取session数据。")
		return false
	}
	res["login-user-id"] = data.UserID
	res["login-last-time"] = strconv.FormatInt(data.LastTime, 10)
	res["login-error-num"] = strconv.Itoa(data.LoginErrorNum)
	b = this.sessionOperate.SessionSet(w, r, "login", res)
	if b == false {
		this.sendLog(IPAddrsGetRequest(r), "User.setSession", "set-session-res", "无法设定session数据。")
	}
	return b
}

//检查权限是否有效
//如果发现有一个权限不存在，则失败
//必须存在至少一个权限
//param permissions []string 要检查的权限列表
//return bool 是否有效
func (this *User) checkPermissions(permissions []string) bool {
	var onlyOne bool = false
	for _, v := range permissions {
		var noExisit bool = true
		for _, v2 := range this.Permissions {
			if v == v2 {
				noExisit = false
				onlyOne = true
				break
			}
		}
		if noExisit == true {
			return false
		}
	}
	if onlyOne == false {
		return false
	}
	return true
}

//将结构体转为map
//强行删除用户密码部分
//param info *UserFields 用户字段
//return map[string]interface{} 转换后的数据数组
func (this *User) fieldsToMap(info *UserFields) map[string]interface{} {
	var res map[string]interface{} = map[string]interface{}{}
	res = map[string]interface{}{
		"ID":          info.ID.Hex(),
		"NiceName":    info.NiceName,
		"UserName":    info.UserName,
		"LastIP":      info.LastIP,
		"LastTime":    info.LastTime,
		"IsDisabled":  info.IsDisabled,
		"Permissions": info.Permissions,
	}
	return res
}

//删除该用户客户端的cookie数据
//param w *http.ResponseWriter Http写入对象
//param r *http.Request Http读取对象
//return bool 是否成功
func (this *User) removeCookie(w http.ResponseWriter, r *http.Request) bool {
	return this.sessionOperate.RemoveCookie(w, r)
}
