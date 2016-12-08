package controller

import (
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"encoding/hex"
)

//The module mainly encapsulates complex file operations.
//But recommended, will reduce the operating efficiency..

//Create a folder
func CreateDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

//Load file
func LoadFile(src string) ([]byte, error) {
	fd, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//Write file
//If the file does not exist, it is created.
func WriteFile(src string, content []byte) error {
	return ioutil.WriteFile(src, content, os.ModeAppend)
}

//Append data to a file
func WriteFileAppend(src string, content []byte, isForward bool) error {
	if IsFile(src) == false {
		return WriteFile(src, content)
	}
	c, err := LoadFile(src)
	if err != nil {
		return err
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

//Modify the file name
//Can be used to cut files.
func CutFile(src string, newSrc string) error {
	return os.Rename(src, newSrc)
}

//Copy file
func CopyFile(src string, dest string) (bool, error) {
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

//Delete file
func DeleteFile(src string) error {
	return os.RemoveAll(src)
}

//Determine whether the file exists
func IsExist(src string) bool {
	_, err := os.Stat(src)
	return err == nil || os.IsExist(err)
}

//Determine whether the file
func IsFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

//To determine whether the folder
func IsFolder(src string) bool {
	info, err := os.Stat(src)
	return err == nil && info.IsDir()
}

//Gets a list of files under the folder
// src string
// filter string eg : gif|jpg
// isSrc bool
func GetFileList(src string,filter string,isSrc bool) ([]string, error) {
	var fs []string
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}
	var filters []string
	if filter != ""{
		filters = strings.Split(filter,"|")
	}
	for _, v := range dir {
		var appendSrc string
		if isSrc == true{
			appendSrc = src + GetPathSep() + v.Name()
		}else{
			appendSrc = v.Name()
		}
		if v.IsDir() == true {
			fs = append(fs, appendSrc)
			continue
		}
		if filter == "" {
			fs = append(fs, appendSrc)
			continue
		}
		names := strings.Split(v.Name(),".")
		if len(names) == 1{
			fs = append(fs, appendSrc)
			continue
		}
		t := names[len(names)-1]
		for _,filterValue := range filters{
			if t != filterValue{
				continue
			}
			fs = append(fs, appendSrc)
		}
	}
	return fs, nil
}

//Gets the number of files in the folder
func GetFileListCount(src string) (int, error) {
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

//Gets the system path separator
func GetPathSep() string {
	return string(os.PathSeparator)
}

//Get the file size
func GetFileSize(src string) int64 {
	info, err := os.Stat(src)
	if err != nil {
		return 0
	}
	return info.Size()
}

//Gets the file name and type
func GetFileNames(src string) (map[string]string, error) {
	info, err := GetFileInfo(src)
	if err != nil {
		return nil, err
	}
	res := map[string]string{
		"name":     info.Name(),
		"type":     "",
		"onlyName": info.Name(),
	}
	if res["name"] == "" {
		return res, nil
	}
	names := strings.Split(res["name"], ".")
	if len(names) < 2 {
		return res, nil
	}
	res["type"] = names[len(names)-1]
	res["onlyName"] = names[0]
	for i := range names {
		if i != 0 && i < len(names)-1 {
			res["onlyName"] = res["onlyName"] + "." + names[i]
		}
	}
	return res, nil
}

//Get file information
func GetFileInfo(src string) (os.FileInfo, error) {
	return os.Stat(src)
}

//Calculates the file sha1 value
func GetFileSha1(src string) (string, error) {
	content, err := LoadFile(src)
	if err != nil {
		return "", err
	}
	if content != nil {
		sha := sha1.New()
		sha.Write(content)
		res := sha.Sum(nil)
		return hex.EncodeToString(res), nil
	}
	return "", nil
}

//Gets the directory path for the time build
//eg : Return and create the path ,"[src]/201611/"
//eg : Return and create the path ,"[src]/201611/2016110102-03[appendFileType]"
func GetTimeDirSrc(src string, appendFileType string) (string, error) {
	t := time.Now()
	sep := GetPathSep()
	newSrc := src + sep + t.Format("200601")
	err = CreateDir(newSrc)
	newSrc = newSrc + sep
	if appendFileType != "" {
		newSrc = newSrc + t.Format("20060102-03") + appendFileType
	}
	return newSrc, err
}
