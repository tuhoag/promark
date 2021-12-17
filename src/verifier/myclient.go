package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	// "encoding/json"
)

const (
	connHost               = "external.promark.com"
	connPort               = "5000"
	connType               = "tcp"
	cryptoServiceSocketURL = "external.promark.com:5000"
)

func main() {
	log.Println("Start running")
	// sendLog("Sending to", cryptoServiceSocketURL)
	conn, err := net.Dial("tcp", cryptoServiceSocketURL)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		log.Println("Error connecting:", err.Error())
		return
	}
	log.Println("Send create command")
	fmt.Fprintf(conn, "create c001\n")

	message, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return
	}
	log.Println("Reiceived: " + message)

	// // sendLog("message:", message)
	// if message == "ok\n" {
	// 	log.Println("Send camId: c001")
	// 	fmt.Fprintf(conn, "c001" + "\n")
	// 	// return
	// } else if message == "nocommand\n" {
	// 	log.Printf("nocommand")
	// 	return
	// }

	// message, err = bufio.NewReader(conn).ReadString('\n')

	// // log.Println("Reiceived: " + message)
	// // log.Println("Reiceived err: " + err.Error())
	// if err != nil {
	// 	// sendLog("Error connecting:", err.Error())
	// 	log.Println("Error after final result:", err.Error())
	// 	return
	// }

	// log.Println("Reiceived: " + message)
	// var cryptoParams CampaignCryptoParams
	// err = json.Unmarshal([]byte(message), &cryptoParams)

	// log.Println(cryptoParams)
}
