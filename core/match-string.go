package core

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//验证模块
//用于验证字符串等验证、过滤类操作
//使用方法：声明后直接调用具体方法即可
//依赖外部模块：无
//依赖内部模块：
// core.LogOperate

//验证模块
type MatchString struct {
}

//检查用户名
//param str string 用户名
//return bool 是否正确
func (this *MatchString) CheckUsername(str string) bool {
	return this.matchStr(`^[a-zA-Z0-9_-]{4,16}$`, str)
}

//检查昵称
//param str string 昵称
//return bool 是否正确
func (this *MatchString) CheckNicename(str string) bool {
	return this.matchStr(`^[\u4e00-\u9fa5_a-zA-Z0-9]+$`, str)
}

//验证邮箱
//param str string 邮箱地址
//return bool 是否正确
func (this *MatchString) CheckEmail(str string) bool {
	return this.matchStr(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, str)
}

//验证密码
//param str string 密码
//return bool 是否正确
func (this *MatchString) CheckPassword(str string) bool {
	return this.matchStr(`^[a-zA-Z0-9_-]{4,16}$`, str)
}

//验证搜索类型的字符串
//param str string 字符串
//return bool 是否正确
func (this *MatchString) CheckSearch(str string) bool {
	return this.matchStr(`^[\u4e00-\u9fa5_a-zA-Z0-9]+$`, str)
}

//验证是否为SHA1
//param str string 字符串
//return bool 是否正确
func (this *MatchString) CheckHexSha1(str string) bool {
	return this.matchStr(`^[a-z0-9]{10,45}$`, str)
}

//验证是否为IP地址
//param str string IP地址
//return bool 是否正确
func (this *MatchString) CheckIP(str string) bool {
	if str == "[::1]" {
		return true
	}
	if this.matchStr(`^$(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, str) == true {
		return true
	}
	return false
}

//获取字符串的SHA1值
//param content string 要计算的字符串
//return string 计算出的SHA1值
func (this *MatchString) GetSha1(content string) string {
	hasher := sha1.New()
	_, err = hasher.Write([]byte(content))
	if err != nil {
		Log.SendLog("core/match-string.go", "0.0.0.0", "MatchString.GetSha1", "write", err.Error())
		return ""
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha)
}

//匹配验证
//param mStr string 验证
//param str string 要验证的字符串
//return bool 是否成功
func (this *MatchString) matchStr(mStr string, str string) bool {
	res, err := regexp.MatchString(mStr, str)
	if err != nil {
		Log.SendLog("core/match-string.go", "0.0.0.0", "MatchString.matchStr", "regexp-match-str", err.Error())
		return false
	}
	return res == true
}

//获取随机字符串
//param n int 随机码
//return string 新随机字符串
func (this *MatchString) GetRandStr(n int) string {
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

//分解URL获取名称和类型
//param sendURL URL地址
//return map[string]string 返回值集合
func (this *MatchString) GetURLNameType(sendURL string) map[string]string {
	res := map[string]string{
		"full-name": "",
		"only-name": "",
		"type":      "",
	}
	urls := strings.Split(sendURL, "/")
	if len(urls) < 1 {
		return res
	}
	res["full-name"] = urls[len(urls)-1]
	if res["full-name"] == "" {
		res["only-name"] = res["full-name"]
		return res
	}
	names := strings.Split(res["full-name"], ".")
	if len(names) < 2 {
		return res
	}
	res["type"] = names[len(names)-1]
	for i := 0; i <= len(names); i++ {
		if i == len(names)-1 {
			break
		}
		if res["only-name"] == "" {
			res["only-name"] = names[i]
		} else {
			res["only-name"] += "." + names[i]
		}
	}
	return res
}

//处理page
//param postPage string 用户提交的page
//return int 过滤后的页数
func (this *MatchString) FilterPage(postPage string) int {
	var res int
	res, err = strconv.Atoi(postPage)
	if err != nil {
		res = 1
	}
	if res < 1 {
		res = 1
	}
	return res
}

//处理max
//限制最小值为1，最大值为999
//param postMax string 用户提交的max
//return int 过滤后的页数
func (this *MatchString) FilterMax(postMax string) int {
	var res int
	res, err = strconv.Atoi(postMax)
	if err != nil {
		res = 1
	}
	if res < 1 {
		res = 1
	}
	if res > 999 {
		res = 999
	}
	return res
}

//过滤非法字符
//param str string 要过滤的字符串
//return string 过滤后的字符串
func (this *MatchString) FilterStr(str string) string{
	//str = strings.Replace(str,"\r","",-1)
	//str = strings.Replace(str,"\n","",-1)
	//str = strings.Replace(str,"\t","",-1)
	str = strings.Replace(str,"~","～",-1)
	str = strings.Replace(str,"<","〈",-1)
	str = strings.Replace(str,">","〉",-1)
	str = strings.Replace(str,"$","￥",-1)
	str = strings.Replace(str,"!","！",-1)
	str = strings.Replace(str,"[","【",-1)
	str = strings.Replace(str,"]","】",-1)
	str = strings.Replace(str,"{","｛",-1)
	str = strings.Replace(str,"}","｝",-1)
	return str
}

//过滤非法字符后判断其长度是否符合标准
//param str string 要过滤的字符串
//param min int 最短，包括该长度
//param max int 最长，包括该长度
//return string 过滤后的字符串，失败返回空字符串
func (this *MatchString) CheckFilterStr(str string,min int,max int) string{
	var newStr string
	newStr = this.FilterStr(str)
	if newStr == ""{
		return ""
	}
	var strLen int
	strLen = len(newStr)
	if strLen >= min && strLen <= max{
		return newStr
	}
	return ""
}
