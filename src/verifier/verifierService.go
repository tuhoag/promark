package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	// "errors"
	"fmt"
	// "log"
	"net"
	"os"

	"github.com/bwesterb/go-ristretto"
	// redis "gopkg.in/redis.v4"
	putils "internal/promark_utils"
	// "strings"
	// "log"
)

// type PromarkRequest struct {
// 	Command string `json:"command"`
// 	Data    string `json:"data"`
// }

// type PromarkResponse struct {
// 	Error string `json:"error"`
// 	Data  string `json:"data"`
// }

// type CampaignCryptoRequest struct {
// 	CamId string
// }

// type NewVerifierCryptoParamsRequest struct {
// 	CamId string `json:"camId"`
// 	H     string `json:"h"`
// }

// type VerifierCryptoParams struct {
// 	CamId string `json:"camId"`
// 	H     string `json:"h"`
// 	S     string `json:"s"`
// }

// type CampaignCryptoParams struct {
// 	CamID string `json:"camId"`
// 	H     string `json:"h"`
// }

func main() {
	port := os.Getenv("API_PORT")
	name := os.Getenv("CORE_PEER_ID")
	// port := "5002"

	fmt.Println("Starting 'tcp' server on " + name + ":" + port)
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
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

	fmt.Println("request:" + buffer)
	var request putils.PromarkRequest
	err = json.Unmarshal([]byte(buffer), &request)

	if err != nil {
		fmt.Println("ERROR:" + err.Error())
		return
	}

	fmt.Println("command:" + request.Command)
	fmt.Println("requestData:" + request.Data)

	if request.Command == "create" {
		// fmt.Println("ok:", string(command))
		CreateVerifierCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "get" {
		// fmt.Println("ok:", string(buffer))
		getCampaignCryptoParamsHandler(conn, request.Data)
	} else {
		putils.SendResponse(conn, "nocommand", "")
		conn.Close()
	}

	conn.Close()
}

func CreateVerifierCampaignCryptoParamsHandler(conn net.Conn, requestData string) {
	var paramsRequest putils.VerifierCryptoParamsRequest
	err := json.Unmarshal([]byte(requestData), &paramsRequest)

	_, err = SetVerifierCryptoParams(paramsRequest)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	vCryptoParams, err := GetVerifierCryptoParams(paramsRequest.CamId)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Println("Got - vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("Got - vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("Got - vCryptoParams.S:" + vCryptoParams.S)

	param, err := json.Marshal(&vCryptoParams)
	putils.SendResponse(conn, "", string(param))
}

func SetVerifierCryptoParams(paramsRequest putils.VerifierCryptoParamsRequest) (bool, error) {
	client := putils.GetRedisConnection()

	var cryptoParams putils.VerifierCryptoParams

	val, err := client.Get(paramsRequest.CamId).Result()
	err = json.Unmarshal([]byte(val), &cryptoParams)
	var s ristretto.Scalar
	if err != nil {
		// params are not existed
		fmt.Println(err)
		s.Rand()
		sBytes := s.Bytes()
		sEnc := b64.StdEncoding.EncodeToString(sBytes)

		cryptoParams = putils.VerifierCryptoParams{
			CamId: paramsRequest.CamId,
			H:     paramsRequest.H,
			S:     sEnc,
		}

		fmt.Println("Store-vCryptoParams.CamId: " + cryptoParams.CamId)
		fmt.Println("Store-vCryptoParams.H: " + cryptoParams.H)
		fmt.Println("Store-vCryptoParams.S: " + cryptoParams.S)

		jsonParam, err := json.Marshal(cryptoParams)
		if err != nil {
			return false, err
		}

		fmt.Println("Store-vCryptoParams JSON: " + string(jsonParam))
		err = client.Set(cryptoParams.CamId, jsonParam, 0).Err()
		if err != nil {
			return false, err
		}

	} else {
		fmt.Printf("The VerifierCryptoParams is existed for id %s", cryptoParams.CamId)
		return false, nil
	}

	return true, nil
}

func GetVerifierCryptoParams(camId string) (*putils.VerifierCryptoParams, error) {
	client := putils.GetRedisConnection()

	var cryptoParams putils.VerifierCryptoParams

	val, err := client.Get(camId).Result()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		// f.WriteString("ERROR: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &cryptoParams)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		// f.WriteString("ERROR: " + err.Error())
		return nil, err
	}

	return &cryptoParams, nil
}

func getCampaignCryptoParamsHandler(conn net.Conn, requestData string) {

}
