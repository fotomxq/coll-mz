package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
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
func (this *MatchString) GetSha1(content string) string {
	hasher := sha1.New()
	_, err = hasher.Write([]byte(content))
	if err != nil {
		return ""
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha)
}

//Gets the regular expression result
func (this *MatchString) matchStr(str string, mSrc string) bool {
	res, err := regexp.MatchString(mSrc, str)
	return res == true && err == nil
}
