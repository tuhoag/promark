package promark_utils

import (
	"bufio"
	// b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	// "math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	ristretto "github.com/bwesterb/go-ristretto"
	redis "github.com/go-redis/redis/v8"
	tecc_utils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

type PromarkRequest struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type PromarkResponse struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

type VerifierCryptoParamsRequest struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
}

type VerifierCryptoParams struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	S     string `json:"s"`
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
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Advertiser   string   `json:"advertiser"`
	Publisher    string   `json:"publisher"`
	VerifierURLs []string `json:"verifierURLs"`
	DeviceIds    []string `json:"deviceIds"`
	StartTime    int64    `json:"startTime"`
	EndTime      int64    `json:"endTime"`
}

type Proof struct {
	Comm string   `json:"comm"`
	Rs   []string `json:"rs"`
}

type CampaignCryptoRequest struct {
	CamId        string `json:"camId"`
	CustomerId   string `json:"customerId"`
	NumVerifiers int    `json:"numVerifiers"`
}

type ProofGenerationRequest struct {
	CamId      string `json:"camId"`
	CustomerId string `json:"customerId`
}

type VerificationRequest struct {
	CamId string `json:"camId"`
	R     string `json:"r"`
}

type VerificationResponse struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	R     string `json:"r"`
	S     string `json:"s"`
	Comm  string `json:"comm"`
}

type CampaignCustomerVerifierProof struct {
	CamId  string `json:"camId"`
	UserId string `json:"userId`
	H      string `json:"h"`
	R      string `json:"r"`
	S      string `json:"s"`
	Comm   string `json:"comm"`
}

type CampaignCustomerProof struct {
	CamId  string   `json:"camId"`
	UserId string   `json:"userId`
	Rs     []string `json:"rs"`
	Comm   string   `json:"comm"`
}

type ProofCustomerCampaign struct {
	Comm    string   `json:"comm"`
	Rs      []string `json:"rs"`
	SubComs []string `json:"subComs"`
}

type PoCProof struct {
	Comm    string   `json:"comm"`
	Rs      []string `json:"rs"`
	SubComs []string `json:"subComs"`
}

type CollectedCustomerProof struct {
	Id            string    `json:"id"`
	CustomerProof Proof     `json:"customerProof"`
	LocationProof Proof     `json:"locationProof"`
	Comm          string    `json:"comm"`
	Rs            []string  `json:"rs"`
	AddedTime     time.Time `json:"addedTime"`
	AddedTimeStr  string    `json:"addedTimeStr"`
}

type VerifierCryptoChannelResult struct {
	URL                  string
	VerifierCryptoParams VerifierCryptoParams
	Error                error
}

type VerifierProofChannelResult struct {
	URL   string
	Proof CampaignCustomerVerifierProof
	Error error
}

type VerifierCommitmentChannelResult struct {
	URL   string
	Comm  string
	Error error
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
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10000,
	})

	return client
}

func GenerateProofFromVerifiers(campaign *Campaign, customerId string) (*PoCProof, error) {
	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)

	var C ristretto.Point
	var subComs, randomValues []string

	C.SetZero()
	fmt.Printf("Init C: %s\n", C)
	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]

		fmt.Println("Call RequestToCreateVerifierCampaignCryptoParamsSocket: " + verifierURL)
		verifierProof, err := RequestCommitmentNoSave(campaign.Id, customerId, verifierURL)

		if err != nil {
			return nil, err
		}

		Ci, err := tecc_utils.ConvertStringToPoint(verifierProof.Comm)

		if err != nil {
			return nil, err
		}

		fmt.Printf("Transfer %s to point %s\n", verifierProof.Comm, Ci)

		C.Add(&C, Ci)
		fmt.Printf("Current C: %s after adding %s\n", C, Ci)

		randomValues = append(randomValues, verifierProof.R)
		subComs = append(subComs, verifierProof.Comm)
	}

	CommEnc := tecc_utils.ConvertPointToString(&C)

	proof := PoCProof{
		Comm:    CommEnc,
		Rs:      randomValues,
		SubComs: subComs,
	}

	fmt.Println("proof.Comm: " + proof.Comm)
	fmt.Printf("proof.Rs: %s\n", proof.Rs)
	fmt.Printf("proof.SubComs: %s\n", proof.SubComs)

	return &proof, nil
}

func GenerateProofFromVerifiersSocketAsync(campaign *Campaign, userId string) (*PoCProof, error) {
	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)
	// putils.SendLog("numVerifiers", string(numVerifiers), LOG_MODE)

	var C ristretto.Point
	C.SetZero()
	vProofChannel := make(chan VerifierProofChannelResult)

	wg := sync.WaitGroup{}
	wg.Add(numVerifiers)

	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]

		fmt.Println("Call RequestToCreateVerifierCampaignCryptoParamsSocket: " + verifierURL)
		go ConcurrentRequestCommitment(campaign.Id, userId, verifierURL, vProofChannel, &wg)
	}

	fmt.Println("Printing results")
	var subComs, randomValues []string

	for i := 0; i < numVerifiers; i++ {
		result := <-vProofChannel
		fmt.Println(result.URL)
		fmt.Println(result.Error)
		fmt.Println(result.Proof)

		// SendLog("result.URL:"+result.URL+"H:"+result.Proof.H+"-R:"+result.Proof.R+"-S:"+result.Proof.S+"-Comm:"+result.Proof.Comm, "", LOG_MODE)
		// putils.SendLog("result.Error", result.Error.Error(), LOG_MODE)
		// putils.SendLog("result.Proof.H", result.Proof.H, LOG_MODE)
		// putils.SendLog("result.Proof.R", result.Proof.R, LOG_MODE)
		// putils.SendLog("result.Proof.S", result.Proof.S, LOG_MODE)
		// putils.SendLog("result.Proof.Comm", result.Proof.Comm, LOG_MODE)

		if result.Error != nil {
			return nil, result.Error
		}

		Ci, err := tecc_utils.ConvertStringToPoint(result.Proof.Comm)

		if err != nil {
			return nil, err
		}

		C.Add(&C, Ci)

		randomValues = append(randomValues, result.Proof.R)
		subComs = append(subComs, result.Proof.Comm)
	}

	close(vProofChannel)

	CommEnc := tecc_utils.ConvertPointToString(&C)

	proof := PoCProof{
		Comm:    CommEnc,
		Rs:      randomValues,
		SubComs: subComs,
	}

	fmt.Println("proof.Comm: " + proof.Comm)
	fmt.Printf("proof.Rs: %s\n", proof.Rs)
	fmt.Printf("proof.SubComs: %s\n", proof.SubComs)

	return &proof, nil
}

func ConcurrentRequestCommitment(camId string, customerId string, url string, results chan VerifierProofChannelResult, wg *sync.WaitGroup) {
	verifierProof, err := RequestCommitment(camId, customerId, url)

	wg.Done()

	fmt.Println("Done with " + url)
	fmt.Println("vCryptoParams.CamId:" + verifierProof.CamId)
	fmt.Println("vCryptoParams.CustomerId:" + verifierProof.UserId)
	fmt.Println("vCryptoParams.H:" + verifierProof.H)
	fmt.Println("vCryptoParams.R:" + verifierProof.R)
	fmt.Println("vCryptoParams.S:" + verifierProof.S)

	// putils.SendLog("Done with", url, LOG_MODE)
	// putils.SendLog("vCryptoParams.CamId:", verifierProof.CamId, LOG_MODE)
	// putils.SendLog("vCryptoParams.CustomerId:", verifierProof.UserId, LOG_MODE)
	// putils.SendLog("vCryptoParams.S:", verifierProof.S, LOG_MODE)

	results <- VerifierProofChannelResult{
		URL:   url,
		Proof: *verifierProof,
		Error: err,
	}
}

func RequestCommitmentNoSave(camId string, customerId string, url string) (*CampaignCustomerVerifierProof, error) {
	// putils.SendLog("RequestCommitment at", url, LOG_MODE)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// putils.SendLog("Error connecting:", err.Error(), LOG_MODE)

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	requestArgs := ProofGenerationRequest{
		CamId:      camId,
		CustomerId: customerId,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	SendRequest(conn, "commit-nosave", string(jsonArgs))
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)
	// SendLog("Reiceived From: "+url+"-Response:", responseStr, LOG_MODE)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}

	var subProof CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(response.Data), &subProof)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	fmt.Println("Returned from " + url + "-subProof.CamId:" + subProof.CamId)

	return &subProof, nil
}

func RequestCommitment(camId string, customerId string, url string) (*CampaignCustomerVerifierProof, error) {
	// putils.SendLog("RequestCommitment at", url, LOG_MODE)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// putils.SendLog("Error connecting:", err.Error(), LOG_MODE)

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	requestArgs := ProofGenerationRequest{
		CamId:      camId,
		CustomerId: customerId,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	SendRequest(conn, "commit", string(jsonArgs))
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)
	// SendLog("Reiceived From: "+url+"-Response:", responseStr, LOG_MODE)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}

	var subProof CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(response.Data), &subProof)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	fmt.Println("Returned from " + url + "-subProof.CamId:" + subProof.CamId)

	return &subProof, nil
}

func StringInSlice(a string, list []string) bool {

	for _, b := range list {
		if strings.Compare(a, b) == 0 {
			return true
		}
	}
	return false
}
