package main

import (
    "fmt"
	"net"
	"bufio"
)

const (
    connHost = "0.0.0.0"
    connPort = "3002"
    connType = "tcp"
)

func main() {
    fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)

	conn, err :=net.Dial(connType, connHost + ":" + connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	fmt.Fprintf(conn, "hello\n")

	message, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Error receiving message:", err.Error())
		return
	}

	fmt.Println("message: " + message)
}