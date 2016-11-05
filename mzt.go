package main

import(
	"ftmplibs"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//日志处理
var log ftmplibs.Log
//配置处理
var config ftmplibs.Config
//数据库
var db *sql.DB
//错误
var err error

//启动脚本
func main(){
	//开始提示
	log.AddLog("* _ * * _ * 脚本开始运行 * _ * * _ *")
	log.AddLog("初始化参数中...")
	//读取配置信息
	err = config.LoadFile("config.json")
	if err != nil{
		log.AddLog("发生一个错误:")
		log.AddErrorLog(err)
	}
	//连接数据库
	db, err = sql.Open(config.Data["databaseType"].(string),config.Data["databaseDNS"].(string))
	if err != nil{
		log.AddLog("发生一个错误:")
		log.AddErrorLog(err)
	}
	defer db.Close()
	//获取秀人数据
	getXRData()
}

func getXRData(){

}