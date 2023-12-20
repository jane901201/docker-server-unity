package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	fmt.Println("Server create.")
	return server
}

func (this *Server) Start() {
	fmt.Println("Start IP ", this.Ip)
	fmt.Println("Port ", this.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err:", err)

		return
	}

	defer listener.Close()

	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		fmt.Println("Start lisetner")

		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		fmt.Println("Have conn? ", conn)
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)

	user.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			len, err := conn.Read(buf)
			if 0 == len {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
			}

			user.DOMessage(buf, len)
		}
	}()
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()

		for _, user := range this.OnlineMap {
			user.Channel <- msg
		}

		this.mapLock.Unlock()
	}
}
