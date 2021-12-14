package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	// "encoding/json"
)

type PromarkRequest struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type PromarkResponse struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

type Campaign struct {
	ID           string   `json:"ID"`
	Name         string   `json:"Name"`
	Advertiser   string   `json:"Advertiser"`
	Business     string   `json:"Business"`
	CommC        string   `json:"CommC"`
	VerifierURLs []string `json:"VerifierURLs"`
}

type VerifierCryptoParams struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	S     string `json:"s"`
}

type CampaignCryptoParams2 struct {
	CamID string `json:camId`
	H     string `json:"h"`
	// R1 [][]byte `json:"r1"`
	// R2 []byte `json:"r2"`
}

const (
	connHost               = "external.promark.com"
	connPort               = "5000"
	connType               = "tcp"
	cryptoServiceSocketURL = "external.promark.com:5001"
)

func main() {
	log.Println("Start running")
	CreateCampaign("c001", "Campaign001", "adv0", "bus0", "0.0.0.0:5002")

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

func ParseResponse(responseStr string) (*PromarkResponse, error) {
	var response PromarkResponse
	err := json.Unmarshal([]byte(responseStr), &response)

	return &response, err
}

func CreateCampaign(camId string, name string, advertiser string, business string, verifierURLStr string) error {
	// sendLog("campaignId", camId)
	// sendLog("campaignName", name)
	// sendLog("advertiser", advertiser)
	// sendLog("business", business)

	// sendLog("camId", camId)

	// existing, err := ctx.GetStub().GetState(camId)
	// if err != nil {
	// 	return errors.New("Unable to read the world state")
	// }

	// if existing != nil {
	// 	return fmt.Errorf("Cannot create asset since its raw id %s is existed", camId)
	// }

	verifierURLs := strings.Split(verifierURLStr, ";")
	numVerifiers := len(verifierURLs)

	cryptoParams, err := RequestCampaignCryptoParamsSocket(camId)
	if err != nil {
		return err
	}

	// sendLog("cryptoParams.H", cryptoParams.H)

	for i := 0; i < numVerifiers; i++ {
		verifierURL := verifierURLs[i]
		requestCreateVerifierCryptoURL := verifierURL
		// sendLog("verifierURL", verifierURL)
		// sendLog("comURL", requestCreateVerifierCryptoURL)

		vCryptoParams, err := RequestToCreateVerifierCampaignCryptoParamsSocket(camId, requestCreateVerifierCryptoURL, cryptoParams)

		if err != nil {
			return err
		}
		fmt.Println(vCryptoParams)
		fmt.Println("vCryptoParams.CamId" + vCryptoParams.CamId)
		fmt.Println("vCryptoParams.H" + vCryptoParams.H)
		fmt.Println("vCryptoParams.S" + vCryptoParams.S)

		// sendLog("vCryptoParams.CamId", vCryptoParams.CamId)
		// sendLog("vCryptoParams.H", vCryptoParams.H)
		// sendLog("vCryptoParams.S", vCryptoParams.S)
	}

	campaign := Campaign{
		ID:           camId,
		Name:         name,
		Advertiser:   advertiser,
		Business:     business,
		CommC:        "",
		VerifierURLs: verifierURLs,
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	// err = ctx.GetStub().PutState(camId, campaignJSON)

	if err != nil {
		return err
	}

	fmt.Println("CampaignJSON:" + string(campaignJSON))

	return nil
}

func RequestCampaignCryptoParamsSocket(camId string) (*CampaignCryptoParams2, error) {
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
	var cryptoParams CampaignCryptoParams2
	err = json.Unmarshal([]byte(response.Data), &cryptoParams)

	return &cryptoParams, nil
}

func RequestToCreateVerifierCampaignCryptoParamsSocket(camId string, url string, cryptoParams *CampaignCryptoParams2) (*VerifierCryptoParams, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	SendRequest(conn, "create", camId)
	// wait for response
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived: " + responseStr)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}

	var vCryptoParams VerifierCryptoParams
	err = json.Unmarshal([]byte(response.Data), &cryptoParams)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	return &vCryptoParams, nil
}
