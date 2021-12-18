package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	// "github.com/bwesterb/go-ristretto"
	// redis "gopkg.in/redis.v4"
	putils "internal/promark_utils"
)

func main() {
	port := os.Getenv("API_PORT")
	// port := "5001"

	fmt.Println("Starting 'tcp' server on external.promark.com:" + port)
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

	fmt.Println("buffer:" + buffer)
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
		CreateCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "get" {
		// fmt.Println("ok:", string(buffer))
		getCampaignCryptoParamsHandler(conn, request.Data)
	} else {
		putils.SendResponse(conn, "nocommand", "")
		conn.Close()
	}

	conn.Close()
}

func CreateCampaignCryptoParamsHandler(conn net.Conn, requestData string) {
	camId := requestData
	fmt.Println("camId:" + camId)
	createCampaignCryptoParams(camId)

	//get param from db
	cryptoParams, err := GetCampaignCryptoParams(camId)

	//temporary return
	param, err := json.Marshal(&cryptoParams)

	// fmt.Println("wrote to file:", n)

	if err != nil {
		panic(err)
	}

	fmt.Println("result: " + string(param))
	putils.SendResponse(conn, "", string(param))
}

func getCampaignCryptoParamsHandler(conn net.Conn, requestData string) {

}

// Campaign function part
func createCampaignCryptoParams(camId string) (*putils.CampaignCryptoParams, error) {
	client := putils.GetRedisConnection()

	//generate campaign param
	_, err := client.Get(camId).Result()
	var cryptoParams putils.CampaignCryptoParams

	if err != nil {
		fmt.Println(err)

		H := putils.GenerateH()
		hEnc := putils.ConvertPointToString(H)
		fmt.Println("hString:.\n", hEnc)

		cryptoParams = putils.CampaignCryptoParams{
			CamID: camId,
			H:     hEnc,
		}

		jsonParam, err := json.Marshal(cryptoParams)
		if err != nil {
			return nil, err
		}

		//store to redis db
		err = client.Set(camId, jsonParam, 0).Err()
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Println("The campaign already existed.\n")
	}

	return &cryptoParams, nil
}

func GetCampaignCryptoParams(camId string) (*putils.CampaignCryptoParams, error) {
	client := putils.GetRedisConnection()

	val, err := client.Get(camId).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(val)

	var campaign putils.CampaignCryptoParams
	err = json.Unmarshal([]byte(val), &campaign)
	if err != nil {
		println(err)
		return nil, err
	}

	//test the value of H
	// H1 := convertBytesToPoint(campaign.H)
	// fmt.Println("H1 point:", H1)

	return &campaign, nil
}
