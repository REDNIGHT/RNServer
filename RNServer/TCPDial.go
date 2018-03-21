package RNServer

import (
	"net"
	"time"
)
import "RNCore"

type TCPDial struct {
	RNCore.Node

	ip   string
	conn net.Conn

	outAddConn chan<- *Name_Conn
}

type Name_Conn struct {
	Name string
	Conn net.Conn
}

func NewTCPDial(name string, ip string) *TCPDial {
	return &TCPDial{Node: RNCore.NewNode(name), ip: ip}
}

func (this *TCPDial) SetOut(outAddConn chan<- *Name_Conn, node_chan_name string) {
	this.outAddConn = outAddConn

	//
	this.SetOutNodeInfos("outAddConn", node_chan_name)
}

func (this *TCPDial) Run() {

	//
	for {
		if this.conn == nil {
			conn, err := net.Dial("tcp", this.ip)

			if err != nil {
				this.Error("err != nil  err=" + err.Error())
				this.conn = nil
			} else {
				this.conn = conn
				this.Log("A new Connection  RemoteAddr=" + conn.RemoteAddr().String())
				this.outAddConn <- &Name_Conn{this.Name(), conn}
			}
		}

		//todo...
		//链接异常断开时 从新链接
		var delta <-chan time.Time
		if this.conn == nil {
			delta = time.After(time.Second * 2)
		} else {
			delta = time.After(time.Second * 30)
		}

		//
		select {
		case <-delta:
			continue

		case <-this.StateSig:
			this.OnState()
			this.StateSig <- true
			continue

		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

//
type _TCPDialStateInfo struct {
	RNCore.StateInfo
	Ip string
}

func (this *TCPDial) OnState() RNCore.IStateInfo {
	return &_TCPDialStateInfo{RNCore.StateInfo{this}, this.ip}
}