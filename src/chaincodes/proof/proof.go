package main

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
	putils "internal/promark_utils"
)

var LOG_MODE = "test"

type ProofSmartContract struct {
	contractapi.Contract
}

var (
	cryptoServiceURL           = "http://external.promark.com:5000"
	cryptoParamsRequestURL     = cryptoServiceURL + "/camp"
	userCryptoParamsRequestURL = cryptoServiceURL + "/usercamp"
	logURL                     = "http://logs.promark.com:5003/log"
)

func (s *ProofSmartContract) GenerateCustomerCampaignProof(ctx contractapi.TransactionContextInterface, camId string, userId string) (*putils.ProofCustomerCampaign, error) {
	putils.SendLog("GenerateCustomerCampaignProof", "", LOG_MODE)
	putils.SendLog("camId:", camId, LOG_MODE)
	putils.SendLog("userId", userId, LOG_MODE)

	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	putils.SendLog("response.Payload", string(response.Payload), LOG_MODE)
	putils.SendLog("response.Status", string(response.Status), LOG_MODE)
	putils.SendLog("response.message", string(response.Message), LOG_MODE)
	// putils.SendLog("response.error", string(response.Error))
	// putils.SendLog("response.message is nil", strconv.FormatBool(response.Message == ""))

	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}

	var campaign putils.Campaign

	err := json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return nil, err
	}

	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)
	putils.SendLog("numVerifiers", string(numVerifiers), LOG_MODE)

	// get crypto params
	// cryptoParams := requestCustomerCampaignCryptoParams(camId, userId, numVerifiers)

	var Ci, C ristretto.Point

	var subComs, randomValues []string

	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]
		comURL := verifierURL + "/camp/" + camId + "/proof/" + userId
		putils.SendLog("verifierURL", verifierURL, LOG_MODE)
		putils.SendLog("comURL", comURL, LOG_MODE)

		// 	testVer(ver)
		// 	putils.SendLog("id", id)
		// 	putils.SendLog("Hvalue", string(cryptoParams.H))
		// 	putils.SendLog("R1value", string(cryptoParams.R1[i]))

		subProof, err := computeCommitment2(camId, userId, comURL)

		if err != nil {
			return nil, err
		}

		// commDec, _ := b64.StdEncoding.DecodeString(comm)
		// Ci = convertStringToPoint(string(commDec))
		putils.SendLog("H"+string(i)+" encoding:", subProof.H, LOG_MODE)
		putils.SendLog("S"+string(i)+" encoding:", subProof.S, LOG_MODE)
		putils.SendLog("R"+string(i)+" encoding:", subProof.R, LOG_MODE)
		putils.SendLog("Comm"+string(i)+" encoding:", subProof.Comm, LOG_MODE)
		Ci = putils.convertStringToPoint(subProof.Comm)
		// CiBytes, _ := b64.StdEncoding.DecodeString(subProof.Comm)
		// Ci = convertBytesToPoint(CiBytes)

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)

			putils.SendLog("Current Comm", b64.StdEncoding.EncodeToString(C.Bytes()))
		}

		randomValues = append(randomValues, subProof.R)
		subComs = append(subComs, subProof.Comm)
	}
	CommEnc := putils.ConvertPointToString(C)
	// CommBytes := C.Bytes()
	// CommEnc := b64.StdEncoding.EncodeToString(CommBytes)

	// get all verifiers URLs

	// calculate commitment
	proof := putils.ProofCustomerCampaign{
		Comm:    CommEnc,
		Rs:      randomValues,
		SubComs: subComs,
	}

	return &proof, nil
}

func (s *ProofSmartContract) GetAllProofs(ctx contractapi.TransactionContextInterface) ([]*putils.CollectedCustomerProof, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var proofs []*putils.CollectedCustomerProof
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		putils.SendLog("queryResponse.Value", string(queryResponse.Value), LOG_MODE)
		var proof putils.CollectedCustomerProof
		err = json.Unmarshal(queryResponse.Value, &proof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, &proof)
	}

	return proofs, nil
}

func (s *ProofSmartContract) AddCustomerProofCampaign(ctx contractapi.TransactionContextInterface, proofId string, comm string, rsStr string) error {
	putils.SendLog("AddCustomerProofCampaign", "", LOG_MODE)
	putils.SendLog("proofId", proofId, LOG_MODE)
	putils.SendLog("comm", comm, LOG_MODE)
	putils.SendLog("rsStr", rsStr, LOG_MODE)

	proofJSON, err := ctx.GetStub().GetState(proofId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if proofJSON != nil {
		return fmt.Errorf("the user proof raw id %s is existed", proofId)
	}

	rs := strings.Split(rsStr, ";")

	collectedProof := putils.CollectedCustomerProof{
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

func (s *ProofSmartContract) GetProofById(ctx contractapi.TransactionContextInterface, proofId string) (*putils.CollectedCustomerProof, error) {
	proofJSON, err := ctx.GetStub().GetState(proofId)
	// backupJSON, err := ctx.GetStub().GetState(backupID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if proofJSON == nil {
		return nil, fmt.Errorf("the campaign raw id %s does not exist", proofId)
	}

	var proof putils.CollectedCustomerProof
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

	var proof putils.CollectedCustomerProof
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
	putils.SendLog("VerifyCampaignProof", "", LOG_MODE)
	putils.SendLog("camId", camId, LOG_MODE)
	putils.SendLog("proofId", proofId, LOG_MODE)

	_, err := s.GetProofById(ctx, proofId)

	if err != nil {
		return false, err
	}

	// get campaign
	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	putils.SendLog("response.Payload", string(response.Payload), LOG_MODE)
	putils.SendLog("response.Status", string(response.Status), LOG_MODE)
	putils.SendLog("response.message", string(response.Message), LOG_MODE)

	if response.Message != "" {
		return false, fmt.Errorf(response.Message)
	}

	var campaign putils.Campaign

	err = json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return false, err
	}

	putils.SendLog("campaign.ID", campaign.ID, LOG_MODE)
	putils.SendLog("campaign.Name", campaign.Name, LOG_MODE)
	putils.SendLog("campaign.VerifierURLs", string(len(campaign.VerifierURLs)), LOG_MODE)

	// get proof
	proof, err := s.GetProofById(ctx, proofId)

	if err != nil {
		return false, err
	}

	// putils.SendLog("proof.H", proof.H)
	putils.SendLog("proof.Comm", proof.Comm)

	for i, R := range proof.Rs {
		putils.SendLog("proof.R["+string(i)+"]", R)
	}

	// callinng verifiers to calculate proof.Comm again based on proof.Rs
	var C ristretto.Point
	for i, verifierURL := range campaign.VerifierURLs {
		// call verifier to compute sub commitment
		comURL := verifierURL + "/camp/" + camId + "/verify"
		ciEnc, err := ComputeCommitment3(camId, proof.Rs[i], comURL)

		if err != nil {
			return false, err
		}

		ciBytes, err := b64.StdEncoding.DecodeString(*ciEnc)

		if err != nil {
			return false, err
		}

		ci := putils.ConvertBytesToPoint(ciBytes)

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

	putils.SendLog("proof.Com", proof.Comm, LOG_MODE)
	putils.SendLog("calculated Com", b64.StdEncoding.EncodeToString(C.Bytes()), LOG_MODE)
	comm := putils.ConvertBytesToPoint(CommBytes)
	if C.Equals(&comm) {
		return true, nil
	} else {
		return false, nil
	}
}

func RequestCustomerCampaignCryptoParams(id string, userId string, numVerifiers int) putils.CampaignCryptoParams {
	var cryptoParams putils.CampaignCryptoParams

	c := &http.Client{}

	// ID             string
	// CustomerId	   stringGet
	// NumOfVerifiers int
	message := putils.CampaignCryptoRequest{
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

func ComputeCommitment3(campId string, rEnc string, url string) (*string, error) {
	putils.SendLog("computeCommitment3 at", url)
	client := &http.Client{}

	message := putils.VerificationRequest{
		CamId: campId,
		R:     rEnc,
	}

	jsonData, err := json.Marshal(message)
	request := string(jsonData)

	putils.SendLog("request", request)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(request))

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	respJSON, err := client.Do(reqData)
	// putils.SendLog("respJSON", *respJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return nil, err
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	putils.SendLog("data", string(data))
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return nil, err
	}

	var verificationResponse putils.VerificationResponse
	err = json.Unmarshal([]byte(data), &verificationResponse)
	if err != nil {
		println(err)
	}

	putils.SendLog("verificationResponse.H:", verificationResponse.H, LOG_MODE)
	putils.SendLog("verificationResponse.s:", verificationResponse.S, LOG_MODE)
	putils.SendLog("verificationResponse.r:", verificationResponse.R, LOG_MODE)
	putils.SendLog("verificationResponse.Comm:", verificationResponse.Comm, LOG_MODE)

	return &verificationResponse.Comm, nil
}

func computeCommitment2(campId string, userId, url string) (*putils.CampaignCustomerVerifierProof, error) {
	putils.SendLog("request calculate proof verifier crypto at", url, LOG_MODE)

	client := &http.Client{}
	// requestArgs := NewVerifierCryptoParamsRequest{
	// 	CamId: camId,
	// 	H:     cryptoParams.H,
	// }

	// jsonArgs, err := json.Marshal(requestArgs)
	// request := string(jsonArgs)
	reqData, err := http.NewRequest("POST", url, strings.NewReader(""))
	// putils.SendLog("request", request)
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
	var subProof putils.CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(data), &subProof)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	return &subProof, nil
}

func computeCommitment(campID string, url string, i int, cryptoParams putils.CampaignCryptoParams) string {
	//connect to verifier: campID,  H , r
	putils.SendLog("Start of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"", LOG_MODE)
	// var param CommRequest
	var rBytes []byte
	var rEnc, hEnc string
	client := &http.Client{}
	// putils.SendLog("create connection of commCompute:", "")

	hBytes := cryptoParams.H
	hEnc = b64.StdEncoding.EncodeToString(hBytes)
	putils.SendLog("Encode H: ", hEnc, LOG_MODE)

	// if url == com1URL {
	// 	rBytes = camParam.R1
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)

	// } else if url == com2URL {
	// 	rBytes = camParam.R2
	// 	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	// }

	putils.SendLog("num r values: ", string(len(cryptoParams.R1)), LOG_MODE)
	rBytes = cryptoParams.R1[i]
	// putils.SendLog("R["+string(i)+"]: ", rBytes)
	rEnc = b64.StdEncoding.EncodeToString(rBytes)
	putils.SendLog("Encode R["+string(i)+"]: ", rEnc, LOG_MODE)

	// jsonData, _ := json.Marshal(param)
	message := fmt.Sprintf("{\"id\": \"%s\", \"H\": \"%s\", \"r\": \"%s\"}", campID, hEnc, rEnc)
	// request := string(jsonData)

	putils.SendLog("commCompute.message", message)

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

	putils.SendLog("commValue:", string(data), LOG_MODE)
	putils.SendLog("end of commCompute:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"", LOG_MODE)

	return string(data)
}

func commVerify(campID string, url string, r string) string {
	//connect to verifier: campID,  H , r
	// putils.SendLog("Start of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	// var hEnc string
	c := &http.Client{}

	// hBytes := camParam.H
	// hEnc = b64.StdEncoding.EncodeToString(hBytes)

	// putils.SendLog("commVerify.r in string", r)

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

	// putils.SendLog("commValue:", string(data))
	// putils.SendLog("end of commVerify:", "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"")

	return string(data)
}
