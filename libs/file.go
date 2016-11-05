//文件操作模块
package libs

import (
	"io/ioutil"
	"os"
	"bytes"
	"crypto/sha1"
)

//文件类结构
type FileOperate struct{
}

//创建新的文件夹
//支持多级创建
func (File *FileOperate) CreateDir(src string) bool {
	err := os.MkdirAll(src,os.ModePerm)
	if err != nil{
		return false
	}
	return true
}

//创建文件
func (File *FileOperate) CreateFile(src string) bool{
	_,err := os.Create(src)
	if err != nil{
		return false
	}
	return true
}

//读取文件
func (File *FileOperate) ReadFile(src string) []byte{
	fd, err := os.Open(src)
	if err != nil {
		return nil
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil
	}
	return c
}

//写入文件
func (File *FileOperate) WriteFile(src string, content []byte) bool {
	err := ioutil.WriteFile(src, content, os.ModeAppend)
	if err != nil {
		return false
	}
	return true
}

//追加写入文件
func (File *FileOperate) WriteFileAppend(src string, content []byte) bool{
	if File.IsFile(src) == false{
		File.WriteFile(src, content)
		return true
	}
	var fileContent []byte = File.ReadFile(src)
	s := [][]byte{
		fileContent,
		content,
	}
	sep := []byte("")
	var newContent []byte = bytes.Join(s,sep)
	File.WriteFile(src,newContent)
	return true
}

//修改文件或文件夹名称
//可用于修改路径，即剪切
func (File *FileOperate) EditFileName(src string, newName string) bool {
	err := os.Rename(src, newName)
	if err != nil {
		return true
	}
	return false
}

//删除文件
func (File *FileOperate) DeleteFile(src string) bool {
	err := os.RemoveAll(src)
	if err != nil {
		return true
	}
	return false
}

//判断路径是否存在
func (File *FileOperate) IsExist(src string) bool{
	_, err := os.Stat(src)
	return err == nil || os.IsExist(err)
}

//判断是否为文件
func (File *FileOperate) IsFile(src string) bool {
	info, err := os.Stat(src)
	if err != nil{
		return false
	}
	return !info.IsDir()
}

//判断是否为文件夹
func (File *FileOperate) IsFolder(src string) bool {
	info, err := os.Stat(src)
	if err != nil{
		return false
	}
	return info.IsDir()
}

//获取文件列表
func (File *FileOperate) GetFileList(src string) string {
	return ""
}

//获取文件大小
func (File *FileOperate) GetFileSize(src string) int64 {
	info, err := os.Stat(src)
	if err != nil{
		return 0
	}
	return info.Size()
}

//获取文件信息
func (File *FileOperate) GetFileInfo(src string) (os.FileInfo ,error) {
	info, err := os.Stat(src)
	return info,err
}

//计算文件sha1值
func (File *FileOperate) GetFileSha1(src string) string{
	content := File.ReadFile(src)
	if content != nil{
		sha := sha1.New()
		sha.Write(content)
		res := sha.Sum(nil)
		return string(res)
	}
	return ""
}