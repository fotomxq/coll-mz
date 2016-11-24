package controller

import (
	"regexp"
	"crypto/sha1"
)

//Authentication and query modules
type MatchString struct {
}

//Verify the 4-16-digit user name
func (this *MatchString) CheckUsername(str string) bool {
	return this.matchStr("^[a-zA-Z][a-zA-Z0-9_]{4,15}$", str)
}

//Verify email
func (this *MatchString) CheckEmail(str string) bool {
	return this.matchStr("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+$", str)
}

//Verify the 6-20 digit password
func (this *MatchString) CheckPassword(str string) bool {
	return this.matchStr("^[a-zA-Z0-9]{5,20}$", str)
}

//Gets the string sha1 value
func (this *MatchString) GetSha1(content string) string{
	sha := sha1.New()
	contentByte := []byte(content)
	_,err := sha.Write(contentByte)
	if err != nil{
		return ""
	}
	res := sha.Sum(nil)
	return string(res)
}

//Gets the regular expression result
func (this *MatchString) matchStr(str string, mSrc string) bool {
	res, err := regexp.MatchString(mSrc, str)
	return res == true && err == nil
}
