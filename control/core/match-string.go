package core

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"time"
	"math/rand"
	"strconv"
)

//Authentication and query modules
//type MatchString struct {
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

//匹配验证
//param str string 要验证的字符串
//param mStr string 验证
~~~~~~~~~~~~~~~~~~~
func (this *MatchString) matchStr(str string, mStr string) bool {
	res, err := regexp.MatchString(mStr, str)
	return res == true && err == nil
}

//获取随机字符串
//param n int 随机码
//return string 新随机字符串
func (this *MatchString) GetRandStr(n int) string{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	re := r.Intn(n)
	return strconv.Itoa(re)
}

//截取字符串
//param str string 要截取的字符串
//param star int 开始位置
//param length int 长度
//return string 新字符串
func (this *MatchString) SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
