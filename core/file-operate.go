package core

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//文件操作模块
//使用方式：直接调用函数即可使用
//依赖内部模块：core.MatchString
//依赖外部库：无

//创建多级文件夹
//param src string 新文件夹路径
//return bool 是否成功
func CreateFolder(src string) bool {
	err = os.MkdirAll(src, os.ModePerm)
	if err != nil {
		SendLog(err.Error())
	}
	return err != nil
}

//读取文件
//param src string 文件路径
//return []byte, bool 文件数据，是否成功
func LoadFile(src string) ([]byte, bool) {
	fd, err := os.Open(src)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	return c, true
}

//写入文件
//param src string 文件路径
//param content []byte 写入内容
//return bool 是否成功
func WriteFile(src string, content []byte) bool {
	err = ioutil.WriteFile(src, content, os.ModeAppend)
	if err != nil {
		SendLog(err.Error())
	}
	return err != nil
}

//追加写入文件
//param src string 文件路径
//param content []byte 写入内容
//param isForward bool 是否追加到最前端
//return bool 是否成功
func WriteFileAppend(src string, content []byte, isForward bool) bool {
	if IsFile(src) == false {
		return WriteFile(src, content)
	}
	c, b := LoadFile(src)
	if b != true {
		return false
	}
	var s [][]byte
	if isForward == false {
		s = [][]byte{
			c,
			content,
		}
	} else {
		s = [][]byte{
			content,
			c,
		}
	}
	sep := []byte("")
	var newContent []byte = bytes.Join(s, sep)
	return WriteFile(src, newContent)
}

//移动文件或文件夹
//param src string 文件路径
//param dest string 新路径
//return bool 是否成功
func MoveF(src string, dest string) bool {
	err = os.Rename(src, dest)
	if err != nil {
		SendLog(err.Error())
	}
	return err != nil
}

//复制文件
//param src string 文件路径
//param dest string 新路径
//return bool 是否成功
func CopyFile(src string, dest string) bool {
	srcF, err := os.Open(src)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	defer srcF.Close()
	destF, err := os.Create(dest)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	defer destF.Close()
	_, err = io.Copy(destF, srcF)
	if err != nil {
		SendLog(err.Error())
		return false
	}
	return true
}

//删除文件或文件夹
//param src string 文件路径
//return bool 是否成功
func DeleteF(src string) bool {
	err = os.RemoveAll(src)
	if err != nil {
		SendLog(err.Error())
	}
	return err != nil
}

//判断文件或文件夹是否存在
//param src string 文件路径
//return bool 是否存在
func IsExist(src string) bool {
	_, err := os.Stat(src)
	return err != nil || os.IsNotExist(err) == false
}

//判断是否为文件
//param src string 文件路径
//return bool 是否为文件
func IsFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

//判断是否为文件夹
//param src string 文件夹路径
//return bool 是否为文件夹
func IsFolder(src string) bool {
	info, err := os.Stat(src)
	return err == nil && info.IsDir()
}

//获取文件列表
//param src string 查询的文件夹路径
//param filtre string 仅保留的文件，文件夹除外，eg : jpg|jpeg|gif
//param isSrc bool 返回是否为文件路径
//return []string,bool 文件列表，是否成功
func GetFileList(src string, filter string, isSrc bool) ([]string, bool) {
	var fs []string
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	var filters []string
	if filter != "" {
		filters = strings.Split(filter, "|")
	}
	for _, v := range dir {
		var appendSrc string
		if isSrc == true {
			appendSrc = src + PathSeparator + v.Name()
		} else {
			appendSrc = v.Name()
		}
		if v.IsDir() == true || filter == "" {
			fs = append(fs, appendSrc)
			continue
		}
		names := strings.Split(v.Name(), ".")
		if len(names) == 1 {
			fs = append(fs, appendSrc)
			continue
		}
		t := names[len(names)-1]
		for _, filterValue := range filters {
			if t != filterValue {
				continue
			}
			fs = append(fs, appendSrc)
		}
	}
	return fs, true
}

//查询文件夹下文件个数
//param src string 文件夹路径
//return int,bool 文件个数，是否成功
func GetFileListCount(src string) (int, bool) {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		SendLog(err.Error())
		return 0, false
	}
	var res int
	for range dir {
		res += 1
	}
	return res, true
}

//获取系统路径分隔符
//return string 分隔符
func GetPathSep() string {
	return string(os.PathSeparator)
}

//获取文件大小
//param src string 文件路径
//return int64,bool 文件大小，是否成功
func GetFileSize(src string) (int64, bool) {
	info, err := os.Stat(src)
	if err != nil {
		SendLog(err.Error())
		return 0, false
	}
	return info.Size(), false
}

//获取文件名称分割序列
//param src string 文件路径
//return map[string]string,bool 文件名称序列，是否成功 eg : {"name","abc.jpg","type":"jpg","only-name":"abc"}
func GetFileNames(src string) (map[string]string, bool) {
	info, err := os.Stat(src)
	if err != nil {
		SendLog(err.Error())
		return nil, false
	}
	res := map[string]string{
		"name":      info.Name(),
		"type":      "",
		"only-name": info.Name(),
	}
	names := strings.Split(res["name"], ".")
	if len(names) < 2 {
		return res, true
	}
	res["type"] = names[len(names)-1]
	res["only-name"] = names[0]
	for i := range names {
		if i != 0 && i < len(names)-1 {
			res["only-name"] = res["only-name"] + "." + names[i]
		}
	}
	return res, true
}

//获取文件信息
//param src string 文件路径
//return os.FileInfo,bool 文件信息，是否成功
func GetFileInfo(src string) (os.FileInfo, bool) {
	var c os.FileInfo
	c, err = os.Stat(src)
	if err != nil {
		SendLog(err.Error())
	}
	return c, err != nil
}

//获取文件SHA1值
//param src string 文件路径
//return string SHA1值
func GetFileSha1(src string) string {
	content, b := LoadFile(src)
	if b != true {
		return ""
	}
	if content != nil {
		sha := sha1.New()
		_, err = sha.Write(content)
		if err != nil {
			SendLog(err.Error())
			return ""
		}
		res := sha.Sum(nil)
		return hex.EncodeToString(res)
	}
	return ""
}

//Gets the directory path for the time build
//eg : Return and create the path ,"[src]/201611/"
//eg : Return and create the path ,"[src]/201611/2016110102-03[appendFileType]"
//获取并创建时间序列创建的多级文件夹
//param src string 文件路径
//param appendFileType string 是否末尾追加文件类型，如果指定值，则返回
//return string 新时间周期目录，失败则返回空字符串
func GetTimeDirSrc(src string, appendFileType string) string {
	sep := GetPathSep()
	newSrc := src + sep + time.Now().Format("200601")
	var b bool
	b = CreateFolder(newSrc)
	if b == false {
		return ""
	}
	newSrc = newSrc + sep
	if appendFileType != "" {
		newSrc = newSrc + time.Now().Format("20060102-03") + appendFileType
	}
	return newSrc
}
