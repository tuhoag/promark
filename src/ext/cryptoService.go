package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
	// "strings"
	// "log"
)

type PromarkRequest struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type PromarkResponse struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

type CampaignCryptoRequest struct {
	CamId string
}

type CampaignCryptoParams struct {
	CamID string `json:"camId"`
	H     string `json:"h"`
}

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
	var request PromarkRequest
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
		SendResponse(conn, "nocommand", "")
		conn.Close()
	}

	conn.Close()
}

func SendData(conn net.Conn, message string) {
	fmt.Fprintf(conn, message+"\n")
}

func SendRequest(conn net.Conn, command string, data string) error {
	request := PromarkRequest{
		Command: command,
		Data:    data,
	}
	requestJSON, err := json.Marshal(&request)

	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return errors.New("ERROR:" + err.Error())
	}

	log.Println("Send create command")
	SendData(conn, string(requestJSON))

	return nil
}

func SendResponse(conn net.Conn, errorStr string, data string) error {
	response := PromarkResponse{
		Error: errorStr,
		Data:  data,
	}
	responseJSON, err := json.Marshal(&response)

	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	SendData(conn, string(responseJSON))

	return nil
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
	SendResponse(conn, "", string(param))
}

func getCampaignCryptoParamsHandler(conn net.Conn, requestData string) {

}

// Campaign function part
func createCampaignCryptoParams(camId string) (*CampaignCryptoParams, error) {
	// var r1, r2 ristretto.Scalar
	// var r ristretto.Scalar
	// var rArr [][]byte

	client := GetRedisConnection()

	//generate campaign param
	_, err := client.Get(camId).Result()
	var cryptoParams CampaignCryptoParams

	if err != nil {
		fmt.Println(err)

		H := generateH()
		hBytes := H.Bytes()
		hEnc := b64.StdEncoding.EncodeToString(hBytes)
		fmt.Println("hString:.\n", hEnc)

		cryptoParams = CampaignCryptoParams{
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

func GetCampaignCryptoParams(camId string) (*CampaignCryptoParams, error) {
	client := GetRedisConnection()

	val, err := client.Get(camId).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(val)

	var campaign CampaignCryptoParams
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

func GetRedisConnection() *redis.Client {
	// pool := redis.ConnectionPool(host="127.0.0.1", port=6379, db=0)
	// client := redis.StrictRedis(connection_pool=pool)
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10000,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Errorf("ERROR: %s", err)
		// f.WriteString("ERROR: " + err.Error())

		return nil
	}
	fmt.Println("pong:" + string(pong))
	// f.WriteString("pong:" + string(pong) + "\n")
	return client
}

func generateH() ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)

	return H
}
