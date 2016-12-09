package core

import "net"

//TCP处理器
//用于服务器和客户端双重通讯，根据业务需求使用即可
//依赖外部库：无
//依赖内部库：
// core.SendLog()

//TCP类
type Tcp struct {
	//服务器监听运行状态
	serverStatus bool
	//客户端监听运行状态
	clientStatus bool
	//TCP对象
	tcpAddr *net.TCPAddr
}

//服务器端监听
//为避免影响主程序，请尽量使用并发运行该函数
//param host string 地址及端口
//param handle func(*net.TCPConn) 调用的函数对象
func (this *Tcp) ServerListen(host string,handle func(*net.TCPConn)){
	if this.serverStatus == true{
		return
	}
	this.tcpAddr, err = net.ResolveTCPAddr("tcp", host)
	if err != nil{
		SendLog(err.Error())
		return
	}
	tcpListener, err := net.ListenTCP("tcp", this.tcpAddr)
	if err != nil{
		SendLog(err.Error())
		return
	}
	defer tcpListener.Close()
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil{
			SendLog(err.Error())
			this.serverStatus = false
			continue
		}
		SendLog("A client connects to the TCP server : " + conn.RemoteAddr().String())
		this.serverStatus = true
		go handle(conn)
	}
}

//客户端监听
//为避免影响主程序，请尽量使用并发运行该函数
//param host string 地址及端口
//param handle func(*net.TCPConn) 调用的函数对象
func (this *Tcp) ClientListen(host string,handle func(*net.TCPConn)){
	if this.clientStatus == true{
		return
	}
	this.tcpAddr, err = net.ResolveTCPAddr("tcp", host)
	if err != nil{
		SendLog(err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, this.tcpAddr)
	if err != nil{
		SendLog(err.Error())
		return
	}
	defer conn.Close()
	SendLog("The client TCP connection to the server is successful.")
	go handle(conn)
}