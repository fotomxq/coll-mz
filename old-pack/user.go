package controller

import (
	"net/http"
	"strconv"
	"time"
)

//The user operates the module
type User struct {
	db                 *Database
	loginMark          string
	loginSessionMark   string
	timeoutSessionMark string
	timeout            int64
	fields             []string
	ip                 string
	matchString        MatchString
}

type UserFieldsStruct struct {
	id         int64
	username   string
	password   string
	last_ip     string
	last_time   string
	is_disabled int
}

//Initialize the user action module
func (this *User) Init(db *Database, timeout int64) {
	this.loginMark = "login"
	this.loginSessionMark = "logged-ok"
	this.timeoutSessionMark = "last-time"
	this.timeout = timeout
	this.db = db
	this.fields = []string{
		"id",
		"username",
		"password",
		"last_ip",
		"last_time",
		"is_disabled",
	}
}

//Update the IP address
func (this *User) UpdateIP() {
	this.ip = IPAddrsGetExternal()
	if this.ip == "" {
		this.ip = IPAddrsGetInternal()
	}
}

//log in
func (this *User) LoginIn(w http.ResponseWriter, r *http.Request, username string, passwd string) bool {
	//debug auto login
	if configData["debug"] == "true" {
		log.NewLog(" ## debug ## auto login.", nil)
		b := this.ChangeLoginSession(w, r, true)
		return b
	}
	//Check the user name and password
	if username == "" || passwd == "" || this.matchString.CheckEmail(username) || this.matchString.CheckPassword(passwd){
		log.NewLog("The user name and password are illegal.",nil)
		return false
	}
	//Check that you are logged in
	if this.CheckLogin(w, r) == true {
		log.NewLog("The user tries to enter the account password to log in, but has actually logged on.", nil)
		return true
	}
	//Calculate the password
	passwdSha1 := this.GetPasswdSha1(passwd)
	if passwdSha1 == "" {
		log.NewLog("The SHA1 value of the password can not be calculated.", nil)
		return false
	}
	//
	t := time.Now()
	t.Format("")
	//Query whether the user exists
	query := "select `id` from `user` where `is_disabled` = 0 and `username` = ? and `password` = ?"
	stmt, err := this.db.db.Prepare(query)
	if err != nil {
		log.NewLog("The user tried to log on to the system, but the user name and password were wrong,E1.", err)
		return false
	}
	defer stmt.Close()
	var id int64
	result := stmt.QueryRow(username,passwdSha1)
	err = result.Scan(&id)
	if err != nil {
		log.NewLog("The user tried to log on to the system, but the user name and password were wrong,E3.", err)
		return false
	}
	if id < 1 {
		log.NewLog("The user tried to log on to the system, but the user name and password were wrong,E4.", nil)
		return false
	}
	//Update table information
	b := this.UpdateLoginInfo(id)
	if b == false {
		return false
	}
	//Update the login status
	if id > 0 {
		log.NewLog("The user successfully logged on to the system,user name : " + username, nil)
		b := this.ChangeLoginSession(w, r, true)
		return b
	}
	log.NewLog("The user can not log on for some unknown reason.", nil)
	return false
}

//sign out
func (this *User) Logout(w http.ResponseWriter, r *http.Request) bool {
	log.NewLog("The user exited the system.", nil)
	return this.ChangeLoginSession(w, r, false)
}

//Check the login status
func (this *User) CheckLogin(w http.ResponseWriter, r *http.Request) bool {
	return this.GetLoginSession(w, r)
}

//Gets the specified user
func (this *User) ViewUser(id int) (UserFieldsStruct, bool) {
	row, err := this.db.GetID("user", this.fields, id)
	var res UserFieldsStruct
	if err != nil {
		log.NewLog("Failed to query user information.", err)
		return res, false
	}
	row.Scan(&res.id,&res.username,&res.password,&res.last_ip,&res.last_time,&res.is_disabled)
	return res, true
}

//Gets the list of users
func (this *User) ViewUserList(searchUser string, page int, max int, sort int, desc bool) (map[int]map[string]interface{}, bool) {
	var result map[int]map[string]interface{}
	query := "select " + this.db.GetFieldsToStr(this.fields) + " from `user` where `username` like ? " + this.db.GetPageSortStr(page, max, this.fields[sort], desc)
	stmt, err := this.db.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.NewLog("", err)
		return result, false
	}
	rows, err := stmt.Query(searchUser)
	defer rows.Close()
	if err != nil {
		log.NewLog("", err)
		return result, false
	}
	result, err = this.db.GetResultToList(this.fields, rows)
	if err != nil {
		log.NewLog("", err)
		return result, false
	}
	return result, true
}

//Create a new user
func (this *User) CreateNewUser(username string, passwd string) int64 {
	//Check if the user exists
	searchUserRes := this.SearchUsername(username)
	if searchUserRes > 0 {
		log.NewLog("Can not create new user because the user already exists.", nil)
		return 0
	}
	//Start creating a new user
	query := "insert into `user`(" + this.db.GetFieldsToStr(this.fields) + ") values(null,?,?,?,?,0)"
	this.UpdateIP()
	passwdSha1 := this.GetPasswdSha1(passwd)
	if passwdSha1 == "" {
		log.NewLog("The SHA1 value of the password can not be calculated.", nil)
		return 0
	}
	stmt, err := this.db.db.Exec(query, username, passwdSha1, this.ip,this.GetNowTimeUnix())
	if err != nil {
		log.NewLog("", err)
		return 0
	}
	id, err := stmt.LastInsertId()
	if err != nil {
		log.NewLog("", err)
		return 0
	}
	return id
}

//Search for the user name get column
func (this *User) SearchUsername(username string) int64 {
	query := "select `id` from `user` where `username` = ?"
	stmt, err := this.db.db.Prepare(query)
	if err != nil {
		return 0
	}
	defer stmt.Close()
	row := stmt.QueryRow(username)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

//Update user information
func (this *User) UpdateUser(id int, username string, passwd string, isDisabled bool) bool {
	query := "update `user` set `username` = ? , `passwd` = ? , `is_disabled` = ? where `id` = ?"
	passwdSha1 := this.GetPasswdSha1(passwd)
	if passwdSha1 == "" {
		log.NewLog("Can not calculate password SHA1 value.", nil)
		return false
	}
	var isDisabledInt int
	if isDisabled == true {
		isDisabledInt = 1
	} else {
		isDisabledInt = 0
	}
	stmt, err := this.db.db.Exec(query, username, passwdSha1, isDisabledInt)
	if err != nil {
		log.NewLog("Unable to update user information E1,user id : "+strconv.Itoa(id), err)
	}
	row, err := stmt.RowsAffected()
	if err == nil && row > 0 {
		log.NewLog("Update the user information successfully.", nil)
		return true
	}
	log.NewLog("Unable to update user information E2,user id : "+strconv.Itoa(id), err)
	return false
}

//delete user
func (this *User) DeleteUser(id int64) bool {
	row, err := this.db.Delete("user", id)
	if err == nil && row > 0 {
		log.NewLog("User successfully deleted.", nil)
		return true
	}
	log.NewLog("An error occurred while deleting the user.", err)
	return false
}

//Change the session state
func (this *User) ChangeLoginSession(w http.ResponseWriter, r *http.Request, b bool) bool {
	s, err := store.Get(r, this.loginMark)
	if err != nil {
		log.NewLog("", err)
		return false
	}
	if b == true {
		s.Values[this.loginSessionMark] = "in"
		s.Values[this.timeoutSessionMark] = this.GetUnixTime()
	} else {
		s.Values[this.loginSessionMark] = "out"
		s.Values[this.timeoutSessionMark] = 0
	}
	err = s.Save(r, w)
	if err != nil {
		log.NewLog("", err)
		return false
	}
	return true
}

//Gets the session state
func (this *User) GetLoginSession(w http.ResponseWriter, r *http.Request) bool {
	data, b := SessionGet(w, r, this.loginMark)
	if b == false {
		return false
	}
	if data[this.loginSessionMark] == nil {
		return false
	}
	if data[this.timeoutSessionMark] == nil {
		data[this.timeoutSessionMark] = this.GetUnixTime()
	}
	if data[this.loginSessionMark].(string) == "in" {
		c := this.GetUnixTime() - data[this.timeoutSessionMark].(int64)
		if c >= this.timeout {
			data[this.loginSessionMark] = "out"
			return false
		}
		return this.ChangeLoginSession(w,r,true)
	}
	return false
}

//Gets the unix timestamp
func (this *User) GetUnixTime() int64 {
	return time.Now().Unix()
}

//Update table information
func (this *User) UpdateLoginInfo(id int64) bool {
	query := "update `user` set `last_ip` = ?,`last_time` = ? where `id` = ?"
	this.UpdateIP()
	smat, err := this.db.db.Exec(query, this.ip,this.GetNowTimeUnix(), id)
	if err != nil {
		log.NewLog("The user information can not be updated when the user logs in E1.", err)
		return false
	}
	row, err := smat.RowsAffected()
	if err == nil && row > 0 {
		return true
	}
	log.NewLog("The user information can not be updated when the user logs in E2.", err)
	return false
}

//Gets the password SHA1 value
func (this *User) GetPasswdSha1(passwd string) string {
	return this.matchString.GetSha1(passwd)
}

//get now time unix
func (this *User) GetNowTimeUnix() int64{
	t := time.Now()
	return t.Unix()
}