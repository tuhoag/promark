package promark_utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
)

type PromarkRequest struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type PromarkResponse struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

type VerifierCryptoParams struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	S     string `json:"s"`
}

type VerifierCryptoParamsRequest struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
}

type DebugLog struct {
	Name  string
	Value string
}

type CampaignCryptoParams struct {
	CamID string `json:camId`
	H     string `json:"h"`
}

type Campaign struct {
	ID           string   `json:"ID"`
	Name         string   `json:"Name"`
	Advertiser   string   `json:"Advertiser"`
	Business     string   `json:"Business"`
	CommC        string   `json:"CommC"`
	VerifierURLs []string `json:"VerifierURLs"`
}

type CampaignCryptoRequest struct {
	CamId        string `json:"camId"`
	CustomerId   string `json:"customerId"`
	NumVerifiers int    `json:"numVerifiers"`
}

type VerifierCryptoChannelResult struct {
	URL                  string
	VerifierCryptoParams VerifierCryptoParams
	Error                error
}

var logURL = "http://logs.promark.com:5003/log"

func SendLog(name, message string, logMode string) {
	if logMode == "test" {
		return
	}

	logmessage := DebugLog{name, message}

	jsonLog, err := json.Marshal(logmessage)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
	}

	logRequest := string(jsonLog)

	c := &http.Client{}

	reqJSON, err := http.NewRequest("POST", logURL, strings.NewReader(logRequest))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	fmt.Println("return data all:", string(data))
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

	fmt.Println("Send create command")
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

func ParseResponse(responseStr string) (*PromarkResponse, error) {
	var response PromarkResponse
	err := json.Unmarshal([]byte(responseStr), &response)

	return &response, err
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

// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

func convertStringToPoint(s string) ristretto.Point {
	var H ristretto.Point
	var hBytes [32]byte

	tmp := []byte(s)
	copy(hBytes[:32], tmp[:])

	H.SetBytes(&hBytes)
	fmt.Println("in convertPointtoBytes hString: ", H)

	return H
}

func convertBytesToPoint(b []byte) ristretto.Point {
	var H ristretto.Point
	var hBytes [32]byte

	copy(hBytes[:32], b[:])

	result := H.SetBytes(&hBytes)
	fmt.Println("in convertBytesToPoint result:", result)

	return H
}

func convertBytesToScalar(b []byte) ristretto.Scalar {
	var r ristretto.Scalar
	var rBytes [32]byte

	copy(rBytes[:32], b[:])

	result := r.SetBytes(&rBytes)
	fmt.Println("in convertBytesToScalar result:", result)

	return r
}

func convertScalarToString(s ristretto.Scalar) string {
	sBytes := s.Bytes()
	sString := string(sBytes)

	return sString
}

func ConvertPointToString(h ristretto.Point) string {
	hBytes := h.Bytes()
	hEnc := b64.StdEncoding.EncodeToString(hBytes)

	return hEnc
}
