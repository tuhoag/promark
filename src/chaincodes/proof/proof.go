package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"

	// "errors"
	"fmt"

	// "log"

	// "strconv"
	"strings"

	putils "internal/promark_utils"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
)

var LOG_MODE = "debug"

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

	proof, err := GenerateCustomerCampaignProofSocket(&campaign, userId)

	return proof, nil
}

func GenerateCustomerCampaignProofSocket(campaign *putils.Campaign, userId string) (*putils.ProofCustomerCampaign, error) {
	// generate a random values for each verifiers
	numVerifiers := len(campaign.VerifierURLs)
	// putils.SendLog("numVerifiers", string(numVerifiers), LOG_MODE)

	// get crypto params
	// cryptoParams := requestCustomerCampaignCryptoParams(camId, userId, numVerifiers)

	var Ci, C ristretto.Point
	vProofChannel := make(chan putils.VerifierProofChannelResult)

	wg := sync.WaitGroup{}
	wg.Add(numVerifiers)

	for i := 0; i < numVerifiers; i++ {
		verifierURL := campaign.VerifierURLs[i]
		// putils.SendLog("verifierURL", verifierURL)
		// putils.SendLog("comURL", requestCreateVerifierCryptoURL)

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

		putils.SendLog("result.URL:"+result.URL+"H:"+result.Proof.H+"-R:"+result.Proof.R+"-S:"+result.Proof.S+"-Comm:"+result.Proof.Comm, "", LOG_MODE)
		// putils.SendLog("result.Error", result.Error.Error(), LOG_MODE)
		// putils.SendLog("result.Proof.H", result.Proof.H, LOG_MODE)
		// putils.SendLog("result.Proof.R", result.Proof.R, LOG_MODE)
		// putils.SendLog("result.Proof.S", result.Proof.S, LOG_MODE)
		// putils.SendLog("result.Proof.Comm", result.Proof.Comm, LOG_MODE)

		if result.Error != nil {
			return nil, result.Error
		}

		Ci = putils.ConvertStringToPoint(result.Proof.Comm)
		// CiBytes, _ := b64.StdEncoding.DecodeString(subProof.Comm)
		// Ci = convertBytesToPoint(CiBytes)

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)

			// putils.SendLog("Current Comm", putils.ConvertPointToString(C), DEBUG_LOG)
		}

		randomValues = append(randomValues, result.Proof.R)
		subComs = append(subComs, result.Proof.Comm)
	}

	close(vProofChannel)

	CommEnc := putils.ConvertPointToString(C)

	proof := putils.ProofCustomerCampaign{
		Comm:    CommEnc,
		Rs:      randomValues,
		SubComs: subComs,
	}

	fmt.Println("proof.Comm: " + proof.Comm)
	fmt.Printf("proof.Rs: %s\n", proof.Rs)
	fmt.Printf("proof.SubComs: %s\n", proof.SubComs)

	// s := fmt.Sprintf("%s is %d years old.\n", name, age)
	putils.SendLog(fmt.Sprintf("proof.Comm:%s- Rs:%s-SubComs:%s", proof.Comm, proof.Rs, proof.SubComs), "", LOG_MODE)

	return &proof, nil
}

func ConcurrentRequestCommitment(camId string, customerId string, url string, results chan putils.VerifierProofChannelResult, wg *sync.WaitGroup) {
	verifierProof, err := RequestCommitment(camId, customerId, url)

	wg.Done()

	fmt.Println("Done with " + url)
	fmt.Println("vCryptoParams.CamId:" + verifierProof.CamId)
	fmt.Println("vCryptoParams.CustomerId:" + verifierProof.UserId)
	fmt.Println("vCryptoParams.H:" + verifierProof.H)
	fmt.Println("vCryptoParams.R:" + verifierProof.R)
	fmt.Println("vCryptoParams.S:" + verifierProof.S)

	putils.SendLog("Done with", url, LOG_MODE)
	putils.SendLog("vCryptoParams.CamId:", verifierProof.CamId, LOG_MODE)
	putils.SendLog("vCryptoParams.CustomerId:", verifierProof.UserId, LOG_MODE)
	putils.SendLog("vCryptoParams.S:", verifierProof.S, LOG_MODE)

	results <- putils.VerifierProofChannelResult{
		URL:   url,
		Proof: *verifierProof,
		Error: err,
	}
}

func RequestCommitment(camId string, customerId string, url string) (*putils.CampaignCustomerVerifierProof, error) {
	putils.SendLog("RequestCommitment at", url, LOG_MODE)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		putils.SendLog("Error connecting:", err.Error(), LOG_MODE)

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	requestArgs := putils.ProofGenerationRequest{
		CamId:      camId,
		CustomerId: customerId,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	putils.SendRequest(conn, "commit", string(jsonArgs))
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return nil, errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)
	putils.SendLog("Reiceived From: "+url+"-Response:", responseStr, LOG_MODE)

	response, err := putils.ParseResponse(responseStr)

	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}

	var subProof putils.CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(response.Data), &subProof)

	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	fmt.Println("Returned from " + url + "-subProof.CamId:" + subProof.CamId)

	return &subProof, nil
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
		Id:   proofId,
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
	// backupJSON, err := ctx.GetStub().GetState(backupId)
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

func (s *ProofSmartContract) DeleteProofById(ctx contractapi.TransactionContextInterface, proofId string) (bool, error) {
	proofJSON, err := ctx.GetStub().GetState(proofId)
	// backupJSON, err := ctx.GetStub().GetState(backupId)
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

	putils.SendLog("campaign.Id", campaign.Id, LOG_MODE)
	putils.SendLog("campaign.Name", campaign.Name, LOG_MODE)
	putils.SendLog("campaign.VerifierURLs", string(len(campaign.VerifierURLs)), LOG_MODE)

	// get proof
	proof, err := s.GetProofById(ctx, proofId)

	if err != nil {
		return false, err
	}

	// putils.SendLog("proof.H", proof.H)
	putils.SendLog("proof.Comm", proof.Comm, LOG_MODE)

	verificationResult, err := VerifyCommitmentSocket(&campaign, proof)

	return verificationResult, err
}

func VerifyCommitmentSocket(campaign *putils.Campaign, proof *putils.CollectedCustomerProof) (bool, error) {
	fmt.Printf("proof.Rs: %s\n", proof.Rs)
	numVerifiers := len(campaign.VerifierURLs)
	vcChannel := make(chan putils.VerifierCommitmentChannelResult)
	wg := sync.WaitGroup{}
	wg.Add(numVerifiers)

	// callinng verifiers to calculate proof.Comm again based on proof.Rs
	for i, verifierURL := range campaign.VerifierURLs {
		// call verifier to compute sub commitment
		go ConcurrentRequestVerification(campaign.Id, proof.Rs[i], verifierURL, vcChannel, &wg)
	}

	fmt.Println("Calculating total comm")
	// var subComs, randomValues []string
	var Ci, C ristretto.Point
	for i := 0; i < numVerifiers; i++ {
		result := <-vcChannel
		fmt.Printf("URL: %s\n", result.URL)
		// fmt.Printf("Error: %v\n" + result.Error)
		fmt.Printf("Comm: %s\n", result.Comm)

		if result.Error != nil {
			return false, result.Error
		}

		Ci = putils.ConvertStringToPoint(result.Comm)
		// CiBytes, _ := b64.StdEncoding.DecodeString(subProof.Comm)
		// Ci = convertBytesToPoint(CiBytes)

		if i == 0 {
			C = Ci
		} else {
			C.Add(&C, &Ci)

			// putils.SendLog("Current Comm", putils.ConvertPointToString(C), DEBUG_LOG)
		}

		// randomValues = append(randomValues, result.Proof.R)
		// subComs = append(subComs, result.Proof.Comm)
	}

	close(vcChannel)

	comm := putils.ConvertStringToPoint(proof.Comm)

	// putils.SendLog("proof.Com", proof.Comm, LOG_MODE)
	// putils.SendLog("calculated Com", b64.StdEncoding.EncodeToString(C.Bytes()), LOG_MODE)
	if C.Equals(&comm) {
		return true, nil
	} else {
		return false, nil
	}
}

func ConcurrentRequestVerification(camId string, r string, url string, results chan putils.VerifierCommitmentChannelResult, wg *sync.WaitGroup) {
	ci, err := RequestVerification(camId, r, url)

	wg.Done()

	fmt.Println("Done with " + url)
	fmt.Printf("ci: %s\n", *ci)

	results <- putils.VerifierCommitmentChannelResult{
		URL:   url,
		Comm:  *ci,
		Error: err,
	}
}

func RequestVerification(camId string, r string, url string) (*string, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("ERROR:" + err.Error())
	}

	requestArgs := putils.VerificationRequest{
		CamId: camId,
		R:     r,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	putils.SendRequest(conn, "verify", string(jsonArgs))
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

	var verificationResponse putils.VerificationResponse
	err = json.Unmarshal([]byte(response.Data), &verificationResponse)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	// putils.SendLog("verificationResponse.H:", verificationResponse.H, LOG_MODE)
	// putils.SendLog("verificationResponse.s:", verificationResponse.S, LOG_MODE)
	// putils.SendLog("verificationResponse.r:", verificationResponse.R, LOG_MODE)
	// putils.SendLog("verificationResponse.Comm:", verificationResponse.Comm, LOG_MODE)

	return &verificationResponse.Comm, nil
}

// func RequestCustomerCampaignCryptoParams(id string, userId string, numVerifiers int) putils.CampaignCryptoParams {
// 	var cryptoParams putils.CampaignCryptoParams

// 	c := &http.Client{}

// 	// Id             string
// 	// CustomerId	   stringGet
// 	// NumOfVerifiers int
// 	message := putils.CampaignCryptoRequest{
// 		CamId:        id,
// 		CustomerId:   userId,
// 		NumVerifiers: numVerifiers,
// 	}

// 	jsonData, err := json.Marshal(message)

// 	request := string(jsonData)

// 	reqJSON, err := http.NewRequest("POST", cryptoParamsRequestURL, strings.NewReader(request))
// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %v\n", err)
// 		return cryptoParams
// 	}

// 	respJSON, err := c.Do(reqJSON)
// 	if err != nil {
// 		fmt.Printf("http.Do() error: %v\n", err)
// 		return cryptoParams
// 	}
// 	defer respJSON.Body.Close()

// 	data, err := ioutil.ReadAll(respJSON.Body)
// 	if err != nil {
// 		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
// 		return cryptoParams
// 	}

// 	fmt.Println("return data all:", string(data))

// 	err = json.Unmarshal([]byte(data), &cryptoParams)
// 	if err != nil {
// 		println(err)
// 	}

// 	return cryptoParams
// }

// func ComputeCommitment3(campId string, rEnc string, url string) (*string, error) {
// 	putils.SendLog("computeCommitment3 at", url, LOG_MODE)
// 	client := &http.Client{}

// 	message := putils.VerificationRequest{
// 		CamId: campId,
// 		R:     rEnc,
// 	}

// 	jsonData, err := json.Marshal(message)
// 	request := string(jsonData)

// 	putils.SendLog("request", request, LOG_MODE)
// 	reqData, err := http.NewRequest("POST", url, strings.NewReader(request))

// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %v\n", err)
// 		return nil, err
// 	}

// 	respJSON, err := client.Do(reqData)
// 	// putils.SendLog("respJSON", *respJSON)
// 	if err != nil {
// 		fmt.Printf("http.Do() error: %v\n", err)
// 		return nil, err
// 	}
// 	defer respJSON.Body.Close()

// 	data, err := ioutil.ReadAll(respJSON.Body)
// 	putils.SendLog("data", string(data), LOG_MODE)
// 	if err != nil {
// 		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
// 		return nil, err
// 	}

// 	var verificationResponse putils.VerificationResponse
// 	err = json.Unmarshal([]byte(data), &verificationResponse)
// 	if err != nil {
// 		println(err)
// 	}

// 	putils.SendLog("verificationResponse.H:", verificationResponse.H, LOG_MODE)
// 	putils.SendLog("verificationResponse.s:", verificationResponse.S, LOG_MODE)
// 	putils.SendLog("verificationResponse.r:", verificationResponse.R, LOG_MODE)
// 	putils.SendLog("verificationResponse.Comm:", verificationResponse.Comm, LOG_MODE)

// 	return &verificationResponse.Comm, nil
// }

// func computeCommitment2(campId string, userId, url string) (*putils.CampaignCustomerVerifierProof, error) {
// 	putils.SendLog("request calculate proof verifier crypto at", url, LOG_MODE)

// 	client := &http.Client{}
// 	// requestArgs := NewVerifierCryptoParamsRequest{
// 	// 	CamId: camId,
// 	// 	H:     cryptoParams.H,
// 	// }

// 	// jsonArgs, err := json.Marshal(requestArgs)
// 	// request := string(jsonArgs)
// 	reqData, err := http.NewRequest("POST", url, strings.NewReader(""))
// 	// putils.SendLog("request", request)
// 	// putils.SendLog("err", err.Error())
// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %v\n", err)
// 		return nil, err
// 	}

// 	respJSON, err := client.Do(reqData)
// 	if err != nil {
// 		fmt.Printf("http.Do() error: %v\n", err)
// 		return nil, err
// 	}
// 	defer respJSON.Body.Close()

// 	data, err := ioutil.ReadAll(respJSON.Body)
// 	putils.SendLog("data", string(data), LOG_MODE)
// 	if err != nil {
// 		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
// 		return nil, err
// 	}

// 	fmt.Println("return data all:", string(data))
// 	var subProof putils.CampaignCustomerVerifierProof
// 	err = json.Unmarshal([]byte(data), &subProof)

// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %v\n", err)
// 		return nil, err
// 	}

// 	return &subProof, nil
// }
