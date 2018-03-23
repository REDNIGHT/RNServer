package RNServer

import (
	"net"
)
import "RNCore"

type TCPAccept struct {
	RNCore.MinNode

	ip       string
	listener net.Listener

	Out func(net.Conn)
}

func NewTCPAccept(name string, ip string) *TCPAccept {
	return &TCPAccept{RNCore.NewMinNode(name), ip, nil, nil}
}

func (this *TCPAccept) Run() {
	l, err := net.Listen("tcp", this.ip)
	if err != nil {
		this.Error("err != nil  err=" + err.Error())
		return
	}
	this.listener = l

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
		this.Out(conn)

	}
}

/*func (this *TCPAccept) Close() {
	this.listener.Close()
}*/
