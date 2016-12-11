package core

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

//数据库操作
//依赖内部库：
// core.SendLog()
//依赖外部库：
// github.com/mattn/go-sqlite3

//数据库类型
type Database struct {
	//数据库句柄
	DB *sql.DB
	//要操作的表
	table string
	//要操作的表字段
	fields []string
	//列字符串
	// 用于加快查询
	fieldStr string
}

//连接数据库
//必须执行该函数后才能使用其他函数
//param dbType string 数据库类型
//param dbDNS string 数据库DNS
//return bool 连接是否成功
func (this *Database) Connect(dbType string, dbDNS string) bool {
	if this.Status() == true {
		return true
	}
	this.DB, err = sql.Open(dbType, dbDNS)
	if err != nil {
		SendLog(err.Error())
		return false
	}
}

//关闭数据库连接
func (this *Database) Close() {
	err = this.DB.Close()
	if err != nil{
		SendLog(err.Error())
	}
}

//查看连接状态
//return bool 是否有连接
func (this *Database) Status() bool{
	stats := this.DB.Stats()
	if stats.OpenConnections > 0{
		return true
	}
	return false
}

//设定要操作的表
//param table string 要操作的表
//param fields []string 要操作的字段组
func (this *Database) Set(table string,fields []string){
	this.table = table
	this.fields = fields
	this.fieldStr = this.GetFieldsToStr()
}

//获取指定ID的单行数据
//param id int64 ID
//return *sql.Row, bool 数据，是否成功
func (this *Database) GetID(id int64) (*sql.Row, bool) {
	if this.Status() == false{
		return nil,false
	}
	query := "select " + this.fieldStr + " from `" + this.table + "` where `id` = ?"
	stmt, err := this.DB.Prepare(query)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	return row, true
}

//获取指定列值单行数据
//param field string 指定列
//param value string 指定值
//return *sql.Row, bool 数据，是否成功
func (this *Database) GetFieldValue(field string,value string) (*sql.Row, bool) {
	if this.Status() == false{
		return nil,false
	}
	query := "select " + this.fieldStr + " from `" + this.table + "` where `" + field + "` = ?"
	stmt, err := this.DB.Prepare(query)
	if err != nil {
		return nil, false
	}
	defer stmt.Close()
	row := stmt.QueryRow(value)
	return row, true
}

//获取多行数据
//param page int 页数
//param max int 页长
//param sort int 排序字段键值
//param desc bool 是否倒序
//return map[int]map[string]interface{}, bool 数据，是否成功
func (this *Database) GetList(page int, max int, sort int, desc bool) (map[int]map[string]interface{}, bool) {
	var returnResult map[int]map[string]interface{}
	if this.Status() == false{
		return returnResult,false
	}
	query := "select " + this.fieldStr + " from `" + this.table + "` " + this.GetPageSortStr(page, max, this.fields[sort], desc)
	stmt, err := this.DB.Prepare(query)
	if err != nil {
		SendLog(err.Error())
		return returnResult,false
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		SendLog(err.Error())
		return returnResult,false
	}
	return this.GetResultToList(rows),true
}

//插入新的数据
//值无过滤提交，建议内部使用
//param values []string 插入的值
//return int64 ID
func (this *Database) Insert(values []string) int64 {
	if this.Status() == false{
		return -1
	}
	query := "insert into `" + this.table + "`(" + this.fieldStr + ") values("
	var valuesLen int = len(values)-1
	for key := range values{
		query += "'" + values[key] + "'"
		if key != valuesLen{
			query += ","
		}
	}
	query += ")"
	stmt, err := this.DB.Exec(query)
	if err != nil {
		SendLog(err.Error())
		return -1
	}
	row,err := stmt.LastInsertId()
	if err != nil{
		SendLog(err.Error())
		return -1
	}
	return row
}

//更新数据
//param setField string 设定的列
//param setValue string 设定的值
//param id int64 要修改的记录ID
//return int64 影响的记录，-1为出错
func (this *Database) Update(setField string, setValue string, id int64) int64 {
	if this.Status() == false{
		return -1
	}
	query := "update `" + this.table + "` set `" + setField + "` = ? where `id` = ?"
	stmt, err := this.DB.Exec(query, setValue, id)
	if err != nil {
		SendLog(err.Error())
		return -1
	}
	row,err := stmt.RowsAffected()
	if err != nil{
		SendLog(err.Error())
		return -1
	}
	return row
}

//删除指定ID记录
//param id int64 ID
//return int64 影响的记录，-1为出错
func (this *Database) Delete(id int64) int64 {
	if this.Status() == false{
		return -1
	}
	query := "delete from `" + this.table + "` where `id` = ?"
	smat, err := this.DB.Exec(query, id)
	if err != nil {
		SendLog(err.Error())
		return -1
	}
	row, err := smat.RowsAffected()
	if err != nil {
		SendLog(err.Error())
		return -1
	}
	return row
}

//获取列字符串
//return string 列字符串
func (this *Database) GetFieldsToStr() string {
	return "`" + strings.Join(this.fields, "`,`") + "`"
}

//通过参数生成页码部分的sql
//param page int 页数
//param max int 页长
//param sort string 排序字段
//param desc bool 是否倒序
//return string sql
func (this *Database) GetPageSortStr(page int, max int, sort string, desc bool) string {
	var descStr string
	if desc == true {
		descStr = "desc"
	} else {
		descStr = "asc"
	}
	newPage := (page-1) * max
	return "order by `" + sort + "` " + descStr + " limit " + strconv.Itoa(newPage) + "," + strconv.Itoa(max)
}

//从多行结果中获取数据
//param result *sql.Rows 数据集
//return map[int]map[string]interface{} 数据数组
func (this *Database) GetResultToList(result *sql.Rows) map[int]map[string]interface{} {
	var resultArray map[int]map[string]interface{}
	defer result.Close()
	var j int = 0
	for {
		b := result.Next()
		if b == false {
			break
		}
		var a map[string]interface{}
		c, err := result.Columns()
		if err != nil {
			break
		}
		for i := range this.fields {
			a[this.fields[i]] = c[i]
		}
		resultArray[j] = a
		j++
	}
	return resultArray
}

