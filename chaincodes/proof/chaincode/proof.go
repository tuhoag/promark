package chaincode

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	// "log"
	"math/big"
	"net/http"
	// "strconv"
	"strings"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

var (
	cryptoServiceURL           = "http://external.promark.com:5000"
	cryptoParamsRequestURL     = cryptoServiceURL + "/camp"
	userCryptoParamsRequestURL = cryptoServiceURL + "/usercamp"
	logURL                     = "http://logs.promark.com:5003/log"
)

var com1URL string
var com2URL string
var comURL string

var camParam CampaignCryptoParams

// Struct of request data to ext service
type CampaignCryptoRequest struct {
	ID           string
	CustomerId   string
	NumVerifiers int
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

func (s *SmartContract) GetCustomerCampaignProof(ctx contractapi.TransactionContextInterface, camId string, userId string) (*ProofCustomerCampaign, error) {
	sendLog("GetCustomerCampaignProof", "")
	sendLog("Campaign:", camId)
	sendLog("userId", userId)

	campaignJSON, err := ctx.GetStub().GetState(camId)

	if err != nil {
		return nil, errors.New("Unable to read the world state")
	}

	if campaignJSON == nil {
		return nil, fmt.Errorf("Cannot get campaign since its raw id %s is unexisted", camId)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, err
	}

	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)
	sendLog("numVerifiers", string(numVerifiers))

	// get crypto params
	cryptoParams := requestCustomerCampaignCryptoParams(camId, userId, numVerifiers)

	var Ci, C ristretto.Point
	var totalCommEnc string

	var randomValues []string
	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]
		comURL = verifierURL + "/comm"
		sendLog("verifierURL", verifierURL)
		sendLog("comURL", comURL)

		// 	testVer(ver)
		// 	sendLog("id", id)
		// 	sendLog("Hvalue", string(cryptoParams.H))
		// 	sendLog("R1value", string(cryptoParams.R1[i]))

		comm := computeCommitment(camId, comURL, i, cryptoParams)
		// commDec, _ := b64.StdEncoding.DecodeString(comm)
		// Ci = convertStringToPoint(string(commDec))
		sendLog("C"+string(i)+" encoding:", comm)
		// sendLog("C"+string(i)+" encoding:", comm)

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)
		}
		CBytes := C.Bytes()
		totalCommEnc = b64.StdEncoding.EncodeToString(CBytes)

		randomValues = append(randomValues, b64.StdEncoding.EncodeToString(cryptoParams.R1[i]))
	}

	// get all verifiers URLs

	// calculate commitment
	proof := ProofCustomerCampaign{
		Comm: totalCommEnc,
		Rs:   randomValues,
	}

	return &proof, nil
}

func (s *SmartContract) CollectCustomerProofCampaign(ctx contractapi.TransactionContextInterface, proofId string, comm string, rsStr string) error {
	sendLog("CollectCustomerProofCampaign", "")
	sendLog("proofId", proofId)
	sendLog("comm", comm)
	sendLog("rsStr", rsStr)

	proofJSON, err := ctx.GetStub().GetState(proofId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if proofJSON != nil {
		return fmt.Errorf("the user proof raw id %s is existed", proofId)
	}

	rs := strings.Split(rsStr, ";")

	collectedProof := CollectedCustomerProof{
		ID:   proofId,
		Comm: comm,
		Rs:   rs,
	}

	collectedProofJSON, err := json.Marshal(collectedProof)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(proofId, collectedProofJSON)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) GetProofById(ctx contractapi.TransactionContextInterface, proofId string) (*CollectedCustomerProof, error) {
	proofJSON, err := ctx.GetStub().GetState(proofId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if proofJSON == nil {
		return nil, fmt.Errorf("the campaign raw id %s does not exist", proofId)
	}

	var proof CollectedCustomerProof
	err = json.Unmarshal(proofJSON, &proof)
	if err != nil {
		return nil, err
	}

	return &proof, nil
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

func requestCustomerCampaignCryptoParams(id string, userId string, numVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	// ID             string
	// CustomerId	   stringGet
	// NumOfVerifiers int
	message := CampaignCryptoRequest{
		ID:           id,
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

func requestCampaignCryptoParams(id string, numVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	message := CampaignCryptoRequest{
		ID:           id,
		CustomerId:   "",
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
