package main

import (
	"bufio"
	// b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	// "github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	putils "internal/promark_utils"
	"io/ioutil"
	"log"
	// "math/big"
	"net"
	"net/http"
	// "os"
	"strings"
	"sync"
)

var LOG_MODE = "DEBUG"

type CampaignSmartContract struct {
	contractapi.Contract
}

var (
	cryptoServiceSocketURL     = "external.promark.com:5000"
	cryptoServiceURL           = "http://external.promark.com:5000"
	cryptoParamsRequestURL     = cryptoServiceURL + "/camp"
	userCryptoParamsRequestURL = cryptoServiceURL + "/usercamp"
	logURL                     = "http://logs.promark.com:5003/log"
)

func (s *CampaignSmartContract) GetAllCampaigns(ctx contractapi.TransactionContextInterface) ([]*putils.Campaign, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var campaigns []*putils.Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// putils.SendLog("queryResponse.Value", string(queryResponse.Value))
		var campaign putils.Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (s *CampaignSmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, camId string, name string, advertiser string, business string, verifierURLStr string) error {
	putils.SendLog("campaignId", camId, LOG_MODE)
	putils.SendLog("campaignName", name, LOG_MODE)
	putils.SendLog("advertiser", advertiser, LOG_MODE)
	putils.SendLog("business", business, LOG_MODE)

	putils.SendLog("camId", camId, LOG_MODE)

	existing, err := ctx.GetStub().GetState(camId)
	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its raw id %s is existed", camId)
	}

	verifierURLs := strings.Split(verifierURLStr, ";")
	campaign, err := CreateCampaignSocket(camId, name, advertiser, business, verifierURLs)

	if err != nil {
		return err
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(camId, campaignJSON)

	if err != nil {
		return err
	}

	fmt.Println("CampaignJSON:" + string(campaignJSON))

	return nil
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
		fmt.Println(result.VerifierCryptoParams)

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

func RequestCampaignCryptoParamsSocket(camId string) (*putils.CampaignCryptoParams, error) {
	conn, err := net.Dial("tcp", cryptoServiceSocketURL)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		log.Println("Error connecting:", err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	putils.SendRequest(conn, "create", camId)

	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	log.Println("Reiceived: " + responseStr)

	response, err := putils.ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}
	var cryptoParams putils.CampaignCryptoParams
	err = json.Unmarshal([]byte(response.Data), &cryptoParams)

	return &cryptoParams, nil
}

func ConcurrentRequestToCreateVerifierCampaignCryptoParamsSocket(camId string, requestCreateVerifierCryptoURL string, cryptoParams *putils.CampaignCryptoParams, results chan putils.VerifierCryptoChannelResult, wg *sync.WaitGroup) {
	vCryptoParams, err := RequestToCreateVerifierCampaignCryptoParamsSocket(camId, requestCreateVerifierCryptoURL, cryptoParams)

	wg.Done()
	fmt.Println("Done with " + requestCreateVerifierCryptoURL)
	fmt.Println("vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("vCryptoParams.S:" + vCryptoParams.S)

	results <- putils.VerifierCryptoChannelResult{
		URL:                  requestCreateVerifierCryptoURL,
		VerifierCryptoParams: *vCryptoParams,
		Error:                err,
	}
}

func RequestToCreateVerifierCampaignCryptoParamsSocket(camId string, url string, cryptoParams *putils.CampaignCryptoParams) (*putils.VerifierCryptoParams, error) {
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

	putils.SendRequest(conn, "create", string(jsonArgs))
	// wait for response
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)

	response, err := putils.ParseResponse(responseStr)

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

func RequestToCreateVerifierCampaignCryptoParamsHandler(camId string, url string, cryptoParams *putils.CampaignCryptoParams) (*putils.VerifierCryptoParams, error) {
	// func computeCommitment(campID string, url string, i int, cryptoParams CampaignCryptoParams) string {
	//connect to verifier: campID,  H , r
	putils.SendLog("request init verifier crypto at", url, LOG_MODE)

	client := &http.Client{}
	requestArgs := putils.VerifierCryptoParamsRequest{
		CamId: camId,
		H:     cryptoParams.H,
	}

	jsonArgs, err := json.Marshal(requestArgs)
	request := string(jsonArgs)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(request))
	putils.SendLog("request", request, LOG_MODE)
	// putils.SendLog("err", err.Error())
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	respJSON, err := client.Do(reqData)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return nil, err
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	putils.SendLog("data", string(data), LOG_MODE)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	fmt.Println("return data all:", string(data))
	var vCryptoParams putils.VerifierCryptoParams
	err = json.Unmarshal([]byte(data), &vCryptoParams)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	return &vCryptoParams, nil
}

func (s *CampaignSmartContract) DeleteCampaignByID(ctx contractapi.TransactionContextInterface, camId string) (bool, error) {
	campaignJSON, err := ctx.GetStub().GetState(camId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if campaignJSON == nil {
		return false, fmt.Errorf("the campaign raw id %s does not exist", camId)
	}

	var campaign putils.Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().DelState(camId)
	if err != nil {
		return false, fmt.Errorf("Failed to delete state:" + err.Error())
	}

	return true, err
}

func (s *CampaignSmartContract) GetCampaignById(ctx contractapi.TransactionContextInterface, camId string) (*putils.Campaign, error) {
	campaignJSON, err := ctx.GetStub().GetState(camId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if campaignJSON == nil {
		return nil, fmt.Errorf("the campaign raw id %s does not exist", camId)
	}

	var campaign putils.Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}
