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
	"strings"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

var (
	extURL        = "http://external.promark.com:5000"
	camRequestURL = extURL + "/camp"
	ver1URL       = "http://verifier1.promark.com:5001"
	com1URL       = ver1URL + "/comm"
	ver2URL       = "http://verifier2.promark.com:5002"
	com2URL       = ver2URL + "/comm"
	logURL        = "http://logs.promark.com:5003/log"
)

var camParam CampaignParam

// Struct of request data to ext service
type Cam struct {
	ID string
	No int
}

type DebugLog struct {
	Name  string
	Value string
}

// Struct of return data from ext service
type CampaignParam struct {
	H  []byte `json:"hvalue"`
	R1 []byte `json:"r1"`
	R2 []byte `json:"r2"`
}

// Struct of data store in Blockchain ledger
type Campaign struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Advertiser string `json:"Advertiser"`
	Business   string `json:"Business"`
	CommC1     string `json:"CommC1"`
	CommC2     string `json:"CommC2`
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
	R2   string `json:"R2"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	campaigns := []Campaign{
		{ID: "id1", Name: "campaign1", Advertiser: "Adv0", Business: "Bus0", CommC1: "", CommC2: ""},
		{ID: "id2", Name: "campaign2", Advertiser: "Adv0", Business: "Bus0", CommC1: "", CommC2: ""},
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
func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, advertiser string, business string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its id %s is existed", id)
	}

	// Request to external service to get params
	var C1, C2, C ristretto.Point

	requestCamParams(id)

	sendLog("id", id)
	sendLog("Hvalue", string(camParam.H))
	sendLog("R1value", string(camParam.R1))
	sendLog("com1URL", string(com1URL))
	sendLog("R2value", string(camParam.R2))
	sendLog("com2URL", string(com2URL))

	comm1 := commCompute(id, com1URL)
	comm1Dec, _ := b64.StdEncoding.DecodeString(comm1)
	C1 = convertStringToPoint(string(comm1Dec))

	comm2 := commCompute(id, com2URL)
	comm2Dec, _ := b64.StdEncoding.DecodeString(comm2)
	C2 = convertStringToPoint(string(comm2Dec))

	C.Add(&C1, &C2)
	CBytes := C.Bytes()
	totalCommEnc := b64.StdEncoding.EncodeToString(CBytes)
	sendLog("total Comm encoding:", totalCommEnc)

	// End of comm computation

	campaign := Campaign{
		ID:         id,
		Name:       name,
		Advertiser: advertiser,
		Business:   business,
		CommC1:     comm1,
		CommC2:     comm2,
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, campaignJSON)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) AddCollectedData(ctx contractapi.TransactionContextInterface, id string, user string, comm string, r1 string, r2 string) error {
	existing, err := ctx.GetStub().GetState(user)

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its id %s is existed", user)
	}

	// request
	requestCamParams(id)

	// comm comuptation - start
	var C1, C2, C, Comm ristretto.Point

	sendLog("id", id)
	sendLog("Hvalue", string(camParam.H))

	// convert encoded Comm value
	commDec, _ := b64.StdEncoding.DecodeString(comm)
	Comm = convertStringToPoint(string(commDec))

	comm1 := commVerify(id, com1URL, r1)
	comm1Dec, _ := b64.StdEncoding.DecodeString(comm1)
	C1 = convertStringToPoint(string(comm1Dec))

	comm2 := commVerify(id, com2URL, r2)
	comm2Dec, _ := b64.StdEncoding.DecodeString(comm2)
	C2 = convertStringToPoint(string(comm2Dec))

	C.Add(&C1, &C2)
	// comm computation - end

	// check the comm conndition
	checkResult := C.Equals(&Comm)

	collectedData := CollectedData{
		User: user,
		Comm: comm,
		R1:   r1,
		R2:   r2,
	}

	dataJSON, err := json.Marshal(collectedData)
	if err != nil {
		return err
	}

	if checkResult {
		err = ctx.GetStub().PutState(id, dataJSON)
		if err != nil {
			return err
		}
	}

	return nil
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
func testVer1() {
	response, err := http.Get(ver1URL)

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

func requestCamParams(id string) {
	c := &http.Client{}

	message := Cam{id, 2}

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", camRequestURL, strings.NewReader(request))
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

func commCompute(campID string, url string) string {
	//connect to verifier: campID,  H , r
	sendLog("Start of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")
	var param CommRequest
	var rBytes []byte
	var rEnc, hEnc string
	c := &http.Client{}

	hBytes := camParam.H
	hEnc = b64.StdEncoding.EncodeToString(hBytes)

	if url == com1URL {
		rBytes = camParam.R1
		rEnc = b64.StdEncoding.EncodeToString(rBytes)

	} else if url == com2URL {
		rBytes = camParam.R2
		rEnc = b64.StdEncoding.EncodeToString(rBytes)
	}

	sendLog("commCompute.R in string", string(rBytes))

	param = CommRequest{ID: campID, H: hBytes, r: rBytes}
	sendLog("commCompute.param in string", param.ID)
	sendLog("commCompute.param in string", string(param.H))
	sendLog("commCompute.param in string", string(param.r))

	// jsonData, _ := json.Marshal(param)
	message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", param.ID, hEnc, rEnc)
	// request := string(jsonData)

	sendLog("commCompute.message", message)

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

	sendLog("commValue:", string(data))
	sendLog("end of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	return string(data)
}

func commVerify(campID string, url string, r string) string {
	//connect to verifier: campID,  H , r
	sendLog("Start of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	var hEnc string
	c := &http.Client{}

	hBytes := camParam.H
	hEnc = b64.StdEncoding.EncodeToString(hBytes)

	sendLog("commVerify.r in string", r)

	// jsonData, _ := json.Marshal(param)
	message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", campID, hEnc, r)
	// request := string(jsonData)

	sendLog("commVerify.message", message)

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

	sendLog("commValue:", string(data))
	sendLog("end of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

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
