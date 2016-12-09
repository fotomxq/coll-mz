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
	db *sql.DB
	//连接状态
	status bool
	//要操作的表
	table string
	//要操作的表字段
	fields []string
}

//连接数据库
//必须执行该函数后才能使用其他函数
//param dbType string 数据库类型
//param dbDNS string 数据库DNS
//return bool 连接是否成功
func (this *Database) Connect(dbType string, dbDNS string) bool {
	if this.status == true {
		return true
	}
	this.db, err = sql.Open(dbType, dbDNS)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	this.status = true
}

//关闭数据库连接
func (this *Database) Close() {
	err = this.db.Close()
	if err != nil{
		SendLog(err.Error())
	}
	this.status = err != nil
}

//设定要操作的表
//param table string 要操作的表
//param fields []string 要操作的字段组
func (this *Database) Set(table string,fields []string){
	this.table = table
	this.fields = fields
}

//获取ID
//param id int64 ID
//return *sql.Row, bool 获取的列

~~~ 尝试直接反射出sql内容
//func (this *Database) GetID(id int64) (*sql.Row, bool) {
	query := "select " + this.GetFieldsToStr(this.fields) + " from `" + this.table + "` where `id` = ?"
	stmt, err := this.db.Prepare(query)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	return row, true
}

//Gets the specified field
func (this *Database) GetField(table string, fields []string, field string,value string) (*sql.Row, error) {
	query := "select " + this.GetFieldsToStr(fields) + " from `" + table + "` where `" + field + "` = ?"
	stmt, err := this.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(value)
	return row, nil
}

//get the list
func (this *Database) GetList(table string, fields []string, page int, max int, sort int, desc bool) (map[int]map[string]interface{}, error) {
	var returnResult map[int]map[string]interface{}
	query := "select " + this.GetFieldsToStr(fields) + " from `" + table + "` " + this.GetPageSortStr(page, max, fields[sort], desc)
	stmt, err := this.db.Prepare(query)
	if err != nil {
		return returnResult, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		return returnResult, err
	}
	return this.GetResultToList(fields, rows)
}

//This is just a placeholder
func (this *Database) Insert(table string, fields []string, values []string) {
	//....
}

//update data
func (this *Database) Update(table string, setField string, setValue string, id int) (int64, error) {
	query := "update `" + table + "` set `" + setField + "` = ? where `id` = ?"
	stmt, err := this.db.Exec(query, setValue, id)
	if err != nil {
		return 0, err
	}
	return stmt.RowsAffected()
}

//Removes the specified id
func (this *Database) Delete(table string, id int64) (int64, error) {
	query := "delete from `" + table + "` where `id` = ?"
	smat, err := this.db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	row, err := smat.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row, nil
}

//Obtains a field string based on the number of fields.
func (this *Database) GetFieldsToStr(fields []string) string {
	return "`" + strings.Join(fields, "`,`") + "`"
}

//Get the sql pagination combo section.
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

//Obtain a list of data according to the field.
func (this *Database) GetResultToList(fields []string, result *sql.Rows) (map[int]map[string]interface{}, error) {
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
		for i := range fields {
			a[fields[i]] = c[i]
		}
		resultArray[j] = a
		j++
	}
	return resultArray, err
}

