package RNServer

import (
	"net"
)
import "RNCore"

type TCPAccept struct {
	RNCore.Node

	ip       string
	listener net.Listener

	out func(net.Conn)
}

func NewTCPAccept(name string, ip string) *TCPAccept {
	return &TCPAccept{Node: RNCore.NewNode(name), ip: ip}
}

func (this *TCPAccept) Out(out func(net.Conn), node_chan_name string) {
	this.out = out

	//
	this.SetOutNodeInfos("out", node_chan_name)

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
		this.out(conn)

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
}

func (this *TCPAccept) OnStateInfo(counts ...*uint) *RNCore.StateInfo {
	si := RNCore.NewStateInfo(this, *counts[0])
	si.StrValues = map[string]string{"ip": this.ip}
	return si
}

func (this *TCPAccept) DebugChanState() {
	//this.OnDebugChanState("In", len(this.In))
}
