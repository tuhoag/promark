package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	putils "internal/promark_utils"
	"log"
	"net"
	// "strings"
	"sync"
)

const (
	connHost               = "external.promark.com"
	connPort               = "5000"
	connType               = "tcp"
	cryptoServiceSocketURL = "external.promark.com:5000"
)

func main() {
	log.Println("Start running")
	CreateCampaignSocket("c001", "Campaign001", "adv0", "bus0", []string{"peer0.adv0.promark.com:5000", "peer0.bus0.promark.com:5000"})
}

func SendData(conn net.Conn, message string) {
	fmt.Fprintf(conn, message+"\n")
}

func SendRequest(conn net.Conn, command string, data string) error {
	request := putils.PromarkRequest{
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
	response := putils.PromarkResponse{
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

func ParseResponse(responseStr string) (*putils.PromarkResponse, error) {
	var response putils.PromarkResponse
	err := json.Unmarshal([]byte(responseStr), &response)

	return &response, err
}

func CreateCampaignSocket(camId string, name string, advertiser string, business string, verifierURLs []string) (*putils.Campaign, error) {
	fmt.Println("Call RequestCampaignCryptoParamsSocket")
	cryptoParams, err := RequestCampaignCryptoParamsSocket(camId)
	if err != nil {
		return nil, err
	}

	// putils.SendLog("cryptoParams.H", cryptoParams.H)
	fmt.Println("cryptoParams.H:" + cryptoParams.H)
	numVerifiers := len(verifierURLs)
	vChannel := make(chan putils.VerifierCryptoChannelResult)
	// defer close(vChannel)

	wg := sync.WaitGroup{}
	wg.Add(numVerifiers)
	for i := 0; i < numVerifiers; i++ {
		verifierURL := verifierURLs[i]
		requestCreateVerifierCryptoURL := verifierURL
		// putils.SendLog("verifierURL", verifierURL)
		// putils.SendLog("comURL", requestCreateVerifierCryptoURL)

		fmt.Println("Call RequestToCreateVerifierCampaignCryptoParamsSocket: " + verifierURL)
		go ConcurrentRequestToCreateVerifierCampaignCryptoParamsSocket(camId, requestCreateVerifierCryptoURL, cryptoParams, vChannel, &wg)
	}

	wg.Wait()

	fmt.Println("Printing results")
	for i := 0; i < numVerifiers; i++ {
		result := <-vChannel
		fmt.Println(result.URL)
		fmt.Println(result.Error)
		fmt.Println(result.verifierCryptoParams)

		if result.Error != nil {
			return nil, result.Error
		}
	}

	close(vChannel)

	campaign := putils.Campaign{
		ID:           camId,
		Name:         name,
		Advertiser:   advertiser,
		Business:     business,
		CommC:        "",
		VerifierURLs: verifierURLs,
	}

	fmt.Println("Closing")
	return &campaign, nil
}

func RequestCampaignCryptoParamsSocket(camId string) (*putils.CampaignCryptoParams2, error) {
	conn, err := net.Dial("tcp", cryptoServiceSocketURL)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		log.Println("Error connecting:", err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	SendRequest(conn, "create", camId)

	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	log.Println("Reiceived: " + responseStr)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}
	var cryptoParams putils.CampaignCryptoParams2
	err = json.Unmarshal([]byte(response.Data), &cryptoParams)

	return &cryptoParams, nil
}

func ConcurrentRequestToCreateVerifierCampaignCryptoParamsSocket(camId string, requestCreateVerifierCryptoURL string, cryptoParams *putils.CampaignCryptoParams2, results chan putils.VerifierCryptoChannelResult, wg *sync.WaitGroup) {
	vCryptoParams, err := RequestToCreateVerifierCampaignCryptoParamsSocket(camId, requestCreateVerifierCryptoURL, cryptoParams)

	wg.Done()
	fmt.Println("Done with " + requestCreateVerifierCryptoURL)
	fmt.Println("vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("vCryptoParams.S:" + vCryptoParams.S)

	results <- putils.VerifierCryptoChannelResult{
		URL:                  requestCreateVerifierCryptoURL,
		verifierCryptoParams: vCryptoParams,
		Error:                err,
	}
}

func RequestToCreateVerifierCampaignCryptoParamsSocket(camId string, url string, cryptoParams *putils.CampaignCryptoParams2) (*putils.VerifierCryptoParams, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	requestArgs := putils.VerifierCryptoParamsRequest{
		CamId: camId,
		H:     cryptoParams.H,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	SendRequest(conn, "create", string(jsonArgs))
	// wait for response
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}

	var vCryptoParams putils.VerifierCryptoParams
	err = json.Unmarshal([]byte(response.Data), &vCryptoParams)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	fmt.Println("Returned-vCryptoParams.CamId:" + vCryptoParams.CamId)
	return &vCryptoParams, nil
}
