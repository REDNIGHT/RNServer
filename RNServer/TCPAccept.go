package RNServer

import (
	"net"
)
import "RNCore"

type TCPAccept struct {
	RNCore.Node

	ip       string
	listener net.Listener

	outAddConn chan<- net.Conn
}

func NewTCPAccept(name string, ip string) *TCPAccept {
	return &TCPAccept{Node: RNCore.NewNode(name), ip: ip}
}

func (this *TCPAccept) SetOut(outAddConn chan<- net.Conn, node_chan_name string) {
	this.outAddConn = outAddConn

	//
	this.SetOutNodeInfos("outAddConn", node_chan_name)

}

func (this *TCPAccept) Run() {
	l, err := net.Listen("tcp", this.ip)
	if err != nil {
		this.Error("err != nil  err=" + err.Error())
		return
	}
	this.listener = l

	//
	this.State()

	//
	for {
		//l.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))

		conn, err := l.Accept()
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			continue
		}
		this.Log("A new Connection  RemoteAddr=" + conn.RemoteAddr().String())

		//
		this.outAddConn <- conn

		//
		select {
		case <-this.CloseSig:
			this.CloseSig <- true
			return
		default:
		}
	}
}

func (this *TCPAccept) Close() {
	this.listener.Close()
	this.Node.Close()
}

//
type _TCPAcceptStateInfo struct {
	RNCore.StateInfo
	Ip string
}

func (this *TCPAccept) OnState() RNCore.IStateInfo {
	return &_TCPAcceptStateInfo{RNCore.StateInfo{this}, this.ip}
}
