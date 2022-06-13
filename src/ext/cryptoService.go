package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	ristretto "github.com/bwesterb/go-ristretto"
	// redis "gopkg.in/redis.v4"
	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
	eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
	putils "internal/promark_utils"
)

var ctx = context.Background()

func main() {
	err := initialize()

	if err != nil {
		fmt.Printf("Cannot start 'tcp' server on external.promark.com because of error: %s", err)
		return
	}

	fmt.Printf("Using H: %s\n", H)

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

var secretFileName = "secret.mycert"
var H ristretto.Point
var HEnc string

func initialize() error {
	// check if the file storing H is existed
	secretFilePath := filepath.Join(".", "certs", secretFileName)
	os.MkdirAll("certs", 0700)

	fmt.Println(secretFilePath)
	data, err := ioutil.ReadFile(secretFilePath)
	if err != nil {
		// generate  and store in a file
		H = *pedersen.GenerateH()
		HEnc = eutils.ConvertPointToString(&H)

		f, err := os.Create(secretFilePath)
		if err != nil {
			return err
		}

		f.WriteString(HEnc)
		f.Close()
	} else {
		HEnc = string(data)
		Hp, err := eutils.ConvertStringToPoint(HEnc)
		if err != nil {
			return err
		}
		H = *Hp
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			putils.SendResponse(conn, fmt.Sprintf("%s", r), "")
		}
	}()

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

	if request.Command == "get" {
		GetDefaultCampaignCryptoParamsHandler(conn, request.Data)
	} else {
		putils.SendResponse(conn, "nocommand", "")
		conn.Close()
	}

	conn.Close()
}

func InitDefaultCampaignCryptoParamsHandler(conn net.Conn, requestData string) {
	// camId := requestData
	cryptoParams, err := initCampaignCryptoParams()

	if err != nil {
		panic(err)
	}

	param, err := json.Marshal(&cryptoParams)

	fmt.Println("result: " + string(param))
	putils.SendResponse(conn, "", string(param))
}

func GetDefaultCampaignCryptoParamsHandler(conn net.Conn, requestData string) {
	//get param from db
	// cryptoParams, err := GetCampaignCryptoParams("")
	cryptoParams := putils.CampaignCryptoParams{
		CamID: "",
		H:     HEnc,
	}

	//temporary return
	param, err := json.Marshal(cryptoParams)

	// fmt.Println("wrote to file:", n)

	if err != nil {
		panic(err)
	}

	fmt.Println("result: " + string(param))
	putils.SendResponse(conn, "", string(param))
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

func generateCryptoParams() (*putils.CampaignCryptoParams, error) {
	var cryptoParams putils.CampaignCryptoParams

	H := pedersen.GenerateH()
	hEnc := eutils.ConvertPointToString(H)
	fmt.Println("hString:.\n", hEnc)

	cryptoParams = putils.CampaignCryptoParams{
		CamID: "",
		H:     hEnc,
	}

	return &cryptoParams, nil
}

func initCampaignCryptoParams() (*putils.CampaignCryptoParams, error) {
	client := putils.GetRedisConnection()

	cryptoParams, err := generateCryptoParams()
	if err != nil {
		return nil, err
	}

	jsonParam, err := json.Marshal(*cryptoParams)
	if err != nil {
		return nil, err
	}

	//store to redis db
	err = client.Set(ctx, "", jsonParam, 0).Err()
	if err != nil {
		return nil, err
	}

	return cryptoParams, nil
}

// Campaign function part
func createCampaignCryptoParams(camId string) (*putils.CampaignCryptoParams, error) {
	client := putils.GetRedisConnection()

	//generate campaign param
	_, err := client.Get(ctx, camId).Result()
	var cryptoParams putils.CampaignCryptoParams

	if err != nil {
		fmt.Println(err)

		H := pedersen.GenerateH()
		hEnc := eutils.ConvertPointToString(H)
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
		err = client.Set(ctx, camId, jsonParam, 0).Err()
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

	val, err := client.Get(ctx, camId).Result()
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
