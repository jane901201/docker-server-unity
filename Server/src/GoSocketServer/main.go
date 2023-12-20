package main

import (
	"fmt"
	"time"
)

func main() {
	go StartServer()

	fmt.Println("Go server.")

	for {
		time.Sleep(1 * time.Second)
	}
}

func StartServer() {
	server := NewServer("0.0.0.0", 8888)
	server.Start()
}
