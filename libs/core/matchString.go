package core

import "regexp"

//验证和查询模块
type MatchString struct {

}

//验证用户名4-15位
func (ms *MatchString) CheckUsername (str string) bool{
	return ms.matchStr("^[a-zA-Z][a-zA-Z0-9_]{4,15}$",str)
}

//验证Email位
func (ms *MatchString) CheckEmail (str string) bool{
	return ms.matchStr("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+$",str)
}

//验证6-20位的密码
func (ms *MatchString) CheckPassword (str string) bool{
	return ms.matchStr("^[a-zA-Z0-9]{5,20}$",str)
}

//获取正则表达式结果
func (ms *MatchString) matchStr(str string,mSrc string) bool{
	res,err := regexp.MatchString(mSrc,str)
	if err != nil || res == false{
		return false
	}
	return true
}