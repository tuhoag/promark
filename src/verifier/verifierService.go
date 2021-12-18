package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	"errors"

	// "errors"
	"fmt"
	// "log"
	"net"
	"os"

	putils "internal/promark_utils"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
	// "strings"
	// "log"
)

func main() {
	port := os.Getenv("API_PORT")
	name := os.Getenv("CORE_PEER_ID")
	// port := "5002"

	fmt.Println("Starting 'tcp' server on " + name + ":" + port)
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		fmt.Println("Client connected.")
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	fmt.Println("request:" + buffer)
	var request putils.PromarkRequest
	err = json.Unmarshal([]byte(buffer), &request)

	if err != nil {
		fmt.Println("ERROR:" + err.Error())
		return
	}

	fmt.Println("command:" + request.Command)
	fmt.Println("requestData:" + request.Data)

	if request.Command == "create" {
		// fmt.Println("ok:", string(command))
		CreateVerifierCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "get" {
		// fmt.Println("ok:", string(buffer))
		GetCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "commit" {
		CalculateCommitmentHandler(conn, request.Data)
	} else if request.Command == "verify" {
		VerifyCommitmentHandler(conn, request.Data)
	} else {
		putils.SendResponse(conn, "nocommand", "")
		conn.Close()
	}

	conn.Close()
}

func CreateVerifierCampaignCryptoParamsHandler(conn net.Conn, requestData string) {
	var paramsRequest putils.VerifierCryptoParamsRequest
	err := json.Unmarshal([]byte(requestData), &paramsRequest)

	_, err = SetVerifierCryptoParams(paramsRequest)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	vCryptoParams, err := GetVerifierCryptoParams(paramsRequest.CamId)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Println("Got - vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("Got - vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("Got - vCryptoParams.S:" + vCryptoParams.S)

	param, err := json.Marshal(&vCryptoParams)
	putils.SendResponse(conn, "", string(param))
}

func SetVerifierCryptoParams(paramsRequest putils.VerifierCryptoParamsRequest) (bool, error) {
	client := putils.GetRedisConnection()

	var cryptoParams putils.VerifierCryptoParams

	val, err := client.Get(paramsRequest.CamId).Result()
	err = json.Unmarshal([]byte(val), &cryptoParams)
	var s ristretto.Scalar
	if err != nil {
		// params are not existed
		fmt.Println(err)
		s.Rand()
		sBytes := s.Bytes()
		sEnc := b64.StdEncoding.EncodeToString(sBytes)

		cryptoParams = putils.VerifierCryptoParams{
			CamId: paramsRequest.CamId,
			H:     paramsRequest.H,
			S:     sEnc,
		}

		fmt.Println("Store-vCryptoParams.CamId: " + cryptoParams.CamId)
		fmt.Println("Store-vCryptoParams.H: " + cryptoParams.H)
		fmt.Println("Store-vCryptoParams.S: " + cryptoParams.S)

		jsonParam, err := json.Marshal(cryptoParams)
		if err != nil {
			return false, err
		}

		fmt.Println("Store-vCryptoParams JSON: " + string(jsonParam))
		err = client.Set(cryptoParams.CamId, jsonParam, 0).Err()
		if err != nil {
			return false, err
		}

	} else {
		fmt.Printf("The VerifierCryptoParams is existed for id %s", cryptoParams.CamId)
		return false, nil
	}

	return true, nil
}

func GetVerifierCryptoParams(camId string) (*putils.VerifierCryptoParams, error) {
	client := putils.GetRedisConnection()

	var cryptoParams putils.VerifierCryptoParams

	val, err := client.Get(camId).Result()

	if err == redis.Nil {
		fmt.Println("camId " + camId + " does not exist")
		return nil, errors.New("camId " + camId + " does not exist")
	} else if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &cryptoParams)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		// f.WriteString("ERROR: " + err.Error())
		return nil, err
	}

	return &cryptoParams, nil
}

func GetCampaignCryptoParamsHandler(conn net.Conn, requestData string) {

}

func CalculateCommitmentHandler(conn net.Conn, requestData string) {
	// calculate commitment
	var paramsRequest putils.ProofGenerationRequest
	err := json.Unmarshal([]byte(requestData), &paramsRequest)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	vCryptoParams, err := GetVerifierCryptoParams(paramsRequest.CamId)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	if vCryptoParams == nil {
		putils.SendResponse(conn, "vCryptoParams is not existed", "")
		return
	}

	fmt.Println("Got - vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("Got - vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("Got - vCryptoParams.S:" + vCryptoParams.S)

	subProof, err := SetCustomerCampaignProof(paramsRequest.CamId, paramsRequest.CustomerId, *vCryptoParams)

	if err != nil && subProof == nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	param, err := json.Marshal(&subProof)
	putils.SendResponse(conn, "", string(param))
}

func SetCustomerCampaignProof(camId string, userId string, vCryptoParams putils.VerifierCryptoParams) (*putils.CampaignCustomerVerifierProof, error) {
	var subProof putils.CampaignCustomerVerifierProof
	proofId := camId + ":" + userId
	client := putils.GetRedisConnection()
	val, err := client.Get(proofId).Result()
	err = json.Unmarshal([]byte(val), &subProof)

	if err != nil {
		// params are not existed
		fmt.Println(err)
		// generate R
		var rScalar ristretto.Scalar
		rScalar.Rand()

		// convert H
		hPoint := putils.ConvertStringToPoint(vCryptoParams.H)

		// convert s
		sScalar := putils.ConvertStringToScalar(vCryptoParams.S)

		// calculate commitment
		comm := putils.CommitTo(&hPoint, &rScalar, &sScalar)

		// return R, Com
		rEnc := putils.ConvertScalarToString(rScalar)
		commEnc := putils.ConvertPointToString(comm)

		subProof = putils.CampaignCustomerVerifierProof{
			CamId:  camId,
			UserId: userId,
			H:      vCryptoParams.H,
			R:      rEnc,
			S:      vCryptoParams.S,
			Comm:   commEnc,
		}

		jsonParam, err := json.Marshal(subProof)
		if err != nil {
			return nil, err
		}

		err = client.Set(proofId, jsonParam, 0).Err()
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Printf("The CampaignCustomerVerifierProof is existed for id %s", proofId)
		return &subProof, fmt.Errorf("subProof with id %s is existed", proofId)
	}

	return &subProof, nil
}

func VerifyCommitmentHandler(conn net.Conn, requestData string) {
	// receive camId, r
	var paramsRequest putils.VerificationRequest
	err := json.Unmarshal([]byte(requestData), &paramsRequest)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Printf("Got - VerificationRequest.CamId: %s\n", paramsRequest.CamId)
	fmt.Printf("Got - VerificationRequest.R: %s\n", paramsRequest.R)

	vCryptoParams, err := GetVerifierCryptoParams(paramsRequest.CamId)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	if vCryptoParams == nil {
		putils.SendResponse(conn, "vCryptoParams is not existed", "")
		return
	}

	fmt.Println("vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("vCryptoParams.S:" + vCryptoParams.S)

	// convert H
	hPoint := putils.ConvertStringToPoint(vCryptoParams.H)
	// hDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.H)
	// hPoint := convertBytesToPoint(hDec)

	// convert s
	sScalar := putils.ConvertStringToScalar(vCryptoParams.S)
	// sDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.S)
	// sScalar := convertBytesToScalar(sDec)

	// convert r
	rScalar := putils.ConvertStringToScalar(paramsRequest.R)
	// rDec, _ := b64.StdEncoding.DecodeString(request.R)
	// rScalar := convertBytesToScalar(rDec)

	// calculate commitment

	comm := putils.CommitTo(&hPoint, &rScalar, &sScalar)
	commEnc := putils.ConvertPointToString(comm)
	// commEnc := b64.StdEncoding.EncodeToString(comm.Bytes())

	response := putils.VerificationResponse{
		CamId: paramsRequest.CamId,
		S:     vCryptoParams.S,
		R:     paramsRequest.R,
		Comm:  commEnc,
		H:     vCryptoParams.H,
	}

	fmt.Println("response.Comm:" + response.Comm)

	responseJSON, err := json.Marshal(&response)
	putils.SendResponse(conn, "", string(responseJSON))
}
