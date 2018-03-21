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
	var inCount uint = 0
	for {
		inCount++

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

		//
		case <-this.StateSig:
			this.OnState(&inCount)
			this.StateSig <- true
			continue

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

	InCount uint
}

func (this *TCPAccept) OnStateInfo(counts ...*uint) RNCore.IStateInfo {
	return &_TCPAcceptStateInfo{RNCore.StateInfo{this}, this.ip, *counts[0]}
}
