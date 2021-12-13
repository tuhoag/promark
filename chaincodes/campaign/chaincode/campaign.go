package chaincode

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"

	// "strconv"
	"strings"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

var com1URL string
var com2URL string
var requestCreateVerifierCryptoURL string

var camParam CampaignCryptoParams

type PromarkRequest struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type PromarkResponse struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

// Struct of request data to ext service
type CampaignCryptoRequest struct {
	CamId        string
	CustomerId   string
	NumVerifiers int
}

type NewVerifierCryptoParamsRequest struct {
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

// Struct of return data from ext service
type CampaignCryptoParams struct {
	H  []byte   `json:"hvalue"`
	R1 [][]byte `json:"r1"`
	// R2 []byte `json:"r2"`
}

type CampaignCryptoParams2 struct {
	CamID string `json:camId`
	H     string `json:"h"`
	// R1 [][]byte `json:"r1"`
	// R2 []byte `json:"r2"`
}

type Campaign struct {
	ID           string   `json:"ID"`
	Name         string   `json:"Name"`
	Advertiser   string   `json:"Advertiser"`
	Business     string   `json:"Business"`
	CommC        string   `json:"CommC"`
	VerifierURLs []string `json:"VerifierURLs"`
}

type ProofCustomerCampaign struct {
	Comm string   `json:"Comm"`
	Rs   []string `json:"Rs"`
}

type CollectedCustomerProof struct {
	ID   string   `json:"ID"`
	Comm string   `json:"Comm"`
	Rs   []string `json:"Rs"`
}

type CommRequest struct {
	ID string `json:"ID"`
	H  []byte `json:"H"`
	r  []byte `json:"r"`
}

type CollectedData struct {
	User string `json:"User"`
	Comm string `json:"Comm"`
	R1   string `json:"R1"`
	// R2   string `json:"R2"`
}

func (s *CampaignSmartContract) GetAllCampaigns(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var campaigns []*Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// sendLog("queryResponse.Value", string(queryResponse.Value))
		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

// func (s *CampaignSmartContract) GetAllCollectedData(ctx contractapi.TransactionContextInterface) ([]*CollectedData, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	// close the resultsIterator when this function is finished
// 	defer resultsIterator.Close()

// 	var allData []*CollectedData
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var data CollectedData
// 		err = json.Unmarshal(queryResponse.Value, &data)
// 		if err != nil {
// 			return nil, err
// 		}
// 		allData = append(allData, &data)
// 	}

// 	return allData, nil
// }
func (s *CampaignSmartContract) CreateCampaignSocket(ctx contractapi.TransactionContextInterface, camId string, name string, advertiser string, business string, verifierURLStr string) error {
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

// Create a new campaign
func (s *CampaignSmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, camId string, name string, advertiser string, business string, verifierURLStr string) error {
	sendLog("campaignId", camId)
	sendLog("campaignName", name)
	sendLog("advertiser", advertiser)
	sendLog("business", business)

	sendLog("camId", camId)

	existing, err := ctx.GetStub().GetState(camId)
	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its raw id %s is existed", camId)
	}

	verifierURLs := strings.Split(verifierURLStr, ";")
	numVerifiers := len(verifierURLs)

	sendLog("numVerifiers", string(numVerifiers))

	cryptoParams, err := requestCampaignCryptoParams(camId)
	if err != nil {
		return err
	}

	sendLog("cryptoParams.H", cryptoParams.H)

	for i := 0; i < numVerifiers; i++ {
		verifierURL := verifierURLs[i]
		requestCreateVerifierCryptoURL = verifierURL + "/camp"
		sendLog("verifierURL", verifierURL)
		sendLog("comURL", requestCreateVerifierCryptoURL)

		vCryptoParams, err := RequestToCreateVerifierCampaignCryptoParamsHandler(camId, requestCreateVerifierCryptoURL, cryptoParams)

		if err != nil {
			return err
		}

		sendLog("vCryptoParams.CamId", vCryptoParams.CamId)
		sendLog("vCryptoParams.H", vCryptoParams.H)
		sendLog("vCryptoParams.S", vCryptoParams.S)
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

	err = ctx.GetStub().PutState(camId, campaignJSON)

	if err != nil {
		return err
	}

	return nil
}

func RequestToCreateVerifierCampaignCryptoParamsHandler(camId string, url string, cryptoParams *CampaignCryptoParams2) (*VerifierCryptoParams, error) {
	// func computeCommitment(campID string, url string, i int, cryptoParams CampaignCryptoParams) string {
	//connect to verifier: campID,  H , r
	sendLog("request init verifier crypto at", url)

	client := &http.Client{}
	requestArgs := NewVerifierCryptoParamsRequest{
		CamId: camId,
		H:     cryptoParams.H,
	}

	jsonArgs, err := json.Marshal(requestArgs)
	request := string(jsonArgs)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(request))
	sendLog("request", request)
	// sendLog("err", err.Error())
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
	sendLog("data", string(data))
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	fmt.Println("return data all:", string(data))
	var vCryptoParams VerifierCryptoParams
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

	var campaign Campaign
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

// func (s *CampaignSmartContract) GetCustomerCampaignProof(ctx contractapi.TransactionContextInterface, camId string, userId string) (*ProofCustomerCampaign, error) {
// 	sendLog("GetCustomerCampaignProof", "")
// 	sendLog("Campaign:", camId)
// 	sendLog("userId", userId)

// 	campaignJSON, err := ctx.GetStub().GetState(camId)

// 	if err != nil {
// 		return nil, errors.New("Unable to read the world state")
// 	}

// 	if campaignJSON == nil {
// 		return nil, fmt.Errorf("Cannot get campaign since its raw id %s is unexisted", camId)
// 	}

// 	var campaign Campaign
// 	err = json.Unmarshal(campaignJSON, &campaign)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// generate a random values for each verifiers
// 	numVerifiers := len(campaign.VerifierURLs)
// 	sendLog("numVerifiers", string(numVerifiers))

// 	// get crypto params
// 	cryptoParams := requestCustomerCampaignCryptoParams(camId, userId, numVerifiers)

// 	var Ci, C ristretto.Point
// 	var totalCommEnc string

// 	var randomValues []string
// 	for i := 0; i < numVerifiers; i++ {
// 		verifierURL := campaign.VerifierURLs[i]
// 		comURL = verifierURL + "/comm"
// 		sendLog("verifierURL", verifierURL)
// 		sendLog("comURL", comURL)

// 		// 	testVer(ver)
// 		// 	sendLog("id", id)
// 		// 	sendLog("Hvalue", string(cryptoParams.H))
// 		// 	sendLog("R1value", string(cryptoParams.R1[i]))

// 		comm := computeCommitment(camId, comURL, i, cryptoParams)
// 		// commDec, _ := b64.StdEncoding.DecodeString(comm)
// 		// Ci = convertStringToPoint(string(commDec))
// 		sendLog("C"+string(i)+" encoding:", comm)
// 		// sendLog("C"+string(i)+" encoding:", comm)

// 		if i == 0 {
// 			C = Ci
// 		} else {
// 			C.Add(&C, &Ci)
// 		}
// 		CBytes := C.Bytes()
// 		totalCommEnc = b64.StdEncoding.EncodeToString(CBytes)

// 		randomValues = append(randomValues, b64.StdEncoding.EncodeToString(cryptoParams.R1[i]))
// 	}

// 	// get all verifiers URLs

// 	// calculate commitment
// 	proof := ProofCustomerCampaign{
// 		Comm: totalCommEnc,
// 		Rs:   randomValues,
// 	}

// 	return &proof, nil
// }

// func (s *CampaignSmartContract) CollectCustomerProofCampaign(ctx contractapi.TransactionContextInterface, proofId string, comm string, rsStr string) error {
// 	sendLog("CollectCustomerProofCampaign", "")
// 	sendLog("proofId", proofId)
// 	sendLog("comm", comm)
// 	sendLog("rsStr", rsStr)

// 	proofJSON, err := ctx.GetStub().GetState(proofId)
// 	if err != nil {
// 		return fmt.Errorf("failed to read from world state: %v", err)
// 	}

// 	if proofJSON != nil {
// 		return fmt.Errorf("the user proof raw id %s is existed", proofId)
// 	}

// 	rs := strings.Split(rsStr, ";")

// 	collectedProof := CollectedCustomerProof{
// 		ID:   proofId,
// 		Comm: comm,
// 		Rs:   rs,
// 	}

// 	collectedProofJSON, err := json.Marshal(collectedProof)
// 	if err != nil {
// 		return err
// 	}

// 	err = ctx.GetStub().PutState(proofId, collectedProofJSON)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *CampaignSmartContract) GetCampaignById(ctx contractapi.TransactionContextInterface, camId string) (*Campaign, error) {
	campaignJSON, err := ctx.GetStub().GetState(camId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if campaignJSON == nil {
		return nil, fmt.Errorf("the campaign raw id %s does not exist", camId)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

// func (s *CampaignSmartContract) GetProofById(ctx contractapi.TransactionContextInterface, proofId string) (*CollectedCustomerProof, error) {
// 	proofJSON, err := ctx.GetStub().GetState(proofId)
// 	// backupJSON, err := ctx.GetStub().GetState(backupID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	if proofJSON == nil {
// 		return nil, fmt.Errorf("the campaign raw id %s does not exist", proofId)
// 	}

// 	var proof CollectedCustomerProof
// 	err = json.Unmarshal(proofJSON, &proof)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &proof, nil
// }

// func (s *CampaignSmartContract) AddCollectedData(ctx contractapi.TransactionContextInterface, id string, user string, n string, comm string, r string, addresses string) error {
// 	existing, err := ctx.GetStub().GetState(user)

// 	if err != nil {
// 		return errors.New("Unable to read the world state")
// 	}

// 	if existing != nil {
// 		return fmt.Errorf("Cannot create asset since its id %s is existed", user)
// 	}

// 	var noVer int
// 	var rList string
// 	noVer, _ = strconv.Atoi(n)
// 	listOfVer := strings.Split(addresses, ";")
// 	listOfR := strings.Split(r, ";")

// 	// comm comuptation - start
// 	var Ci, C, Comm ristretto.Point

// 	for i := 0; i < noVer; i++ {
// 		comURL = listOfVer[i] + "/verify"

// 		commi := commVerify(id, comURL, listOfR[i])
// 		sendLog("AddCollectedData - listOfR[i]:", listOfR[i])
// 		sendLog("AddCollectedData - commi:", commi)

// 		tempCommDec, _ := b64.StdEncoding.DecodeString(commi)
// 		Ci = convertStringToPoint(string(tempCommDec))

// 		if i == 0 {
// 			C = Ci
// 		} else {
// 			C.Add(&C, &Ci)
// 		}

// 		rList += listOfR[i] + ";"
// 	}

// 	// convert encoded Comm value
// 	commDec, _ := b64.StdEncoding.DecodeString(comm)
// 	Comm = convertStringToPoint(string(commDec))

// 	// check the comm conndition
// 	checkResult := C.Equals(&Comm)
// 	sendLog("AddCollectedData - rList:", string(rList))

// 	collectedData := CollectedData{
// 		User: user,
// 		Comm: comm,
// 		R1:   rList,
// 		// R2:   r2,
// 	}

// 	dataJSON, err := json.Marshal(collectedData)
// 	if err != nil {
// 		return err
// 	}

// 	if checkResult {
// 		sendLog("AddCollectedData - check comm result value:", "true")
// 		err = ctx.GetStub().PutState(user, dataJSON)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *CampaignSmartContract) DeleteDataByUserId(ctx contractapi.TransactionContextInterface, userId string) (bool, error) {
// 	dataJSON, err := ctx.GetStub().GetState(userId)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	if dataJSON == nil {
// 		return false, fmt.Errorf("the backup %s does not exist", userId)
// 	}

// 	var collectedData CollectedData
// 	err = json.Unmarshal(dataJSON, &collectedData)
// 	if err != nil {
// 		return false, err
// 	}

// 	err = ctx.GetStub().DelState(userId)
// 	if err != nil {
// 		return false, fmt.Errorf("Failed to delete state:" + err.Error())
// 	}

// 	return true, err
// }

// func (s *CampaignSmartContract) QueryLedgerById(ctx contractapi.TransactionContextInterface, id string) ([]*Campaign, error) {
// 	queryString := fmt.Sprintf(`{"selector":{"id":{"$lte": "%s"}}}`, id)

// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resultsIterator.Close()

// 	var campaigns []*Campaign

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		var campaign Campaign
// 		err = json.Unmarshal(queryResponse.Value, &campaign)
// 		if err != nil {
// 			return nil, err
// 		}

// 		campaigns = append(campaigns, &campaign)
// 	}

// 	resultsIterator.Close()

// 	return campaigns, nil
// }

/////////////////// External service functions //////////////////////////////////
func testVer(url string) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	sendLog("testVer1", string(responseData))
	fmt.Println(string(responseData))
}

func sendLog(name, message string) {
	if LOG_MODE == "test" {
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

func requestCustomerCampaignCryptoParams(camId string, userId string, numVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	// ID             string
	// CustomerId	   string
	// NumOfVerifiers int
	message := CampaignCryptoRequest{
		CamId:        camId,
		CustomerId:   userId,
		NumVerifiers: numVerifiers,
	}

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", cryptoParamsRequestURL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return cryptoParams
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return cryptoParams
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return cryptoParams
	}

	fmt.Println("return data all:", string(data))

	err = json.Unmarshal([]byte(data), &cryptoParams)
	if err != nil {
		println(err)
	}

	return cryptoParams
}

func requestCampaignCryptoParams(camId string) (*CampaignCryptoParams2, error) {
	var cryptoParams CampaignCryptoParams2

	c := &http.Client{}

	message := CampaignCryptoRequest{
		CamId:        camId,
		CustomerId:   "",
		NumVerifiers: 0,
	}

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", cryptoParamsRequestURL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return nil, err
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	fmt.Println("return data all:", string(data))

	err = json.Unmarshal([]byte(data), &cryptoParams)
	if err != nil {
		return nil, err
	}

	return &cryptoParams, nil
}

// func requestCamParams(id string, n int) {
// 	c := &http.Client{}

// 	message := CampaignCryptoRequest{id, "", n}

// 	jsonData, err := json.Marshal(message)

// 	request := string(jsonData)

// 	reqJSON, err := http.NewRequest("POST", cryptoParamsRequestURL, strings.NewReader(request))
// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %v\n", err)
// 		return
// 	}

// 	respJSON, err := c.Do(reqJSON)
// 	if err != nil {
// 		fmt.Printf("http.Do() error: %v\n", err)
// 		return
// 	}
// 	defer respJSON.Body.Close()

// 	data, err := ioutil.ReadAll(respJSON.Body)
// 	if err != nil {
// 		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
// 		return
// 	}

// 	fmt.Println("return data all:", string(data))

// 	err = json.Unmarshal([]byte(data), &camParam)
// 	if err != nil {
// 		println(err)
// 	}

// 	sendLog("returnedH", convertBytesToPoint(camParam.H).String())
// }

func computeCommitment(campID string, url string, i int, cryptoParams CampaignCryptoParams) string {
	//connect to verifier: campID,  H , r
	sendLog("Start of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")
	// var param CommRequest
	var rBytes []byte
	var rEnc, hEnc string
	client := &http.Client{}
	// sendLog("create connection of commCompute:", "")

	hBytes := cryptoParams.H
	hEnc = b64.StdEncoding.EncodeToString(hBytes)
	sendLog("Encode H: ", hEnc)

	// if url == com1URL {
	// 	rBytes = camParam.R1
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)

	// } else if url == com2URL {
	// 	rBytes = camParam.R2
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	// }

	sendLog("num r values: ", string(len(cryptoParams.R1)))
	rBytes = cryptoParams.R1[i]
	// sendLog("R["+string(i)+"]: ", rBytes)
	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	sendLog("Encode R["+string(i)+"]: ", rEnc)

	// jsonData, _ := json.Marshal(param)
	message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", campID, hEnc, rEnc)
	// request := string(jsonData)

	sendLog("commCompute.message", message)

	reqJSON, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
	}

	respJSON, err := client.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
	}

	sendLog("commValue:", string(data))
	sendLog("end of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	return string(data)
}

func commVerify(campID string, url string, r string) string {
	//connect to verifier: campID,  H , r
	// sendLog("Start of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	// var hEnc string
	c := &http.Client{}

	// hBytes := camParam.H
	// hEnc = b64.StdEncoding.EncodeToString(hBytes)

	// sendLog("commVerify.r in string", r)

	// jsonData, _ := json.Marshal(param)
	// message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", campID, hEnc, r)
	message := fmt.Sprintf("{\"id\": \"%s\", \"r\": \"%s\"}", campID, r)

	// request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
	}

	// sendLog("commValue:", string(data))
	// sendLog("end of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	return string(data)
}

/////////////////// Pedersen functions //////////////////////////////////
// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

// Generate a random point on the curve
func generateH() ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)

	return H
}

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
