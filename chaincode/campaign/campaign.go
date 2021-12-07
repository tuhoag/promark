package campaign

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

var (
	cryptoServiceURL       = "http://external.promark.com:5000"
	cryptoParamsRequestURL = cryptoServiceURL + "/camp"
	logURL                 = "http://logs.promark.com:5003/log"
)

var com1URL string
var com2URL string
var comURL string

var camParam CampaignCryptoParams

// Struct of request data to ext service
type CampaignCryptoRequest struct {
	ID             string
	NumOfVerifiers int
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

type Campaign struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Advertiser string `json:"Advertiser"`
	Business   string `json:"Business"`
	CommC      string `json:"CommC"`
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

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	campaigns := []Campaign{
		{ID: "id1", Name: "campaign1", Advertiser: "Adv0", Business: "Bus0", CommC: ""},
		// {ID: "id2", Name: "campaign2", Advertiser: "Adv0", Business: "Bus0", CommC1: "", CommC2: ""},
	}

	for _, campaign := range campaigns {
		campaignJSON, err := json.Marshal(campaign)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(campaign.ID, campaignJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
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

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (s *SmartContract) GetAllCollectedData(ctx contractapi.TransactionContextInterface) ([]*CollectedData, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var allData []*CollectedData
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var data CollectedData
		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, err
		}
		allData = append(allData, &data)
	}

	return allData, nil
}

// Create a new campaign
func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, advertiser string, business string, verifierURLStr string) error {
	existing, err := ctx.GetStub().GetState(id)

	sendLog("campaignId", id)
	sendLog("campaignName", name)
	sendLog("advertiser", advertiser)
	sendLog("business", business)

	// split the verifier address
	// var verifierURL, commStr, totalCommEnc string
	// var ver1, ver2 string

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its id %s is existed", id)
	}

	verifierURLs := strings.Split(verifierURLStr, ";")
	numVerifiers := len(verifierURLs)

	sendLog("numVerifiers", string(numVerifiers))

	cryptoParams := requestCampaignCryptoParams(id, numVerifiers)
	sendLog("cryptoParams.H", convertBytesToPoint(cryptoParams.H).String())
	// sendLog("R1-1 value", string(cryptoParams.R1[0]))
	// sendLog("R1-2 value", string(cryptoParams.R1[1]))

	// // Request to external service to get params
	// var Ci, C ristretto.Point

	for i := 0; i < numVerifiers; i++ {
		verifierURL := verifierURLs[i]
		comURL = verifierURL + "/comm"
		sendLog("verifierURL", verifierURL)
		sendLog("comURL", comURL)

		// 	testVer(ver)
		// 	sendLog("id", id)
		// 	sendLog("Hvalue", string(cryptoParams.H))
		// 	sendLog("R1value", string(cryptoParams.R1[i]))

		comm := computeCommitment(id, comURL, i, cryptoParams)
		// commDec, _ := b64.StdEncoding.DecodeString(comm)
		// Ci = convertStringToPoint(string(commDec))
		sendLog("C"+string(i)+" encoding:", comm)
		// sendLog("C"+string(i)+" encoding:", comm)

		// 	commStr += comm + ";"

		// 	if i == 0 {
		// 		C = Ci
		// 	} else {
		// 		C.Add(&C, &Ci)
		// 	}
		// 	CBytes := C.Bytes()
		// 	totalCommEnc = b64.StdEncoding.EncodeToString(CBytes)
	}

	// sendLog("total Comm encoding:", totalCommEnc)
	// // End of comm computation

	// campaign := Campaign{
	// 	ID:         id,
	// 	Name:       name,
	// 	Advertiser: advertiser,
	// 	Business:   business,
	// 	CommC:      commStr,
	// 	// CommC2:     comm2,
	// }

	// campaignJSON, err := json.Marshal(campaign)
	// if err != nil {
	// 	return err
	// }

	// err = ctx.GetStub().PutState(id, campaignJSON)

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *SmartContract) DeleteCampaignByID(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	campaignJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if campaignJSON == nil {
		return false, fmt.Errorf("the backup %s does not exist", id)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to delete state:" + err.Error())
	}

	return true, err
}

func (s *SmartContract) AddCollectedData(ctx contractapi.TransactionContextInterface, id string, user string, n string, comm string, r string, addresses string) error {
	existing, err := ctx.GetStub().GetState(user)

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its id %s is existed", user)
	}

	var noVer int
	var rList string
	noVer, _ = strconv.Atoi(n)
	listOfVer := strings.Split(addresses, ";")
	listOfR := strings.Split(r, ";")

	// comm comuptation - start
	var Ci, C, Comm ristretto.Point

	for i := 0; i < noVer; i++ {
		comURL = listOfVer[i] + "/verify"

		commi := commVerify(id, comURL, listOfR[i])
		sendLog("AddCollectedData - listOfR[i]:", listOfR[i])
		sendLog("AddCollectedData - commi:", commi)

		tempCommDec, _ := b64.StdEncoding.DecodeString(commi)
		Ci = convertStringToPoint(string(tempCommDec))

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)
		}

		rList += listOfR[i] + ";"
	}

	// convert encoded Comm value
	commDec, _ := b64.StdEncoding.DecodeString(comm)
	Comm = convertStringToPoint(string(commDec))

	// check the comm conndition
	checkResult := C.Equals(&Comm)
	sendLog("AddCollectedData - rList:", string(rList))

	collectedData := CollectedData{
		User: user,
		Comm: comm,
		R1:   rList,
		// R2:   r2,
	}

	dataJSON, err := json.Marshal(collectedData)
	if err != nil {
		return err
	}

	if checkResult {
		sendLog("AddCollectedData - check comm result value:", "true")
		err = ctx.GetStub().PutState(user, dataJSON)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SmartContract) DeleteDataByUserId(ctx contractapi.TransactionContextInterface, userId string) (bool, error) {
	dataJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataJSON == nil {
		return false, fmt.Errorf("the backup %s does not exist", userId)
	}

	var collectedData CollectedData
	err = json.Unmarshal(dataJSON, &collectedData)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().DelState(userId)
	if err != nil {
		return false, fmt.Errorf("Failed to delete state:" + err.Error())
	}

	return true, err
}

func (s *SmartContract) QueryLedgerById(ctx contractapi.TransactionContextInterface, id string) ([]*Campaign, error) {
	queryString := fmt.Sprintf(`{"selector":{"id":{"$lte": "%s"}}}`, id)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var campaigns []*Campaign

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, &campaign)
	}

	resultsIterator.Close()

	return campaigns, nil
}

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

func requestCampaignCryptoParams(id string, numOfVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	message := CampaignCryptoRequest{id, numOfVerifiers}

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

func requestCamParams(id string, n int) {
	c := &http.Client{}

	message := CampaignCryptoRequest{id, n}

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", cryptoParamsRequestURL, strings.NewReader(request))
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

	err = json.Unmarshal([]byte(data), &camParam)
	if err != nil {
		println(err)
	}

	sendLog("returnedH", convertBytesToPoint(camParam.H).String())
}

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

	rBytes = cryptoParams.R1[i]
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
