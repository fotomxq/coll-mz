//日志模块
package libs

//引用模块
import(
	"fmt"
	"time"
)

//日志类结构
type Log struct {
	//新日志发送方式
	// 0 - 向控制台发送信息 ; 1 - 保存日志信息 ; 2 - 全部发送
	newLogType int
	//日志数据保存的路径目录
	//如果不存在则自动建立为“log”目录下
	dirSrc string
	//日志保存结构
	// 0 - 年月/日.log ; 1 - 年/月/日.log ; 2 - 年月日.log
	dirType int
	//file模块调用
	file FileOperate
}

//设定发送方式
func (log *Log) SetNewLogType(num int) {
	if num != 1 && num != 2{
		num = 0
	}
	log.newLogType = num
}

//设定存储路径
func (log *Log) SetDirSrc(src string) {
	if log.file.IsFile(src) == false {
		log.dirSrc = src
	}else{
		log.dirSrc = "log"
	}
}

//设定日志保存结构
func (log *Log) SetDirType(t int){
	if t != 1 && t != 2{
		t = 0
	}
	log.dirType = t
}

//添加新的日志
func (log *Log) AddLog(content string) {
	switch log.newLogType {
		case 1:
			log.postFileLog(content)
			break
		case 2:
			log.postFmtLog(content)
			log.postFileLog(content)
			break
		default:
			log.postFmtLog(content)
			break
	}
}

//系统级别错误日志
func (log *Log) AddErrorLog(err error){
	errMsg := err.Error()
	log.AddLog(errMsg)
}

//向控制台发送日志信息
func (log *Log) postFmtLog(content string) {
	fmt.Println(content)
}

//向日志文件发送日志信息
func (log *Log) postFileLog(content string){
	//检查或初始化变量
	if log.dirSrc == ""{
		log.SetDirSrc("log")
	}
	if log.newLogType != 0 && log.newLogType != 1 && log.newLogType != 2{
		log.SetNewLogType(2)
	}
	if log.dirType != 0 && log.dirType != 1 && log.dirType != 2{
		log.SetDirType(0)
	}
	//构建存储位置
	var dir string
	var logSrc string
	t := time.Now()
	switch log.dirType{
		case 1:
			dir = log.dirSrc + "/" + t.Format("2006/01")
			logSrc = dir + "/" + t.Format("20060102") + ".log"
			break
		case 2:
			dir = log.dirSrc
			logSrc = dir + "/" + t.Format("20060102") + ".log"
			break
		default:
			dir = log.dirSrc + "/" + t.Format("200601")
			logSrc = dir + "/" + t.Format("20060102") + ".log"
			break

	}
	if log.file.CreateDir(dir) == false {
		log.postFmtLog("ERROR : Cannot create log dir.")
		return
	}
	//构建日志
	var nowTime string = log.getNowDateString()
	var logContent string = nowTime + " " + content + "\n"
	logContentByte := []byte(logContent)
	//向日志文件添加日志
	var _ bool = log.file.WriteFileAppend(logSrc,logContentByte)
}

//获取当前系统日期
func (log *Log) getNowDateString() (string){
	return time.Now().String()
}
