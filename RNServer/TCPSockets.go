package RNServer

import (
	"encoding/binary"
	"io"
	"net"
	"time"
	"unsafe"
)
import "RNCore"

type TCPSockets struct {
	RNCore.Node

	MaxConnNum int

	sockets       map[uintptr]*Socket
	socketsByName map[string]*Socket

	InAddConn            chan net.Conn
	InAddConnWithName    chan *Name_Conn
	InRemoveSocketByName chan string
	InRemoveSocket       chan uintptr

	InSendBuffer       chan *SocketBuffer
	InSendBufferByName chan *SocketBufferByName
	InBroadcast        chan []byte

	//
	outAddSocket    chan<- *Socket
	outRemoveSocket chan<- *Socket

	outSocketsBuffer chan<- *SocketBuffer
}

type Socket struct {
	Name      string
	Conn      net.Conn
	OutBuffer chan<- *SocketBuffer
	InBuffer  chan []byte
}

type SocketBuffer struct {
	SocketID uintptr
	Buffer   []byte
}

type SocketBufferByName struct {
	Name   string
	Buffer []byte
}

func NewTCPSockets(name string, maxConnNum int) *TCPSockets {
	return &TCPSockets{
		Node:          RNCore.NewNode(name),
		MaxConnNum:    maxConnNum,
		sockets:       make(map[uintptr]*Socket),
		socketsByName: make(map[string]*Socket),

		InAddConn:            make(chan net.Conn, RNCore.InChanCount),
		InAddConnWithName:    make(chan *Name_Conn, RNCore.InChanCount),
		InRemoveSocketByName: make(chan string, RNCore.InChanCount),
		InRemoveSocket:       make(chan uintptr, RNCore.InChanCount),

		InSendBuffer:       make(chan *SocketBuffer, RNCore.InChanCount),
		InSendBufferByName: make(chan *SocketBufferByName, RNCore.InChanCount),
		InBroadcast:        make(chan []byte, RNCore.InChanCount)}
}

func (this *TCPSockets) SetOut(
	outAddSocket chan<- *Socket,
	outRemoveSocket chan<- *Socket,
	outSocketsBuffer chan *SocketBuffer,
	outNodeInfos ...string) {

	this.outAddSocket = outAddSocket
	this.outRemoveSocket = outRemoveSocket

	this.outSocketsBuffer = outSocketsBuffer

	//
	this.OutNodeInfos = outNodeInfos

	this.OutNodeInfos[0] = "outAddSocket." + this.OutNodeInfos[0]
	this.OutNodeInfos[1] = "outRemoveSocket." + this.OutNodeInfos[1]
	this.OutNodeInfos[2] = "outSocketsBuffer." + this.OutNodeInfos[2]
}

func (this *TCPSockets) Run() {
	//
	for {
		select {
		case conn := <-this.InAddConn:
			if len(this.sockets) >= this.MaxConnNum {
				conn.Close()
				continue
			}
			this.addSocket(conn, "")

		case c_n := <-this.InAddConnWithName:
			if len(this.sockets) >= this.MaxConnNum {
				c_n.Conn.Close()
				continue
			}
			this.addSocket(c_n.Conn, c_n.Name)

		case sendBuffer := <-this.InSendBuffer:
			if socket, have := this.sockets[sendBuffer.SocketID]; have == true {
				socket.InBuffer <- sendBuffer.Buffer
			} else {
				this.Error("this.sockets have = false  SocketID=%v", sendBuffer.SocketID)
			}

		case name_buffer := <-this.InSendBufferByName:
			if socket, have := this.socketsByName[name_buffer.Name]; have == true {
				socket.InBuffer <- name_buffer.Buffer
			} else {
				this.Error("this.socketsByName have = false  name=" + name_buffer.Name)
			}

		case buffer := <-this.InBroadcast:
			for _, s := range this.sockets {
				s.InBuffer <- buffer
			}

		case s := <-this.InRemoveSocket:
			this.removeSocket(s)

		case name := <-this.InRemoveSocketByName:
			this.removeSocketByName(name)

			//
		case <-time.After(time.Second * RNCore.StateTime):
			this.State()
			continue

		case <-this.CloseSig:
			for _, s := range this.sockets {
				this.removeSocket(uintptr(unsafe.Pointer(s)))
			}
			this.CloseSig <- true
			return
		}
	}
}

func (this *TCPSockets) addSocket(conn net.Conn, name string) {
	socket := &Socket{name, conn, this.outSocketsBuffer, make(chan []byte)}

	this.sockets[uintptr(unsafe.Pointer(socket))] = socket

	if name != "" {
		this.socketsByName[name] = socket
	}

	if this.outAddSocket != nil {
		this.outAddSocket <- socket
	}

	go this.readConnection(socket)
	go this.writeConnection(socket)
}

func (this *TCPSockets) removeSocketByName(name string) {
	if socket, have := this.socketsByName[name]; have == true {
		this.removeSocket(uintptr(unsafe.Pointer(socket)))
	} else {
		this.Error("this.socketsByName have = false  name=" + name)
	}
}

func (this *TCPSockets) removeSocket(socketId uintptr) {
	if socket, have := this.sockets[socketId]; have == true {
		delete(this.sockets, socketId)

		if _, have2 := this.socketsByName[socket.Name]; have2 == true {
			delete(this.socketsByName, socket.Name)
		} else {
			this.Error("this.socketsByName have = false  RemoteAddr=" + socket.Conn.RemoteAddr().String())
		}

		socket.Conn.Close()
		this.Log("Now, %d connections is alve.\n", len(this.sockets))

		if this.outRemoveSocket != nil {
			this.outRemoveSocket <- socket
		}

	} else {
		this.Error("this.sockets have = false  RemoteAddr=" + socket.Conn.RemoteAddr().String())
	}
}

func (this *TCPSockets) readConnection(socket *Socket) {
	for {
		buffer, err := read(socket.Conn)
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			break
		}

		socket.OutBuffer <- &SocketBuffer{uintptr(unsafe.Pointer(socket)), buffer}
	}

	this.InRemoveSocket <- uintptr(unsafe.Pointer(socket))
}

func (this *TCPSockets) writeConnection(socket *Socket) {
	for b := range socket.InBuffer {
		if b == nil {
			break
		}

		err := write(socket.Conn, b)
		if err != nil {
			break
		}
	}

	this.InRemoveSocket <- uintptr(unsafe.Pointer(socket))
}

//
func read(conn net.Conn) (buffer []byte, err error) {

	buffer = make([]byte, 2)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}
	bufferLength := uint32(binary.BigEndian.Uint16(buffer))

	// data
	buffer = make([]byte, bufferLength)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func write(conn net.Conn, buffer []byte) error {

	_buffer := make([]byte, len(buffer)+2)

	binary.BigEndian.PutUint16(_buffer, uint16(len(buffer)))
	copy(_buffer[2:], buffer)

	_, e := conn.Write(_buffer)
	return e
}

//
type _TCPSocketsStateInfo struct {
	RNCore.StateInfo

	MaxConnNum int

	sockets       int
	socketsByName int

	InAddConn            int
	InAddConnWithName    int
	InRemoveSocketByName int
	InRemoveSocket       int

	InSendBuffer       int
	InSendBufferByName int
	InBroadcast        int
}

func (this *TCPSockets) OnState() RNCore.IStateInfo {
	return &_TCPSocketsStateInfo{
		RNCore.StateInfo{this},

		this.MaxConnNum,

		len(this.sockets),
		len(this.socketsByName),

		len(this.InAddConn),
		len(this.InAddConnWithName),
		len(this.InRemoveSocketByName),
		len(this.InRemoveSocket),

		len(this.InSendBuffer),
		len(this.InSendBufferByName),
		len(this.InBroadcast)}
}
