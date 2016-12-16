package core

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//该模块用于通过外网、内部获取IP地址
//使用IPAddrsGetExternal()前请确保联网
//使用方式：
// 直接调用函数即可获取IP地址，失败将返回0.0.0.0
//依赖内部模块：core.LogOperate
//依赖外部库：无

//通过外部网络获取IP地址
//return string IP地址
func IPAddrsGetExternal() string {
	var url string = "http://myexternalip.com/raw"
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		Log.SendLog("core/ip-addrs.go","0.0.0.0","IPAddrsGetExternal","http-get",err.Error())
		return "0.0.0.0"
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.SendLog("core/ip-addrs.go","0.0.0.0","IPAddrsGetExternal","ioutil-read-all",err.Error())
		return "0.0.0.0"
	}
	var html string
	html = string(body)
	html = strings.Replace(html, " ", "", -1)
	html = strings.Replace(html, "\n", "", -1)
	return html
}

//通过内部获取IP地址
//return string IP地址
func IPAddrsGetInternal() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Log.SendLog("core/ip-addrs.go","0.0.0.0","IPAddrsGetInternal","get-interface-addrs",err.Error())
		return "0.0.0.0"
	}
	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}

//通过r获取客户端IP地址
//自动剔除端口部分，只获取IP地址
//param r *http.Request Http读取对象
//return string IP地址
func IPAddrsGetRequest(r *http.Request) string{
	var ipAddr string
	ipAddr = r.RemoteAddr
	if ipAddr == ""{
		return "0.0.0.0"
	}
	var ipAddrs []string
	ipAddrs = strings.Split(ipAddr,":")
	if len(ipAddrs) > 1{
		var res string
		res = strings.Join(ipAddrs[0:len(ipAddrs)-1],":")
		return res
	}
	return ipAddrs[0]
}

///////////////////////////////////////////////////////////////////////////////
// IP黑名单和白名单模块
///////////////////////////////////////////////////////////////////////////////

//初始化后即可使用

//IP黑名单结构
type IPAddrBan struct {
	//是否启动黑名单模式
	banBool bool
	//是否启动白名单模式
	whiteBool bool
	//数据库操作句柄
	db *mgo.Database
}

//IP黑名单数据库结构
type IPAddrBanFields struct {
	//IP地址
	IPAddr string
	//是否拉黑
	// 如果改为白名单模式，则会判断仅为false的放行
	IsBan bool
}

//初始化
//param db *mgo.Database 数据库句柄
//param banBool bool 是否启动黑名单
//param whiteBool bool 是否启动白名单
func (this *IPAddrBan) Init(db *mgo.Database,banBool bool, whiteBool bool){
	this.db = db
	this.banBool = banBool
	this.whiteBool = whiteBool
}

//检查IP地址是否可通行
//param ipAddr string IP地址
//return bool 是否可以通行
func (this *IPAddrBan) CheckList(ipAddr string) bool{
	//如果 白名单和黑名单模式关闭，则返回可通行
	if this.whiteBool == false && this.banBool == false{
		return true
	}
	//获取集合
	var result IPAddrBanFields
	var dbColl *mgo.Collection
	dbColl = this.db.C("ip")
	err = dbColl.Find(bson.M{"ipaddr":ipAddr}).One(&result)
	//如果 没有数据 && 白名单未启动，返回可通行，否则不允许
	if err != nil && this.whiteBool == true{
		Log.SendLog("core/ip-addrs.go",ipAddr,"IPAddrBan.IPAddrsCheck","ip-ban","未在白名单内的IP，尝试访问。")
		return false
	}
	//如果 ban=false，无论是否启动黑白名单，均返回可通行
	if result.IsBan == false{
		return true
	}
	//其他情况返回失败
	Log.SendLog("core/ip-addrs.go",ipAddr,"IPAddrBan.IPAddrsCheck","ip-ban","未在白名单或在黑名单内的IP，尝试访问。")
	return false
}

//将IP地址记录到数据库
//param ipaddr string IP地址
//param isBan bool 是否禁用
//return bool 是否保存成功
func (this *IPAddrBan) SaveToList(ipAddr string,isBan bool) bool{
	//获取集合
	var result IPAddrBanFields
	var dbColl *mgo.Collection
	dbColl = this.db.C("ip")
	err = dbColl.Find(bson.M{"ipaddr":ipAddr}).One(&result)
	if err != nil{
		//不存在则创建
		err = dbColl.Insert(&IPAddrBanFields{ipAddr,isBan})
		if err != nil{
			Log.SendLog("core/ip-addrs.go",ipAddr,"IPAddrBan.IPAddrsSaveBan","create-db",err.Error())
			return false
		}
		return true
	}
	//如果存在，则更新数据
	err = dbColl.Update(bson.M{"ipaddr":ipAddr},bson.M{"$set":bson.M{"isban":isBan}})
	if err != nil{
		Log.SendLog("core/ip-addrs.go",ipAddr,"IPAddrBan.IPAddrsSaveBan","update-db",err.Error())
		return false
	}
	return true
}