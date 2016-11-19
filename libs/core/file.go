//文件操作模块
//该包直接调用函数即可
package core

import (
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"
)

//文件操作类
type FileOperate struct {
}

//创建新的文件夹
//支持多级创建
func (f *FileOperate) CreateDir(src string) (bool, error) {
	b := f.IsFolder(src)
	if b == true {
		return true, nil
	}
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

//创建文件
func (f *FileOperate) CreateFile(src string) bool {
	_, err := os.Create(src)
	if err != nil {
		return true
	}
	return false
}

//读取文件
func (f *FileOperate) ReadFile(src string) ([]byte, error) {
	fd, fdErr := os.Open(src)
	if fdErr != nil {
		return nil, fdErr
	}
	defer fd.Close()
	c, cErr := ioutil.ReadAll(fd)
	if cErr != nil {
		return nil, cErr
	}
	return c, nil
}

//写入文件
func (f *FileOperate) WriteFile(src string, content []byte) (bool, error) {
	err := ioutil.WriteFile(src, content, os.ModeAppend)
	if err != nil {
		return false, err
	}
	return true, nil
}

//追加写入文件
func (f *FileOperate) WriteFileAppend(src string, content []byte) (bool, error) {
	if f.IsFile(src) == false {
		writeBool, writeErr := f.WriteFile(src, content)
		return writeBool, writeErr
	}
	fileContent, fcErr := f.ReadFile(src)
	if fcErr != nil {
		return false, fcErr
	}
	s := [][]byte{
		fileContent,
		content,
	}
	sep := []byte("")
	var newContent []byte = bytes.Join(s, sep)
	writeBool2, writeErr2 := f.WriteFile(src, newContent)
	return writeBool2, writeErr2
}

//向前追加写入文件
func (f *FileOperate) WriteFileForward(src string, content []byte) (bool, error) {
	if f.IsFile(src) == false {
		writeBool, writeErr := f.WriteFile(src, content)
		return writeBool, writeErr
	}
	fileContent, fcErr := f.ReadFile(src)
	if fcErr != nil {
		return false, fcErr
	}
	s := [][]byte{
		content,
		fileContent,
	}
	sep := []byte("")
	var newContent []byte = bytes.Join(s, sep)
	writeBool2, writeErr2 := f.WriteFile(src, newContent)
	return writeBool2, writeErr2
}

//修改文件或文件夹名称
//可用于修改路径，即剪切
func (f *FileOperate) EditFileName(src string, newName string) (bool, error) {
	err := os.Rename(src, newName)
	if err != nil {
		return true, err
	}
	return false, nil
}

//复制文件
func (f *FileOperate) CopyFile(src string, dest string) (bool, error) {
	srcF, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer srcF.Close()
	destF, err := os.Create(dest)
	if err != nil {
		return false, err
	}
	defer destF.Close()
	_, err = io.Copy(destF, srcF)
	if err != nil {
		return false, err
	}
	return true, err
}

//删除文件
func (f *FileOperate) DeleteFile(src string) bool {
	err := os.RemoveAll(src)
	if err != nil {
		return true
	}
	return false
}

//判断路径是否存在
func (f *FileOperate) IsExist(src string) bool {
	_, err := os.Stat(src)
	return err == nil || os.IsExist(err)
}

//判断是否为文件
func (f *FileOperate) IsFile(src string) bool {
	info, err := os.Stat(src)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//判断是否为文件夹
func (f *FileOperate) IsFolder(src string) bool {
	info, err := os.Stat(src)
	if err != nil {
		return false
	}
	return info.IsDir()
}

//获取文件夹下文件和目录列表
func (f *FileOperate) GetFileList(src string) ([]string, error) {
	var fs []string
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}
	sep := f.GetPathSep()
	for _, v := range dir {
		fs = append(fs, src+sep+v.Name())
	}
	return fs, nil
}

//获取文件夹下文件和目录个数
func (f *FileOperate) GetFileListCount(src string) (int, error) {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return 0, err
	}
	var res int
	for range dir {
		res += 1
	}
	return res, nil
}

//获取系统路径分隔符
func (f *FileOperate) GetPathSep() string {
	return string(os.PathSeparator)
}

//获取文件大小
func (f *FileOperate) GetFileSize(src string) int64 {
	info, err := os.Stat(src)
	if err != nil {
		return 0
	}
	return info.Size()
}

//获取文件信息
func (f *FileOperate) GetFileInfo(src string) (os.FileInfo, error) {
	info, err := os.Stat(src)
	return info, err
}

//计算文件sha1值
func (f *FileOperate) GetFileSha1(src string) (string, error) {
	content, err := f.ReadFile(src)
	if err != nil {
		return "", err
	}
	if content != nil {
		sha := sha1.New()
		sha.Write(content)
		res := sha.Sum(nil)
		return string(res), nil
	}
	return "", nil
}
