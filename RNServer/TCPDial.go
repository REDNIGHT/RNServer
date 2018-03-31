package RNServer

import (
	"net"
	"time"

	"../RNCore"
)

type TCPDial struct {
	RNCore.Node

	ip   string
	conn net.Conn

	Out func(*Name_Conn)
}

type Name_Conn struct {
	Name string
	Conn net.Conn
}

func NewTCPDial(name string, ip string) *TCPDial {
	return &TCPDial{RNCore.NewNode(name + " " + ip), ip, nil, nil}
}

func (this *TCPDial) Run() {
	defer RNCore.CatchPanic(nil)

	for {
		//
		if this.conn == nil {
			conn, err := net.Dial("tcp", this.ip)

			if err != nil {
				this.Error("err != nil  err=" + err.Error())
				this.conn = nil
			} else {
				this.conn = conn
				this.Log("A new Connection  RemoteAddr=" + conn.RemoteAddr().String())
				this.Out(&Name_Conn{this.Name(), conn})
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
		}
	}
}
