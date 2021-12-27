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

	"github.com/bwesterb/go-ristretto"
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

	campaign := putils.Campaign{
		Id:           "c54",
		Name:         "Campaign001",
		Advertiser:   "adv0",
		Business:     "bus0",
		CommC:        "",
		VerifierURLs: []string{"peer0.adv0.promark.com:5000", "peer0.bus0.promark.com:5000"},
	}
	proof, _ := GenerateCustomerCampaignProofSocket(&campaign, "u92")

	collectedProof := putils.CollectedCustomerProof{
		Id:   "p001",
		Comm: proof.Comm,
		Rs:   proof.Rs,
	}
	VerifyCommitmentSocket(&campaign, &collectedProof)
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
		Id:           camId,
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

	results <- putils.VerifierProofChannelResult{
		URL:   url,
		Proof: *verifierProof,
		Error: err,
	}
}

func RequestCommitment(camId string, customerId string, url string) (*putils.CampaignCustomerVerifierProof, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

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

	fmt.Println("Returned-vCryptoParams.CamId:" + subProof.CamId)
	return &subProof, nil
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
