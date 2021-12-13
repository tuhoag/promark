package main

import (
    "fmt"
	"net"
	"os"
	"bufio"
	"log"
)

type CampaignCryptoRequest struct {
	CamId        string
}

type CampaignCryptoParams struct {
	CamID string `json:"camId"`
	H     string `json:"h"`
}


const (
    connHost = "localhost"
    connPort = "8080"
    connType = "tcp"
)

func main() {
    fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		fmt.Println("Client connected.")
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer))

	if buffer == "create" {
		createCampaignCryptoParamsHandler(conn)
	} else if buffer == "get" {
		getCampaignCryptoParamsHandler(conn)
	} else {
		fmt.Write(conn, "nocommand")
		conn.Close()
	}

	conn.Write(buffer)

    conn.Close()
}

func createCampaignCryptoParamsHandler(conn net.Conn) {

}

func getCampaignCryptoParamsHandler(conn net.Conn) {

}