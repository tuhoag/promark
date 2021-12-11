package chaincode

import (
	b64 "encoding/base64"
	"encoding/json"

	// "errors"
	"fmt"
	"io/ioutil"

	// "log"
	"math/big"
	"net/http"

	// "strconv"
	"strings"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
)

type ProofSmartContract struct {
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
	CamId        string `json:"camId"`
	CustomerId   string `json:"userId"`
	NumVerifiers int    `json:"numVerifiers"`
}

type VerificationRequest struct {
	CamId string `json:"camId"`
	R     string `json:"R"`
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
	Comm    string   `json:"Comm"`
	Rs      []string `json:"Rs"`
	SubComs []string `json:"SubComs`
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

func (s *ProofSmartContract) GenerateCustomerCampaignProof(ctx contractapi.TransactionContextInterface, camId string, userId string) (*ProofCustomerCampaign, error) {
	SendLog("GenerateCustomerCampaignProof", "")
	SendLog("camId:", camId)
	SendLog("userId", userId)

	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	SendLog("response.Payload", string(response.Payload))
	SendLog("response.Status", string(response.Status))
	SendLog("response.message", string(response.Message))
	// sendLog("response.error", string(response.Error))
	// sendLog("response.message is nil", strconv.FormatBool(response.Message == ""))

	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}

	var campaign Campaign

	err := json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return nil, err
	}

	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)
	SendLog("numVerifiers", string(numVerifiers))

	// get crypto params
	// cryptoParams := requestCustomerCampaignCryptoParams(camId, userId, numVerifiers)

	var Ci, C ristretto.Point

	var subComs, randomValues []string

	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]
		comURL = verifierURL + "/camp/" + camId + "/proof/" + userId
		SendLog("verifierURL", verifierURL)
		SendLog("comURL", comURL)

		// 	testVer(ver)
		// 	sendLog("id", id)
		// 	sendLog("Hvalue", string(cryptoParams.H))
		// 	sendLog("R1value", string(cryptoParams.R1[i]))

		subProof, err := computeCommitment2(camId, userId, comURL)

		if err != nil {
			return nil, err
		}

		// commDec, _ := b64.StdEncoding.DecodeString(comm)
		// Ci = convertStringToPoint(string(commDec))
		SendLog("H"+string(i)+" encoding:", subProof.H)
		SendLog("S"+string(i)+" encoding:", subProof.S)
		SendLog("R"+string(i)+" encoding:", subProof.R)
		SendLog("Comm"+string(i)+" encoding:", subProof.Comm)
		CiBytes, _ := b64.StdEncoding.DecodeString(subProof.Comm)
		Ci = convertBytesToPoint(CiBytes)

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)

			SendLog("Current Comm", b64.StdEncoding.EncodeToString(C.Bytes()))
		}

		randomValues = append(randomValues, subProof.R)
		subComs = append(subComs, subProof.Comm)
	}
	CommBytes := C.Bytes()
	CommEnc := b64.StdEncoding.EncodeToString(CommBytes)

	// get all verifiers URLs

	// calculate commitment
	proof := ProofCustomerCampaign{
		Comm:    CommEnc,
		Rs:      randomValues,
		SubComs: subComs,
	}

	return &proof, nil
}

func (s *ProofSmartContract) GetAllProofs(ctx contractapi.TransactionContextInterface) ([]*CollectedCustomerProof, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var proofs []*CollectedCustomerProof
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		SendLog("queryResponse.Value", string(queryResponse.Value))
		var proof CollectedCustomerProof
		err = json.Unmarshal(queryResponse.Value, &proof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, &proof)
	}

	return proofs, nil
}

func (s *ProofSmartContract) AddCustomerProofCampaign(ctx contractapi.TransactionContextInterface, proofId string, comm string, rsStr string) error {
	SendLog("AddCustomerProofCampaign", "")
	SendLog("proofId", proofId)
	SendLog("comm", comm)
	SendLog("rsStr", rsStr)

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

func (s *ProofSmartContract) GetProofById(ctx contractapi.TransactionContextInterface, proofId string) (*CollectedCustomerProof, error) {
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

func (s *ProofSmartContract) DeleteProofByID(ctx contractapi.TransactionContextInterface, proofId string) (bool, error) {
	proofJSON, err := ctx.GetStub().GetState(proofId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if proofJSON == nil {
		return false, fmt.Errorf("the proof id %s does not exist", proofId)
	}

	var proof CollectedCustomerProof
	err = json.Unmarshal(proofJSON, &proof)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().DelState(proofId)
	if err != nil {
		return false, fmt.Errorf("Failed to delete state:" + err.Error())
	}

	return true, err
}

func (s *ProofSmartContract) VerifyCampaignProof(ctx contractapi.TransactionContextInterface, camId string, proofId string) (bool, error) {
	SendLog("VerifyCampaignProof", "")
	SendLog("camId", camId)
	SendLog("proofId", proofId)

	_, err := s.GetProofById(ctx, proofId)

	if err != nil {
		return false, err
	}

	// get campaign
	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	SendLog("response.Payload", string(response.Payload))
	SendLog("response.Status", string(response.Status))
	SendLog("response.message", string(response.Message))

	if response.Message != "" {
		return false, fmt.Errorf(response.Message)
	}

	var campaign Campaign

	err = json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return false, err
	}

	SendLog("campaign.ID", campaign.ID)
	SendLog("campaign.Name", campaign.Name)
	SendLog("campaign.VerifierURLs", string(len(campaign.VerifierURLs)))

	// get proof
	proof, err := s.GetProofById(ctx, proofId)

	if err != nil {
		return false, err
	}

	// SendLog("proof.H", proof.H)
	SendLog("proof.Comm", proof.Comm)

	for i, R := range proof.Rs {
		SendLog("proof.R["+string(i)+"]", R)
	}

	// callinng verifiers to calculate proof.Comm again based on proof.Rs
	var C ristretto.Point
	for i, verifierURL := range campaign.VerifierURLs {
		// call verifier to compute sub commitment
		comURL := verifierURL + "/camp/" + camId + "/verify"
		ciEnc, err := computeCommitment3(camId, proof.Rs[i], comURL)

		if err != nil {
			return false, err
		}

		ciBytes, err := b64.StdEncoding.DecodeString(*ciEnc)

		if err != nil {
			return false, err
		}

		ci := convertBytesToPoint(ciBytes)

		if i == 0 {
			C = ci
		} else {
			C.Add(&C, &ci)
		}
	}

	CommBytes, err := b64.StdEncoding.DecodeString(proof.Comm)
	if err != nil {
		return false, err
	}

	SendLog("proof.Com", proof.Comm)
	SendLog("calculated Com", b64.StdEncoding.EncodeToString(C.Bytes()))
	comm := convertBytesToPoint(CommBytes)
	if C.Equals(&comm) {
		return true, nil
	} else {
		return false, nil
	}
}

func SendLog(name, message string) {
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

func RequestCustomerCampaignCryptoParams(id string, userId string, numVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	// ID             string
	// CustomerId	   stringGet
	// NumOfVerifiers int
	message := CampaignCryptoRequest{
		CamId:        id,
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

func RequestCampaignCryptoParams(id string, numVerifiers int) CampaignCryptoParams {
	var cryptoParams CampaignCryptoParams

	c := &http.Client{}

	message := CampaignCryptoRequest{
		CamId:        id,
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

func computeCommitment3(campId string, rEnc string, url string) (*string, error) {
	SendLog("computeCommitment3 at", url)
	client := &http.Client{}

	message := VerificationRequest{
		CamId: campId,
		R:     rEnc,
	}

	jsonData, err := json.Marshal(message)
	request := string(jsonData)

	SendLog("request", request)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(request))

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	respJSON, err := client.Do(reqData)
	// SendLog("respJSON", *respJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return nil, err
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	SendLog("data", string(data))
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	var verificationResponse VerificationResponse
	err = json.Unmarshal([]byte(data), &verificationResponse)
	if err != nil {
		println(err)
	}

	SendLog("verificationResponse.H:", verificationResponse.H)
	SendLog("verificationResponse.s:", verificationResponse.S)
	SendLog("verificationResponse.r:", verificationResponse.R)
	SendLog("verificationResponse.Comm:", verificationResponse.Comm)

	return &verificationResponse.Comm, nil
}

func computeCommitment2(campId string, userId, url string) (*CampaignCustomerVerifierProof, error) {
	SendLog("request calculate proof verifier crypto at", url)

	client := &http.Client{}
	// requestArgs := NewVerifierCryptoParamsRequest{
	// 	CamId: camId,
	// 	H:     cryptoParams.H,
	// }

	// jsonArgs, err := json.Marshal(requestArgs)
	// request := string(jsonArgs)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(""))
	// sendLog("request", request)
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
	SendLog("data", string(data))
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	fmt.Println("return data all:", string(data))
	var subProof CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(data), &subProof)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	return &subProof, nil
}

func computeCommitment(campID string, url string, i int, cryptoParams CampaignCryptoParams) string {
	//connect to verifier: campID,  H , r
	SendLog("Start of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")
	// var param CommRequest
	var rBytes []byte
	var rEnc, hEnc string
	client := &http.Client{}
	// sendLog("create connection of commCompute:", "")

	hBytes := cryptoParams.H
	hEnc = b64.StdEncoding.EncodeToString(hBytes)
	SendLog("Encode H: ", hEnc)

	// if url == com1URL {
	// 	rBytes = camParam.R1
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)

	// } else if url == com2URL {
	// 	rBytes = camParam.R2
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	// }

	SendLog("num r values: ", string(len(cryptoParams.R1)))
	rBytes = cryptoParams.R1[i]
	// sendLog("R["+string(i)+"]: ", rBytes)
	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	SendLog("Encode R["+string(i)+"]: ", rEnc)

	// jsonData, _ := json.Marshal(param)
	message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", campID, hEnc, rEnc)
	// request := string(jsonData)

	SendLog("commCompute.message", message)

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

	SendLog("commValue:", string(data))
	SendLog("end of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

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
