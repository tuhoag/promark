package promark_utils

import (
	"bufio"
	// b64 "encoding/base64"
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	// "math/big"
	"net"
	// ristretto "github.com/bwesterb/go-ristretto"
	// eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

func RequestClearVerifierData(url string, camId string) error {
	// putils.SendLog("RequestCommitment at", url, LOG_MODE)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// putils.SendLog("Error connecting:", err.Error(), LOG_MODE)

		fmt.Println("Error connecting:" + err.Error())
		return errors.New("ERROR:" + err.Error())
	}

	SendRequest(conn, "clear", camId)
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error clearning verifier data:", err.Error())
		return errors.New("Error clearning verifier data:" + err.Error())
	}

	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)
	// SendLog("Reiceived From: "+url+"-Response:", responseStr, LOG_MODE)

	_, err = ParseResponse(responseStr)

	if err != nil {
		return err
	}

	return nil
}
