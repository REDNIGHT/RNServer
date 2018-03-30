package RNServer

import (
	"encoding/binary"
	"io"
	"net"
	//"time"
	"math"
	"unsafe"

	"../RNCore"
)

type SocketBuffer struct {
	SocketId uintptr
	Buffer   []byte
}
type TCPSockets struct {
	RNCore.MNode

	MaxSocketCount int

	sockets       map[uintptr]*Socket
	socketsByName map[string]*Socket

	Out func(*SocketBuffer)
}

type Socket struct {
	Name      string
	Conn      net.Conn
	OutBuffer func(*SocketBuffer)
	InBuffer  chan []byte
}

func NewTCPSockets(name string, maxSocketCount int) *TCPSockets {
	return &TCPSockets{
		MNode:          RNCore.NewMNode(name),
		MaxSocketCount: maxSocketCount,
		sockets:        make(map[uintptr]*Socket),
		socketsByName:  make(map[string]*Socket)}
}

func (this *TCPSockets) AddSocket(conn net.Conn, name string) {
	this.InCall() <- func(_ RNCore.IMessage) {

		if len(this.sockets) >= this.MaxSocketCount {
			conn.Close()
			return
		}

		//
		socket := &Socket{name, conn, this.Out, make(chan []byte)}

		this.sockets[uintptr(socket)] = socket

		if name != "" {
			this.socketsByName[name] = socket
		}

		go this.readConnection(socket)
		go this.writeConnection(socket)
	}
}

func (this *TCPSockets) RemoveSocketByName(name string) {
	this.InCall() <- func(_ RNCore.IMessage) {
		if socket, have := this.socketsByName[name]; have == true {
			this.removeSocket(uintptr(unsafe.Pointer(socket)))
		} else {
			this.Error("this.socketsByName have = false  name=" + name)
		}
	}
}

func (this *TCPSockets) RemoveSocket(socketId uintptr) {
	this.InCall() <- func(_ RNCore.IMessage) {
		this.removeSocket(socketId)
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

	} else {
		this.Error("this.sockets have = false  RemoteAddr=" + socket.Conn.RemoteAddr().String())
	}
}

//
func (this *TCPSockets) SendBuffer(socketId uintptr, buffer []byte) {
	if socket, have := this.sockets[socketId]; have == true {
		socket.InBuffer <- buffer
	} else {
		this.Error("this.sockets have = false  SocketID=%v", socketId)
	}
}

func (this *TCPSockets) SendBufferByName(socketName string, buffer []byte) {
	if socket, have := this.socketsByName[socketName]; have == true {
		socket.InBuffer <- buffer
	} else {
		this.Error("this.socketsByName have = false  socketName=" + socketName)
	}
}

func (this *TCPSockets) Broadcast(buffer []byte) {
	for _, s := range this.sockets {
		s.InBuffer <- buffer
	}
}

//
func (this *TCPSockets) readConnection(socket *Socket) {
	for {
		buffer, err := this.read(socket.Conn)
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			break
		}

		socket.OutBuffer(&SocketBuffer{uintptr(unsafe.Pointer(socket)), buffer})
	}

	this.RemoveSocket(uintptr(unsafe.Pointer(socket)))
}

func (this *TCPSockets) writeConnection(socket *Socket) {
	for b := range socket.InBuffer {
		if b == nil {
			break
		}

		err := this.write(socket.Conn, b)
		if err != nil {
			break
		}
	}

	this.RemoveSocket(uintptr(unsafe.Pointer(socket)))
}

//
func (this *TCPSockets) read(conn net.Conn) (buffer []byte, err error) {

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

func (this *TCPSockets) write(conn net.Conn, buffer []byte) error {
	buferLen := len(buffer) + 2
	if buferLen >= math.MaxUint16 {
		this.Error("buferLen >= math.MaxUint16  buferLen=%v", buferLen)
		return nil
	}
	_buffer := make([]byte, buferLen)

	binary.BigEndian.PutUint16(_buffer, uint16(len(buffer)))
	copy(_buffer[2:], buffer)

	_, e := conn.Write(_buffer)
	return e
}

//
type _TCPSocketsStateInfo struct {
	RNCore.StateInfo

	MaxSocketCount     uint
	socketsCount       uint
	socketsByNameCount uint

	InCount uint
}

func (this *TCPSockets) GetStateInfo() *RNCore.StateInfo {
	si := this.MNode.GetStateInfo()

	si.Values = map[string]uint{
		"maxSocketCount":     uint(this.MaxSocketCount),
		"socketCount":        uint(len(this.sockets)),
		"socketsByNameCount": uint(len(this.socketsByName))}
	return si
}
