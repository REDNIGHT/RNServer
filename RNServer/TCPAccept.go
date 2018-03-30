package RNServer

import (
	"net"

	"../RNCore"
)

type TCPAccept struct {
	RNCore.Node
	listener net.Listener

	Out func(net.Conn)
}

func NewTCPAccept(name, ip string) *TCPAccept {
	this := &TCPAccept{RNCore.NewNode(name + " " + ip), nil, nil}

	l, err := net.Listen("tcp", ip)
	if err != nil {
		this.Panic("err != nil  err=" + err.Error())
	}
	this.listener = l

	return this
}

func (this *TCPAccept) Run() {
	//
	for {
		//l.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
		if this.listener == nil {
			return
		}

		conn, err := this.listener.Accept()
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			continue
		}
		this.Log("A new Connection  RemoteAddr=" + conn.RemoteAddr().String())

		//
		this.Out(conn)
	}
}

func (this *TCPAccept) Close() {
	this.listener.Close()
	this.listener = nil
}
