package main

import (
	"bufio"
	"context"
	// "crypto/rsa"
	b64 "encoding/base64"
	"encoding/json"
	"errors"

	// "errors"
	"fmt"
	// "log"
	"io/ioutil"
	"net"
	"os"

	putils "internal/promark_utils"

	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
	eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"

	ristretto "github.com/bwesterb/go-ristretto"
	redis "github.com/go-redis/redis/v8"

	// "strings"
	// "log"
	"path/filepath"
)

var ctx = context.Background()
var LOG_MODE = "debug"
var secretFileName string
var H ristretto.Point
var HEnc string
var s ristretto.Scalar
var SEnc string
var cryptoServiceSocketURL = "external.promark.com:5000"

func main() {
	port := os.Getenv("API_PORT")
	name := os.Getenv("CORE_PEER_ID")
	secretFileName = filepath.Join("certs", fmt.Sprintf("%s.mycert", name))

	err := initialize()
	if err != nil {
		fmt.Printf("Cannot generate secret number and get H:%s\n", err)
	}
	fmt.Printf("verifier:%s - H:%s - s:%s\n", name, H, s)

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

func getHFromCryptoService() (*putils.CampaignCryptoParams, error) {
	conn, err := net.Dial("tcp", cryptoServiceSocketURL)
	if err != nil {
		// sendLog("Error connecting:", err.Error())

		fmt.Println("Error connecting:" + err.Error())
		return nil, errors.New("Error connecting:" + err.Error())
	}

	putils.SendRequest(conn, "get", "")
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		fmt.Println("Error in getting H:", err.Error())
		return nil, errors.New("Error in getting H:" + err.Error())
	}
	fmt.Println("Reiceived From: " + cryptoServiceSocketURL + "-Response:" + responseStr)

	response, err := putils.ParseResponse(responseStr)

	if err != nil {
		return nil, err
	}
	var cryptoParams putils.CampaignCryptoParams
	err = json.Unmarshal([]byte(response.Data), &cryptoParams)

	if err != nil {
		return nil, err
	}
	return &cryptoParams, nil
}

func initialize() error {
	// check if the file storing H is existed
	data, err := ioutil.ReadFile(secretFileName)
	if err != nil {
		// generate  and store in a file
		// var s ristretto.Scalar
		s.Rand()
		sBytes := s.Bytes()
		sEnc := b64.StdEncoding.EncodeToString(sBytes)

		f, err := os.Create(secretFileName)
		if err != nil {
			return err
		}

		f.WriteString(sEnc)
		f.Close()
	} else {

		sp, err := eutils.ConvertStringToScalar(string(data))
		if err != nil {
			return err
		}
		s = *sp
	}

	// contact cryptoService for H
	for {
		cryptoParams, err := getHFromCryptoService()

		if err != nil {
			fmt.Println(err)
			continue
		}

		HEnc = cryptoParams.H
		Hp, err := eutils.ConvertStringToPoint(HEnc)

		if err != nil {
			return err
		}

		H = *Hp
		break
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			putils.SendResponse(conn, fmt.Sprintf("%s", r), "")
		}
	}()

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
	putils.SendLog(os.Getenv("CORE_PEER_ID")+"command", request.Command, LOG_MODE)
	putils.SendLog(os.Getenv("CORE_PEER_ID")+"requestData", request.Data, LOG_MODE)

	if request.Command == "create" {
		// fmt.Println("ok:", string(command))
		CreateVerifierCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "get" {
		// fmt.Println("ok:", string(buffer))
		GetCampaignCryptoParamsHandler(conn, request.Data)
	} else if request.Command == "commit" {
		CalculateCommitmentHandler(conn, request.Data)
	} else if request.Command == "commit-nosave" {
		CalculateCommitmentNoSaveHandler(conn, request.Data)
	} else if request.Command == "genpoc" {
		GenerateSubPoCProof(conn, request.Data)
	} else if request.Command == "commit-nocam" {
		CalculateCommitmentNoCamHandler(conn, request.Data)
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

	val, err := client.Get(ctx, paramsRequest.CamId).Result()
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
		err = client.Set(ctx, cryptoParams.CamId, jsonParam, 0).Err()
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

	val, err := client.Get(ctx, camId).Result()

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

	err = SetCustomerCampaignProofTransaction(paramsRequest.CamId, paramsRequest.CustomerId, vCryptoParams)

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}
	subProof, err := GetCustomerCampaignProof(paramsRequest.CamId, paramsRequest.CustomerId)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Println("Calculated Proof - Comm:" + subProof.Comm)
	// response := putils.CampaignCustomerVerifierProof{}
	responseData, err := json.Marshal(subProof)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	putils.SendResponse(conn, "", string(responseData))
}

func GenerateSubPoCProof(conn net.Conn, requestData string) {
	// calculate commitment
	subPoCProof, err := GeneratePoCProof()

	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Printf("Calculated Proof - Comm:%s", subPoCProof.Comm)
	// response := putils.CampaignCustomerVerifierProof{}
	responseData, err := json.Marshal(subPoCProof)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	putils.SendResponse(conn, "", string(responseData))
}

func CalculateCommitmentNoCamHandler(conn net.Conn, requestData string) {
	// calculate commitment
	var paramsRequest putils.VerificationRequest
	err := json.Unmarshal([]byte(requestData), &paramsRequest)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Printf("Got - paramsRequest.CamId:%s\n", paramsRequest.CamId)

	commEnc, err := CalculateCommitment(paramsRequest.R)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
		return
	}

	fmt.Printf("Calculated Proof - Comm:%s", commEnc)
	response := putils.VerificationResponse{
		CamId: paramsRequest.CamId,
		S:     SEnc,
		R:     paramsRequest.R,
		Comm:  commEnc,
		H:     HEnc,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	putils.SendResponse(conn, "", string(responseData))
}

func CalculateCommitmentNoSaveHandler(conn net.Conn, requestData string) {
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

	subProof, err := GenerateRAndCalculateCommitment(paramsRequest.CamId, paramsRequest.CustomerId, vCryptoParams)
	// err = SetCustomerCampaignProofTransaction(paramsRequest.CamId, paramsRequest.CustomerId, *vCryptoParams)

	// if err != nil {
	// 	putils.SendResponse(conn, err.Error(), "")
	// }
	// subProof, err := GetCustomerCampaignProof(paramsRequest.CamId, paramsRequest.CustomerId)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	fmt.Printf("Calculated Proof - Comm:%s", subProof.Comm)
	// response := putils.CampaignCustomerVerifierProof{}
	responseData, err := json.Marshal(subProof)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	putils.SendResponse(conn, "", string(responseData))
}

func GetCustomerCampaignProof(camId string, userId string) (*putils.CampaignCustomerVerifierProof, error) {
	client := putils.GetRedisConnection()
	proofId := camId + ":" + userId
	subProofJSON, err := client.Get(ctx, proofId).Result()

	if err != nil {
		// proof is not existed
		return nil, err
	}
	var subProof putils.CampaignCustomerVerifierProof
	err = json.Unmarshal([]byte(subProofJSON), &subProof)

	if err != nil {
		return nil, err
	}

	return &subProof, nil
}

func GeneratePoCProof() (*putils.PoCProof, error) {
	// generate r
	var rScalar ristretto.Scalar
	rScalar.Rand()
	fmt.Printf("Generated R: %s\n", rScalar)

	// calculate commitment
	comm := pedersen.CommitTo(&H, &rScalar, &s)
	fmt.Printf("Calculated comm: %s\n", comm)

	commEnc := eutils.ConvertPointToString(comm)
	rEnc := eutils.ConvertScalarToString(&rScalar)

	poc := putils.PoCProof{
		Comm: commEnc,
		R:    rEnc,
	}

	return &poc, nil
}

func CalculateCommitment(rEnc string) (string, error) {
	r, err := eutils.ConvertStringToScalar(rEnc)
	if err != nil {
		return "", err
	}

	comm := pedersen.CommitTo(&H, r, &s)
	cEnc := eutils.ConvertPointToString(comm)

	return cEnc, nil
}

func GenerateRAndCalculateCommitment(camId string, userId string, vCryptoParams *putils.VerifierCryptoParams) (*putils.CampaignCustomerVerifierProof, error) {
	var rScalar ristretto.Scalar
	rScalar.Rand()
	fmt.Printf("Generated R: %s\n", rScalar)

	// convert s
	sScalar, err := eutils.ConvertStringToScalar(vCryptoParams.S)
	if err != nil {
		return nil, err
	}
	fmt.Printf("s: %s\n", sScalar)

	// convert H
	hPoint, err := eutils.ConvertStringToPoint(vCryptoParams.H)
	if err != nil {
		return nil, err
	}
	fmt.Printf("H: %s\n", hPoint)

	// calculate commitment
	comm := pedersen.CommitTo(hPoint, &rScalar, sScalar)

	fmt.Printf("Calculated comm: %s\n", comm)

	// return R, Com
	rEnc := eutils.ConvertScalarToString(&rScalar)
	commEnc := eutils.ConvertPointToString(comm)

	subProof := putils.CampaignCustomerVerifierProof{
		CamId:  camId,
		UserId: userId,
		H:      vCryptoParams.H,
		R:      rEnc,
		S:      vCryptoParams.S,
		Comm:   commEnc,
	}

	return &subProof, nil
}

func SetCustomerCampaignProofTransaction(camId string, userId string, vCryptoParams *putils.VerifierCryptoParams) error {
	fmt.Println(os.Getenv("CORE_PEER_ID") + ":SetCustomerCampaignProofTransaction:" + camId + ":" + userId)
	putils.SendLog(os.Getenv("CORE_PEER_ID")+"SetCustomerCampaignProofTransaction", camId+":"+userId, LOG_MODE)

	maxRetries := 3
	proofId := camId + ":" + userId

	// var subProof putils.CampaignCustomerVerifierProof

	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		// val, err := rdb.Get(ctx, "key").Result()
		_, err := tx.Get(ctx, proofId).Result()
		if err != nil && err != redis.Nil {
			// proof is not existed
			return err
		}

		// Actual operation (local in optimistic lock).
		if err == redis.Nil {
			subProof, err := GenerateRAndCalculateCommitment(camId, userId, vCryptoParams)
			// var rScalar ristretto.Scalar
			// rScalar.Rand()
			// fmt.Println("Generated random scalar")

			// // convert H
			// hPoint := putils.ConvertStringToPoint(vCryptoParams.H)

			// // convert s
			// sScalar := putils.ConvertStringToScalar(vCryptoParams.S)

			// // calculate commitment
			// comm := putils.CommitTo(&hPoint, &rScalar, &sScalar)

			// fmt.Println("Calculated comm")

			// // return R, Com
			// rEnc := putils.ConvertScalarToString(rScalar)
			// commEnc := putils.ConvertPointToString(comm)

			// subProof := putils.CampaignCustomerVerifierProof{
			// 	CamId:  camId,
			// 	UserId: userId,
			// 	H:      vCryptoParams.H,
			// 	R:      rEnc,
			// 	S:      vCryptoParams.S,
			// 	Comm:   commEnc,
			// }

			fmt.Println("Initialized subproof")
			fmt.Println("subproof.CamId:" + subProof.CamId)

			jsonParam, err := json.Marshal(subProof)
			if err != nil {
				return err
			}

			fmt.Println("Converted subproof to JSON:" + string(jsonParam))

			// Operation is commited only if the watched keys remain unchanged.
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, proofId, jsonParam, 0)

				fmt.Println("Added subproof to db")

				return nil
			})

			return err
		}

		return nil
	}

	// Retry if the key has been changed.
	client := putils.GetRedisConnection()
	for i := 0; i < maxRetries; i++ {
		err := client.Watch(ctx, txf, proofId)
		if err == nil {
			// Success.
			putils.SendLog(os.Getenv("CORE_PEER_ID")+"Success", "", LOG_MODE)

			// subProofJSON, _ := client.Get(ctx, proofId).Result()
			// subProof
			// json.Unmarshal([]byte(subProofJSON), &subProof)
			// // // json.Unmarshal(client.Get(ctx, proofId), &subProof)
			// fmt.Println("Updated subproof successfully")
			// fmt.Println("subproof.comm:" + subProof.Comm)
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			fmt.Println("There are some modifications on subproof")
			continue
		}
		// Return any other error.
		return err
	}

	return errors.New("SetCustomerCampaignProofTransaction reached maximum number of retries")
}

func SetCustomerCampaignProof2(camId string, userId string, vCryptoParams putils.VerifierCryptoParams) (*putils.CampaignCustomerVerifierProof, error) {
	// generate random values

	var subProof putils.CampaignCustomerVerifierProof
	proofId := camId + ":" + userId
	client := putils.GetRedisConnection()
	val, err := client.Get(ctx, proofId).Result()
	err = json.Unmarshal([]byte(val), &subProof)

	if err != nil {
		// params are not existed
		fmt.Println(err)
		// generate R
		var rScalar *ristretto.Scalar
		rScalar.Rand()

		// convert H
		hPoint, err := eutils.ConvertStringToPoint(vCryptoParams.H)
		if err != nil {
			return nil, err
		}

		// convert s
		sScalar, err := eutils.ConvertStringToScalar(vCryptoParams.S)
		if err != nil {
			return nil, err
		}

		// calculate commitment
		comm := pedersen.CommitTo(hPoint, rScalar, sScalar)

		// return R, Com
		rEnc := eutils.ConvertScalarToString(rScalar)
		commEnc := eutils.ConvertPointToString(comm)

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

		err = client.Set(ctx, proofId, jsonParam, 0).Err()
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Printf("The CampaignCustomerVerifierProof is existed for id %s", proofId)
		return &subProof, fmt.Errorf("subProof with id %s is existed", proofId)
	}

	return &subProof, nil
}

func SetCustomerCampaignProof(camId string, userId string, vCryptoParams putils.VerifierCryptoParams) (*putils.CampaignCustomerVerifierProof, error) {
	// generate random values

	var subProof putils.CampaignCustomerVerifierProof
	proofId := camId + ":" + userId
	client := putils.GetRedisConnection()
	val, err := client.Get(ctx, proofId).Result()
	err = json.Unmarshal([]byte(val), &subProof)

	if err != nil {
		// params are not existed
		fmt.Println(err)
		// generate R
		var rScalar *ristretto.Scalar
		rScalar.Rand()

		// convert H
		hPoint, err := eutils.ConvertStringToPoint(vCryptoParams.H)
		if err != nil {
			return nil, err
		}

		// convert s
		sScalar, err := eutils.ConvertStringToScalar(vCryptoParams.S)
		if err != nil {
			return nil, err
		}

		// calculate commitment
		comm := pedersen.CommitTo(hPoint, rScalar, sScalar)

		// return R, Com
		rEnc := eutils.ConvertScalarToString(rScalar)
		commEnc := eutils.ConvertPointToString(comm)

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

		err = client.Set(ctx, proofId, jsonParam, 0).Err()
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
	}

	fmt.Println("vCryptoParams.CamId:" + vCryptoParams.CamId)
	fmt.Println("vCryptoParams.H:" + vCryptoParams.H)
	fmt.Println("vCryptoParams.S:" + vCryptoParams.S)

	// convert H
	hPoint, err := eutils.ConvertStringToPoint(vCryptoParams.H)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}
	// hDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.H)
	// hPoint := convertBytesToPoint(hDec)

	// convert s
	sScalar, err := eutils.ConvertStringToScalar(vCryptoParams.S)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}
	// sDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.S)
	// sScalar := convertBytesToScalar(sDec)

	// convert r
	rScalar, err := eutils.ConvertStringToScalar(paramsRequest.R)
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}
	// rDec, _ := b64.StdEncoding.DecodeString(request.R)
	// rScalar := convertBytesToScalar(rDec)

	// calculate commitment

	comm := pedersen.CommitTo(hPoint, rScalar, sScalar)
	commEnc := eutils.ConvertPointToString(comm)
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
	if err != nil {
		putils.SendResponse(conn, err.Error(), "")
	}

	putils.SendResponse(conn, "", string(responseJSON))
}
