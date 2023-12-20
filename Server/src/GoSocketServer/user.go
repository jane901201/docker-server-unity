package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
	level   int
}

func NewUser(conn net.Conn, server *Server) *User {
	fmt.Println("Create new user. ")
	userAddr := conn.RemoteAddr().String()
	fmt.Println(userAddr)

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
		level:   0,
	}

	CreateNewSqlUser(*user)
	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
}

func (this *User) DOMessage(buf []byte, len int) {
	msg := string(buf[:len-1])
	this.server.BroadCast(this, msg)
}

// 接受訊息
func (this *User) ListenMessage() {
	for {
		msg := <-this.Channel
		fmt.Println("Send msg to client: ", msg, ", len: ", int16(len(msg)))

		re := regexp.MustCompile(`:(\d+)$`)

		matches := re.FindStringSubmatch(msg)

		// 抓取 Level 的數字
		if len(matches) >= 2 {
			result := matches[1]
			fmt.Println("Level:", result)
			level, err := strconv.ParseInt(result, 10, 64)

			if err == nil {
				this.level = int(level)
				UpdateSqlUser(*this)
			}

		} else {
			fmt.Println("Don't find level.")
		}

		bytebuf := bytes.NewBuffer([]byte{})
		binary.Write(bytebuf, binary.BigEndian, int16(len(msg)))
		binary.Write(bytebuf, binary.BigEndian, []byte(msg))
		this.conn.Write(bytebuf.Bytes())
	}
}
