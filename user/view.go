package user

import "gopkg.in/mgo.v2/bson"

//该文件定义查询用户数据

//根据ID查询用户信息
//param id string 用户ID
//return *UserFields,bool 用户信息组，是否成功
func GetID(id string) (*UserFields,bool){
	//初始化变量
	var result UserFields
	//获取数据
	var err error
	err = dbColl.FindId(bson.M{"_id":bson.ObjectIdHex(id)}).One(&result)
	if err != nil{
		return &result,false
	}
	//返回数据
	return &result,true
}

//查询用户列表
//param search string 搜索昵称或用户名
//param page int 页数
//param max int 页码
//param sort int 排序字段键值
//param desc bool 是否倒序
//return []UserFields,bool 数据结果，是否成功
func GetList(search string,page int,max int,sort int,desc bool) (*[]UserFields,bool){
	//分析sortStr
	var sortStr string
	if fields[sort] != ""{
		sortStr = fields[sort]
	}else{
		sortStr = fields[0]
	}
	//分析desc
	if desc == true{
		sortStr = "-" + sortStr
	}
	//获取数据
	var result []UserFields
	var err error
	var skip int
	skip = (page - 1) * max
	if search == ""{
		err = dbColl.Find(nil).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	}else{
		err = dbColl.Find(bson.M{"$or":bson.M{"nicename":search,"username":search}}).Sort(sortStr).Skip(skip).Limit(max).All(&result)
	}
	if err != nil{
		return &result,false
	}
	//返回结果
	return &result,true
}

