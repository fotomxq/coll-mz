package core

import (
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"bytes"
)

//日志处理模块
//使用方式：
// 可预先配置好LogOperateSrc变量，指定其他存储日志的路径
// 直接调用函数即可使用
//依赖内部模块：无
//依赖外部库：无

//日志存储文件夹路径
//默认存储到程序所在目录的log文件夹下
var LogOperateSrc = "log"

//向控制台发送一个日志
//自动向控制台和日志文件输出信息
//日志同时将存储到"LogOperateSrc/200601/02/2006010215.log"文件内
//param message string 日志信息
func SendLog(message string){
	//加入时间
	var t time.Time
	t = time.Now()
	message = t.Format("2006-01-02 15:04:05.999999999") + " " + message + "\n"
	//向控制台输出日志
	fmt.Println(message)
	//生成日志文件路径
	var logDir string
	logDir = LogOperateSrc + PathSeparator +t.Format("200601") + PathSeparator + t.Format("02")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil{
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " Failed to create the folder : " + logDir)
		return
	}
	var logSrc string
	logSrc = logDir + PathSeparator + t.Format("2006010215") + ".log"
	//向文件输出日志
	var messageByte []byte
	messageByte = []byte(message)
	//如果文件不存在，则直接创建写入日志
	_,err = os.Stat(logSrc)
	if err != nil || os.IsNotExist(err) == true{
		err = ioutil.WriteFile(logSrc, messageByte, os.ModeAppend)
		if err != nil{
			fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		}
		return
	}
	//如果文件存在，则读出文件内容，添加新日志
	var fd *os.File
	fd, err = os.Open(logSrc)
	if err != nil{
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
	defer fd.Close()
	var c []byte
	c,err = ioutil.ReadAll(fd)
	if err != nil {
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
	var s [][]byte
	s = [][]byte{
		c,
		messageByte,
	}
	var sep []byte
	messageByte = bytes.Join(s, sep)
	err = ioutil.WriteFile(logSrc, messageByte, os.ModeAppend)
	if err != nil{
		fmt.Println(t.Format("2006-01-02 15:04:05.999999999") + " " + err.Error())
		return
	}
}