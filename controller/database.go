package controller

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

//Database operation module
type Database struct {
	db        *sql.DB
	isConnect bool
}

//Connect to the database
func (this *Database) Connect(dbType string, dbDNS string) error {
	if this.isConnect == true {
		return nil
	}
	this.db, err = sql.Open(dbType, dbDNS)
	if err != nil {
		this.isConnect = true
	}
	return err
}

//Close the database connection
func (this *Database) Close() error {
	return this.db.Close()
}

//Gets the specified ID
func (this *Database) GetID(table string, fields []string, id int) (*sql.Row, error) {
	query := "select " + this.GetFieldsToStr(fields) + " from `" + table + "` where `id` = ?"
	stmt, err := this.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	return row, nil
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
		descStr = "ase"
	}
	return "limit " + strconv.Itoa(page) + "," + strconv.Itoa(max) + " order by `" + sort + "` " + descStr
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
